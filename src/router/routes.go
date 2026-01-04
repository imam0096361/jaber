package router

import (
	"net/http"
	"os"
	"path/filepath"

	"news-portal/src/controller"
	"news-portal/src/service"
)

// SetupRoutes initializes all application routes
func SetupRoutes(mux *http.ServeMux) {
	// Initialize services and controllers
	articleService := &service.ArticleService{}
	articleController := controller.NewArticleController(articleService)

	userService := &service.UserService{}
	authController := controller.NewAuthController(userService)

	// Setup all routes
	setupStaticFiles(mux)
	setupAuthRoutes(mux, authController)
	setupArticleRoutes(mux, articleController)
	setupProtectedRoutes(mux, authController)
	setupHomeRoute(mux)
}

// setupStaticFiles serves static files from /frontend directory
func setupStaticFiles(mux *http.ServeMux) {
	staticDir := "frontend"
	if wd, err := os.Getwd(); err == nil {
		staticDir = filepath.Join(wd, "frontend")
	}
	mux.Handle("/frontend/", http.StripPrefix("/frontend/", http.FileServer(http.Dir(staticDir))))
}

// setupProtectedRoutes setup protected routes (JWT middleware removed)
func setupProtectedRoutes(mux *http.ServeMux, authController *controller.AuthController) {
	mux.HandleFunc("/api/auth/profile", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			authController.GetProfile(w, r)
		}
	})
}

// setupHomeRoute serves the home page
func setupHomeRoute(mux *http.ServeMux) {
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "frontend/index.html")
		} else {
			http.NotFound(w, r)
		}
	})
}
