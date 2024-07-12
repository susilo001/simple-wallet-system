package main

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	userpb "github.com/susilo001/simple-wallet-system/user/proto/user/v1"
	walletpb "github.com/susilo001/simple-wallet-system/wallet/proto/wallet/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	userConn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to User service: %v", err)
	}
	defer userConn.Close()
	userClient := userpb.NewUserServiceClient(userConn)

	walletConn, err := grpc.NewClient("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to Wallet service: %v", err)
	}
	defer walletConn.Close()
	walletClient := walletpb.NewWalletServiceClient(walletConn)

	r := gin.Default()

	r.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")

		userId, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Call User service
		userResp, err := userClient.GetUser(context.Background(), &userpb.GetUserRequest{Id: int32(userId)})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Call Wallet service
		walletResp, err := walletClient.GetWallet(context.Background(), &walletpb.GetWalletRequest{WalletId: int32(userId)})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"user":   userResp.User,
			"wallet": walletResp.Wallet,
		})
	})

	r.GET("/users/:id/transactions", func(c *gin.Context) {
		id := c.Param("id")

		userId, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Call Wallet service to get transaction history
		walletResp, err := walletClient.GetTransactions(context.Background(), &walletpb.GetTransactionsRequest{WalletId: int32(userId)})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"transactions": walletResp.Transactions,
		})
	})

	r.POST("/wallets/:id/topup", func(c *gin.Context) {
		id := c.Param("id")
		walletId, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var req struct {
			Amount float64 `json:"amount" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Call Wallet service to perform top-up
		_, err = walletClient.TopUpWallet(context.Background(), &walletpb.TopUpRequest{WalletId: int32(walletId), Amount: req.Amount})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Wallet top-up successful"})
	})

	r.POST("/wallets/:id/transfers", func(c *gin.Context) {
		id := c.Param("id")
		senderId, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var req struct {
			RecipientId int     `json:"recipient_id" binding:"required"`
			Amount      float64 `json:"amount" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Call Wallet service to perform transfer
		_, err = walletClient.Transfer(context.Background(), &walletpb.TransferRequest{
			FromWalletId:    int32(senderId),
			ToWalletId: int32(req.RecipientId),
			Amount:      req.Amount,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Wallet transfer successful"})
	})

	r.Run(":8080")

}
