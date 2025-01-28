package stock

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/wafi04/golang-backend/grpc/pb"
)

type Database struct {
	db          *sqlx.DB
	redisClient *redis.Client
}

func NewDB(db *sqlx.DB, redisClient *redis.Client) *Database {
	return &Database{
		db:          db,
		redisClient: redisClient,
	}
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
		fmt.Printf("Failed to save data to Redis: %v\n", err)
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
