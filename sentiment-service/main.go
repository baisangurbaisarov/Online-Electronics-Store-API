package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Request struct {
	Text string `json:"text"`
}

type Response struct {
	Sentiment string `json:"sentiment"`
}

func analyzeSentiment(text string) string {
	if text == "" {
		return "neutral"
	}

	lower := strings.ToLower(text)

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

	pos, neg := 0, 0
	for _, w := range positiveWords {
		if strings.Contains(lower, w) {
			pos++
		}
	}
	for _, w := range negativeWords {
		if strings.Contains(lower, w) {
			neg++
		}
	}

	switch {
	case pos > neg:
		return "positive"
	case neg > pos:
		return "negative"
	default:
		return "neutral"
	}
}

func main() {
	r := gin.Default()

	r.POST("/analyze", func(c *gin.Context) {
		var req Request
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}
		sentiment := analyzeSentiment(req.Text)
		c.JSON(http.StatusOK, Response{Sentiment: sentiment})
	})

	r.Run(":9090")
}
