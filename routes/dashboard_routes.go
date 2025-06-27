package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/maiscrianca/backend/controllers"
	"github.com/maiscrianca/backend/middleware"
)

// SetupDashboardRoutes configura as rotas do dashboard
func SetupDashboardRoutes(app *fiber.App) {
	dashboard := app.Group("/api/dashboard", middleware.AuthMiddleware)
	dashboard.Get("/metrics", controllers.GetDashboardMetrics)
}
