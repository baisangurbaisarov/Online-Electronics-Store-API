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

func TestGetProducts_Empty(t *testing.T) {
	setupTestDB(t)

	w := performJSONRequest(handlers.GetProducts, http.MethodGet, "/products", nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var products []models.Product
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &products))
	assert.Empty(t, products)
}

func TestGetProductByID_NotFound(t *testing.T) {
	setupTestDB(t)

	w := performRequestWithParam(handlers.GetProductByID, http.MethodGet, "/products/99", "id", "99", nil)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCreateProduct_MissingName(t *testing.T) {
	setupTestDB(t)

	brand := models.Brand{Name: "Sony"}
	category := models.Category{Name: "TV"}
	require.NoError(t, database.DB.Create(&brand).Error)
	require.NoError(t, database.DB.Create(&category).Error)

	w := performJSONRequest(handlers.CreateProduct, http.MethodPost, "/products", map[string]any{
		"name":        "",
		"price":       99.99,
		"stock":       5,
		"brand_id":    brand.ID,
		"category_id": category.ID,
	})

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateProduct_Success(t *testing.T) {
	setupTestDB(t)

	brand := models.Brand{Name: "Apple"}
	category := models.Category{Name: "Phones"}
	require.NoError(t, database.DB.Create(&brand).Error)
	require.NoError(t, database.DB.Create(&category).Error)

	w := performJSONRequest(handlers.CreateProduct, http.MethodPost, "/products", map[string]any{
		"name":        "iPhone",
		"price":       999.99,
		"stock":       10,
		"brand_id":    brand.ID,
		"category_id": category.ID,
	})

	assert.Equal(t, http.StatusCreated, w.Code)

	var product models.Product
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &product))
	assert.Equal(t, "iPhone", product.Name)
}
