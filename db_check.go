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
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("Arquivo .env não encontrado, usando variáveis de ambiente")
	}

	// Obter string de conexão do DB
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("Variável de ambiente DATABASE_URL não definida")
	}

	// Conectar ao banco de dados
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Erro ao conectar ao PostgreSQL: %v", err)
	}
	defer db.Close()

	// Verificar conexão
	err = db.Ping()
	if err != nil {
		log.Fatalf("Erro ao verificar conexão com PostgreSQL: %v", err)
	}
	fmt.Println("Conexão com o PostgreSQL estabelecida com sucesso!")

	// Listar todas as tabelas no banco de dados
	rows, err := db.Query(`
		SELECT tablename FROM pg_catalog.pg_tables 
		WHERE schemaname != 'pg_catalog' AND schemaname != 'information_schema'
	`)
	if err != nil {
		log.Fatalf("Erro ao listar tabelas: %v", err)
	}
	defer rows.Close()

	fmt.Println("\nTabelas no banco de dados:")
	var tabelas []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			log.Fatalf("Erro ao ler nome da tabela: %v", err)
		}
		tabelas = append(tabelas, tableName)
		fmt.Println("-", tableName)
	}

	// Se encontrou a tabela User, listar suas colunas
	for _, tabela := range tabelas {
		if tabela == "User" || tabela == "user" {
			fmt.Printf("\nColunas da tabela %s:\n", tabela)
			colRows, err := db.Query(fmt.Sprintf(`
				SELECT column_name, data_type 
				FROM information_schema.columns 
				WHERE table_name = '%s'
			`, tabela))
			if err != nil {
				log.Fatalf("Erro ao listar colunas: %v", err)
			}
			defer colRows.Close()

			for colRows.Next() {
				var colName, dataType string
				if err := colRows.Scan(&colName, &dataType); err != nil {
					log.Fatalf("Erro ao ler coluna: %v", err)
				}
				fmt.Printf("- %s (%s)\n", colName, dataType)
			}
		}
	}
}
