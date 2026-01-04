package router

import (
	"net/http"

	"news-portal/src/controller"
)

// setupAuthRoutes setup all authentication related routes
func setupAuthRoutes(mux *http.ServeMux, authController *controller.AuthController) {
	mux.HandleFunc("/api/auth/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			authController.Register(w, r)
		}
	})

	mux.HandleFunc("/api/auth/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			authController.Login(w, r)
		}
	})
}
