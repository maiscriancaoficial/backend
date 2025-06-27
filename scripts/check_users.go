package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Carregar variáveis de ambiente
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Println("Arquivo .env não encontrado, usando variáveis de ambiente")
	}

	// Obter a string de conexão do banco de dados
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("Variável de ambiente DATABASE_URL não encontrada")
	}

	// Conectar ao banco de dados
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	defer db.Close()

	// Verificar conexão com o banco
	err = db.Ping()
	if err != nil {
		log.Fatalf("Erro ao testar conexão com o banco de dados: %v", err)
	}
	fmt.Println("Conexão com o PostgreSQL estabelecida com sucesso!")

	// Verificar se existem usuários no banco
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM \"User\"").Scan(&count)
	if err != nil {
		log.Fatalf("Erro ao contar usuários: %v", err)
	}
	fmt.Printf("Total de usuários no banco: %d\n", count)

	// Listar usuários existentes
	rows, err := db.Query("SELECT id, name, email, role FROM \"User\" LIMIT 10")
	if err != nil {
		log.Fatalf("Erro ao consultar usuários: %v", err)
	}
	defer rows.Close()

	fmt.Println("\nLista de usuários:")
	fmt.Println("----------------------------------")
	fmt.Printf("%-36s | %-20s | %-30s | %-10s\n", "ID", "Nome", "Email", "Papel")
	fmt.Println("----------------------------------")

	for rows.Next() {
		var id, name, email, role string
		if err := rows.Scan(&id, &name, &email, &role); err != nil {
			log.Fatalf("Erro ao ler usuário: %v", err)
		}
		fmt.Printf("%-36s | %-20s | %-30s | %-10s\n", id, name, email, role)
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Erro durante iteração dos usuários: %v", err)
	}
}
