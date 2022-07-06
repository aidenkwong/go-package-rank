package model

import "gorm.io/gorm"

type StandardPackage struct {
	gorm.Model
	Name         string `gorm:"not null;unique_index"`
	NumOfImports int    `gorm:"not null"`
}

type APIStandardPackage struct {
	Name         string `json:"name"`
	NumOfImports int    `json:"num_of_imports"`
}

func GetAllStandardPackages() []APIStandardPackage {
	var stdPkgs []APIStandardPackage
	db := GetDB()
	db.Model(&StandardPackage{}).Select([]string{"name", "num_of_imports"}).Find(&stdPkgs)
	return stdPkgs
}
