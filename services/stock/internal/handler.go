package internal

import (
	"context"

	"github.com/wafi04/golang-backend/grpc/pb"
)

type Stockhandler struct {
	pb.UnimplementedStockServiceServer
	db *Database
}

func NewStockService(db *Database) *Stockhandler {
	return &Stockhandler{
		db: db,
	}
}

func (h *Stockhandler) ChangeStock(ctx context.Context, req *pb.ChangeStockRequest) (*pb.Stock, error) {
	return h.db.ChangeStock(ctx, req)
}

func (h *Stockhandler) GetStock(ctx context.Context, req *pb.GetStockRequest) (*pb.Stock, error) {
	return h.db.GetStock(ctx, req)
}
func (h *Stockhandler) CheckStockAvailability(ctx context.Context, req *pb.CheckStockAvailabilityRequest) (*pb.CheckStockAvailabilityResponse, error) {
	return h.db.CheckStockAvailability(ctx, req)
}
