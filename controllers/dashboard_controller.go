package controllers

import (
	"database/sql"
	"time"

	"github.com/gofiber/fiber/v2"
)

// DashboardMetrics representa os dados de métricas do dashboard
type DashboardMetrics struct {
	TotalVendas     float64       `json:"totalVendas"`
	LivrosVendidos  int           `json:"livrosVendidos"`
	ProdutosVendidos int          `json:"produtosVendidos"`
	Visualizacoes   int           `json:"visualizacoes"`
	Cliques         int           `json:"cliques"`
	GraficoMetricas []MetricaData `json:"graficoMetricas"`
}

// MetricaData representa os dados do gráfico de métricas por período
type MetricaData struct {
	Data      string  `json:"data"`
	Vendas    float64 `json:"vendas"`
	Livros    int     `json:"livros"`
	Produtos  int     `json:"produtos"`
	Views     int     `json:"views"`
	Cliques   int     `json:"cliques"`
}

// GetDashboardMetrics retorna as métricas para o dashboard
func GetDashboardMetrics(c *fiber.Ctx) error {
	db := c.Locals("db").(*sql.DB)

	// Busca total de vendas
	var totalVendas sql.NullFloat64
	err := db.QueryRow("SELECT COALESCE(SUM(valor_total), 0) FROM pedidos WHERE status != 'cancelado'").Scan(&totalVendas)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro ao buscar total de vendas: " + err.Error(),
		})
	}

	// Busca total de livros vendidos
	var livrosVendidos sql.NullInt64
	err = db.QueryRow(`
		SELECT COALESCE(SUM(ip.quantidade), 0) 
		FROM itens_pedido ip 
		JOIN produtos p ON ip.produto_id = p.id
		JOIN pedidos pe ON ip.pedido_id = pe.id
		WHERE p.categoria = 'livro' AND pe.status != 'cancelado'
	`).Scan(&livrosVendidos)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro ao buscar livros vendidos: " + err.Error(),
		})
	}

	// Busca total de produtos vendidos
	var produtosVendidos sql.NullInt64
	err = db.QueryRow(`
		SELECT COALESCE(SUM(quantidade), 0) 
		FROM itens_pedido ip
		JOIN pedidos p ON ip.pedido_id = p.id
		WHERE p.status != 'cancelado'
	`).Scan(&produtosVendidos)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro ao buscar produtos vendidos: " + err.Error(),
		})
	}

	// Busca visualizações
	var visualizacoes sql.NullInt64
	err = db.QueryRow("SELECT COALESCE(SUM(visualizacoes), 0) FROM metricas_site").Scan(&visualizacoes)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro ao buscar visualizações: " + err.Error(),
		})
	}

	// Busca cliques
	var cliques sql.NullInt64
	err = db.QueryRow("SELECT COALESCE(SUM(cliques), 0) FROM metricas_site").Scan(&cliques)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro ao buscar cliques: " + err.Error(),
		})
	}

	// Busca dados para o gráfico dos últimos 7 dias
	graficoMetricas := []MetricaData{}

	rows, err := db.Query(`
		WITH dias AS (
			SELECT date(generate_series(current_date - interval '6 days', current_date, '1 day')) AS data
		)
		SELECT 
			to_char(d.data, 'DD/MM/YYYY') as data,
			COALESCE(SUM(p.valor_total), 0) as vendas,
			COALESCE(SUM(CASE WHEN pr.categoria = 'livro' THEN ip.quantidade ELSE 0 END), 0) as livros,
			COALESCE(SUM(ip.quantidade), 0) as produtos,
			COALESCE(SUM(m.visualizacoes), 0) as views,
			COALESCE(SUM(m.cliques), 0) as cliques
		FROM 
			dias d
		LEFT JOIN 
			pedidos p ON date(p.data_criacao) = d.data AND p.status != 'cancelado'
		LEFT JOIN 
			itens_pedido ip ON p.id = ip.pedido_id
		LEFT JOIN 
			produtos pr ON ip.produto_id = pr.id
		LEFT JOIN 
			metricas_site m ON date(m.data) = d.data
		GROUP BY 
			d.data
		ORDER BY 
			d.data
	`)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro ao buscar dados para o gráfico: " + err.Error(),
		})
	}
	defer rows.Close()

	for rows.Next() {
		var metrica MetricaData
		err := rows.Scan(&metrica.Data, &metrica.Vendas, &metrica.Livros, &metrica.Produtos, &metrica.Views, &metrica.Cliques)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Erro ao ler dados do gráfico: " + err.Error(),
			})
		}
		graficoMetricas = append(graficoMetricas, metrica)
	}

	// Tratamento de valores nulos
	totalVendasValue := 0.0
	if totalVendas.Valid {
		totalVendasValue = totalVendas.Float64
	}

	livrosVendidosValue := 0
	if livrosVendidos.Valid {
		livrosVendidosValue = int(livrosVendidos.Int64)
	}

	produtosVendidosValue := 0
	if produtosVendidos.Valid {
		produtosVendidosValue = int(produtosVendidos.Int64)
	}

	visualizacoesValue := 0
	if visualizacoes.Valid {
		visualizacoesValue = int(visualizacoes.Int64)
	}

	cliquesValue := 0
	if cliques.Valid {
		cliquesValue = int(cliques.Int64)
	}

	// Retorna os dados
	metrics := DashboardMetrics{
		TotalVendas:      totalVendasValue,
		LivrosVendidos:   livrosVendidosValue,
		ProdutosVendidos: produtosVendidosValue,
		Visualizacoes:    visualizacoesValue,
		Cliques:          cliquesValue,
		GraficoMetricas:  graficoMetricas,
	}

	return c.JSON(metrics)
}
