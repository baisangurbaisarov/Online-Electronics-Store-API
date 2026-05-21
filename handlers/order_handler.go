package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"electronicsStore/database"
	"electronicsStore/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type orderItemInput struct {
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
}

func CreateOrder(c *gin.Context) {
	userID, ok := requireUserID(c)
	if !ok {
		return
	}

	var input struct {
		Items []orderItemInput `json:"items"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	if len(input.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order must contain at least one item"})
		return
	}

	var order models.Order
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		var total float64
		lineItems := make([]models.OrderItem, 0, len(input.Items))

		for _, item := range input.Items {
			if item.ProductID == 0 || item.Quantity <= 0 {
				return errInvalidOrderItem
			}

			var product models.Product
			if err := tx.First(&product, item.ProductID).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return errProductNotFound
				}
				return err
			}
			if product.Stock < item.Quantity {
				return errInsufficientStock
			}

			lineTotal := product.Price * float64(item.Quantity)
			total += lineTotal

			lineItems = append(lineItems, models.OrderItem{
				ProductID: product.ID,
				Quantity:  item.Quantity,
				Price:     product.Price,
			})

			if err := tx.Model(&product).Update("stock", product.Stock-item.Quantity).Error; err != nil {
				return err
			}
		}

		order = models.Order{
			UserID:    userID,
			Status:    "placed",
			Total:     total,
			CreatedAt: time.Now(),
		}
		if err := tx.Create(&order).Error; err != nil {
			return err
		}

		for i := range lineItems {
			lineItems[i].OrderID = order.ID
		}
		if err := tx.Create(&lineItems).Error; err != nil {
			return err
		}

		order.Items = lineItems
		return nil
	})

	if err != nil {
		switch err {
		case errInvalidOrderItem:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Each item needs a valid product_id and quantity > 0"})
		case errProductNotFound:
			c.JSON(http.StatusBadRequest, gin.H{"error": "One or more products not found"})
		case errInsufficientStock:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient stock for one or more products"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		}
		return
	}

	database.DB.Preload("Items").Preload("Items.Product").First(&order, order.ID)
	c.JSON(http.StatusCreated, order)
}

func GetMyOrders(c *gin.Context) {
	userID, ok := requireUserID(c)
	if !ok {
		return
	}

	var orders []models.Order
	database.DB.
		Preload("Items").
		Preload("Items.Product").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&orders)

	c.JSON(http.StatusOK, orders)
}

func GetOrderByID(c *gin.Context) {
	userID, ok := requireUserID(c)
	if !ok {
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	var order models.Order
	result := database.DB.
		Preload("Items").
		Preload("Items.Product").
		Where("id = ? AND user_id = ?", id, userID).
		First(&order)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}

var (
	errInvalidOrderItem  = errors.New("invalid order item")
	errProductNotFound   = errors.New("product not found")
	errInsufficientStock = errors.New("insufficient stock")
)
