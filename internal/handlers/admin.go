package handlers

import (
	"cardapio-digital/internal/models"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var validate = validator.New()

// CreateRestaurant gerencia o POST do novo estabelecimento
func CreateRestaurant(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var restaurant models.Restaurant

		if err := c.BodyParser(&restaurant); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"erro": "Dados inválidos"})
		}

		// Validação de entrada
		if err := validate.Struct(&restaurant); err != nil {
			validationErrors := formatValidationErrors(err)
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"erro": "Validação falhou", "detalhes": validationErrors})
		}

		if err := db.Create(&restaurant).Error; err != nil {
			log.Printf("Erro ao criar restaurante: %v", err)
			status, errMsg := mapDatabaseError(err)
			return c.Status(status).JSON(fiber.Map{"erro": errMsg})
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

		// Validação de entrada
		if err := validate.Struct(&category); err != nil {
			validationErrors := formatValidationErrors(err)
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"erro": "Validação falhou", "detalhes": validationErrors})
		}

		// Valida se o restaurante existe
		var restaurant models.Restaurant
		if err := db.First(&restaurant, category.RestaurantID).Error; err != nil {
			if strings.Contains(err.Error(), "record not found") {
				return c.Status(http.StatusBadRequest).JSON(fiber.Map{"erro": "Restaurante não encontrado"})
			}
			log.Printf("Erro ao validar restaurante: %v", err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"erro": "Erro ao validar restaurante"})
		}

		if err := db.Create(&category).Error; err != nil {
			log.Printf("Erro ao salvar categoria: %v", err)
			status, errMsg := mapDatabaseError(err)
			return c.Status(status).JSON(fiber.Map{"erro": errMsg})
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

		// Validação de entrada
		if err := validate.Struct(&item); err != nil {
			validationErrors := formatValidationErrors(err)
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"erro": "Validação falhou", "detalhes": validationErrors})
		}

		// Valida se a categoria existe
		var category models.Category
		if err := db.First(&category, item.CategoryID).Error; err != nil {
			if strings.Contains(err.Error(), "record not found") {
				return c.Status(http.StatusBadRequest).JSON(fiber.Map{"erro": "Categoria não encontrada"})
			}
			log.Printf("Erro ao validar categoria: %v", err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"erro": "Erro ao validar categoria"})
		}

		if err := db.Create(&item).Error; err != nil {
			log.Printf("Erro ao salvar item: %v", err)
			status, errMsg := mapDatabaseError(err)
			return c.Status(status).JSON(fiber.Map{"erro": errMsg})
		}

		return c.Status(http.StatusCreated).JSON(item)
	}
}

// formatValidationErrors retorna um mapa dos erros de validação
func formatValidationErrors(err error) map[string]string {
	validationErrors := make(map[string]string)
	if validationErr, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range validationErr {
			validationErrors[fieldErr.Field()] = fmt.Sprintf("Tag: %s", fieldErr.Tag())
		}
	}
	return validationErrors
}

// mapDatabaseError mapeia erros de banco de dados para status HTTP apropriados
func mapDatabaseError(err error) (int, string) {
	errMsg := err.Error()

	// Violação de constraint de unicidade
	if strings.Contains(errMsg, "UNIQUE constraint failed") {
		return http.StatusConflict, "Este registro já existe (restrição de unicidade violada)"
	}

	// Violação de chave estrangeira
	if strings.Contains(errMsg, "FOREIGN KEY constraint failed") {
		return http.StatusBadRequest, "Referência inválida (chave estrangeira)"
	}

	// Violação de erro not null
	if strings.Contains(errMsg, "NOT NULL constraint failed") {
		return http.StatusBadRequest, "Campo obrigatório ausente"
	}

	// Erro de conexão ou timeout
	if strings.Contains(errMsg, "connection refused") || strings.Contains(errMsg, "timeout") {
		return http.StatusServiceUnavailable, "Serviço temporariamente indisponível"
	}

	// Erro genérico inesperado
	return http.StatusInternalServerError, "Erro ao processar solicitação"
}
