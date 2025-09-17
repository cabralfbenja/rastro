package main

import (
	"rastro/cmd/api/requests"
	"rastro/cmd/api/services"
	"rastro/common"
)

func main() {
	db, err := common.NewMySql()

	if err != nil {
		panic("Failed to connect to database")
	}

	// Initialize the CategoryService with the database connection
	service := services.CategoryService{
		DB: db,
	}
	categories := []string{"Food", "Gifts", "Transport", "Utilities", "Entertainment", "Health", "Education", "Travel", "Shopping", "Miscellaneous"}

	for _, name := range categories {
		_, err := service.Create(&requests.CreateCategoryRequest{Name: name, IsCustom: false})
		if err != nil {
			panic("Failed to seed categories: " + err.Error())
		}
		println("Seeded category ", name, " created successfully")
	}
}
