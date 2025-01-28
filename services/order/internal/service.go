package internal

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/wafi04/golang-backend/grpc/pb"
	"github.com/wafi04/golang-backend/services/common"
)

type OrderService struct {
	db          *sqlx.DB
	redisClient *redis.Client
}

func NewOrderService(db *sqlx.DB, rediclient *redis.Client) *OrderService {
	return &OrderService{
		db:          db,
		redisClient: rediclient,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.Order, error) {
	orderID := common.GenerateRandomId("order")
	var order pb.Order

	query := `
	INSERT INTO orders (id, variants_id, quantity, total, user_id, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id, variants_id, quantity, total, user_id, created_at, updated_at
	`
	now := time.Now()

	err := s.db.QueryRowContext(ctx, query, orderID, req.VariantsId, req.Quantity, req.Total, req.UserId, now, now).
		Scan(&order.OrderId, &order.VariantsId, &order.Quantity, &order.Total, &order.UserId, &order.CreatedAt, &order.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no rows were returned after insert")
		}
		return nil, fmt.Errorf("failed to create order: %v", err)
	}

	return &order, nil
}
