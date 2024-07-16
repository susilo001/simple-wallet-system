package main

import (
	"log"
	"net"

	"github.com/susilo001/simple-wallet-system/user/entity"
	"github.com/susilo001/simple-wallet-system/user/handler"
	pb "github.com/susilo001/simple-wallet-system/user/proto/user/v1"
	"github.com/susilo001/simple-wallet-system/user/repository"
	"github.com/susilo001/simple-wallet-system/user/service"

	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// setup gorm connection
	dsn := "postgresql://postgres:password@localhost:5432/postgres"
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		log.Fatalln(err)
	}

	gormDB.AutoMigrate(&entity.User{})

	repository := repository.NewUserRepository(gormDB)
	service := service.NewUserService(repository)
	handler := handler.NewUserHandler(service)

	// Run the grpc server
	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, handler)
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("Running grpc server in port :50051")
	_ = grpcServer.Serve(lis)
}
