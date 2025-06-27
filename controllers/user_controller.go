package controllers

import (
	"database/sql"
	"log"

	"github.com/WBianchi/maiscrianca/models"
	"github.com/gofiber/fiber/v2"
)

// UserController gerencia as operações relacionadas ao usuário
type UserController struct {
	DB *sql.DB
}

// NewUserController cria uma nova instância de UserController
func NewUserController(db *sql.DB) *UserController {
	return &UserController{
		DB: db,
	}
}

// GetUserProfile recupera o perfil do usuário autenticado
func (c *UserController) GetUserProfile(ctx *fiber.Ctx) error {
	userId := ctx.Locals("userId").(string)
	
	var user models.User
	query := `SELECT id, email, name, role, profile_avatar, created_at, updated_at 
			  FROM "User" WHERE id = $1`
	
	err := c.DB.QueryRow(query, userId).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Role,
		&user.ProfileAvatar,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Usuário não encontrado",
			})
		}
		log.Printf("Erro ao buscar usuário: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro interno do servidor",
		})
	}

	return ctx.JSON(fiber.Map{
		"user": user,
	})
}

// UpdateUserProfile atualiza o perfil do usuário
func (c *UserController) UpdateUserProfile(ctx *fiber.Ctx) error {
	userId := ctx.Locals("userId").(string)
	
	// Verifica se o usuário existe
	var exists bool
	err := c.DB.QueryRow(`SELECT EXISTS(SELECT 1 FROM "User" WHERE id = $1)`, userId).Scan(&exists)
	
	if err != nil {
		log.Printf("Erro ao verificar usuário: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro interno do servidor",
		})
	}
	
	if !exists {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Usuário não encontrado",
		})
	}
	
	// Extrair dados da requisição
	type UpdateRequest struct {
		Name          string `json:"name"`
		ProfileAvatar string `json:"profileAvatar"`
	}
	
	var req UpdateRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Dados inválidos",
		})
	}
	
	// Atualiza apenas os campos fornecidos
	if req.Name != "" || req.ProfileAvatar != "" {
		query := `UPDATE "User" SET updated_at = NOW()`
		args := []interface{}{}
		paramCount := 1
		
		if req.Name != "" {
			query += `, name = $` + string(rune('0'+paramCount))
			args = append(args, req.Name)
			paramCount++
		}
		
		if req.ProfileAvatar != "" {
			query += `, profile_avatar = $` + string(rune('0'+paramCount))
			args = append(args, req.ProfileAvatar)
			paramCount++
		}
		
		query += ` WHERE id = $` + string(rune('0'+paramCount))
		args = append(args, userId)
		
		_, err = c.DB.Exec(query, args...)
		if err != nil {
			log.Printf("Erro ao atualizar usuário: %v", err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Erro ao atualizar perfil",
			})
		}
	}
	
	return ctx.JSON(fiber.Map{
		"message": "Perfil atualizado com sucesso",
	})
}
