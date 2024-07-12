package service

import (
	"context"
	"fmt"

	"github.com/susilo001/simple-wallet-system/wallet/entity"
)

type IWalletService interface {
	CreateWallet(ctx context.Context, wallet *entity.Wallet) (entity.Wallet, error)
	GetWalletByID(ctx context.Context, id int) (entity.Wallet, error)
	UpdateWallet(ctx context.Context, id int, wallet entity.Wallet) (entity.Wallet, error)
	TopUpWallet(ctx context.Context, walletID int, amount float64) error
	Transfer(ctx context.Context, fromWalletID int, toWalletID int, amount float64) error
	GetTransactions(ctx context.Context, walletID int) ([]entity.Transaction, error)
}

type IWalletRepository interface {
	CreateWallet(ctx context.Context, wallet *entity.Wallet) (entity.Wallet, error)
	GetWalletByID(ctx context.Context, id int) (entity.Wallet, error)
	UpdateWallet(ctx context.Context, id int, wallet entity.Wallet) (entity.Wallet, error)
	DeleteWallet(ctx context.Context, id int) error
	GetAllWallets(ctx context.Context) ([]entity.Wallet, error)
	TopupWallet(ctx context.Context, walletID int, amount float64) error
	Transfer(ctx context.Context, fromWalletID int, toWalletID int, amount float64) error
	GetTransactions(ctx context.Context, walletID int) ([]entity.Transaction, error)
}

type walletService struct {
	walletRepo IWalletRepository
}

func NewWalletService(walletRepo IWalletRepository) IWalletService {
	return &walletService{walletRepo: walletRepo}
}

func (s *walletService) CreateWallet(ctx context.Context, wallet *entity.Wallet) (entity.Wallet, error) {
	createdWallet, err := s.walletRepo.CreateWallet(ctx, wallet)
	if err != nil {
		return entity.Wallet{}, fmt.Errorf("failed to create wallet: %v", err)
	}
	return createdWallet, nil
}

func (s *walletService) GetWalletByID(ctx context.Context, id int) (entity.Wallet, error) {
	wallet, err := s.walletRepo.GetWalletByID(ctx, id)
	if err != nil {
		return entity.Wallet{}, fmt.Errorf("failed to get wallet by ID: %v", err)
	}
	return wallet, nil
}

func (s *walletService) UpdateWallet(ctx context.Context, id int, wallet entity.Wallet) (entity.Wallet, error) {
	updatedWallet, err := s.walletRepo.UpdateWallet(ctx, id, wallet)
	if err != nil {
		return entity.Wallet{}, fmt.Errorf("failed to update wallet: %v", err)
	}
	return updatedWallet, nil
}

func (s *walletService) TopUpWallet(ctx context.Context, walletID int, amount float64) error {
	err := s.walletRepo.TopupWallet(ctx, walletID, amount)
	if err != nil {
		return fmt.Errorf("failed to top up wallet: %v", err)
	}
	return nil
}

func (s *walletService) Transfer(ctx context.Context, fromWalletID int, toWalletID int, amount float64) error {
	err := s.walletRepo.Transfer(ctx, fromWalletID, toWalletID, amount)
	if err != nil {
		return fmt.Errorf("failed to transfer amount: %v", err)
	}
	return nil
}

func (s *walletService) GetTransactions(ctx context.Context, walletID int) ([]entity.Transaction, error) {
	transactions, err := s.walletRepo.GetTransactions(ctx, walletID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %v", err)
	}
	return transactions, nil
}

func (s *walletService) DeleteWallet(ctx context.Context, id int) error {
	err := s.walletRepo.DeleteWallet(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete wallet: %v", err)
	}
	return nil
}

func (s *walletService) GetAllWallets(ctx context.Context) ([]entity.Wallet, error) {
	wallets, err := s.walletRepo.GetAllWallets(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all wallets: %v", err)
	}
	return wallets, nil
}
