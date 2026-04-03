package main

import (
	"electronicsStore/database"
	"electronicsStore/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	
	database.Connect()

	r := gin.Default()

	r.GET("/products", handlers.GetProducts)
	r.GET("/products/:id", handlers.GetProductByID)
	r.POST("/products", handlers.CreateProduct)
	r.PUT("/products/:id", handlers.UpdateProduct)
	r.DELETE("/products/:id", handlers.DeleteProduct)

	r.GET("/brands", handlers.GetBrands)
	r.POST("/brands", handlers.CreateBrand)
	r.DELETE("/brands/:id", handlers.DeleteBrand)

	r.GET("/categories", handlers.GetCategories)
	r.POST("/categories", handlers.CreateCategory)
	r.DELETE("/categories/:id", handlers.DeleteCategory)

	r.Run(":8080")
}
