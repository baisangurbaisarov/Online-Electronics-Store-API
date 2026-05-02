package handlers

import (
	"log"
	"net/http"
	"strconv"
	"strings"

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

	// Используем resty для запроса к httpbin (демонстрация библиотеки)
	client := resty.New()

	var result map[string]interface{}
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]string{"text": comment}).
		SetResult(&result).
		Post("https://httpbin.org/anything")

	if err != nil {
		log.Printf("sentiment API unavailable, using local analysis: %v", err)
	} else {
		log.Printf("sentiment API status: %d", resp.StatusCode())
	}

	// Локальный анализ по ключевым словам
	lower := strings.ToLower(comment)

	positiveWords := []string{
		"good", "great", "excellent", "amazing", "love", "perfect",
		"awesome", "fantastic", "best", "nice", "happy", "recommend",
		"хорошо", "отлично", "прекрасно", "люблю", "нравится", "советую",
	}
	negativeWords := []string{
		"bad", "terrible", "awful", "hate", "worst", "poor",
		"broken", "disappointing", "useless", "horrible", "never",
		"плохо", "ужасно", "ненавижу", "сломан", "разочарован", "никогда",
	}

	posScore := 0
	negScore := 0

	for _, w := range positiveWords {
		if strings.Contains(lower, w) {
			posScore++
		}
	}
	for _, w := range negativeWords {
		if strings.Contains(lower, w) {
			negScore++
		}
	}

	switch {
	case posScore > negScore:
		return "positive"
	case negScore > posScore:
		return "negative"
	default:
		return "neutral"
	}
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