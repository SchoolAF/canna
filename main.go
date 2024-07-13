package main

import (
	"api/config"
	_ "api/docs"
	"api/internal/router"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"log"
	"os"
	"strings"
)

// @title Canna Got That Swag - API Documentation
// @version 1.0
// @description Canna API Documentation with Swagger

// @contact.name Halcyon Team
// @contact.url https://github.com/halcyonproject
// @contact.email dev@hlcyn.co

// @host canna.hlcyn.co
// @BasePath /
// @schemes https http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter your bearer token in the format "Bearer {token}"

func main() {
	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New())
	router.SetupRoutes(app)

	app.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Join(config.GetAllowedOrigins(), ","),
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))

	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404)
	})

	port := os.Getenv("PORT")

	err := app.Listen(":" + port)
	if err != nil {
		log.Fatalf("Error starting the server: %v", err)
	}
}
