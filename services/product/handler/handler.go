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
	log            logger.Logger
}

func NewProductHandler(service *productrepo.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: service,
	}
}

func (h *ProductHandler) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.Product, error) {
	h.log.Log(logger.InfoLevel, "incoming request ")
	return h.productService.CreateProduct(ctx, req)
}

func (h *ProductHandler) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.Product, error) {
	h.log.Log(logger.InfoLevel, "incoming request ")
	return h.productService.GetProduct(ctx, req)
}
func (h *ProductHandler) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	h.log.Log(logger.InfoLevel, "incoming request list")
	return h.productService.ListProducts(ctx, req)
}

func (h *ProductHandler) CreateProductVariant(ctx context.Context, req *pb.CreateProductVariantRequest) (*pb.ProductVariant, error) {
	h.log.Log(logger.InfoLevel, "Incoming Request Create Varinat")
	return h.productService.CreateProductVariant(ctx, req)
}

func (h *ProductHandler) UpdateProductVariant(ctx context.Context, req *pb.UpdateProductVariantRequest) (*pb.ProductVariant, error) {
	h.log.Log(logger.InfoLevel, "Incoming Request Update Varinat")
	return h.productService.UpdateProductVariant(ctx, req)
}

func (h *ProductHandler) DeleteProductVariant(ctx context.Context, req *pb.DeleteProductVariantRequest) (*pb.DeleteProductResponse, error) {
	h.log.Log(logger.InfoLevel, "Incoming Request Delete Varinat")
	return h.productService.DeleteProductVariant(ctx, req)
}

func (h *ProductHandler) AddProductImage(ctx context.Context, req *pb.AddProductImageRequest) (*pb.ProductImage, error) {
	h.log.Log(logger.InfoLevel, "Incoming Request Create Image")
	return h.productService.AddProductImage(ctx, req)
}

func (h *ProductHandler) UpdateProductImage(ctx context.Context, req *pb.UpdateProductImageRequest) (*pb.ProductImage, error) {
	h.log.Log(logger.InfoLevel, "Incoming Request Update image")
	return h.productService.UpdateProductImage(ctx, req)
}

func (h *ProductHandler) DeleteProductImage(ctx context.Context, req *pb.DeleteProductImageRequest) (*pb.DeleteProductResponse, error) {
	h.log.Log(logger.InfoLevel, "Incoming Request Delete image")
	return h.productService.DeleteProductImage(ctx, req)
}
