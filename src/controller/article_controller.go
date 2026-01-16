package controller

import (
	"app/src/model"
	"app/src/response"
	"app/src/service"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type ArticleController struct {
	_ArticleService service.ArticleService
}

func NewArticleController(s service.ArticleService) *ArticleController {
	return &ArticleController{_ArticleService: s}
}

func (a *ArticleController) CreateArticle(c *fiber.Ctx) error {
	req := new(model.Article)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "fail",
			Message: "invalid json",
			Errors:  err.Error(),
		})
	}

	req.Created = time.Now()

	created, err := a._ArticleService.CreateArticle(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorDetails{
			Code:    fiber.StatusInternalServerError,
			Status:  "error",
			Message: "could not create article",
			Errors:  err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(created)
}

func (a *ArticleController) GetAll(c *fiber.Ctx) error {
	category := c.Query("category")
	search := c.Query("search")
	limit, _ := strconv.Atoi(c.Query("limit"))

	// Debug logging
	println("=== Controller GetAll ===")
	println("Raw category from query:", category)
	println("Category bytes:", []byte(category))
	println("=========================")

	items, err := a._ArticleService.GetAllArticles(category, search, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorDetails{
			Code:    fiber.StatusInternalServerError,
			Status:  "error",
			Message: "could not fetch articles",
			Errors:  err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(items)
}

func (a *ArticleController) GetByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	if idStr == "" {
		idStr = c.Query("id")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "fail",
			Message: "invalid id",
			Errors:  err.Error(),
		})
	}

	item, err := a._ArticleService.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.ErrorDetails{
			Code:    fiber.StatusNotFound,
			Status:  "fail",
			Message: "article not found",
			Errors:  err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(item)
}

func (a *ArticleController) DeleteByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	if idStr == "" {
		idStr = c.Query("id")
	}

	if err := a._ArticleService.DeleteByID(idStr); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorDetails{
			Code:    fiber.StatusInternalServerError,
			Status:  "error",
			Message: "could not delete article",
			Errors:  err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.Common{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "deleted",
	})
}

func (a *ArticleController) UpdateArticle(c *fiber.Ctx) error {
	req := new(model.Article)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "fail",
			Message: "invalid json",
			Errors:  err.Error(),
		})
	}

	updated, err := a._ArticleService.UpdateArticle(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorDetails{
			Code:    fiber.StatusInternalServerError,
			Status:  "error",
			Message: "could not update article",
			Errors:  err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(updated)
}

func (a *ArticleController) GetFeatured(c *fiber.Ctx) error {
	category := c.Query("category")
	limit, _ := strconv.Atoi(c.Query("limit"))

	items, err := a._ArticleService.GetFeatured(category, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorDetails{
			Code:    fiber.StatusInternalServerError,
			Status:  "error",
			Message: "could not fetch featured articles",
			Errors:  err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(items)
}

func (a *ArticleController) GetCategories(c *fiber.Ctx) error {
	items, err := a._ArticleService.GetCategories()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorDetails{
			Code:    fiber.StatusInternalServerError,
			Status:  "error",
			Message: "could not fetch categories",
			Errors:  err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(items)
}

func (a *ArticleController) UploadImage(c *fiber.Ctx) error {
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "fail",
			Message: "image is required",
			Errors:  err.Error(),
		})
	}

	uploadsDir := "frontend/uploads"
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorDetails{
			Code:    fiber.StatusInternalServerError,
			Status:  "error",
			Message: "could not create uploads directory",
			Errors:  err.Error(),
		})
	}

	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
	filePath := filepath.Join(uploadsDir, filename)

	if err := c.SaveFile(file, filePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorDetails{
			Code:    fiber.StatusInternalServerError,
			Status:  "error",
			Message: "could not save file",
			Errors:  err.Error(),
		})
	}

	url := "/frontend/uploads/" + filename
	return c.Status(fiber.StatusOK).JSON(model.ImageUploadResponse{
		Message: "ছবি সফলভাবে আপলোড হয়েছে",
		URL:     url,
	})
}
