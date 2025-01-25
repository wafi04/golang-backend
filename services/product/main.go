package main

import (
	"net"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/wafi04/golang-backend/configs/database"
	"github.com/wafi04/golang-backend/grpc/pb"
	"github.com/wafi04/golang-backend/services/common"
	"github.com/wafi04/golang-backend/services/product/handler"
	productrepo "github.com/wafi04/golang-backend/services/product/repository"
	"google.golang.org/grpc"
)


func main(){
	log :=  common.NewLogger()
	db ,err :=  database.NewDB(common.LoadEnv("DATABASE_PRODUCT"))
		if err != nil {
		log.Log(common.ErrorLevel, "Failed to initialize database : %v: ", err)
	}
	defer db.Close()

	health := db.Health()
	log.Log(common.InfoLevel, "Database health : %v", health["status"])


	productService := productrepo.NewProductService(db.DB)
    productHandler := handler.NewProductHandler(productService)

	grpcServer := grpc.NewServer()

	port :=  common.LoadEnv("PRODUCT_PORT")

    pb.RegisterProductServiceServer(grpcServer, productHandler)

    lis, err := net.Listen("tcp", port)
    if err != nil {
        log.Log(common.ErrorLevel, "Failed to listen: %v", err)
        return
    }

    log.Log(common.InfoLevel, "Product service starting on port %s", port)

    if err := grpcServer.Serve(lis); err != nil {
        log.Log(common.ErrorLevel, "Failed to serve: %v", err)
        return
    }

}