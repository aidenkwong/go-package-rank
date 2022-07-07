package model

import "gorm.io/gorm"

type StandardPackage struct {
	gorm.Model
	Name       string `gorm:"not null;unique"`
	ImportedBy int    `gorm:"not null"`
}

type APIStandardPackage struct {
	Name       string `json:"name"`
	ImportedBy int    `json:"imported_by"`
}

func GetAllStandardPackages() []APIStandardPackage {
	var stdPkgs []APIStandardPackage
	db := GetDB()
	db.Model(&StandardPackage{}).Select([]string{"name", "imported_by"}).Find(&stdPkgs)
	return stdPkgs
}
