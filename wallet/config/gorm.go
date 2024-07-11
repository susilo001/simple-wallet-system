package config

import (
	"log"

	"github.com/susilo001/simple-wallet-system/wallet/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDatabases() {
	// setup gorm connection
	dsn := "postgresql://postgres:password@localhost:5432/postgres"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		log.Fatalln(err)
	}

	db.AutoMigrate(&entity.Wallet{})
}