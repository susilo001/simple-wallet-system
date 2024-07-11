package service

import (
	"context"
	"fmt"

	"github.com/susilo001/simple-wallet-system/wallet/entity"
)

// IWalletService defines the interface for the wallet service
type IWalletService interface {
	CreateWallet(ctx context.Context, wallet *entity.Wallet) (entity.Wallet, error)
	GetWalletByID(ctx context.Context, id int) (entity.Wallet, error)
	UpdateWallet(ctx context.Context, id int, wallet entity.Wallet) (entity.Wallet, error)
	DeleteWallet(ctx context.Context, id int) error
	GetAllWallets(ctx context.Context) ([]entity.Wallet, error)
}

// IWalletRepository defines the interface for the wallet repository
type IWalletRepository interface {
	CreateWallet(ctx context.Context, wallet *entity.Wallet) (entity.Wallet, error)
	GetWalletByID(ctx context.Context, id int) (entity.Wallet, error)
	UpdateWallet(ctx context.Context, id int, wallet entity.Wallet) (entity.Wallet, error)
	DeleteWallet(ctx context.Context, id int) error
	GetAllWallets(ctx context.Context) ([]entity.Wallet, error)
}

// walletService is an implementation of IWalletService that uses IWalletRepository
type walletService struct {
	walletRepo IWalletRepository
}

// NewWalletService creates a new instance of walletService
func NewWalletService(walletRepo IWalletRepository) IWalletService {
	return &walletService{walletRepo: walletRepo}
}

// CreateWallet creates a new wallet
func (s *walletService) CreateWallet(ctx context.Context, wallet *entity.Wallet) (entity.Wallet, error) {
	// Calls CreateWallet from the repository to create a new wallet
	createdWallet, err := s.walletRepo.CreateWallet(ctx, wallet)
	if err != nil {
		return entity.Wallet{}, fmt.Errorf("failed to create wallet: %v", err)
	}
	return createdWallet, nil
}

// GetWalletByID gets a wallet by ID
func (s *walletService) GetWalletByID(ctx context.Context, id int) (entity.Wallet, error) {
	// Calls GetWalletByID from the repository to get a wallet by ID
	wallet, err := s.walletRepo.GetWalletByID(ctx, id)
	if err != nil {
		return entity.Wallet{}, fmt.Errorf("failed to get wallet by ID: %v", err)
	}
	return wallet, nil
}

// UpdateWallet updates wallet data
func (s *walletService) UpdateWallet(ctx context.Context, id int, wallet entity.Wallet) (entity.Wallet, error) {
	// Calls UpdateWallet from the repository to update wallet data
	updatedWallet, err := s.walletRepo.UpdateWallet(ctx, id, wallet)
	if err != nil {
		return entity.Wallet{}, fmt.Errorf("failed to update wallet: %v", err)
	}
	return updatedWallet, nil
}

// DeleteWallet deletes a wallet by ID
func (s *walletService) DeleteWallet(ctx context.Context, id int) error {
	// Calls DeleteWallet from the repository to delete a wallet by ID
	err := s.walletRepo.DeleteWallet(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete wallet: %v", err)
	}
	return nil
}

// GetAllWallets gets all wallets
func (s *walletService) GetAllWallets(ctx context.Context) ([]entity.Wallet, error) {
	// Calls GetAllWallets from the repository to get all wallets
	wallets, err := s.walletRepo.GetAllWallets(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all wallets: %v", err)
	}
	return wallets, nil
}
