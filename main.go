package main

import (
	"electronicsStore/database"
	"electronicsStore/handlers"
	"electronicsStore/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	database.Connect()
	database.RunMigrations()

	r := gin.Default()

	auth := r.Group("/auth")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
	}

	api := r.Group("/")
	api.Use(middleware.AuthRequired())
	{
		// Products
		api.GET("/products", handlers.GetProducts)
		api.GET("/products/:id", handlers.GetProductByID)
		api.POST("/products", handlers.CreateProduct)
		api.PUT("/products/:id", handlers.UpdateProduct)
		api.DELETE("/products/:id", handlers.DeleteProduct)

		// Reviews (nested under products)
		api.GET("/products/:id/reviews", handlers.GetReviews)
		api.POST("/products/:id/reviews", handlers.CreateReview)
		api.DELETE("/reviews/:id", handlers.DeleteReview)

		// Brands
		api.GET("/brands", handlers.GetBrands)
		api.POST("/brands", handlers.CreateBrand)
		api.DELETE("/brands/:id", handlers.DeleteBrand)

		// Categories
		api.GET("/categories", handlers.GetCategories)
		api.POST("/categories", handlers.CreateCategory)
		api.DELETE("/categories/:id", handlers.DeleteCategory)
	}

	r.Run(":8080")
}