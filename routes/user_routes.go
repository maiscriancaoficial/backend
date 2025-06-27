package routes

import (
	"github.com/WBianchi/maiscrianca/configs"
	"github.com/WBianchi/maiscrianca/controllers"
	"github.com/WBianchi/maiscrianca/middleware"
	"github.com/WBianchi/maiscrianca/models"
	"github.com/gofiber/fiber/v2"
)

// SetupUserRoutes configura as rotas de usuário
func SetupUserRoutes(app *fiber.App, userController *controllers.UserController, config *configs.Config) {
	api := app.Group("/api")
	
	// Rotas protegidas por autenticação
	user := api.Group("/user", middleware.AuthMiddleware(config))
	
	// Rotas para todos os usuários autenticados
	user.Get("/profile", userController.GetUserProfile)
	user.Put("/profile", userController.UpdateUserProfile)
	
	// Rotas protegidas por role
	admin := api.Group("/admin", middleware.AuthMiddleware(config), middleware.RoleGuard(models.ADMIN))
	admin.Get("/dashboard-data", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Dados do dashboard administrativo",
		})
	})
	
	employee := api.Group("/employee", middleware.AuthMiddleware(config), middleware.RoleGuard(models.EMPLOYEE))
	employee.Get("/dashboard-data", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Dados do dashboard de funcionário",
		})
	})
	
	client := api.Group("/client", middleware.AuthMiddleware(config), middleware.RoleGuard(models.CLIENT))
	client.Get("/dashboard-data", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Dados do dashboard de cliente",
		})
	})
	
	affiliate := api.Group("/affiliate", middleware.AuthMiddleware(config), middleware.RoleGuard(models.AFFILIATE))
	affiliate.Get("/dashboard-data", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Dados do dashboard de afiliado",
		})
	})
}
