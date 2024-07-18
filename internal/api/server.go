package api

import (
	"go-ecommerce-app/config"
	"go-ecommerce-app/internal/api/rest"
	"go-ecommerce-app/internal/api/rest/handlers"
	"go-ecommerce-app/internal/domain"
	"go-ecommerce-app/internal/helper"
	"log"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func StartServer(config config.AppConfig) {
	app := fiber.New()

	db, err := gorm.Open(postgres.Open(config.Dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln("Database connection error: ", err)
	}

	log.Println("Database connected")
	// run migrations
	err = db.AutoMigrate(&domain.User{}, &domain.BankAccount{})
	if err != nil {
		log.Fatalf("Error on running migration: %v", err.Error())
	}
	log.Println("migration was succefull")

	auth := helper.SetupAuth(config.AppSecret)

	rh := &rest.RestHandler{
		App:    app,
		DB:     db,
		Auth:   auth,
		Config: config,
	}
	setupRoutes(rh)

	app.Listen(config.ServerPort)
}

func setupRoutes(rh *rest.RestHandler) {
	// user handler
	handlers.SetupUserRoutes(rh)
	// transactions
	// catalouges

}
