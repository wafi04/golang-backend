package producthandler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/wafi04/golang-backend/grpc/pb"
	"github.com/wafi04/golang-backend/services/common"
)

func (h *ProductHandler) HandleCreateVariants(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Color string `json:"color"`
		Sku   string `json:"sku"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Sku == "" {
		req.Sku = common.GenerateSku(req.Color)
	} else if !common.IsSkuValid(req.Sku) {
		return
	}

	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, "Product Id is required", http.StatusBadRequest)
		return
	}

	variants, err := h.productclient.CreateProductVariant(r.Context(), &pb.CreateProductVariantRequest{
		ProductId: id,
		Color:     req.Color,
		Sku:       req.Sku,
	})

	if err != nil {
		common.SendErrorResponseWithDetails(w, http.StatusBadRequest, "Failed to create variants", err.Error())
		return
	}
	common.SendSuccessResponse(w, http.StatusCreated, "Created variants successfully", variants)

}

func (p *ProductHandler) HandleUpdateVariants(w http.ResponseWriter, r *http.Request) {

	var req struct {
		Color     string `json:"color"`
		ProductID string `json:"product_id"`
		Sku       string `json:"sku"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, "Product Variants Id is required", http.StatusBadRequest)
		return
	}

	if req.Sku == "" {
		req.Sku = common.GenerateSku(req.Color)
	} else if !common.IsSkuValid(req.Sku) {
		return
	}

	update, err := p.productclient.UpdateProductVariant(r.Context(), &pb.UpdateProductVariantRequest{
		Variant: &pb.ProductVariant{
			Id:        id,
			Color:     req.Color,
			Sku:       req.Sku,
			ProductId: req.ProductID,
		},
	})

	if err != nil {
		common.SendErrorResponseWithDetails(w, http.StatusBadRequest, "Failed to update variants", err.Error())
		return
	}
	common.SendSuccessResponse(w, http.StatusOK, "Update variants successfully", update)
}

func (p *ProductHandler) HandleDeleteVariants(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, "Product Variants Id is required", http.StatusBadRequest)
		return
	}

	delete, err := p.productclient.DeleteProductVariant(r.Context(), &pb.DeleteProductVariantRequest{
		Id: id,
	})
	if err != nil {
		common.SendErrorResponseWithDetails(w, http.StatusBadRequest, "Failed to update variants", err.Error())
		return
	}

	common.SendSuccessResponse(w, http.StatusOK, "Update variants successfully", delete)

}
