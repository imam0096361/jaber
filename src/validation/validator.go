package validation

import (
	"fmt"
	"news-portal/src/model"
	"regexp"
	"strings"
)

// Article Validation
func ValidateArticle(article *model.Article) []string {
	var errors []string

	if article.Title == "" {
		errors = append(errors, "শিরোনাম প্রয়োজন")
	}
	if len(article.Title) < 3 {
		errors = append(errors, "শিরোনাম কমপক্ষে ৩ অক্ষরের হতে হবে")
	}
	// Word count check
	wordCount := len(strings.Fields(article.Title))
	if wordCount > 500 {
		errors = append(errors, "শিরোনাম ৫০০ শব্দের বেশি হতে পারবে না")
	}

	if article.Content == "" {
		errors = append(errors, "খবর এর বিস্তারিত প্রয়োজন")
	}
	if len(article.Content) < 10 {
		errors = append(errors, "খবর কমপক্ষে ১০ অক্ষরের হতে হবে")
	}

	if article.Category == "" {
		errors = append(errors, "ক্যাটাগরি প্রয়োজন")
	}

	if article.Author == "" {
		errors = append(errors, "লেখক এর নাম প্রয়োজন")
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
