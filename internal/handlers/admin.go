package handlers

import (
	"cardapio-digital/internal/models"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// CreateRestaurant gerencia o POST do novo estabelecimento
func CreateRestaurant(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var restaurant models.Restaurant
		
		if err := c.BodyParser(&restaurant); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"erro": "Dados inválidos"})
		}

		if err := db.Create(&restaurant).Error; err != nil {
			return c.Status(http.StatusConflict).JSON(fiber.Map{"erro": "Nao foi possível criar o restaurante (verifique o nome ou URL)"})
		}

		return c.Status(http.StatusCreated).JSON(restaurant)
	}
}

// CreateCategory adiciona seções como "Bebidas" ou "Sobremesas"
func CreateCategory(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var category models.Category
		if err := c.BodyParser(&category); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"erro": "Dados inválidos"})
		}

		if err := db.Create(&category).Error; err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"erro": "Erro ao salvar categoria"})
		}

		return c.Status(http.StatusCreated).JSON(category)
	}
}

// CreateItem adiciona o produto final atrelado a uma categoria
func CreateItem(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var item models.Item
		if err := c.BodyParser(&item); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"erro": "Dados inválidos"})
		}

		if err := db.Create(&item).Error; err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"erro": "Erro ao salvar item"})
		}

		return c.Status(http.StatusCreated).JSON(item)
	}
}