package handler

import (
	"context"
	"fmt"
	"log"

	"github.com/susilo001/simple-wallet-system/wallet/entity"
	pb "github.com/susilo001/simple-wallet-system/wallet/proto/wallet/v1"
	"github.com/susilo001/simple-wallet-system/wallet/service"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type WalletHandler struct {
	pb.UnimplementedWalletServiceServer
	walletService service.IWalletService
}

func NewWalletHandler(walletService service.IWalletService) *WalletHandler {
	return &WalletHandler{
		walletService: walletService,
	}
}

func (h *WalletHandler) GetWallet(ctx context.Context, req *pb.GetWalletRequest) (*pb.GetWalletResponse, error) {
	wallet, err := h.walletService.GetWalletByID(ctx, int(req.GetWalletId()))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	res := &pb.GetWalletResponse{
		Wallet: &pb.Wallet{
			Id:        int32(wallet.ID),
			UserId:    int32(wallet.UserID),
			Balance:   wallet.Balance,
			CreatedAt: timestamppb.New(wallet.CreatedAt),
			UpdatedAt: timestamppb.New(wallet.UpdatedAt),
		},
	}
	return res, nil
}

func (h *WalletHandler) CreateWallet(ctx context.Context, req *pb.CreateWalletRequest) (*pb.MutationResponse, error) {
	createdWallet, err := h.walletService.CreateWallet(ctx, &entity.Wallet{
		UserID: int(req.GetUserId()),
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &pb.MutationResponse{
		Message: fmt.Sprintf("Success created wallet with ID %d", createdWallet.ID),
	}, nil
}

func (h *WalletHandler) UpdateWallet(ctx context.Context, req *pb.UpdateWalletRequest) (*pb.MutationResponse, error) {
	updatedWallet, err := h.walletService.UpdateWallet(ctx, int(req.UserId), entity.Wallet{
		UserID:  int(req.UserId),
		Balance: req.Balance,
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &pb.MutationResponse{
		Message: fmt.Sprintf("Success updated wallet with ID %d", updatedWallet.ID),
	}, nil
}

func (h *WalletHandler) GetBalance(ctx context.Context, req *pb.GetBalanceRequest) (*pb.GetBalanceResponse, error) {
	wallet, err := h.walletService.GetWalletByID(ctx, int(req.GetWalletId()))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &pb.GetBalanceResponse{
		Balance: wallet.Balance,
	}, nil
}

func (h *WalletHandler) TopupWallet(ctx context.Context, req *pb.TopupRequest) (*pb.MutationResponse, error) {
	err := h.walletService.TopupWallet(ctx, int(req.GetWalletId()), req.GetAmount())
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &pb.MutationResponse{
		Message: fmt.Sprintf("Successfully topped up wallet with ID %d", req.GetWalletId()),
	}, nil
}

func (h *WalletHandler) Transfer(ctx context.Context, req *pb.TransferRequest) (*pb.MutationResponse, error) {
	err := h.walletService.Transfer(ctx, int(req.GetFromWalletId()), int(req.GetToWalletId()), req.GetAmount())
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &pb.MutationResponse{
		Message: fmt.Sprintf("Successfully transferred from wallet ID %d to wallet ID %d", req.GetFromWalletId(), req.GetToWalletId()),
	}, nil
}

func (h *WalletHandler) GetTransactions(ctx context.Context, req *pb.GetTransactionsRequest) (*pb.GetTransactionsResponse, error) {
	transactions, err := h.walletService.GetTransactions(ctx, int(req.GetWalletId()))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var pbTransactions []*pb.Transaction
	for _, transaction := range transactions {
		pbTransactions = append(pbTransactions, &pb.Transaction{
			Id:           int32(transaction.ID),
			FromWalletId: int32(transaction.FromWalletID),
			ToWalletId:   int32(transaction.ToWalletID),
			Amount:       transaction.Amount,
			CreatedAt:    timestamppb.New(transaction.CreatedAt),
		})
	}
	return &pb.GetTransactionsResponse{
		Transactions: pbTransactions,
	}, nil
}
