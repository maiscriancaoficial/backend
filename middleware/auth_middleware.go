package middleware

import (
	"strings"

	"github.com/WBianchi/maiscrianca/auth"
	"github.com/WBianchi/maiscrianca/configs"
	"github.com/WBianchi/maiscrianca/models"
	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware verifica se o usuário está autenticado
func AuthMiddleware(config *configs.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Não autorizado: token não fornecido",
			})
		}
		
		// Verificar formato "Bearer {token}"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Formato de autorização inválido",
			})
		}
		
		tokenString := parts[1]
		claims, err := auth.ValidateToken(tokenString, config)
		
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token inválido: " + err.Error(),
			})
		}
		
		// Adicionar dados do usuário ao contexto
		c.Locals("userId", claims.UserID)
		c.Locals("userRole", claims.Role)
		
		return c.Next()
	}
}

// RoleGuard verifica se o usuário possui a role necessária
func RoleGuard(roles ...models.Role) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole, ok := c.Locals("userRole").(models.Role)
		
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Informações de autenticação ausentes",
			})
		}
		
		for _, role := range roles {
			if userRole == role {
				return c.Next()
			}
		}
		
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Acesso negado: permissão insuficiente",
		})
	}
}
