package controllers

import (
	"database/sql"
	"log"

	"github.com/WBianchi/maiscrianca/auth"
	"github.com/WBianchi/maiscrianca/configs"
	"github.com/WBianchi/maiscrianca/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// AuthController gerencia a autenticação de usuários
type AuthController struct {
	DB     *sql.DB
	Config *configs.Config
}

// NewAuthController cria uma nova instância de AuthController
func NewAuthController(db *sql.DB, config *configs.Config) *AuthController {
	return &AuthController{
		DB:     db,
		Config: config,
	}
}

// Login autentica um usuário
func (c *AuthController) Login(ctx *fiber.Ctx) error {
	log.Println("Endpoint de login chamado")
	var loginRequest models.LoginRequest

	if err := ctx.BodyParser(&loginRequest); err != nil {
		log.Printf("Erro ao fazer parse do body: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Dados de login inválidos",
		})
	}

	log.Printf("Tentativa de login para email: %s", loginRequest.Email)

	// Validação básica
	if loginRequest.Email == "" || loginRequest.Password == "" {
		log.Println("Email ou senha vazios")
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email e senha são obrigatórios",
		})
	}

	// Verificar conexão com o DB primeiro
	if err := c.DB.Ping(); err != nil {
		log.Printf("Erro na conexão com o banco de dados: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro na conexão com o banco de dados",
		})
	}

	log.Println("Conexão com o banco de dados verificada com sucesso")

	// Buscar usuário pelo email
	var user models.User
	var hashedPassword string
	// Usamos NullString para o campo que pode ser nulo
	var profileAvatar sql.NullString

	query := `SELECT id, email, password, name, role, "profileAvatar", "createdAt", "updatedAt" 
              FROM "User" WHERE email = $1`
	
	log.Printf("Executando query: %s com email: %s", query, loginRequest.Email)
	
	err := c.DB.QueryRow(query, loginRequest.Email).Scan(
		&user.ID,
		&user.Email,
		&hashedPassword,
		&user.Name,
		&user.Role,
		&profileAvatar,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	
	// Se o valor for válido (não nulo), atribuímos à estrutura user
	if profileAvatar.Valid {
		user.ProfileAvatar = profileAvatar.String
	} else {
		user.ProfileAvatar = ""
	}

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Usuário não encontrado para o email: %s", loginRequest.Email)
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Credenciais inválidas",
			})
		}
		log.Printf("Erro ao buscar usuário: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro interno do servidor: " + err.Error(),
		})
	}

	// Verificar senha
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(loginRequest.Password))
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Credenciais inválidas",
		})
	}

	// Gerar token JWT
	token, err := auth.GenerateToken(&user, c.Config)
	if err != nil {
		log.Printf("Erro ao gerar token: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro ao gerar token de autenticação",
		})
	}

	// Verificar a página de redirecionamento com base na role
	var redirectUrl string
	switch user.Role {
	case models.ADMIN:
		redirectUrl = "/dashboard"
	case models.CLIENT:
		redirectUrl = "/cliente/visao-geral"
	case models.EMPLOYEE:
		redirectUrl = "/funcionario/visao-geral"
	case models.AFFILIATE:
		redirectUrl = "/afiliado/visao-geral"
	default:
		redirectUrl = "/cliente/visao-geral"
	}

	// Resposta com token e informações do usuário
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"user": models.UserResponse{
			ID:            user.ID,
			Email:         user.Email,
			Name:          user.Name,
			Role:          user.Role,
			ProfileAvatar: user.ProfileAvatar,
			Token:         token,
		},
		"redirectUrl": redirectUrl,
	})
}

// Register registra um novo usuário
func (c *AuthController) Register(ctx *fiber.Ctx) error {
	log.Println("Endpoint de registro chamado")
	var registerRequest models.RegisterRequest

	if err := ctx.BodyParser(&registerRequest); err != nil {
		log.Printf("Erro ao fazer parse do body: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Dados de registro inválidos",
		})
	}
	log.Printf("Dados de registro recebidos: %+v", registerRequest)

	// Validação básica
	if registerRequest.Email == "" || registerRequest.Password == "" || registerRequest.Name == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email, senha e nome são obrigatórios",
		})
	}

	// Verificar conexão com o DB primeiro
	if err := c.DB.Ping(); err != nil {
		log.Printf("Erro na conexão com o banco de dados: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro na conexão com o banco de dados",
		})
	}
	log.Println("Conexão com o banco de dados verificada com sucesso")

	// Verificar se o email já está em uso
	var count int
	emailCheckQuery := `SELECT COUNT(*) FROM "User" WHERE email = $1`
	log.Printf("Executando verificação de email com query: %s e email: %s", emailCheckQuery, registerRequest.Email)
	err := c.DB.QueryRow(emailCheckQuery, registerRequest.Email).Scan(&count)
	if err != nil {
		log.Printf("Erro ao verificar email: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro interno do servidor",
		})
	}

	if count > 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Este email já está em uso",
		})
	}

	// Hash da senha
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Erro ao gerar hash da senha: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro ao processar senha",
		})
	}

	// Definir role (default CLIENT)
	role := models.CLIENT
	if registerRequest.Role != "" {
		role = registerRequest.Role
	}

	// Inserir novo usuário
	id := uuid.New().String()
	log.Printf("ID gerado para novo usuário: %s", id)
	
	insertQuery := `INSERT INTO "User" (id, email, password, name, role, "createdAt", "updatedAt") 
              VALUES ($1, $2, $3, $4, $5, NOW(), NOW()) 
              RETURNING id, email, name, role, "profileAvatar", "createdAt", "updatedAt"`

	log.Printf("Executando query de inserção: %s", insertQuery)
	log.Printf("Parâmetros: ID=%s, Email=%s, Nome=%s, Role=%s", 
		id, registerRequest.Email, registerRequest.Name, role)

	var user models.User
	// Usamos um ponteiro para NullString para o campo que pode ser nulo
	var profileAvatar sql.NullString

	err = c.DB.QueryRow(
		insertQuery, 
		id, 
		registerRequest.Email, 
		string(hashedPassword), 
		registerRequest.Name, 
		string(role),
	).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Role,
		&profileAvatar,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	
	// Se o valor for válido (não nulo), atribuímos à estrutura user
	if profileAvatar.Valid {
		user.ProfileAvatar = profileAvatar.String
	} else {
		user.ProfileAvatar = ""
	}

	if err != nil {
		log.Printf("Erro ao inserir usuário: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro ao criar usuário: " + err.Error(),
		})
	}
	
	log.Printf("Usuário criado com sucesso: %+v", user)

	// Gerar token JWT
	token, err := auth.GenerateToken(&user, c.Config)
	if err != nil {
		log.Printf("Erro ao gerar token: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro ao gerar token de autenticação",
		})
	}

	// Verificar a página de redirecionamento com base na role
	var redirectUrl string
	switch user.Role {
	case models.ADMIN:
		redirectUrl = "/dashboard"
	case models.CLIENT:
		redirectUrl = "/cliente/visao-geral"
	case models.EMPLOYEE:
		redirectUrl = "/funcionario/visao-geral"
	case models.AFFILIATE:
		redirectUrl = "/afiliado/visao-geral"
	default:
		redirectUrl = "/cliente/visao-geral"
	}

	// Resposta com token e informações do usuário
	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"user": models.UserResponse{
			ID:            user.ID,
			Email:         user.Email,
			Name:          user.Name,
			Role:          user.Role,
			ProfileAvatar: user.ProfileAvatar,
			Token:         token,
		},
		"redirectUrl": redirectUrl,
	})
}
