package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("Iniciando depuração do banco de dados...")

	// Carregar variáveis de ambiente
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Println("Erro ao carregar .env:", err)
	} else {
		fmt.Println("Arquivo .env carregado com sucesso")
	}

	// Verificar DATABASE_URL
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		fmt.Println("ERRO: Variável DATABASE_URL não definida ou vazia")
		return
	} else {
		fmt.Println("DATABASE_URL encontrada (não exibindo por segurança)")
	}

	// Tentar conexão
	fmt.Println("Tentando conectar ao PostgreSQL...")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println("ERRO ao abrir conexão:", err)
		return
	}
	defer db.Close()

	// Testar conexão
	fmt.Println("Testando conexão com Ping()...")
	err = db.Ping()
	if err != nil {
		fmt.Println("ERRO no ping ao banco:", err)
		return
	}
	fmt.Println("Conexão com PostgreSQL estabelecida com sucesso!")

	// Listar tabelas
	fmt.Println("Obtendo lista de tabelas...")
	rows, err := db.Query(`
		SELECT tablename 
		FROM pg_catalog.pg_tables 
		WHERE schemaname != 'pg_catalog' 
		AND schemaname != 'information_schema'
	`)
	if err != nil {
		fmt.Println("ERRO ao listar tabelas:", err)
		return
	}
	defer rows.Close()

	fmt.Println("\nTabelas no banco de dados:")
	var temTabela bool
	for rows.Next() {
		temTabela = true
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			fmt.Println("ERRO ao ler nome da tabela:", err)
			continue
		}
		fmt.Println("-", tableName)
	}

	if !temTabela {
		fmt.Println("Nenhuma tabela encontrada!")
	}
}
