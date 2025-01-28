package productrepo

import (
	"context"

	"github.com/google/uuid"
	"github.com/wafi04/golang-backend/grpc/pb"
)

func (pr *ProductService) CreateProductVariant(ctx context.Context, req *pb.CreateProductVariantRequest) (*pb.ProductVariant, error) {
	variantsID := uuid.New().String()
	var variants pb.ProductVariant
	query := `
		INSERT INTO product_variants (id,color,sku,product_id)
		VALUES ($1,$2,$3,$4)
		RETURNING id, color, sku, product_id
	`

	err := pr.db.QueryRowContext(ctx, query, variantsID, req.Color, req.Sku, req.ProductId).Scan(
		&variants.Id,
		&variants.Color,
		&variants.Sku,
		&variants.ProductId,
	)

	if err != nil {
		pr.log.Error("Failed to Create Variants : %v ", err)
	}

	return &pb.ProductVariant{
		Id:        variants.Id,
		Color:     variants.Color,
		Sku:       variants.Sku,
		ProductId: variants.ProductId,
	}, nil
}

func (pr *ProductService) UpdateProductVariant(ctx context.Context, req *pb.UpdateProductVariantRequest) (*pb.ProductVariant, error) {
	query := `
        UPDATE product_variants
        SET color = $1, sku = $2
        WHERE id = $3
        RETURNING id, color, sku, product_id
    `

	var variant pb.ProductVariant
	err := pr.db.QueryRowContext(ctx, query, req.Variant.Color, req.Variant.Sku, req.Variant.Id).Scan(
		&variant.Id,
		&variant.Color,
		&variant.Sku,
		&variant.ProductId,
	)

	if err != nil {
		pr.log.Error("Failed to update variant: %v", err)
		return nil, err
	}

	return &variant, nil
}

func (pr *ProductService) DeleteProductVariant(ctx context.Context, req *pb.DeleteProductVariantRequest) (*pb.DeleteProductResponse, error) {
	query := `
	DELETE FROM product_variants WHERE id = $1
	`

	_, err := pr.db.ExecContext(ctx, query, req.Id)

	if err != nil {
		pr.log.Error("Failed to Delete Variants : %v", err)
		return nil, err
	}

	return &pb.DeleteProductResponse{
		Success: true,
	}, nil

}
