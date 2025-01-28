package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/wafi04/golang-backend/services/common"
	"github.com/wafi04/golang-backend/services/gateway/server"
	authhandler "github.com/wafi04/golang-backend/services/gateway/server/auth"
	categoryhandler "github.com/wafi04/golang-backend/services/gateway/server/category"
	filehandler "github.com/wafi04/golang-backend/services/gateway/server/files"
	orderhandler "github.com/wafi04/golang-backend/services/gateway/server/order"
	producthandler "github.com/wafi04/golang-backend/services/gateway/server/product"
	stockhandler "github.com/wafi04/golang-backend/services/gateway/server/stock"
)

func main() {
	logs := common.NewLogger()

	logs.Info("Staring Server gateway ")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	gateway, err := authhandler.NewGateway(ctx)
	if err != nil {
		logs.Log(common.ErrorLevel, "Failed to connect Auth Service : %v", err)
	}
	categorygateway, err := categoryhandler.NewCategoryGateway(ctx)
	if err != nil {
		logs.Log(common.ErrorLevel, "Failed to connect Category Service : %v", err)
	}

	fileGateway, err := filehandler.NewFilesGateway(ctx)
	if err != nil {
		logs.Log(common.ErrorLevel, "Failed to connect File Service : %v", err)
	}
	productGateway, err := producthandler.NewProductGateway(ctx)
	if err != nil {
		logs.Log(common.ErrorLevel, "Failed to co	nnect product Service : %v", err)
	}
	stockGateway, err := stockhandler.NewStockGateway(ctx)
	if err != nil {
		logs.Log(common.ErrorLevel, "Failed to conect	nnect stock  Service : %v", err)
	}

	orderGateway, err := orderhandler.NewProductGateway(ctx)
	if err != nil {
		logs.Log(common.ErrorLevel, "Failed to conec order Service : %v", err)
	}

	r := mux.NewRouter()

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logs.Log(common.InfoLevel, "Incoming request: %s %s", r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
		})
	})

	r = server.SetupRoutes(gateway, categorygateway, fileGateway, productGateway, stockGateway, orderGateway)
	r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, _ := route.GetPathTemplate()
		methods, _ := route.GetMethods()
		logs.Log(common.InfoLevel, "Registered route: %s (Methods: %v)", path, methods)
		return nil
	})

	srv := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:4000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logs.Log(common.InfoLevel, "Gateway server starting on http://localhost:4000")
	log.Print(srv.ListenAndServe())
}
