package model

import (
	"log"
	"os"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db  *gorm.DB
	dsn string
)

func init() {

	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if pair[0] == "DATABASE_URL" {
			dsn = pair[1]
			break
		}
	}
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
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
