package handlers_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"electronicsStore/handlers"
	"electronicsStore/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateCategory_Success(t *testing.T) {
	setupTestDB(t)

	w := performJSONRequest(handlers.CreateCategory, http.MethodPost, "/categories", map[string]string{
		"name": "Laptops",
	})

	assert.Equal(t, http.StatusCreated, w.Code)

	var category models.Category
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &category))
	assert.Equal(t, "Laptops", category.Name)
}
