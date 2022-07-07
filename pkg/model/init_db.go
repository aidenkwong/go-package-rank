package model

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db  *gorm.DB
	dsn string
)

func init() {

	dsn = os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABAASE_URL must be set")
	}
	d, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		// Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}
	db = d
	db.AutoMigrate(&StandardPackage{})

}

func GetDB() *gorm.DB {
	return db
}
