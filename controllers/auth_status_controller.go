package controllers

import (
	"github.com/WBianchi/maiscrianca/configs"
	"github.com/WBianchi/maiscrianca/models"
	"github.com/gofiber/fiber/v2"
)

// AuthStatusController fornece endpoints para verificar o status da autenticação
type AuthStatusController struct {
	Config *configs.Config
}

// NewAuthStatusController cria uma nova instância de AuthStatusController
func NewAuthStatusController(config *configs.Config) *AuthStatusController {
	return &AuthStatusController{
		Config: config,
	}
}

// GetAuthStatus retorna o status de autenticação do usuário e para onde redirecionar com base na role
func (c *AuthStatusController) GetAuthStatus(ctx *fiber.Ctx) error {
	// Recupera a role do usuário a partir do middleware de autenticação
	userRole, ok := ctx.Locals("userRole").(models.Role)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"authenticated": false,
			"message":       "Usuário não autenticado",
		})
	}

	// Determina a URL para redirecionamento com base na role
	var redirectUrl string
	switch userRole {
	case models.ADMIN:
		redirectUrl = "/dashboard"
	case models.CLIENT:
		redirectUrl = "/cliente/visao-geral"
	case models.EMPLOYEE:
		redirectUrl = "/funcionario/visao-geral"
	case models.AFFILIATE:
		redirectUrl = "/afiliado/visao-geral"
	default:
		redirectUrl = "/"
	}

	// Retorna as informações de autenticação
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"authenticated": true,
		"redirectUrl":   redirectUrl,
	})
}
