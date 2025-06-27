package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/WBianchi/maiscrianca/controllers"
	"github.com/WBianchi/maiscrianca/middleware"
)

// SetupLivrosRoutes configura as rotas para gestão de livros
func SetupLivrosRoutes(app *fiber.App) {
	livros := app.Group("/api/livros", middleware.AuthMiddleware)
	
	// Rotas de livros
	livros.Get("/", controllers.GetLivros)
	livros.Post("/", controllers.CreateLivro)
	livros.Get("/:id", controllers.GetLivro)
	livros.Put("/:id", controllers.UpdateLivro)
	livros.Delete("/:id", controllers.DeleteLivro)
	
	// Rotas para upload de imagens e arquivos usando vercel blob
	livros.Post("/upload/capa", controllers.UploadCapa)
	livros.Post("/upload/arquivo", controllers.UploadArquivo)
	livros.Post("/upload/pagina", controllers.UploadPagina)
	
	// Rotas de categorias
	categorias := app.Group("/api/categorias", middleware.AuthMiddleware)
	categorias.Get("/", controllers.GetCategorias)
	categorias.Post("/", controllers.CreateCategoria)
	categorias.Put("/:id", controllers.UpdateCategoria)
	categorias.Delete("/:id", controllers.DeleteCategoria)
}
