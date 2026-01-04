package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"news-portal/src/model"
	"news-portal/src/service"
)

type AuthController struct {
	userService *service.UserService
}

func NewAuthController(userService *service.UserService) *AuthController {
	return &AuthController{userService: userService}
}

func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req model.RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validation
	if req.Name == "" || req.Email == "" || req.Password == "" {
		http.Error(w, "Name, Email, and Password are required", http.StatusBadRequest)
		return
	}

	user, token, err := c.userService.Register(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := model.AuthResponse{
		Token:   token,
		User:    user,
		Message: "Registration successful",
	}

	json.NewEncoder(w).Encode(response)
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req model.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validation
	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and Password are required", http.StatusBadRequest)
		return
	}

	user, token, err := c.userService.Login(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	response := model.AuthResponse{
		Token:   token,
		User:    user,
		Message: "Login successful",
	}

	json.NewEncoder(w).Encode(response)
}

func (c *AuthController) GetProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := c.userService.GetUserByID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}
