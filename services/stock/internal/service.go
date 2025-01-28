package internal

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/wafi04/golang-backend/grpc/pb"
)

type Database struct {
	db          *sqlx.DB
	redisClient *redis.Client
}

func NewDB(db *sqlx.DB, redis *redis.Client) *Database {
	return &Database{
		db:          db,
		redisClient: redis,
	}
}

func (r *Database) ChangeStock(ctx context.Context, req *pb.ChangeStockRequest) (*pb.Stock, error) {
	now := time.Now()

	query := `
	INSERT INTO stock (quantity, variant_id, created_at, updated_at)
	VALUES ($1, $2, $3, $4)
	RETURNING quantity, variant_id
	`

	var stock pb.Stock
	err := r.db.QueryRowContext(ctx, query, req.Quantity, req.VariantId, now, now).Scan(&stock.Quantity, &stock.VariantId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no rows were returned after insert")
		}
		return nil, fmt.Errorf("failed to create stock: %v", err)
	}

	stock.CreatedAt = now.Unix()
	stock.UpdatedAt = now.Unix()

	redisKey := fmt.Sprintf("stock-%s", req.VariantId)

	customStock := CustomStock{
		Quantity:  stock.Quantity,
		VariantId: stock.VariantId,
		CreatedAt: stock.CreatedAt,
		UpdatedAt: stock.UpdatedAt,
	}

	stockJSON, err := json.Marshal(customStock)
	if err != nil {
		fmt.Printf("Failed to marshal stock data for Redis: %v\n", err)
	} else {
		// Store the stock data in Redis
		err = r.redisClient.Set(ctx, redisKey, stockJSON, 0).Err()
		if err != nil {
			fmt.Printf("Failed to update Redis cache: %v\n", err)
		} else {
			fmt.Printf("Successfully updated Redis cache for key: %s\n", redisKey)
		}
	}

	return &stock, nil
}

func (r *Database) GetStock(ctx context.Context, req *pb.GetStockRequest) (*pb.Stock, error) {
	now := time.Now()
	query := `
	SELECT quantity,variant_id,created_at,updated_at  FROM stock WHERE variant_id = 1
	`
	var stock pb.Stock
	err := r.db.QueryRowContext(ctx, query, req.VariantId).Scan(&stock.Quantity, &stock.VariantId, &stock.CreatedAt, &stock.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no rows were returned after insert")
		}
		return nil, fmt.Errorf("failed to create stock: %v", err)
	}

	stock.CreatedAt = now.Unix()
	stock.UpdatedAt = now.Unix()

	return &stock, nil
}

func (r *Database) CheckStockAvailability(ctx context.Context, req *pb.CheckStockAvailabilityRequest) (*pb.CheckStockAvailabilityResponse, error) {
	redisKey := fmt.Sprintf("stock:%s", req.VariantId)
	log.Printf("Variant : %s  , qty : %d", req.VariantId, req.RequestedQuantity)

	availableQuantity, err := r.redisClient.HGet(ctx, redisKey, "quantity").Int64()
	if err == nil {
		if availableQuantity >= req.RequestedQuantity {
			return &pb.CheckStockAvailabilityResponse{
				IsAvailable:       true,
				AvailableQuantity: availableQuantity,
			}, nil
		} else {
			return &pb.CheckStockAvailabilityResponse{
				IsAvailable:       false,
				AvailableQuantity: availableQuantity,
			}, nil
		}
	}

	query := `
        SELECT quantity 
        FROM stock 
        WHERE variant_id = $1
    `
	var dbQuantity int64
	err = r.db.QueryRowContext(ctx, query, req.VariantId).Scan(&dbQuantity)
	if err != nil {
		if err == sql.ErrNoRows {
			return &pb.CheckStockAvailabilityResponse{
				IsAvailable:       false,
				AvailableQuantity: 0,
			}, nil
		}
		return nil, err
	}

	err = r.redisClient.HSet(ctx, redisKey, "quantity", dbQuantity).Err()
	if err != nil {
		fmt.Printf("Gagal menyimpan data ke Redis: %v\n", err)
	}

	if dbQuantity >= req.RequestedQuantity {
		return &pb.CheckStockAvailabilityResponse{
			IsAvailable:       true,
			AvailableQuantity: dbQuantity,
		}, nil
	} else {
		return &pb.CheckStockAvailabilityResponse{
			IsAvailable:       false,
			AvailableQuantity: dbQuantity,
		}, nil
	}
}
