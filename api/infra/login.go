package infra

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"pyg.com/api/app"
	apperrors "pyg.com/api/app/errors"
)

func Login(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		var params struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		err := json.NewDecoder(request.Body).Decode(&params)
		if err != nil {
			log.Println(err.Error())
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		sessionToken, err := app.Login(params.Username, params.Password, db)

		var responseJson []byte
		if err == nil {
			response := struct {
				SessionToken string `json:"sessionToken"`
			}{
				SessionToken: string(sessionToken),
			}

			responseJson, err = json.Marshal(response)
		}

		switch {
		case errors.Is(err, apperrors.ErrAccountNotFound), errors.Is(err, apperrors.ErrInvalidPassword):
			http.Error(writer, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)

		case err != nil:
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		}

		if err != nil {
			log.Println(err.Error())
			return
		}

		if responseJson == nil {
			panic("login: jsonResponse expected not nil (assertion failure)")
		}

		writer.Write(responseJson)
	})
}
