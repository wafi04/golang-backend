package handler

import (
	"context"

	"github.com/wafi04/common/pkg/logger"
	"github.com/wafi04/golang-backend/grpc/pb"
	productrepo "github.com/wafi04/golang-backend/services/product/repository"
)

type ProductHandler struct {
	pb.UnimplementedProductServiceServer
	productService *productrepo.ProductService
	log  logger.Logger
}


func NewProductHandler(service  *productrepo.ProductService)  *ProductHandler{
	return  &ProductHandler{
		productService: service,
	}
}

func (h *ProductHandler)  CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.Product,error){
	h.log.Log(logger.InfoLevel, "incoming request ")
	return h.productService.CreateProduct(ctx, req)
}

func (h *ProductHandler)  GetProduct(ctx context.Context,req *pb.GetProductRequest)  (*pb.Product,error){
	h.log.Log(logger.InfoLevel, "incoming request ")
	return h.productService.GetProduct(ctx, req)
}
func (h *ProductHandler)  ListProducts(ctx context.Context,req *pb.ListProductsRequest)  (*pb.ListProductsResponse,error){
	h.log.Log(logger.InfoLevel, "incoming request ")
	return h.productService.ListProducts(ctx, req)
}