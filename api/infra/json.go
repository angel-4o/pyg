package infra

import (
	"encoding/json"
	"net/http"
)

func parseBody[T any](request *http.Request, output *T) error {
	return json.NewDecoder(request.Body).Decode(&output)
}
