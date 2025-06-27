package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("Iniciando diagnóstico do banco de dados...")

	// Carregar variáveis de ambiente
	err := godotenv.Load("../../.env")
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
	fmt.Println("\nListando tabelas no banco de dados:")
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

	var tabelas []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			fmt.Println("ERRO ao ler nome da tabela:", err)
			continue
		}
		tabelas = append(tabelas, tableName)
		fmt.Println("-", tableName)
	}

	if len(tabelas) == 0 {
		fmt.Println("ATENÇÃO: Nenhuma tabela encontrada no banco de dados!")
		return
	}

	// Para cada tabela que pode ser "User", verificar colunas
	fmt.Println("\nProcurando e analisando tabela de usuários...")
	encontrouTabelaUser := false

	for _, tabela := range tabelas {
		if tabela == "User" || tabela == "user" || tabela == "users" {
			encontrouTabelaUser = true
			fmt.Printf("\nColunas da tabela %s:\n", tabela)
			
			colRows, err := db.Query(fmt.Sprintf(`
				SELECT column_name, data_type 
				FROM information_schema.columns 
				WHERE table_name = '%s'
			`, tabela))
			
			if err != nil {
				fmt.Println("ERRO ao listar colunas:", err)
				continue
			}

			for colRows.Next() {
				var colName, dataType string
				if err := colRows.Scan(&colName, &dataType); err != nil {
					fmt.Println("ERRO ao ler coluna:", err)
					continue
				}
				fmt.Printf("- %s (%s)\n", colName, dataType)
			}
			colRows.Close()
		}
	}

	if !encontrouTabelaUser {
		fmt.Println("ATENÇÃO: Não foi encontrada nenhuma tabela de usuários (User/user/users)!")
		fmt.Println("Certifique-se de que o banco de dados contém as tabelas necessárias.")
	}
}
