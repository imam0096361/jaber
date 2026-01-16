package service

import (
	"app/src/model"
	"time"

	"gorm.io/gorm"
)

type ArticleService interface {
	CreateArticle(a *model.Article) (*model.Article, error)
	GetAllArticles(category, search string, limit int) ([]model.Article, error)
	GetByID(id int) (*model.Article, error)
	DeleteByID(id string) error
	GetFeatured(category string, limit int) ([]model.Article, error)
	GetCategories() ([]string, error)
	UpdateArticle(a *model.Article) (*model.Article, error)
}

type articleService struct {
	db *gorm.DB
}

func NewArticleService(db *gorm.DB) ArticleService {
	return &articleService{db: db}
}

func (s *articleService) CreateArticle(a *model.Article) (*model.Article, error) {
	a.UpdatedAt = time.Now()
	if err := s.db.Create(a).Error; err != nil {
		return nil, err
	}
	return a, nil
}

func (s *articleService) UpdateArticle(a *model.Article) (*model.Article, error) {
	a.UpdatedAt = time.Now()
	// Using Save or Select("*") is necessary to update boolean fields to false (zero values)
	if err := s.db.Model(&model.Article{}).Where("id = ?", a.ID).Select("*").Updates(a).Error; err != nil {
		return nil, err
	}
	return s.GetByID(a.ID)
}

func (s *articleService) GetAllArticles(category, search string, limit int) ([]model.Article, error) {
	var items []model.Article
	query := s.db.Order("updated_at desc")

	// Debug logging
	println("=== GetAllArticles Debug ===")
	println("Category received:", category)
	println("Category length:", len(category))
	println("Category bytes:", []byte(category))

	if category != "" && category != "সব" {
		println("Applying category filter:", category)
		query = query.Where("category = ?", category)
	} else {
		println("No category filter applied (showing all)")
	}

	if search != "" {
		// Use ILIKE for case-insensitive search in Postgres
		query = query.Where("title ILIKE ? OR content ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&items).Error; err != nil {
		return nil, err
	}

	println("Articles found:", len(items))
	for i, item := range items {
		println("Article", i, "- Category:", item.Category)
	}
	println("=========================")

	return items, nil
}

func (s *articleService) GetByID(id int) (*model.Article, error) {
	var a model.Article
	if err := s.db.First(&a, id).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

func (s *articleService) DeleteByID(id string) error {
	// GORM Delete needs the primary key to be specified if not using a pointer with ID
	return s.db.Where("id = ?", id).Delete(&model.Article{}).Error
}

func (s *articleService) GetFeatured(category string, limit int) ([]model.Article, error) {
	var items []model.Article
	query := s.db.Where("featured = ?", true).Order("updated_at desc")

	// Debug logging
	println("=== GetFeatured Debug ===")
	println("Category received:", category)

	if category != "" && category != "সব" {
		println("Applying category filter:", category)
		query = query.Where("category = ?", category)
	} else {
		println("No category filter applied (showing all featured)")
	}

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&items).Error; err != nil {
		return nil, err
	}

	println("Featured articles found:", len(items))
	println("======================")

	return items, nil
}

func (s *articleService) GetCategories() ([]string, error) {
	// Syncing with the Admin Panel list as requested
	categories := []string{
		"জাতীয়",
		"রাজনীতি",
		"খেলাধুলা",
		"প্রযুক্তি",
		"স্বাস্থ্য",
		"বিনোদন",
		"আন্তর্জাতিক",
		"অন্যান্য",
	}
	return categories, nil
}
