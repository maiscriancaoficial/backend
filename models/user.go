package models

import (
	"time"
)

// Role representa o papel do usuário no sistema
type Role string

// Constantes para os diferentes papéis
const (
	CLIENT    Role = "CLIENT"
	EMPLOYEE  Role = "EMPLOYEE"
	ADMIN     Role = "ADMIN"
	AFFILIATE Role = "AFFILIATE"
)

// User representa o modelo de usuário no banco de dados
type User struct {
	ID            string    `json:"id"`
	Email         string    `json:"email"`
	Password      string    `json:"-"` // Não retornamos a senha nas respostas JSON
	Name          string    `json:"name"`
	Role          Role      `json:"role"`
	ProfileAvatar string    `json:"profileAvatar,omitempty"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// UserResponse representa a resposta para o cliente após autenticação
type UserResponse struct {
	ID            string    `json:"id"`
	Email         string    `json:"email"`
	Name          string    `json:"name"`
	Role          Role      `json:"role"`
	ProfileAvatar string    `json:"profileAvatar,omitempty"`
	Token         string    `json:"token"`
}

// LoginRequest representa a requisição de login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterRequest representa a requisição de registro
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Role     Role   `json:"role,omitempty"`
}
