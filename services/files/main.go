package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/wafi04/golang-backend/grpc/pb"
	"github.com/wafi04/golang-backend/services/common"
	"google.golang.org/grpc"
)

func main() {
	log := common.NewLogger()

	// Initialize Cloudinary
	cld, err := cloudinary.NewFromParams(
		common.LoadEnv("CLOUDINARY_CLOUD_NAME"),
		common.LoadEnv("CLOUDINARY_API_KEY"),
		common.LoadEnv("CLOUDINARY_API_SECRET"),
	)
	if err != nil {
		log.Log(common.ErrorLevel, "Failed to initialize Cloudinary: %v", err)
		return
	}

	// Create gRPC server with Prometheus interceptors
	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
		grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor),
	)

	fileService := NewCloudinaryService(cld)
	pb.RegisterFileServiceServer(grpcServer, fileService)

	http.Handle("/metrics", promhttp.Handler())
	httpServer := &http.Server{
		Addr:    ":5054",
		Handler: nil,
	}

	go func() {
		log.Log(common.InfoLevel, "Starting HTTP server for metrics on :8084")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Log(common.ErrorLevel, "HTTP server error: %v", err)
		}
	}()

	port := common.LoadEnv("FILES_PORT")
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Log(common.ErrorLevel, "Failed to listen: %v", err)
		return
	}

	go func() {
		log.Log(common.InfoLevel, "gRPC server starting on port %s", port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Log(common.ErrorLevel, "Failed to serve gRPC: %v", err)
		}
	}()

	// Graceful shutdown
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
