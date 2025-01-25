package main

import (
	"net"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/wafi04/common/pkg/logger"
	"github.com/wafi04/golang-backend/configs/database"
	"github.com/wafi04/golang-backend/grpc/pb"
	"github.com/wafi04/golang-backend/services/auth/handler"
	"github.com/wafi04/golang-backend/services/auth/repository"
	"github.com/wafi04/golang-backend/services/auth/service"
	"github.com/wafi04/golang-backend/services/common"
	"google.golang.org/grpc"
)

func main(){
	log := logger.NewLogger()

	db, err := database.NewDB(common.LoadEnv("DATABASE_AUTH"))
	if err != nil {
		log.Log(logger.ErrorLevel, "Failed to initialize database : %v: ", err)
	}
	defer db.Close()
	health := db.Health()
	log.Log(logger.InfoLevel, "Database health : %v", health["status"])

	userRepo := repository.NewUserRepository(db.DB)
	userService := &service.UserService{
		UserRepository: userRepo,
	}
	authHandler := &handler.AuthHandler{
		UserService: userService,

	}
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Log(logger.ErrorLevel, "failed to listen: %v", err)
	}

	
	s := grpc.NewServer()
	pb.RegisterAuthServiceServer(s, authHandler)
	log.Log(logger.InfoLevel, "Auth server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Log(logger.ErrorLevel, "failed to serve: %v", err)
	}

}