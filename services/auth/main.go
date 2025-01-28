package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/wafi04/common/pkg/logger"
	"github.com/wafi04/golang-backend/configs/database"
	"github.com/wafi04/golang-backend/grpc/pb"
	"github.com/wafi04/golang-backend/services/auth/handler"
	"github.com/wafi04/golang-backend/services/auth/repository"
	"github.com/wafi04/golang-backend/services/auth/service"
	"github.com/wafi04/golang-backend/services/common"
	"google.golang.org/grpc"
)

type Config struct {
	DatabaseURL string
	Port        string
}

func loadConfig() Config {
	return Config{
		DatabaseURL: common.LoadEnv("DATABASE_AUTH"),
		Port:        common.LoadEnv("AUTH_PORT"),
	}
}

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"method", "endpoint"},
	)
	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)
)

func init() {
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
}
func handlertesting(w http.ResponseWriter, r *http.Request) {
	timer := prometheus.NewTimer(httpRequestDuration.WithLabelValues(r.Method, r.URL.Path))
	defer timer.ObserveDuration()

	httpRequestsTotal.WithLabelValues(r.Method, r.URL.Path).Inc()
	w.Write([]byte("Hello, world!"))
}

func main() {
	log := logger.NewLogger()
	config := loadConfig()

	db, err := database.NewDB(config.DatabaseURL)
	if err != nil {
		log.Log(logger.ErrorLevel, "Failed to initialize database: %v", err)
		return

	}
	defer db.Close()

	health := db.Health()
	log.Log(logger.InfoLevel, "Database health: %v", health["status"])

	userRepo := repository.NewUserRepository(db.DB)
	userService := &service.UserService{
		UserRepository: userRepo,
	}
	authHandler := &handler.AuthHandler{
		UserService: userService,
	}

	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
		grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor),
	)
	pb.RegisterAuthServiceServer(grpcServer, authHandler)

	http.HandleFunc("/", handlertesting)
	http.Handle("/metrics", promhttp.Handler())
	httpServer := &http.Server{
		Addr:    ":5051",
		Handler: nil,
	}
	

	go func() {
		log.Log(logger.InfoLevel, "Starting HTTP server for metrics on :5051")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Log(logger.ErrorLevel, "HTTP server error: %v", err)
		}
	}()

	lis, err := net.Listen("tcp", ":5001")
	if err != nil {
		log.Log(logger.ErrorLevel, "Failed to listen: %v", err)
		return
	}

	go func() {
		log.Log(logger.InfoLevel, "gRPC server starting on port %s", ":5001")
		if err := grpcServer.Serve(lis); err != nil {
			log.Log(logger.ErrorLevel, "Failed to serve gRPC: %v", err)
		}
	}()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)
	<-stopChan

	log.Log(logger.InfoLevel, "Shutting down servers...")

	grpcServer.GracefulStop()
	log.Log(logger.InfoLevel, "gRPC server stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Log(logger.ErrorLevel, "HTTP server shutdown error: %v", err)
	}
	log.Log(logger.InfoLevel, "HTTP server stopped")
}
