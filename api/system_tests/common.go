package systemtest

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func sendRequest(t *testing.T, handler http.Handler, method, api string, params any, sessionToken *string) *httptest.ResponseRecorder {
	requestBody, err := json.Marshal(params)
	require.Nil(t, err)

	request, err := http.NewRequest(method, api, strings.NewReader(string(requestBody)))
	require.Nil(t, err)

	if sessionToken != nil {
		request.AddCookie(&http.Cookie{
			Name:  "sessionToken",
			Value: *sessionToken,
		})
	}

	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	return response
}

func parseResponse[T any](response *httptest.ResponseRecorder, output *T) error {
	return json.NewDecoder(response.Body).Decode(&output)
}
