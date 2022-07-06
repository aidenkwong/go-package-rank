package model

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db  *gorm.DB
	dsn = "host=localhost user=postgres password=Highfive5+ dbname=go_package_rank port=5432 sslmode=disable"
)

func init() {
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
