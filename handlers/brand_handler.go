package handlers

import (
	"net/http"
	"strconv"

	"electronicsStore/database"
	"electronicsStore/models"

	"github.com/gin-gonic/gin"
)

func GetBrands(c *gin.Context) {
	var brands []models.Brand
	database.DB.Find(&brands)
	c.JSON(http.StatusOK, brands)
}

func CreateBrand(c *gin.Context) {
	var brand models.Brand
	if err := c.ShouldBindJSON(&brand); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	if brand.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name is required"})
		return
	}

	if err := database.DB.Create(&brand).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create brand"})
		return
	}
	c.JSON(http.StatusCreated, brand)
}

func DeleteBrand(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	result := database.DB.Delete(&models.Brand{}, id)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Brand not found"})
		return
	}
	c.Status(http.StatusNoContent)
}
