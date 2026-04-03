package handlers

import (
	"net/http"
	"strconv"

	"electronicsStore/database"
	"electronicsStore/models"

	"github.com/gin-gonic/gin"
)

func GetProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))
	categoryStr := c.Query("category_id")
	brandStr := c.Query("brand_id")

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 5
	}

	offset := (page - 1) * limit

	query := database.DB.Preload("Brand").Preload("Category")

	if categoryStr != "" {
		if catID, err := strconv.Atoi(categoryStr); err == nil {
			query = query.Where("category_id = ?", catID)
		}
	}
	if brandStr != "" {
		if brandID, err := strconv.Atoi(brandStr); err == nil {
			query = query.Where("brand_id = ?", brandID)
		}
	}

	var products []models.Product
	query.Offset(offset).Limit(limit).Find(&products)

	c.JSON(http.StatusOK, products)
}

func GetProductByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var product models.Product
	result := database.DB.Preload("Brand").Preload("Category").First(&product, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

func CreateProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	if product.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name is required"})
		return
	}
	if product.Price < 0.01 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Price must be at least 0.01"})
		return
	}
	if product.Stock < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Stock cannot be negative"})
		return
	}

	var brand models.Brand
	if err := database.DB.First(&brand, product.BrandID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid brand_id"})
		return
	}

	var category models.Category
	if err := database.DB.First(&category, product.CategoryID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category_id"})
		return
	}

	if err := database.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	database.DB.Preload("Brand").Preload("Category").First(&product, product.ID)
	c.JSON(http.StatusCreated, product)
}

func UpdateProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var product models.Product
	if err := database.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	var updated models.Product
	if err := c.ShouldBindJSON(&updated); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	if updated.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name is required"})
		return
	}
	if updated.Price < 0.01 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Price must be at least 0.01"})
		return
	}
	if updated.Stock < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Stock cannot be negative"})
		return
	}

	var brand models.Brand
	if err := database.DB.First(&brand, updated.BrandID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid brand_id"})
		return
	}

	var category models.Category
	if err := database.DB.First(&category, updated.CategoryID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category_id"})
		return
	}

	updated.ID = uint(id)
	database.DB.Save(&updated)

	database.DB.Preload("Brand").Preload("Category").First(&updated, id)
	c.JSON(http.StatusOK, updated)
}

func DeleteProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	result := database.DB.Delete(&models.Product{}, id)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	c.Status(http.StatusNoContent)
}
