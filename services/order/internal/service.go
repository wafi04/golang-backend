package handler

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/wafi04/golang-backend/grpc/pb"
	"github.com/wafi04/golang-backend/services/common"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	lockTimeout   = 10 * time.Second
	retryTimeout  = 100 * time.Millisecond
	maxRetries    = 50
	lockKeyPrefix = "order-lock:"
)

type OrderRepository struct {
	db          *sqlx.DB
	redisClient *redis.Client
}

func NewOrderService(db *sqlx.DB, redisClient *redis.Client) *OrderRepository {
	return &OrderRepository{
		db:          db,
		redisClient: redisClient,
	}
}

func (s *OrderRepository) acquireLock(ctx context.Context, variantID string) (bool, error) {
	lockKey := fmt.Sprintf("%s%s", lockKeyPrefix, variantID)

	for i := 0; i < maxRetries; i++ {
		success, err := s.redisClient.SetNX(ctx, lockKey, "locked", lockTimeout).Result()
		if err != nil {
			return false, fmt.Errorf("error acquiring lock: %v", err)
		}

		if success {
			return true, nil
		}

		time.Sleep(retryTimeout)
	}

	return false, fmt.Errorf("timeout acquiring lock for variant %s", variantID)
}

func (s *OrderRepository) releaseLock(ctx context.Context, variantID string) error {
	lockKey := fmt.Sprintf("%s%s", lockKeyPrefix, variantID)
	return s.redisClient.Del(ctx, lockKey).Err()
}

func (s *OrderRepository) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.Order, error) {
	log.Printf("Incoming order request for variant %s", req.VariantsId)

	locked, err := s.acquireLock(ctx, req.VariantsId)
	if err != nil || !locked {
		return nil, fmt.Errorf("failed to acquire lock: %v", err)
	}
	defer func() {
		if unlockErr := s.releaseLock(ctx, req.VariantsId); unlockErr != nil {
			log.Printf("failed to release lock: %v", unlockErr)
		}
	}()

	var result *pb.Order
	var orderErr error

	defer func() {
		if unlockErr := s.releaseLock(ctx, req.VariantsId); unlockErr != nil {
			log.Printf("Warning: Failed to release lock for variant %s: %v", req.VariantsId, unlockErr)
			if orderErr == nil {
				orderErr = fmt.Errorf("failed to release lock: %v", unlockErr)
			}
			if orderErr != nil {
				log.Printf("Multiple errors occurred: Order error: %v, Unlock error: %v", orderErr, unlockErr)
			}
		}
	}()

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		orderErr = fmt.Errorf("failed to begin transaction: %v", err)
		return nil, orderErr
	}
	defer tx.Rollback()

	var currentStock int64
	err = tx.QueryRowContext(ctx, `
		SELECT quantity 
		FROM stock 
		WHERE variant_id = $1 
		FOR UPDATE`,
		req.VariantsId).Scan(&currentStock)

	if err != nil {
		orderErr = fmt.Errorf("failed to check current stock: %v", err)
		return nil, orderErr
	}

	if currentStock < req.Quantity {
		orderErr = fmt.Errorf("insufficient stock for variant %s: available %d, requested %d",
			req.VariantsId, currentStock, req.Quantity)
		return nil, orderErr
	}

	_, err = tx.ExecContext(ctx, `
		UPDATE stock 
		SET quantity = quantity - $1,
		    updated_at = $2
		WHERE variant_id = $3`,
		req.Quantity, time.Now(), req.VariantsId)

	if err != nil {
		orderErr = fmt.Errorf("failed to update stock: %v", err)
		return nil, orderErr
	}

	orderID := common.GenerateRandomId("order")
	var order pb.Order
	now := time.Now()
	var createdAt, updatedAt time.Time

	err = tx.QueryRowContext(ctx, `
		INSERT INTO orders (id, variants_id, quantity, total, user_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, variants_id, quantity, total, user_id, created_at, updated_at`,
		orderID, req.VariantsId, req.Quantity, req.Total, req.UserId, now, now).
		Scan(&order.OrderId, &order.VariantsId, &order.Quantity, &order.Total, &order.UserId, &createdAt, &updatedAt)

	if err != nil {
		orderErr = fmt.Errorf("failed to create order: %v", err)
		return nil, orderErr
	}

	redisKey := fmt.Sprintf("stock:%s", req.VariantsId)
	stockData := map[string]interface{}{
		"quantity":  currentStock - req.Quantity,
		"updatedAt": now.Unix(),
	}

	if err := s.redisClient.HMSet(ctx, redisKey, stockData).Err(); err != nil {
		log.Printf("Warning: Failed to update Redis cache: %v", err)
	}

	if err := tx.Commit(); err != nil {
		orderErr = fmt.Errorf("failed to commit transaction: %v", err)
		return nil, orderErr
	}

	order.CreatedAt = timestamppb.New(createdAt)
	order.UpdatedAt = timestamppb.New(updatedAt)

	log.Printf("Successfully created order %s for variant %s", orderID, req.VariantsId)
	result = &order
	return result, orderErr
}
