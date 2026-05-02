package handlers

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"electronicsStore/database"
	"electronicsStore/models"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

type sentimentResult struct {
	Sentiment string `json:"sentiment"`
}

func analyzeSentiment(comment string) string {
	if comment == "" {
		return "neutral"
	}

	sentimentURL := os.Getenv("SENTIMENT_URL")
	if sentimentURL == "" {
		sentimentURL = "http://localhost:9090"
	}

	client := resty.New()

	var result sentimentResult
	_, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]string{"text": comment}).
		SetResult(&result).
		Post(sentimentURL + "/analyze")

	if err != nil {
		log.Printf("sentiment service error: %v", err)
		return "unknown"
	}

	return result.Sentiment
}

func GetReviews(c *gin.Context) {
	productID, err := strconv.Atoi(c.Param("id"))
	if err != nil || productID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var product models.Product
	if err := database.DB.First(&product, productID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	var reviews []models.Review
	database.DB.Preload("User").Where("product_id = ?", productID).Find(&reviews)
	c.JSON(http.StatusOK, reviews)
}

func CreateReview(c *gin.Context) {
	productID, err := strconv.Atoi(c.Param("id"))
	if err != nil || productID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var product models.Product
	if err := database.DB.First(&product, productID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	var input struct {
		Rating  int    `json:"rating"`
		Comment string `json:"comment"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	if input.Rating < 1 || input.Rating > 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Rating must be between 1 and 5"})
		return
	}

	userIDRaw, _ := c.Get("userID")
	userIDFloat, ok := userIDRaw.(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Could not identify user"})
		return
	}
	userID := uint(userIDFloat)

	sentiment := analyzeSentiment(input.Comment)

	review := models.Review{
		ProductID: uint(productID),
		UserID:    userID,
		Rating:    input.Rating,
		Comment:   input.Comment,
		Sentiment: sentiment,
	}

	if err := database.DB.Create(&review).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create review"})
		return
	}

	database.DB.Preload("User").First(&review, review.ID)
	c.JSON(http.StatusCreated, review)
}

func DeleteReview(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	userIDRaw, _ := c.Get("userID")
	userIDFloat, ok := userIDRaw.(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Could not identify user"})
		return
	}
	userID := uint(userIDFloat)

	var review models.Review
	if err := database.DB.First(&review, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
		return
	}

	if review.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own reviews"})
		return
	}

	database.DB.Delete(&review)
	c.Status(http.StatusNoContent)
}
