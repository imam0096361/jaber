package router

import (
	"net/http"

	"news-portal/src/controller"
)

// setupArticleRoutes setup all article related routes
func setupArticleRoutes(mux *http.ServeMux, articleController *controller.ArticleController) {
	// Get endpoints
	mux.HandleFunc("/api/articles", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			articleController.GetArticles(w, r)
		}
	})

	mux.HandleFunc("/api/featured", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			articleController.GetFeatured(w, r)
		}
	})

	mux.HandleFunc("/api/article", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			articleController.GetArticleByID(w, r)
		}
	})

	mux.HandleFunc("/api/categories", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			articleController.GetCategories(w, r)
		}
	})

	// Create/Update endpoints
	mux.HandleFunc("/api/add-article", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			articleController.AddArticle(w, r)
		}
	})

	mux.HandleFunc("/api/upload-image", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			articleController.UploadImage(w, r)
		}
	})
}
