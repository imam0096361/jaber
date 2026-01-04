package validation

import (
	"fmt"
	"news-portal/src/model"
	"regexp"
)

// Article Validation
func ValidateArticle(article *model.Article) []string {
	var errors []string

	if article.Title == "" {
		errors = append(errors, "Title is required")
	}
	if len(article.Title) < 3 {
		errors = append(errors, "Title must be at least 3 characters")
	}
	if len(article.Title) > 200 {
		errors = append(errors, "Title must not exceed 200 characters")
	}

	if article.Content == "" {
		errors = append(errors, "Content is required")
	}
	if len(article.Content) < 10 {
		errors = append(errors, "Content must be at least 10 characters")
	}

	if article.Category == "" {
		errors = append(errors, "Category is required")
	}

	if article.Author == "" {
		errors = append(errors, "Author is required")
	}

	return errors
}

// User Validation
func ValidateRegister(req *model.RegisterRequest) []string {
	var errors []string

	if req.Name == "" {
		errors = append(errors, "Name is required")
	}
	if len(req.Name) < 2 {
		errors = append(errors, "Name must be at least 2 characters")
	}

	if err := ValidateEmail(req.Email); err != nil {
		errors = append(errors, err.Error())
	}

	if req.Password == "" {
		errors = append(errors, "Password is required")
	}
	if len(req.Password) < 6 {
		errors = append(errors, "Password must be at least 6 characters")
	}

	return errors
}

func ValidateLogin(req *model.LoginRequest) []string {
	var errors []string

	if err := ValidateEmail(req.Email); err != nil {
		errors = append(errors, err.Error())
	}

	if req.Password == "" {
		errors = append(errors, "Password is required")
	}

	return errors
}

// Email Validation
func ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email is required")
	}

	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	if !matched {
		return fmt.Errorf("invalid email format")
	}

	return nil
}
