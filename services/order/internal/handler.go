package handler

import (
	"context"

	"github.com/wafi04/golang-backend/grpc/pb"
)

type OrderHandler struct {
	pb.UnimplementedOrderServiceServer
	orderRepository *OrderRepository
}

func NewOrderHandler(orderRepository *OrderRepository) *OrderHandler {
	return &OrderHandler{
		orderRepository: orderRepository,
	}
}

func (h *OrderHandler) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.Order, error) {
	return h.orderRepository.CreateOrder(ctx, req)
}
