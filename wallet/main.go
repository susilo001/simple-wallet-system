package main

import (
	"log"
	"net"

	"github.com/susilo001/simple-wallet-system/wallet/entity"
	"github.com/susilo001/simple-wallet-system/wallet/handler"
	pb "github.com/susilo001/simple-wallet-system/wallet/proto/wallet/v1"
	"github.com/susilo001/simple-wallet-system/wallet/repository"
	"github.com/susilo001/simple-wallet-system/wallet/service"

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

	gormDB.AutoMigrate(&entity.Wallet{})
	// setup service

	// uncomment to use postgres gorm
	walletRepo := repository.NewWalletRepository(gormDB)
	walletService := service.NewWalletService(walletRepo)
	//walletHandler := ginHandler.NewWalletHandler(walletService)
	walletHandler := handler.NewWalletHandler(walletService)

	// Run the grpc server
	grpcServer := grpc.NewServer()
	pb.RegisterWalletServiceServer(grpcServer, walletHandler)
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("Running grpc server in port :50052")
	_ = grpcServer.Serve(lis)
}
