package main

import (
	"net"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/wafi04/golang-backend/grpc/pb"
	"github.com/wafi04/golang-backend/services/common"

	"google.golang.org/grpc"
)

func main() {
    log :=  common.NewLogger()
    cld, err := cloudinary.NewFromParams(
        common.LoadEnv("CLOUDINARY_CLOUD_NAME"),
       common.LoadEnv("CLOUDINARY_API_KEY"),
       common.LoadEnv("CLOUDINARY_API_SECRET"),
    )
	
    if err != nil {
        log.Log(common.ErrorLevel,"Failed to initialize Cloudinary: %v", err)
    }

	grpcServer := grpc.NewServer()
    fileService := NewCloudinaryService(cld)

    port :=   common.LoadEnv("FILES_PORT")
    pb.RegisterFileServiceServer(grpcServer, fileService)

    lis, err := net.Listen("tcp",port)
    if err != nil {
        log.Log(common.ErrorLevel,"Failed to listen: %v", err)
    }

    log.Log(common.InfoLevel,"Server is running on port %s",port)
    if err := grpcServer.Serve(lis); err != nil {
        log.Log(common.ErrorLevel,"Failed to serve: %v", err)
    }
}