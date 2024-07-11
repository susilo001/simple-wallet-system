package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	userpb "github.com/susilo001/simple-wallet-system/user/proto/user/v1"
	walletpb "github.com/susilo001/simple-wallet-system/wallet/proto/wallet/v1"
	"google.golang.org/grpc"
)

func main() {
	// Initialize gRPC clients
	userConn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to User service: %v", err)
	}
	defer userConn.Close()
	userClient := userpb.NewUserServiceClient(userConn)

	walletConn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to Wallet service: %v", err)
	}
	defer walletConn.Close()
	walletClient := walletpb.NewWalletServiceClient(walletConn)

	r := gin.Default()

	r.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		// Call User service
		userResp, err := userClient.GetUser(context.Background(), &userpb.GetUserRequest{Id: id})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// Call Wallet service
		walletResp, err := walletClient.GetWallet(context.Background(), &walletpb.GetWalletRequest{WalletId: id})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"user":   userResp.User,
			"wallet": walletResp.Wallet,
		})
	})

	r.Run(":8080")
}
