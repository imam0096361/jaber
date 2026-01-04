package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"news-portal/src/config"
	"news-portal/src/model"
	"news-portal/src/service"
	"news-portal/src/validation"
)

type ArticleController struct {
	service *service.ArticleService
}

func NewArticleController(service *service.ArticleService) *ArticleController {
	return &ArticleController{service: service}
}

func (c *ArticleController) GetArticles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	category := r.URL.Query().Get("category")
	search := r.URL.Query().Get("search")

	articles, err := c.service.GetAllArticles(category, search)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if articles == nil {
		articles = []model.Article{}
	}

	json.NewEncoder(w).Encode(articles)
}

func (c *ArticleController) GetFeatured(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	articles, err := c.service.GetFeaturedArticles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if articles == nil {
		articles = []model.Article{}
	}

	json.NewEncoder(w).Encode(articles)
}

func (c *ArticleController) GetArticleByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID required", http.StatusBadRequest)
		return
	}

	article, err := c.service.GetArticleByID(id)
	if err != nil {
		http.Error(w, "Article not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(article)
}

func (c *ArticleController) GetCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	categories, err := c.service.GetCategories()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(categories)
}

func (c *ArticleController) AddArticle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var a model.Article
	err := json.NewDecoder(r.Body).Decode(&a)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate article
	if validationErrors := validation.ValidateArticle(&a); len(validationErrors) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Validation failed",
			"errors":  validationErrors,
		})
		return
	}

	article, err := c.service.CreateArticle(&a)
	if err != nil {
		http.Error(w, "Error inserting article: "+err.Error(), http.StatusInternalServerError)
		return
	}

	article.Created = time.Now()

	response := model.Response{
		Message: "Article added successfully",
		Article: article,
	}

	json.NewEncoder(w).Encode(response)
}

func (c *ArticleController) UploadImage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Create uploads directory if it doesn't exist
	uploadsDir := config.UPLOADS_DIR
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		http.Error(w, "Error creating uploads directory", http.StatusInternalServerError)
		return
	}

	// Parse file upload
	err := r.ParseMultipartForm(config.MAX_FILE_SIZE)
	if err != nil {
		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create unique filename
	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), handler.Filename)
	filePath := filepath.Join(uploadsDir, filename)

	// Save file
	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy file content
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Error writing file", http.StatusInternalServerError)
		return
	}

	imageURL := "/static/uploads/" + filename

	response := model.ImageUploadResponse{
		Message: "Image uploaded successfully",
		URL:     imageURL,
	}

	json.NewEncoder(w).Encode(response)
}
