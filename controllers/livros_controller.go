package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/maiscrianca/backend/models"
	"github.com/maiscrianca/backend/repository"
	"os"
	"strconv"
	"github.com/vercel/blob"
)

// GetLivros retorna todos os livros
func GetLivros(c *fiber.Ctx) error {
	// Obter o ID do espaço do usuário do middleware de autenticação
	espacoId := c.Locals("espacoId").(string)

	// Buscar livros do repositório
	livros, err := repository.GetLivrosByEspacoId(espacoId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro ao buscar livros: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    livros,
	})
}

// GetLivro retorna um livro específico
func GetLivro(c *fiber.Ctx) error {
	id := c.Params("id")
	espacoId := c.Locals("espacoId").(string)

	livro, err := repository.GetLivroById(id, espacoId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Livro não encontrado",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    livro,
	})
}

// CreateLivro cria um novo livro
func CreateLivro(c *fiber.Ctx) error {
	// Obter o ID do espaço do usuário do middleware de autenticação
	espacoId := c.Locals("espacoId").(string)

	// Criar uma instância de livro para receber os dados JSON
	livro := new(models.Livro)

	// Parsear o body da requisição
	if err := c.BodyParser(livro); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Erro ao processar dados: " + err.Error(),
		})
	}

	// Adicionar o espacoId ao livro
	livro.EspacoId = espacoId

	// Inserir o livro no banco de dados
	if err := repository.CreateLivro(livro); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro ao criar livro: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Livro criado com sucesso",
		"data":    livro,
	})
}

// UpdateLivro atualiza um livro existente
func UpdateLivro(c *fiber.Ctx) error {
	id := c.Params("id")
	espacoId := c.Locals("espacoId").(string)

	// Verificar se o livro existe e pertence ao espaço
	_, err := repository.GetLivroById(id, espacoId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Livro não encontrado",
		})
	}

	// Criar uma instância de livro para receber os dados JSON
	livroUpdate := new(models.Livro)

	// Parsear o body da requisição
	if err := c.BodyParser(livroUpdate); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Erro ao processar dados: " + err.Error(),
		})
	}

	// Garantir que o ID e o espacoId corretos sejam usados
	livroUpdate.ID = id
	livroUpdate.EspacoId = espacoId

	// Atualizar o livro no banco de dados
	if err := repository.UpdateLivro(livroUpdate); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro ao atualizar livro: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Livro atualizado com sucesso",
		"data":    livroUpdate,
	})
}

// DeleteLivro exclui um livro
func DeleteLivro(c *fiber.Ctx) error {
	id := c.Params("id")
	espacoId := c.Locals("espacoId").(string)

	// Verificar se o livro existe e pertence ao espaço
	_, err := repository.GetLivroById(id, espacoId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Livro não encontrado",
		})
	}

	// Excluir o livro
	if err := repository.DeleteLivro(id, espacoId); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro ao excluir livro: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Livro excluído com sucesso",
	})
}

// UploadCapa faz upload da capa do livro usando Vercel Blob
func UploadCapa(c *fiber.Ctx) error {
	return handleBlobUpload(c, "book-covers")
}

// UploadArquivo faz upload do arquivo PDF do livro usando Vercel Blob
func UploadArquivo(c *fiber.Ctx) error {
	return handleBlobUpload(c, "book-files")
}

// UploadPagina faz upload de uma imagem de página do livro usando Vercel Blob
func UploadPagina(c *fiber.Ctx) error {
	return handleBlobUpload(c, "book-pages")
}

// handleBlobUpload gerencia o upload para o Vercel Blob
func handleBlobUpload(c *fiber.Ctx, folder string) error {
	// Obter token do Vercel Blob
	blobToken := os.Getenv("BLOB_READ_WRITE_TOKEN")
	if blobToken == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Token do Vercel Blob não configurado",
		})
	}

	// Obter o arquivo do multipart/form-data
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Arquivo não fornecido: " + err.Error(),
		})
	}

	// Abrir o arquivo
	fileContent, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Não foi possível ler o arquivo: " + err.Error(),
		})
	}
	defer fileContent.Close()

	// Preparar cliente do Vercel Blob
	client := blob.NewClient(blobToken)

	// Upload do arquivo para o Vercel Blob
	uploadResult, err := client.Upload(c.Context(), &blob.UploadOptions{
		Data:        fileContent,
		Filename:    file.Filename,
		ContentType: file.Header["Content-Type"][0],
		Folder:      folder,
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro ao fazer upload do arquivo: " + err.Error(),
		})
	}

	// Retornar a URL do arquivo no Blob
	return c.JSON(fiber.Map{
		"success": true,
		"url":     uploadResult.URL,
	})
}
