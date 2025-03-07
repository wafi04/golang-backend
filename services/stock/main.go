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
	"github.com/redis/go-redis/v9"
	"github.com/wafi04/golang-backend/configs/database"
	"github.com/wafi04/golang-backend/grpc/pb"
	"github.com/wafi04/golang-backend/services/common"
	"github.com/wafi04/golang-backend/services/stock/internal"
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

	log := common.NewLogger()
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "P@ssw0rd*1",
		DB:       0,
	})

	db, err := database.NewDB(common.LoadEnv("DATABASE_STOCK"))
	if err != nil {
		log.Log(common.ErrorLevel, "Failed to initialize database: %v", err)
		return
	}
	defer db.Close()

	database := internal.NewDB(db.DB, redisClient)
	stockService := internal.NewStockService(database)

	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
		grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor),
	)
	pb.RegisterStockServiceServer(grpcServer, stockService)

	http.Handle("/metrics", promhttp.Handler())
	httpServer := &http.Server{
		Addr:    ":5055",
		Handler: nil,
	}

	go func() {
		log.Log(common.InfoLevel, "Starting HTTP server for metrics on :5055")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Log(common.ErrorLevel, "HTTP server error: %v", err)
		}
	}()

	lis, err := net.Listen("tcp", ":5005")
	if err != nil {
		log.Log(common.ErrorLevel, "Failed to listen: %v", err)
		return
	}

	go func() {
		log.Log(common.InfoLevel, "gRPC server starting on port %s", ":5005")
		if err := grpcServer.Serve(lis); err != nil {
			log.Log(common.ErrorLevel, "Failed to serve gRPC: %v", err)
		}
	}()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)
	<-stopChan

	log.Log(common.InfoLevel, "Shutting down servers...")

	grpcServer.GracefulStop()
	log.Log(common.InfoLevel, "gRPC server stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Log(common.ErrorLevel, "HTTP server shutdown error: %v", err)
	}
	log.Log(common.InfoLevel, "HTTP server stopped")

}
