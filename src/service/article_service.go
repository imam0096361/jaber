package service

import (
	"fmt"
	"log"
	"news-portal/src/database"
	"news-portal/src/model"
)

type ArticleService struct{}

func (s *ArticleService) GetAllArticles(category, search string) ([]model.Article, error) {
	query := "SELECT id, title, content, category, author, image, created, featured FROM articles WHERE 1=1"
	var args []interface{}
	placeholder := 1

	if category != "" && category != "সব" {
		query += fmt.Sprintf(" AND category = $%d", placeholder)
		args = append(args, category)
		placeholder++
	}

	if search != "" {
		query += fmt.Sprintf(" AND (title LIKE $%d OR content LIKE $%d)", placeholder, placeholder+1)
		searchTerm := "%" + search + "%"
		args = append(args, searchTerm, searchTerm)
		placeholder += 2
	}

	query += " ORDER BY created DESC LIMIT 50"

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []model.Article
	for rows.Next() {
		var a model.Article
		err := rows.Scan(&a.ID, &a.Title, &a.Content, &a.Category, &a.Author, &a.Image, &a.Created, &a.Featured)
		if err != nil {
			log.Println("Error scanning article:", err)
			continue
		}
		articles = append(articles, a)
	}

	return articles, nil
}

func (s *ArticleService) GetFeaturedArticles() ([]model.Article, error) {
	rows, err := database.DB.Query("SELECT id, title, content, category, author, image, created, featured FROM articles WHERE featured = true ORDER BY created DESC LIMIT 5")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []model.Article
	for rows.Next() {
		var a model.Article
		err := rows.Scan(&a.ID, &a.Title, &a.Content, &a.Category, &a.Author, &a.Image, &a.Created, &a.Featured)
		if err != nil {
			log.Println("Error scanning article:", err)
			continue
		}
		articles = append(articles, a)
	}

	return articles, nil
}

func (s *ArticleService) GetArticleByID(id string) (*model.Article, error) {
	var a model.Article
	err := database.DB.QueryRow("SELECT id, title, content, category, author, image, created, featured FROM articles WHERE id = $1", id).Scan(
		&a.ID, &a.Title, &a.Content, &a.Category, &a.Author, &a.Image, &a.Created, &a.Featured,
	)
	if err != nil {
		return nil, err
	}

	return &a, nil
}

func (s *ArticleService) CreateArticle(article *model.Article) (*model.Article, error) {
	log.Println("DEBUG: Starting CreateArticle service")
	query := `INSERT INTO articles (title, content, category, author, image, featured) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	log.Printf("DEBUG: Executing query: %s", query)
	
	err := database.DB.QueryRow(query, article.Title, article.Content, article.Category, article.Author, article.Image, article.Featured).Scan(&article.ID)
	if err != nil {
		log.Printf("DEBUG: Database error: %v", err)
		return nil, err
	}

	log.Printf("DEBUG: Article created with ID: %d", article.ID)
	return article, nil
}

func (s *ArticleService) GetCategories() ([]string, error) {
	rows, err := database.DB.Query("SELECT DISTINCT category FROM articles")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []string
	categories = append(categories, "সব")
	for rows.Next() {
		var cat string
		err := rows.Scan(&cat)
		if err != nil {
			log.Println("Error scanning category:", err)
			continue
		}
		categories = append(categories, cat)
	}

	return categories, nil
}
