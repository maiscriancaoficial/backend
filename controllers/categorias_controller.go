package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/maiscrianca/backend/models"
	"github.com/maiscrianca/backend/repository"
)

// GetCategorias retorna todas as categorias
func GetCategorias(c *fiber.Ctx) error {
	// Obter o ID do espaço do usuário do middleware de autenticação
	espacoId := c.Locals("espacoId").(string)

	// Buscar categorias do repositório
	categorias, err := repository.GetCategoriasByEspacoId(espacoId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro ao buscar categorias: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    categorias,
	})
}

// CreateCategoria cria uma nova categoria
func CreateCategoria(c *fiber.Ctx) error {
	// Obter o ID do espaço do usuário do middleware de autenticação
	espacoId := c.Locals("espacoId").(string)

	// Criar uma instância de categoria para receber os dados JSON
	categoria := new(models.Categoria)

	// Parsear o body da requisição
	if err := c.BodyParser(categoria); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Erro ao processar dados: " + err.Error(),
		})
	}

	// Adicionar o espacoId à categoria
	categoria.EspacoId = espacoId

	// Inserir a categoria no banco de dados
	if err := repository.CreateCategoria(categoria); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro ao criar categoria: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Categoria criada com sucesso",
		"data":    categoria,
	})
}

// UpdateCategoria atualiza uma categoria existente
func UpdateCategoria(c *fiber.Ctx) error {
	id := c.Params("id")
	espacoId := c.Locals("espacoId").(string)

	// Verificar se a categoria existe e pertence ao espaço
	_, err := repository.GetCategoriaById(id, espacoId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Categoria não encontrada",
		})
	}

	// Criar uma instância de categoria para receber os dados JSON
	categoriaUpdate := new(models.Categoria)

	// Parsear o body da requisição
	if err := c.BodyParser(categoriaUpdate); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Erro ao processar dados: " + err.Error(),
		})
	}

	// Garantir que o ID e o espacoId corretos sejam usados
	categoriaUpdate.ID = id
	categoriaUpdate.EspacoId = espacoId

	// Atualizar a categoria no banco de dados
	if err := repository.UpdateCategoria(categoriaUpdate); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro ao atualizar categoria: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Categoria atualizada com sucesso",
		"data":    categoriaUpdate,
	})
}

// DeleteCategoria exclui uma categoria
func DeleteCategoria(c *fiber.Ctx) error {
	id := c.Params("id")
	espacoId := c.Locals("espacoId").(string)

	// Verificar se a categoria existe e pertence ao espaço
	_, err := repository.GetCategoriaById(id, espacoId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Categoria não encontrada",
		})
	}

	// Verificar se existem livros associados à categoria
	hasBooks, err := repository.HasLivrosByCategoria(id, espacoId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro ao verificar livros associados: " + err.Error(),
		})
	}

	if hasBooks {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Não é possível excluir esta categoria pois existem livros associados a ela",
		})
	}

	// Excluir a categoria
	if err := repository.DeleteCategoria(id, espacoId); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro ao excluir categoria: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Categoria excluída com sucesso",
	})
}
