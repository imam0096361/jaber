package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
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
		http.Error(w, "আইডি প্রয়োজন", http.StatusBadRequest)
		return
	}

	article, err := c.service.GetArticleByID(id)
	if err != nil {
		http.Error(w, "খবরটি পাওয়া যায়নি", http.StatusNotFound)
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
		http.Error(w, "মেথড অনুমোদিত নয়", http.StatusMethodNotAllowed)
		return
	}

	var a model.Article
	err := json.NewDecoder(r.Body).Decode(&a)
	if err != nil {
		log.Printf("Error decoding JSON: %v", err)
		http.Error(w, "ভুল JSON ফরম্যাট", http.StatusBadRequest)
		return
	}

	log.Printf("Adding article: %+v", a)

	// Validate article
	if validationErrors := validation.ValidateArticle(&a); len(validationErrors) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "ভ্যালিডেশন ব্যর্থ হয়েছে",
			"errors":  validationErrors,
		})
		return
	}

	article, err := c.service.CreateArticle(&a)
	if err != nil {
		log.Printf("Error inserting article: %v", err)
		http.Error(w, "খবর যোগ করতে সমস্যা হয়েছে: "+err.Error(), http.StatusInternalServerError)
		return
	}

	article.Created = time.Now()

	response := model.Response{
		Message: "খবর সফলভাবে যোগ করা হয়েছে",
		Article: article,
	}

	json.NewEncoder(w).Encode(response)
}

func (c *ArticleController) UploadImage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "মেথড অনুমোদিত নয়", http.StatusMethodNotAllowed)
		return
	}

	// Create uploads directory if it doesn't exist
	uploadsDir := config.UPLOADS_DIR
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		log.Printf("Error creating uploads directory: %v", err)
		http.Error(w, "আপলোড ডিরেক্টরি তৈরি করা যায়নি", http.StatusInternalServerError)
		return
	}

	// Parse file upload
	err := r.ParseMultipartForm(config.MAX_UPLOAD_SIZE)
	if err != nil {
		log.Printf("Error parsing multipart form: %v", err)
		http.Error(w, "ফাইলের আকার অনেক বড়", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("image")
	if err != nil {
		log.Printf("Error retrieving file: %v", err)
		http.Error(w, "ফাইল পাওয়া যায়নি", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create unique filename
	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), handler.Filename)
	filePath := filepath.Join(uploadsDir, filename)

	// Save file
	dst, err := os.Create(filePath)
	if err != nil {
		log.Printf("Error creating file: %v", err)
		http.Error(w, "ফাইল সেভ করা যায়নি", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy file content
	_, err = io.Copy(dst, file)
	if err != nil {
		log.Printf("Error writing file: %v", err)
		http.Error(w, "ফাইল রাইট করা যায়নি", http.StatusInternalServerError)
		return
	}

	imageURL := "/frontend/uploads/" + filename
	log.Printf("Image uploaded successfully: %s", imageURL)

	response := model.ImageUploadResponse{
		Message: "ছবি সফলভাবে আপলোড হয়েছে",
		URL:     imageURL,
	}

	json.NewEncoder(w).Encode(response)
}
