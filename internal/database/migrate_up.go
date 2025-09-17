package main

import (
	"rastro/common"
	"rastro/internal/models"
)

func main() {
	db, err := common.NewMySql()

	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&models.UserModel{}, &models.AppTokenModel{}, &models.CategoryModel{})

	if err != nil {
		panic(err)
	}
}
