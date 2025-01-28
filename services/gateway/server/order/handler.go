package orderhandler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/wafi04/golang-backend/grpc/pb"
	"github.com/wafi04/golang-backend/services/common"
	"github.com/wafi04/golang-backend/services/common/middleware"
	"github.com/wafi04/golang-backend/services/gateway/server/config"
)

type OrderHandler struct {
	stockService pb.StockServiceClient
	orderService pb.OrderServiceClient
}

func containsString(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}
func isInsufficientStockError(err error) bool {
	return err != nil && containsString(err.Error(), "insufficient stock")
}
func isValidationError(err error) bool {
	return err != nil && (containsString(err.Error(), "invalid") ||
		containsString(err.Error(), "required") ||
		containsString(err.Error(), "must be"))
}
func NewProductGateway(ctx context.Context) (*OrderHandler, error) {

	connOrder, err := config.ConnectWithRetry(config.Load().OrderServiceURL, "order")
	log.Printf("ORDER : %s", common.LoadEnv("ORDER_SERVICE_URL"))
	if err != nil {
		return nil, err
	}

	connStock, err := config.ConnectWithRetry(config.Load().StockServiceURL, "stock")
	log.Printf("stock : %s", common.LoadEnv("STOCK_SERVICE_URL"))
	if err != nil {
		return nil, err
	}

	return &OrderHandler{
		orderService: pb.NewOrderServiceClient(connOrder),
		stockService: pb.NewStockServiceClient(connStock),
	}, nil
}

func (h *OrderHandler) HandleCreateOrder(w http.ResponseWriter, r *http.Request) {
	log.Printf("Incoming Request from : %s", r.URL.Path)
	user, _ := middleware.GetUserFromContext(r.Context())
	var req struct {
		Qty       int64  `json:"quantity"`
		Varinatid string `json:"variant_id"`
	}
	log.Printf("cahange")

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if req.Qty <= 0 {
		common.SendErrorResponse(w, http.StatusBadRequest, "Quantity must be a positive number")
		return
	}

	order, err := h.orderService.CreateOrder(r.Context(), &pb.CreateOrderRequest{
		Quantity:   req.Qty,
		VariantsId: req.Varinatid,
		Total:      float64(req.Qty) * 8,
		UserId:     user.UserId,
	})
	if err != nil {
		log.Printf("Error creating order: %v", err)
		switch {
		case isInsufficientStockError(err):
			common.SendErrorResponseWithDetails(w, http.StatusConflict, "Insufficient stock ", err.Error())
		case isValidationError(err):
			common.SendErrorResponseWithDetails(w, http.StatusBadRequest, "Validation error", err.Error())
		default:
			common.SendErrorResponseWithDetails(w, http.StatusInternalServerError, "Failed to create order", err.Error())
		}
		return
	}

	common.SendSuccessResponse(w, http.StatusCreated, "Created success", order)
}
