package handlers_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"electronicsStore/database"
	"electronicsStore/handlers"
	"electronicsStore/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateOrder_Success(t *testing.T) {
	setupTestDB(t)

	brand := models.Brand{Name: "Dell"}
	category := models.Category{Name: "Laptops"}
	require.NoError(t, database.DB.Create(&brand).Error)
	require.NoError(t, database.DB.Create(&category).Error)

	product := models.Product{
		Name: "XPS", Price: 1200, Stock: 5,
		BrandID: brand.ID, CategoryID: category.ID,
	}
	require.NoError(t, database.DB.Create(&product).Error)

	user := models.User{Username: "buyer", Password: "hashed"}
	require.NoError(t, database.DB.Create(&user).Error)

	w := performJSONRequestAsUser(handlers.CreateOrder, http.MethodPost, "/orders", map[string]any{
		"items": []map[string]any{
			{"product_id": product.ID, "quantity": 2},
		},
	}, user.ID)

	assert.Equal(t, http.StatusCreated, w.Code)

	var order models.Order
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &order))
	assert.Equal(t, user.ID, order.UserID)
	assert.Equal(t, 2400.0, order.Total)

	var updated models.Product
	require.NoError(t, database.DB.First(&updated, product.ID).Error)
	assert.Equal(t, 3, updated.Stock)
}

func TestGetMyOrders_OnlyOwnOrders(t *testing.T) {
	setupTestDB(t)

	user1 := models.User{Username: "u1", Password: "x"}
	user2 := models.User{Username: "u2", Password: "x"}
	require.NoError(t, database.DB.Create(&user1).Error)
	require.NoError(t, database.DB.Create(&user2).Error)

	require.NoError(t, database.DB.Create(&models.Order{UserID: user1.ID, Status: "placed", Total: 10}).Error)
	require.NoError(t, database.DB.Create(&models.Order{UserID: user2.ID, Status: "placed", Total: 99}).Error)

	w := performJSONRequestAsUser(handlers.GetMyOrders, http.MethodGet, "/orders", nil, user1.ID)
	assert.Equal(t, http.StatusOK, w.Code)

	var orders []models.Order
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &orders))
	require.Len(t, orders, 1)
	assert.Equal(t, user1.ID, orders[0].UserID)
	assert.Equal(t, 10.0, orders[0].Total)
}
