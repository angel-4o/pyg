package systemtest

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type registerParams struct {
	Email     string `json:"email"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func newRegisterParams() *registerParams {
	return &registerParams{
		Email:     "user@email.com",
		Username:  "username",
		Password:  "password",
		FirstName: "firstName",
		LastName:  "lastName",
	}
}

func TestRegister(t *testing.T) {
	db, handler := newDbAndHandler()
	defer db.Close()

	t.Run("registraion of new user succeeds", func(t *testing.T) {
		response := sendRequest(t, handler, "POST", "/v1/register", newRegisterParams(), nil)
		assert.Equal(t, http.StatusCreated, response.Result().StatusCode)
	})

	t.Run("duplicate users are not allowed", func(t *testing.T) {
		response := sendRequest(t, handler, "POST", "/v1/register", newRegisterParams(), nil)
		assert.Equal(t, http.StatusConflict, response.Result().StatusCode)
	})
}
