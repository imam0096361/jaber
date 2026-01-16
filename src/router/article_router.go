package router

import (
	"app/src/controller"
	"app/src/service"

	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
)

func ArticleRoutes(router fiber.Router, s service.ArticleService) {
	articleController := controller.NewArticleController(s)

	// Cache Config - Include full URI (with query params) in cache key
	cacheConfig := cache.Config{
		Expiration:   30 * time.Second,
		CacheControl: true,
		KeyGenerator: func(c *fiber.Ctx) string {
			// Include the full path AND query string in the cache key
			// This ensures different query parameters get different cached responses
			return c.OriginalURL()
		},
	}

	router.Get("/articles", cache.New(cacheConfig), articleController.GetAll)
	router.Get("/article", articleController.GetByID) // Don't cache detailed view to ensure hit count/freshness if needed
	router.Get("/featured", cache.New(cacheConfig), articleController.GetFeatured)
	router.Get("/categories", cache.New(cacheConfig), articleController.GetCategories)
	
	router.Post("/add-article", articleController.CreateArticle)
	router.Post("/upload-image", articleController.UploadImage)
	router.Post("/delete-article", articleController.DeleteByID)
	router.Post("/update-article", articleController.UpdateArticle)
}
