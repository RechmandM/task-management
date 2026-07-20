package config

import (
	"fmt"
	"log"

	"github.com/rechmand/task-management/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		Env("DB_USER"),
		Env("DB_PASS"),
		Env("DB_HOST"),
		Env("DB_PORT"),
		Env("DB_NAME"),
	)

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Failed to connect database : ", err)
	}

	DB = database

	log.Println("✅ MySQL Connected")

	err = DB.AutoMigrate(&models.Task{})
	if err != nil {
		log.Fatal("❌ Auto Migration Failed : ", err)
	}

	log.Println("✅ Database Migrated")
}