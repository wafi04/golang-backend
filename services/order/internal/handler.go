package internal

import (
	"context"

	"github.com/wafi04/golang-backend/grpc/pb"
)

type OrderHandler struct {
	pb.UnimplementedOrderServiceServer
	orderservice *OrderService
}

func NewOrderHandler(orderService *OrderService) *OrderHandler {
	return &OrderHandler{
		orderservice: orderService,
	}
}

func (h *OrderHandler) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.Order, error) {
	return h.orderservice.CreateOrder(ctx, req)
}
