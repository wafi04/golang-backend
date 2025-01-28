package stockhandler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/wafi04/golang-backend/grpc/pb"
	"github.com/wafi04/golang-backend/services/common"
	authhandler "github.com/wafi04/golang-backend/services/gateway/server/auth"
)

type StockHandler struct {
	stockHandler pb.StockServiceClient
	logger       common.Logger
}

func NewStockGateway(ctx context.Context) (*StockHandler, error) {
	conn, err := authhandler.ConnectWithRetry("stock:5005", "stock")
	log.Printf("STOCK : %s", common.LoadEnv("STOCK_SERVICE_URL"))
	if err != nil {
		return nil, err
	}

	return &StockHandler{
		stockHandler: pb.NewStockServiceClient(conn),
	}, nil
}
func (h *StockHandler) HandleCreateStock(w http.ResponseWriter, r *http.Request) {
	start := time.Now() // Start time of the request
	h.logger.Log(common.InfoLevel, "Incoming request client")

	var req struct {
		VariantId string `json:"variant_id"`
		Quantity  int64  `json:"quantity"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	stock, err := h.stockHandler.ChangeStock(r.Context(), &pb.ChangeStockRequest{
		VariantId: req.VariantId,
		Quantity:  req.Quantity,
	})

	if err != nil {
		log.Printf("Failed to change stock: %v", err) // Log the error
		common.SendErrorResponse(w, http.StatusBadRequest, "Failed to change stock")
		return
	}
	duration := time.Since(start)
	h.logger.Log(common.InfoLevel, "Response time for HandleCreateStock: %v\n", duration)

	common.SendSuccessResponse(w, http.StatusAccepted, "Change stock available", stock)

}

func (h *StockHandler) HandleCheckAvaibility(w http.ResponseWriter, r *http.Request) {
	start := time.Now() // Start time of the request
	h.logger.Log(common.InfoLevel, "Incoming request from client")

	var req struct {
		Qty int64 `json:"quantity"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	variant_id, ok := vars["id"]
	if !ok {
		http.Error(w, "Category ID is required", http.StatusBadRequest)
		return
	}
	available, err := h.stockHandler.CheckStockAvailability(r.Context(), &pb.CheckStockAvailabilityRequest{
		VariantId:         variant_id,
		RequestedQuantity: req.Qty,
	})

	if err != nil {
		common.SendErrorResponse(w, 400, "Errorr Check avaibility")
		return
	}

	common.SendSuccessResponse(w, 200, "Stock is available", available)

	// Log response time
	duration := time.Since(start)
	log.Printf("Response time for HandleCheckAvaibility: %v\n", duration)
}
