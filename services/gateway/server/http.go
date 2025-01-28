package server

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/wafi04/golang-backend/services/common"
	"github.com/wafi04/golang-backend/services/common/middleware"
	authhandler "github.com/wafi04/golang-backend/services/gateway/server/auth"
	categoryhandler "github.com/wafi04/golang-backend/services/gateway/server/category"
	filehandler "github.com/wafi04/golang-backend/services/gateway/server/files"
	orderhandler "github.com/wafi04/golang-backend/services/gateway/server/order"
	producthandler "github.com/wafi04/golang-backend/services/gateway/server/product"
	stockhandler "github.com/wafi04/golang-backend/services/gateway/server/stock"
)

var (
	counter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "devbuls_counter",
		Help: "Counting the total number of request handled",
	})

	gauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "devbulls_gauge",
		Help: "Monitoring node usage",
	}, []string{"node", "namespace"})
)

func RecordMetrics() {
	go func() {
		for {
			counter.Inc()
			gauge.WithLabelValues("node-1", "namespace-b").Set(rand.Float64())
			time.Sleep(time.Second * 5)
		}
	}()

}
func init() {
	prometheus.MustRegister(gauge)
}

func SetupRoutes(
	authGateway *authhandler.AuthHandler,
	categoryGateway *categoryhandler.CategoryHandler,
	fileGateway *filehandler.FileHandler,
	productGateway *producthandler.ProductHandler,
	stockGateway *stockhandler.StockHandler,
	orderGateway *orderhandler.OrderHandler,
) *mux.Router {
	r := mux.NewRouter()
	RecordMetrics()

	// CORS middleware
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "http://192.168.100.81:3000")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	})

	r.Use(common.ResponseTimeMiddleware)

	// API routes
	api := r.PathPrefix("/api/v1").Subrouter()
	r.Handle("/metrics", promhttp.Handler())
	public := api.PathPrefix("").Subrouter()
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		healthStatus := true
		timestamp := time.Now().Format(time.RFC3339)

		data := struct {
			Health    bool   `json:"health"`
			Timestamp string `json:"time"`
		}{
			Health:    healthStatus,
			Timestamp: timestamp,
		}

		if healthStatus {
			common.SendSuccessResponse(w, http.StatusOK, "Connection Ready", data)
		} else {
			common.SendErrorResponse(w, http.StatusServiceUnavailable, "Service Unhealthy")
		}
	}).Methods("GET")
	public.HandleFunc("/auth/register", authGateway.HandleCreateUser).Methods("POST", "OPTIONS")
	public.HandleFunc("/auth/login", authGateway.HandleLogin).Methods("POST", "OPTIONS")

	// Protected routes
	protected := api.PathPrefix("").Subrouter()
	protected.Use(middleware.AuthMiddleware)
	// Auth protected routes
	protected.HandleFunc("/auth/profile", authGateway.HandleGetProfile).Methods("GET", "OPTIONS")
	protected.HandleFunc("/auth/logout", authGateway.HandleLogout).Methods("POST", "OPTIONS")
	protected.HandleFunc("/auth/verification-email", authGateway.HandleVerifyEmail).Methods("POST", "OPTIONS")
	protected.HandleFunc("/auth/refresh-token", authGateway.HandleRefreshToken).Methods("POST", "OPTIONS")
	protected.HandleFunc("/auth/resend-verification", authGateway.HandleResendVerification).Methods("POST", "OPTIONS")
	protected.HandleFunc("/auth/list-sessions", authGateway.HandlerListSessions).Methods("GET", "OPTIONS")
	protected.HandleFunc("/auth/revoke-session/{id}", authGateway.HandleRevokeSessions).Methods("DELETE", "OPTIONS")

	// Category protected routes
	public.HandleFunc("/category", categoryGateway.HandleCreateCategory).Methods("POST", "OPTIONS")
	protected.HandleFunc("/category", categoryGateway.HandleGetCategories).Methods("GET", "OPTIONS")
	public.HandleFunc("/list-categories", categoryGateway.HandleListCategories).Methods("GET", "OPTIONS")
	protected.HandleFunc("/category/{id}", categoryGateway.HandleUpdateCategory).Methods("PUT", "OPTIONS")
	protected.HandleFunc("/category/{id}", categoryGateway.HandleDeleteCategory).Methods("DELETE", "OPTIONS")

	// UploadFile
	protected.HandleFunc("/upload", fileGateway.HandleUploadFile).Methods("POST", "OPTIONS")

	// products
	protected.HandleFunc("/create-product", productGateway.HandleCreateProduct).Methods("POST", "OPTIONS")
	protected.HandleFunc("/product/{id}", productGateway.HandleGetProduct).Methods("GET", "OPTIONS")
	protected.HandleFunc("/product", productGateway.HandleListProducts).Methods("GET", "OPTIONS")
	// variants
	protected.HandleFunc("/product/{id}/variants", productGateway.HandleCreateVariants).Methods("POST", "OPTIONS")
	protected.HandleFunc("/product/variants/{id}", productGateway.HandleUpdateVariants).Methods("PUT", "OPTIONS")
	protected.HandleFunc("/product/variants/{id}", productGateway.HandleDeleteVariants).Methods("DELETE", "OPTIONS")

	// images
	protected.HandleFunc("/product/images/{id}", productGateway.HandleAddProductImage).Methods("POST", "OPTIONS")
	protected.HandleFunc("/product/images/{id}", productGateway.HandleDeleteProductImage).Methods("DELETE", "OPTIONS")

	protected.HandleFunc("/stock", stockGateway.HandleCreateStock).Methods("POST", "OPTIONS")
	protected.HandleFunc("/stock/{id}", stockGateway.HandleCheckAvaibility).Methods("GET", "OPTIONS")

	protected.HandleFunc("/order", orderGateway.HandleCreateOrder).Methods("POST", "OPTIONS")

	return r
}
