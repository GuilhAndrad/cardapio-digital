package handlers

import (
	"cardapio-digital/internal/models"
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// GetMenuBySlug busca o restaurante, suas categorias e seus respectivos itens de uma só vez
func GetMenuBySlug(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		slug := c.Params("slug")
		var restaurant models.Restaurant

		// O pulo do gato: aninhamos o Preload para trazer Categorias -> Itens
		err := db.Preload("Categories.Items").Where("slug = ?", slug).First(&restaurant).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(http.StatusNotFound).JSON(fiber.Map{"erro": "Cardápio não encontrado"})
			}
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"erro": "Erro interno no servidor"})
		}

		return c.JSON(restaurant)
	}
}
