package service

import (
	"fmt"
	"time"

	"news-portal/src/database"
	"news-portal/src/middleware"
	"news-portal/src/model"
	"news-portal/src/utils"
)

type UserService struct{}

func (s *UserService) Register(req *model.RegisterRequest) (*model.User, string, error) {
	// Check if email already exists
	var existingID int
	err := database.DB.QueryRow("SELECT id FROM users WHERE email = ?", req.Email).Scan(&existingID)
	if err == nil {
		return nil, "", fmt.Errorf("email already registered")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, "", err
	}

	// Insert user
	result, err := database.DB.Exec(
		"INSERT INTO users (name, email, password_hash) VALUES (?, ?, ?)",
		req.Name, req.Email, hashedPassword,
	)
	if err != nil {
		return nil, "", err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return nil, "", err
	}

	user := &model.User{
		ID:        int(userID),
		Name:      req.Name,
		Email:     req.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Generate JWT token
	token, err := middleware.GenerateJWT(int(userID), req.Email)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *UserService) Login(req *model.LoginRequest) (*model.User, string, error) {
	var user model.User
	var passwordHash string

	// Find user by email
	err := database.DB.QueryRow(
		"SELECT id, name, email, password_hash, created_at, updated_at FROM users WHERE email = ?",
		req.Email,
	).Scan(&user.ID, &user.Name, &user.Email, &passwordHash, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, "", fmt.Errorf("user not found")
	}

	// Check password
	if !utils.CheckPasswordHash(req.Password, passwordHash) {
		return nil, "", fmt.Errorf("invalid password")
	}

	// Generate JWT token
	token, err := middleware.GenerateJWT(user.ID, user.Email)
	if err != nil {
		return nil, "", err
	}

	return &user, token, nil
}

func (s *UserService) GetUserByID(userID int) (*model.User, error) {
	var user model.User
	err := database.DB.QueryRow(
		"SELECT id, name, email, created_at, updated_at FROM users WHERE id = ?",
		userID,
	).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return &user, nil
}

func (s *UserService) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	err := database.DB.QueryRow(
		"SELECT id, name, email, created_at, updated_at FROM users WHERE email = ?",
		email,
	).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return &user, nil
}
