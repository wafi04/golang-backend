package main

import (
	"net"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/wafi04/common/pkg/logger"
	"github.com/wafi04/golang-backend/configs/database"
	"github.com/wafi04/golang-backend/grpc/pb"
	"github.com/wafi04/golang-backend/services/category/handler"
	"github.com/wafi04/golang-backend/services/category/service"
	"github.com/wafi04/golang-backend/services/common"

	"google.golang.org/grpc"
)

func main(){
	log :=  logger.NewLogger()
	db ,err :=  database.NewDB(common.LoadEnv("DATABASE_CATEGORY"))
		if err != nil {
		log.Log(logger.ErrorLevel, "Failed to initialize database : %v: ", err)
	}
	defer db.Close()

	health := db.Health()
	log.Log(logger.InfoLevel, "Database health : %v", health["status"])


	categoryService := service.NewCategoryService(db.DB)
    categoryHandler := handler.NewCategoryHandler(categoryService)

	grpcServer := grpc.NewServer()
	port :=  common.LoadEnv("CATEGORY_PORT")
    pb.RegisterCategoryServiceServer(grpcServer, categoryHandler)

    lis, err := net.Listen("tcp", port)
    if err != nil {
        log.Log(logger.ErrorLevel, "Failed to listen: %v", err)
        return
    }

    log.Log(logger.InfoLevel, "gRPC server starting on port %s", port)

    if err := grpcServer.Serve(lis); err != nil {
        log.Log(logger.ErrorLevel, "Failed to serve: %v", err)
        return
    }

}