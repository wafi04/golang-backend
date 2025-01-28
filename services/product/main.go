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
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/wafi04/common/pkg/logger"
	"github.com/wafi04/golang-backend/configs/database"
	"github.com/wafi04/golang-backend/grpc/pb"
	"github.com/wafi04/golang-backend/services/category/handler"
	"github.com/wafi04/golang-backend/services/category/service"
	"github.com/wafi04/golang-backend/services/common"
	"google.golang.org/grpc"
)

type Config struct {
	DatabaseURL string
	Port        string
}

func loadConfig() Config {
	return Config{
		DatabaseURL: common.LoadEnv("DATABASE_PRODUCT"),
		Port:        common.LoadEnv("PRODUCT_PORT"),
	}
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

	categoryService := service.NewCategoryService(db.DB)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
		grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor),
	)
	pb.RegisterCategoryServiceServer(grpcServer, categoryHandler)

	http.Handle("/metrics", promhttp.Handler())
	httpServer := &http.Server{
		Addr:    ":5053",
		Handler: nil,
	}

	go func() {
		log.Log(logger.InfoLevel, "Starting HTTP server for metrics on :8083")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Log(logger.ErrorLevel, "HTTP server error: %v", err)
		}
	}()

	lis, err := net.Listen("tcp", config.Port)
	if err != nil {
		log.Log(logger.ErrorLevel, "Failed to listen: %v", err)
		return
	}

	go func() {
		log.Log(logger.InfoLevel, "gRPC server starting on port %s", config.Port)
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
