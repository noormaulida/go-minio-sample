package main

import (
	"log"

	"go-minio-sample/pkg/config"
	"go-minio-sample/pkg/database"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// initiate config
	err := config.Load(".")
	if err != nil {
		log.Fatalf("Could not load environment variables: %s", err)
	}

	postgresErr := database.NewDB(config.ConfigData)
	if err != postgresErr {
		log.Fatalf("Postgresql initialization error: %s", err)
	}

	server := fiber.New()
	server.Static("/", "./static")

	server.Use(recover.New())
	server.Use(logger.New())

	server.Get("/healthchecker", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  true,
			"message": "🚀 Go Minio Sample is up",
		})
	})

	log.Fatal(server.Listen(":" + config.ConfigData.ServerPort))
}
