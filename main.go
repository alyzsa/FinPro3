package main

import (
	"github.com/alyzsa/FinPro3/database"
	"github.com/alyzsa/FinPro3/entity"
	"github.com/alyzsa/FinPro3/router"
	"fmt"
	"log"
	"os"

	"gorm.io/gorm"
)

func main() {
	database.StartDB()
	var db = database.GetDB()
	seedAdminData(db)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := router.StartApp()
	r.Run(":" + port)
}

func seedAdminData(db *gorm.DB) {
	var adminCount int64
	var User []entity.User
	db.Model(&User).Where("role = ?", "admin").Count(&adminCount)

	if adminCount == 0 {
		admin := entity.User{
			Full_Name: "admin",
			Email:     "admin@gmail.com",
			Password:  "admin123",
			Role:      "admin",
		}

		if err := db.Create(&admin).Error; err != nil {
			log.Fatal(err)
		}

		fmt.Println("Admin user seeded successfully.")
	} else {
		fmt.Println("Admin user already exists.")
	}
}
