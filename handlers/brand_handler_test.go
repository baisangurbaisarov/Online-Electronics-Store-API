package handlers

import (
	"encoding/json"
	"net/http"
	"testing"

	"electronicsStore/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateBrand_Success(t *testing.T) {
	setupTestDB(t)

	w := performJSONRequest(CreateBrand, http.MethodPost, "/brands", map[string]string{
		"name": "Samsung",
	})

	assert.Equal(t, http.StatusCreated, w.Code)

	var brand models.Brand
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &brand))
	assert.Equal(t, "Samsung", brand.Name)
}

func TestDeleteBrand_NotFound(t *testing.T) {
	setupTestDB(t)

	w := performRequestWithParam(DeleteBrand, http.MethodDelete, "/brands/42", "id", "42", nil)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
