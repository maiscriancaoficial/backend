package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/WBianchi/maiscrianca/configs"
	"github.com/WBianchi/maiscrianca/controllers"
	"github.com/WBianchi/maiscrianca/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Carregar variáveis de ambiente
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Erro ao carregar arquivo .env:", err)
	}

	// Carregar configurações
	config := configs.LoadConfig()

	// Configurar conexão com o PostgreSQL
	db, err := sql.Open("postgres", config.DatabaseURL)
	if err != nil {
		log.Fatal("Erro ao conectar ao banco de dados:", err)
	}
	defer db.Close()

	// Verificar conexão com o banco
	err = db.Ping()
	if err != nil {
		log.Fatal("Erro ao testar conexão com o banco de dados:", err)
	}
	fmt.Println("Conexão com o PostgreSQL estabelecida com sucesso!")

	// Inicializar controladores
	authController := controllers.NewAuthController(db, config)
	userController := controllers.NewUserController(db)

	// Inicializar o aplicativo Fiber
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// Tratamento de erros padrão
			code := fiber.StatusInternalServerError

			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(logger.New())
	
	// Configuração de CORS
	allowOrigins := config.AllowedOrigins
	if allowOrigins == "" {
		allowOrigins = "http://localhost:3000"
	}
	
	app.Use(cors.New(cors.Config{
		AllowOrigins: allowOrigins,
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE",
		AllowCredentials: true,
	}))

	// Rota de healthcheck
	app.Get("/api/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"message": "API do Mais Criança está funcionando corretamente",
		})
	})

	// Configurar rotas
	routes.SetupAuthRoutes(app, authController)
	routes.SetupUserRoutes(app, userController, config)

	// Iniciar o servidor
	port := config.Port
	fmt.Printf("Servidor iniciado na porta %s\n", port)
	log.Fatal(app.Listen(":" + port))
}
