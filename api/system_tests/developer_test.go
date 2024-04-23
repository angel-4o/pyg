package systemtest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"pyg.com/api/domain"
)

type createDeveloperParams struct {
	Name string
}

func newCreateDeveloperParams() *createDeveloperParams {
	return &createDeveloperParams{
		Name: "name",
	}
}

type updateDeveloperParams struct {
	Id         int64
	Name       *string
	MemberList domain.MemberList
}

func newUpdateDeveloperParams(id int64) *updateDeveloperParams {
	newName := "newName"
	return &updateDeveloperParams{
		Id:   id,
		Name: &newName,
	}
}

type getDeveloperParams struct {
	Id int64
}

func newGetDeveloperParams(id int64) *getDeveloperParams {
	return &getDeveloperParams{
		Id: id,
	}
}

func TestDeveloper(t *testing.T) {
	db, handler := newDbAndHandler()
	defer db.Close()

	registerResponse := sendRequest(t, handler, "POST", "/v1/register", newRegisterParams(), nil)
	require.Equal(t, http.StatusCreated, registerResponse.Result().StatusCode)

	loginResponse := sendRequest(t, handler, "POST", "/v1/login", newLoginParams(), nil)
	require.Equal(t, http.StatusOK, loginResponse.Result().StatusCode)

	loginResult := struct {
		SessionToken string `json:"sessionToken"`
	}{}
	err := json.NewDecoder(loginResponse.Body).Decode(&loginResult)
	require.Nil(t, err)

	sessionToken := loginResult.SessionToken

	t.Run("anonymous user cannot create developer", func(t *testing.T) {
		response := sendRequest(t, handler, "POST", "/v1/developer", newCreateDeveloperParams(), nil)
		assert.Equal(t, http.StatusUnauthorized, response.Result().StatusCode)
	})

	var developerId int64

	t.Run("creation of new developer succeeds", func(t *testing.T) {
		response := sendRequest(t, handler, "POST", "/v1/developer", newCreateDeveloperParams(), &sessionToken)
		require.Equal(t, http.StatusCreated, response.Result().StatusCode)

		var creationResponse struct {
			Id int64
		}
		parseResponse(response, &creationResponse)
		developerId = creationResponse.Id

		response = sendRequest(t, handler, "GET", fmt.Sprintf("/v1/developer/%d", developerId), newGetDeveloperParams(developerId), nil)
		assert.Equal(t, http.StatusOK, response.Result().StatusCode)

		var developer domain.Developer
		parseResponse(response, &developer)
		assert.Equal(t, "name", developer.Name)
	})

	registerParamsForNewUser := newRegisterParams()
	registerParamsForNewUser.Email = "new@email.com"
	registerParamsForNewUser.Username = "newUsername"
	registerResponse = sendRequest(t, handler, "POST", "/v1/register", registerParamsForNewUser, nil)
	require.Equal(t, http.StatusCreated, registerResponse.Result().StatusCode)

	t.Run("developer update is successful", func(t *testing.T) {
		response := sendRequest(t, handler, "PUT", "/v1/developer", newUpdateDeveloperParams(developerId), &sessionToken)
		assert.Equal(t, http.StatusOK, response.Result().StatusCode)

		response = sendRequest(t, handler, "GET", fmt.Sprintf("/v1/developer/%d", developerId), newGetDeveloperParams(developerId), nil)
		assert.Equal(t, http.StatusOK, response.Result().StatusCode)

		var developer domain.Developer
		parseResponse(response, &developer)
		assert.Equal(t, "newName", developer.Name)
	})
}
