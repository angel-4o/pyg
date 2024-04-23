package systemtest

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type loginParams struct {
	Username string
	Password string
}

func newLoginParams() *loginParams {
	return &loginParams{
		Username: "username",
		Password: "password",
	}
}

func TestLogin(t *testing.T) {
	db, handler := newDbAndHandler()
	defer db.Close()

	registerResponse := sendRequest(t, handler, "POST", "/v1/register", newRegisterParams(), nil)
	require.Equal(t, http.StatusCreated, registerResponse.Result().StatusCode)

	t.Run("login succeeds with correct credentials", func(t *testing.T) {
		response := sendRequest(t, handler, "POST", "/v1/login", newLoginParams(), nil)
		assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	})

	t.Run("login fails with invalid password", func(t *testing.T) {
		params := newLoginParams()
		params.Password = "wrongPassword"
		response := sendRequest(t, handler, "POST", "/v1/login", params, nil)
		assert.Equal(t, http.StatusUnauthorized, response.Result().StatusCode)
	})

	t.Run("login fails with nonexisting user", func(t *testing.T) {
		params := newLoginParams()
		params.Username = "nonexisting"
		response := sendRequest(t, handler, "POST", "/v1/login", params, nil)
		assert.Equal(t, http.StatusUnauthorized, response.Result().StatusCode)
	})
}
