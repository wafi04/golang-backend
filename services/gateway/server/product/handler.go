package producthandler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/wafi04/golang-backend/grpc/pb"
	"github.com/wafi04/golang-backend/services/common"
	"github.com/wafi04/golang-backend/services/gateway/server/config"
)

type ProductHandler struct {
	productclient pb.ProductServiceClient
	filesClient   pb.FileServiceClient
}

func NewProductGateway(ctx context.Context) (*ProductHandler, error) {

	conn, err := config.ConnectWithRetry(config.Load().ProductServiceURL, "product")
	log.Printf("produdtc : %s", config.Load().ProductServiceURL)
	if err != nil {
		return nil, err
	}

	connFile, err := config.ConnectWithRetry(config.Load().FilesServiceURL, "files")
	log.Printf("produdtc : %s", common.LoadEnv("FILES_SERVICE_URL"))
	if err != nil {
		return nil, err
	}

	return &ProductHandler{
		productclient: pb.NewProductServiceClient(conn),
		filesClient:   pb.NewFileServiceClient(connFile),
	}, nil
}

func (h *ProductHandler) HandleCreateProduct(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received create product request: %s %s", r.Method, r.URL.Path)

	var req ProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("Decoded request: %+v", &req)

	if req.Sku == "" {
		req.Sku = common.GenerateSku(req.Name)
	} else if !common.IsSkuValid(req.Sku) {
		return
	}

	resp, err := h.productclient.CreateProduct(r.Context(), &pb.CreateProductRequest{
		Product: &pb.Product{
			Name:        req.Name,
			Description: req.Description,
			SubTitle:    req.SubTitle,
			Price:       req.Price,
			Sku:         req.Sku,
			CategoryId:  req.CategoryId,
		},
	})

	if err != nil {
		log.Printf("Error from auth service: %v", err.Error())
		http.Error(w, fmt.Sprintf("Error creating user: %v", err.Error()), http.StatusInternalServerError)
		return
	}

	log.Printf("Received response from auth service: %+v", resp)

	w.Header().Set("Content-Type", "application/json")

	response := common.Success(resp, "Product Created Successfully")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func (h *ProductHandler) HandleGetProduct(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received get product request: %s %s", r.Method, r.URL.Path)

	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, ok := vars["id"]

	if !ok {
		http.Error(w, "Category ID is required", http.StatusBadRequest)
		return
	}

	res, err := h.productclient.GetProduct(r.Context(), &pb.GetProductRequest{
		Id: id,
	})

	if err != nil {
		log.Printf("Failed to get Product: %v", err)
		common.SendErrorResponseWithDetails(w, http.StatusInternalServerError, "Failed to get products", err.Error())
		return
	}
	resp := common.Success(res, "Get product successfully")

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

}

func (h *ProductHandler) HandleListProducts(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received get product request: %s %s", r.Method, r.URL.Path)
	w.Header().Set("Content-Type", "application/json")

	_ = r.URL.Query().Get("page")
	// limit, err := strconv.Atoi(r.URL.Query().Get("limit"))

	req := &pb.ListProductsRequest{
		PageSize:  10,
		PageToken: "0",
	}

	res, err := h.productclient.ListProducts(r.Context(), req)

	if err != nil {
		log.Printf("Failed to get Product: %v", err)
		common.SendErrorResponseWithDetails(w, http.StatusInternalServerError, "Failed to get products", err.Error())
		return
	}

	// Success response
	resp := common.Success(res, "Pagination product successfully")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

}
