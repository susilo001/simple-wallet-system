package repository

import (
	"context"
	"errors"
	"log"

	"github.com/susilo001/simple-wallet-system/wallet/entity"
	"github.com/susilo001/simple-wallet-system/wallet/service"
	"gorm.io/gorm"
)

type GormDBIface interface {
	WithContext(ctx context.Context) *gorm.DB
	Create(value interface{}) *gorm.DB
	First(dest interface{}, conds ...interface{}) *gorm.DB
	Save(value interface{}) *gorm.DB
	Delete(value interface{}, conds ...interface{}) *gorm.DB
	Find(dest interface{}, conds ...interface{}) *gorm.DB
}

type walletRepository struct {
	db GormDBIface
}

func NewWalletRepository(db GormDBIface) service.IWalletRepository {
	return &walletRepository{db: db}
}

func (r *walletRepository) CreateWallet(ctx context.Context, wallet *entity.Wallet) (entity.Wallet, error) {
	if err := r.db.WithContext(ctx).Create(wallet).Error; err != nil {
		log.Printf("Error creating wallet: %v\n", err)
		return entity.Wallet{}, err
	}
	return *wallet, nil
}

func (r *walletRepository) GetWalletByID(ctx context.Context, id int) (entity.Wallet, error) {
	var wallet entity.Wallet
	if err := r.db.WithContext(ctx).First(&wallet, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Wallet{}, nil
		}
		log.Printf("Error getting wallet by ID: %v\n", err)
		return entity.Wallet{}, err
	}
	return wallet, nil
}

func (r *walletRepository) UpdateWallet(ctx context.Context, id int, wallet entity.Wallet) (entity.Wallet, error) {
	var existingWallet entity.Wallet
	if err := r.db.WithContext(ctx).First(&existingWallet, id).Error; err != nil {
		log.Printf("Error finding wallet to update: %v\n", err)
		return entity.Wallet{}, err
	}

	existingWallet.UserID = wallet.UserID
	existingWallet.Balance = wallet.Balance
	if err := r.db.WithContext(ctx).Save(&existingWallet).Error; err != nil {
		log.Printf("Error updating wallet: %v\n", err)
		return entity.Wallet{}, err
	}
	return existingWallet, nil
}

func (r *walletRepository) DeleteWallet(ctx context.Context, id int) error {
	if err := r.db.WithContext(ctx).Delete(&entity.Wallet{}, id).Error; err != nil {
		log.Printf("Error deleting wallet: %v\n", err)
		return err
	}
	return nil
}

func (r *walletRepository) GetAllWallets(ctx context.Context) ([]entity.Wallet, error) {
	var wallets []entity.Wallet
	if err := r.db.WithContext(ctx).Find(&wallets).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return wallets, nil
		}
		log.Printf("Error getting all wallets: %v\n", err)
		return nil, err
	}
	return wallets, nil
}
