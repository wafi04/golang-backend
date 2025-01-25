package handler

import (
	"context"

	"github.com/wafi04/common/pkg/logger"
	"github.com/wafi04/golang-backend/grpc/pb"
	"github.com/wafi04/golang-backend/services/category/service"
)

type CategoryHandler struct {
	pb.UnimplementedCategoryServiceServer
	logger logger.Logger
	categoryservice  *service.CategoryService
}

func  NewCategoryHandler (service *service.CategoryService) *CategoryHandler{
	return &CategoryHandler{
		categoryservice: service,
	}
}


func  (h *CategoryHandler)  CreateCategory(ctx context.Context,req *pb.CreateCategoryRequest ) (*pb.Category,error){
	h.logger.Log(logger.InfoLevel, "Received Create Category ")
	return h.categoryservice.CreateCategory(ctx, req)
}

func (h *CategoryHandler) GetCategories(ctx context.Context, req *pb.ListCategoriesRequest) (*pb.ListCategoriesResponse, error) {
    h.logger.Log(logger.InfoLevel, "Received Get Categories request")
    return h.categoryservice.GetCategories(ctx, req)
}

func (h *CategoryHandler)  ListCategories(ctx context.Context, req *pb.ListCategoriesRequest) (*pb.ListCategoriesResponse, error) {
    h.logger.Log(logger.InfoLevel, "Incoming  Request Categories Called")
    return h.categoryservice.ListCategories(ctx, req)
}

func (h *CategoryHandler)  UpdateCategory(ctx context.Context, req *pb.UpdateCategoryRequest) (*pb.Category, error){
	h.logger.Log(logger.InfoLevel, "Update Category Callled")
	return h.categoryservice.UpdateCategory(ctx,req)
}
func (h *CategoryHandler) DeleteCategory(ctx context.Context, req *pb.DeleteCategoryRequest) (*pb.DeleteCategoryResponse, error){
	h.logger.Log(logger.InfoLevel, "Delete  Category Callled")
	return h.categoryservice.DeleteCategory(ctx,req)
}