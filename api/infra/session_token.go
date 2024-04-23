package infra

import (
	"net/http"

	"pyg.com/api/app/data"
)

func getSessionToken(request* http.Request) data.SessionToken {
	cookie, err := request.Cookie("sessionToken")
	if err != nil {
		return ""
	}

	return data.SessionToken(cookie.Value)
}