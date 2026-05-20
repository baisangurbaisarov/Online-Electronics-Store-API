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
	"golang.org/x/crypto/bcrypt"
)

func TestRegister_Success(t *testing.T) {
	setupTestDB(t)

	w := performJSONRequest(handlers.Register, http.MethodPost, "/auth/register", map[string]string{
		"username": "alice",
		"password": "secret1",
	})

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "alice", resp["username"])
}

func TestRegister_DuplicateUsername(t *testing.T) {
	setupTestDB(t)

	body := map[string]string{"username": "bob", "password": "secret1"}
	performJSONRequest(handlers.Register, http.MethodPost, "/auth/register", body)
	w := performJSONRequest(handlers.Register, http.MethodPost, "/auth/register", body)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestRegister_ShortPassword(t *testing.T) {
	setupTestDB(t)

	w := performJSONRequest(handlers.Register, http.MethodPost, "/auth/register", map[string]string{
		"username": "carol",
		"password": "123",
	})

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLogin_InvalidCredentials(t *testing.T) {
	setupTestDB(t)

	hashed, _ := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.DefaultCost)
	require.NoError(t, database.DB.Create(&models.User{Username: "dave", Password: string(hashed)}).Error)

	w := performJSONRequest(handlers.Login, http.MethodPost, "/auth/login", map[string]string{
		"username": "dave",
		"password": "wrong",
	})

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLogin_Success(t *testing.T) {
	setupTestDB(t)

	performJSONRequest(handlers.Register, http.MethodPost, "/auth/register", map[string]string{
		"username": "eve",
		"password": "secret1",
	})

	w := performJSONRequest(handlers.Login, http.MethodPost, "/auth/login", map[string]string{
		"username": "eve",
		"password": "secret1",
	})

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.NotEmpty(t, resp["token"])
}
