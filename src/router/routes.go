package router

import (
	"net/http"
	"os"
	"path/filepath"

	"news-portal/src/controller"
	"news-portal/src/middleware"
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

// setupStaticFiles serves static files from /static directory
func setupStaticFiles(mux *http.ServeMux) {
	staticDir := "static"
	if wd, err := os.Getwd(); err == nil {
		staticDir = filepath.Join(wd, "static")
	}
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))
}

// setupProtectedRoutes setup protected routes with JWT middleware
func setupProtectedRoutes(mux *http.ServeMux, authController *controller.AuthController) {
	mux.HandleFunc("/api/auth/profile", middleware.JWTAuthMiddleware(
		http.HandlerFunc(authController.GetProfile),
	).ServeHTTP)
}

// setupHomeRoute serves the home page
func setupHomeRoute(mux *http.ServeMux) {
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "static/index.html")
		} else {
			http.NotFound(w, r)
		}
	})
}
