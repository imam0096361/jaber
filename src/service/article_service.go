package service

import (
	"log"
	"news-portal/src/database"
	"news-portal/src/model"
)

type ArticleService struct{}

func (s *ArticleService) GetAllArticles(category, search string) ([]model.Article, error) {
	query := "SELECT id, title, content, category, author, image, created, featured FROM articles WHERE 1=1"
	var args []interface{}

	if category != "" && category != "সব" {
		query += " AND category = ?"
		args = append(args, category)
	}

	if search != "" {
		query += " AND (title LIKE ? OR content LIKE ?)"
		searchTerm := "%" + search + "%"
		args = append(args, searchTerm, searchTerm)
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
	rows, err := database.DB.Query("SELECT id, title, content, category, author, image, created, featured FROM articles WHERE featured = 1 ORDER BY created DESC LIMIT 5")
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
	err := database.DB.QueryRow("SELECT id, title, content, category, author, image, created, featured FROM articles WHERE id = ?", id).Scan(
		&a.ID, &a.Title, &a.Content, &a.Category, &a.Author, &a.Image, &a.Created, &a.Featured,
	)
	if err != nil {
		return nil, err
	}

	return &a, nil
}

func (s *ArticleService) CreateArticle(article *model.Article) (*model.Article, error) {
	query := `INSERT INTO articles (title, content, category, author, image, featured) VALUES (?, ?, ?, ?, ?, ?)`
	result, err := database.DB.Exec(query, article.Title, article.Content, article.Category, article.Author, article.Image, article.Featured)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	article.ID = int(id)
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
