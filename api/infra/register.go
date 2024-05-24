package infra

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"pyg.com/api/app"
	"pyg.com/api/domain"
)

func Register(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		var params struct {
			Email     string `json:"email"`
			Username  string `json:"username"`
			FirstName string `json:"firstName"`
			LastName  string `json:"lastName"`
			Password  string `json:"password"`
			Type      string `json:"type"`
		}

		err := json.NewDecoder(request.Body).Decode(&params)
		if err != nil {
			log.Println(err.Error())
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		err = app.Register(params.Email, params.Username, params.FirstName, params.LastName, params.Password, db)
		if err != nil {
			log.Println(err.Error())
		}

		switch {
		case errors.Is(err, domain.ErrAccountExists):
			http.Error(writer, http.StatusText(http.StatusConflict), http.StatusConflict)
		case errors.Is(err, domain.ErrValidationFailed):
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		case err != nil:
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		default:
			writer.WriteHeader(http.StatusCreated)
		}
	})
}
