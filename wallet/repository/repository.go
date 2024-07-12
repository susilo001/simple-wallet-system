package repository

import (
	"context"
	"errors"
	"log"
	"time"

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

func (r *walletRepository) TopUpWallet(ctx context.Context, walletID int, amount float64) error {
	var wallet entity.Wallet
	if err := r.db.WithContext(ctx).First(&wallet, walletID).Error; err != nil {
		log.Printf("Error finding wallet for top-up: %v\n", err)
		return err
	}

	wallet.Balance += amount
	if err := r.db.WithContext(ctx).Save(&wallet).Error; err != nil {
		log.Printf("Error updating wallet balance for top-up: %v\n", err)
		return err
	}

	return r.createTransaction(ctx, walletID, 0, amount)
}

func (r *walletRepository) Transfer(ctx context.Context, senderID int, recipientID int, amount float64) error {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		log.Printf("Error beginning transaction: %v\n", tx.Error)
		return tx.Error
	}

	var senderWallet entity.Wallet
	if err := tx.First(&senderWallet, recipientID).Error; err != nil {
		log.Printf("Error finding sender's wallet for transfer: %v\n", err)
		tx.Rollback()
		return err
	}

	if senderWallet.Balance < amount {
		tx.Rollback()
		return errors.New("insufficient balance")
	}

	senderWallet.Balance -= amount
	if err := tx.Save(&senderWallet).Error; err != nil {
		log.Printf("Error updating sender's wallet balance: %v\n", err)
		tx.Rollback()
		return err
	}

	var toWallet entity.Wallet
	if err := tx.First(&toWallet, recipientID).Error; err != nil {
		log.Printf("Error finding receiver's wallet for transfer: %v\n", err)
		tx.Rollback()
		return err
	}

	toWallet.Balance += amount
	if err := tx.Save(&toWallet).Error; err != nil {
		log.Printf("Error updating receiver's wallet balance: %v\n", err)
		tx.Rollback()
		return err
	}

	if err := r.createTransaction(ctx, senderID, recipientID, amount); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *walletRepository) GetTransactions(ctx context.Context, walletID int) ([]entity.Transaction, error) {
	var transactions []entity.Transaction
	if err := r.db.WithContext(ctx).Where("sender_id = ? OR recipient_id = ?", walletID, walletID).Find(&transactions).Error; err != nil {
		log.Printf("Error getting transactions: %v\n", err)
		return nil, err
	}
	return transactions, nil
}

func (r *walletRepository) createTransaction(ctx context.Context, senderID int, recipientID int, amount float64) error {
	transaction := entity.Transaction{
		SenderID:    senderID,
		RecipientID: recipientID,
		Amount:      amount,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := r.db.WithContext(ctx).Create(&transaction).Error; err != nil {
		log.Printf("Error creating transaction: %v\n", err)
		return err
	}

	return nil
}
