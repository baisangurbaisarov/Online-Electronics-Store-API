package router

import (
	"electronicsStore/handlers"
	"electronicsStore/middleware"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORS())

	auth := r.Group("/auth")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
	}

	api := r.Group("/")
	api.Use(middleware.AuthRequired())
	{
		api.GET("/products", handlers.GetProducts)
		api.GET("/products/:id", handlers.GetProductByID)
		api.POST("/products", handlers.CreateProduct)
		api.PUT("/products/:id", handlers.UpdateProduct)
		api.DELETE("/products/:id", handlers.DeleteProduct)

		api.GET("/products/:id/reviews", handlers.GetReviews)
		api.POST("/products/:id/reviews", handlers.CreateReview)
		api.DELETE("/reviews/:id", handlers.DeleteReview)

		api.GET("/brands", handlers.GetBrands)
		api.POST("/brands", handlers.CreateBrand)
		api.DELETE("/brands/:id", handlers.DeleteBrand)

		api.GET("/categories", handlers.GetCategories)
		api.POST("/categories", handlers.CreateCategory)
		api.DELETE("/categories/:id", handlers.DeleteCategory)

		api.POST("/orders", handlers.CreateOrder)
		api.GET("/orders", handlers.GetMyOrders)
		api.GET("/orders/:id", handlers.GetOrderByID)
	}

	return r
}
