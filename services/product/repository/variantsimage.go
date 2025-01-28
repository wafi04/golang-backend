package productrepo

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/wafi04/common/pkg/logger"
	"github.com/wafi04/golang-backend/grpc/pb"
)

//   rpc AddProductImage (AddProductImageRequest) returns (ProductImage);
//     rpc UpdateProductImage (UpdateProductImageRequest) returns (ProductImage);
//     rpc DeleteProductImage (DeleteProductImageRequest) returns (DeleteProductResponse);

func (pr *ProductService) AddProductImage(ctx context.Context, req *pb.AddProductImageRequest) (*pb.ProductImage, error) {
    imageID := uuid.New().String()

    querySelect := `
        SELECT id FROM product_variants WHERE id = $1
    `
    var variantID string
    err := pr.db.QueryRowContext(ctx, querySelect, req.VariantId).Scan(&variantID)
    if err != nil {
        pr.log.Log(logger.ErrorLevel, "Variant not found: %v", err)
        return nil, fmt.Errorf("variant not found: %v", err)
    }

    queryInsert := `
        INSERT INTO product_images (id, url, variant_id, is_main)
        VALUES ($1, $2, $3, $4)
        RETURNING id, url, variant_id, is_main
    `

    var image pb.ProductImage
    err = pr.db.QueryRowContext(ctx, queryInsert, imageID, req.Url, req.VariantId, req.IsMain).Scan(
        &image.Id,
        &image.Url,
        &image.VariantId,
        &image.IsMain,
    )
    if err != nil {
        pr.log.Log(logger.ErrorLevel, "Failed to create product image: %v", err)
        return nil, fmt.Errorf("failed to create product image: %v", err)
    }

    return &image, nil
}


func (pr *ProductService) UpdateProductImage(ctx context.Context, req *pb.UpdateProductImageRequest) (*pb.ProductImage, error) {
    querySelect := `
        SELECT id FROM product_images WHERE id = $1
    `
    var imageID string
    err := pr.db.QueryRowContext(ctx, querySelect, req.Image.Id).Scan(&imageID)
    if err != nil {
        pr.log.Log(logger.ErrorLevel, "Image not found: %v", err)
        return nil, fmt.Errorf("image not found: %v", err)
    }

    queryUpdate := `
        UPDATE product_images
        SET url = $1, is_main = $2
        WHERE id = $3
        RETURNING id, url, variant_id, is_main
    `

    var image pb.ProductImage
    err = pr.db.QueryRowContext(ctx, queryUpdate, req.Image.Url, req.Image.IsMain ,req.Image.Id).Scan(
        &image.Id,
        &image.Url,
        &image.VariantId,
        &image.IsMain,
    )
    if err != nil {
        pr.log.Log(logger.ErrorLevel, "Failed to update product image: %v", err)
        return nil, fmt.Errorf("failed to update product image: %v", err)
    }

    return &image, nil
}


func (pr *ProductService) DeleteProductImage(ctx context.Context, req *pb.DeleteProductImageRequest) (*pb.DeleteProductResponse, error) {
    queryDelete := `
        DELETE FROM product_images
        WHERE id = $1
    `

    result, err := pr.db.ExecContext(ctx, queryDelete, req.Id)
    if err != nil {
        pr.log.Log(logger.ErrorLevel, "Failed to delete product image: %v", err)
        return nil, fmt.Errorf("failed to delete product image: %v", err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        pr.log.Log(logger.ErrorLevel, "Failed to get rows affected: %v", err)
        return nil, fmt.Errorf("failed to get rows affected: %v", err)
    }

    if rowsAffected == 0 {
        return &pb.DeleteProductResponse{
            Success: false,
        }, nil
    }

    return &pb.DeleteProductResponse{
        Success: true,
    }, nil
}