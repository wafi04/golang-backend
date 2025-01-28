package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/wafi04/golang-backend/services/common"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Config struct {
	AuthServiceURL     string
	CategoryServiceURL string
	ProductServiceURL  string
	FilesServiceURL    string
	StockServiceURL    string
	OrderServiceURL    string
}

func Load() *Config {
	cfg := &Config{
		AuthServiceURL:     common.LoadEnv("AUTH_SERVICE_URL"),
		CategoryServiceURL: common.LoadEnv("CATEGORY_SERVICE_URL"),
		ProductServiceURL:  common.LoadEnv("PRODUCT_SERVICE_URL"),
		FilesServiceURL:    common.LoadEnv("FILES_SERVICE_URL"),
		StockServiceURL:    common.LoadEnv("STOCK_SERVICE_URL"),
		OrderServiceURL:    common.LoadEnv("ORDER_SERVICE_URL"),
	}

	// Validate URLs
	validateURL("AUTH_SERVICE_URL", cfg.AuthServiceURL)
	validateURL("CATEGORY_SERVICE_URL", cfg.CategoryServiceURL)
	validateURL("PRODUCT_SERVICE_URL", cfg.ProductServiceURL)
	validateURL("FILES_SERVICE_URL", cfg.FilesServiceURL)
	validateURL("STOCK_SERVICE_URL", cfg.StockServiceURL)
	validateURL("ORDER_SERVICE_URL", cfg.OrderServiceURL)

	return cfg
}

func validateURL(envVarName, url string) {
	if url == "" {
		log.Fatalf("%s environment variable is not set", envVarName)
	}

}

func ConnectWithRetry(target string, service string) (*grpc.ClientConn, error) {
	maxAttempts := 5
	var conn *grpc.ClientConn
	var err error

	for i := 0; i < maxAttempts; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		log.Printf("Attempting to connect to %s service (attempt %d/%d)...", service, i+1, maxAttempts)

		conn, err = grpc.DialContext(ctx,
			target,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithBlock(),
		)

		if err == nil {
			log.Printf("Successfully connected to %s service", service)
			return conn, nil
		}

		log.Printf("Failed to connect to %s service: %v. Retrying...", service, err)
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("failed to connect to %s service after %d attempts: %v", service, maxAttempts, err)
}
