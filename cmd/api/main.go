package main

import (
	"cardapio-digital/internal/database"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {

	db, err := database.Connect()
	if err != nil {
		log.Fatal("Fogo no parquinho: não foi possível inicializar o banco de dados: ", err)
	}

	_ = db

	app := fiber.New(fiber.Config{
		AppName: "Cardapio Digital API v1.0",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Requesição inválida ou mal formatada",
			})
		},
	})

	app.Use(logger.New())
	app.Use(recover.New())

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":    "Servidor rodando normalmente",
			"version":   "1.0",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	log.Fatal(app.Listen(":3000"))
}
