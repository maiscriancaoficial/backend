package routes

import (
	"github.com/WBianchi/maiscrianca/controllers"
	"github.com/gofiber/fiber/v2"
)

// SetupAuthRoutes configura as rotas de autenticação
func SetupAuthRoutes(app *fiber.App, authController *controllers.AuthController) {
	auth := app.Group("/api/auth")
	
	// Rotas públicas
	auth.Post("/login", authController.Login)
	auth.Post("/register", authController.Register)
}
