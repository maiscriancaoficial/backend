package configs

import (
	"os"
	"time"
)

// Config armazena todas as configurações da aplicação
type Config struct {
	JWTSecret           string
	JWTExpirationHours  time.Duration
	DatabaseURL         string
	AllowedOrigins      string
	Port                string
}

// LoadConfig carrega as configurações do ambiente
func LoadConfig() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "maiscrianca_secret_key" // Valor padrão, deve ser substituído em produção
	}

	return &Config{
		JWTSecret:          jwtSecret,
		JWTExpirationHours: 24, // Token válido por 24 horas
		DatabaseURL:        os.Getenv("DATABASE_URL"),
		AllowedOrigins:     os.Getenv("ALLOWED_ORIGINS"),
		Port:               port,
	}
}
