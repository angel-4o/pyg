package infra

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"pyg.com/api/app"
	apperrors "pyg.com/api/app/errors"
	"pyg.com/api/domain"
	"pyg.com/api/persistence"

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

func Login2(db *sql.DB) http.Handler {
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

		account, sessionToken, err := app.Login2(params.Username, params.Password, db)

		var responseJson []byte
		if err == nil {
			// Create a new struct without the password field
			response := struct {
				SessionToken string                 `json:"sessionToken"`
				Account      domain.AccountResponse `json:"account"`
			}{
				SessionToken: string(sessionToken),
				Account: domain.AccountResponse{
					Id:        account.Id,
					Username:  account.Username,
					Email:     account.Email,
					FirstName: account.FirstName,
					LastName:  account.LastName,
					//Type:      account.Type,

					// Add other fields you want to include in the response
				},
			}
			if account.LastName == "developer" {
				developerRepo := persistence.MakeDeveloperRepo(db)

				//persistence.MakeAccountRepo()
				//account_repo
				//developerRepo.find§
				developer, err := developerRepo.FindByName(params.Username)
				if (developer == nil || err != nil){
					app.CreateDeveloper(sessionToken, params.Username, db)
				}

				// dev, er := persistence.FindByName(params.Username, db);
				// developerId, err := 
				// if err != nil {

				// }
				//fmt.Print(developerId)
				// fmt.Printf("developerId: %v\n", developerId)
			}
			if account.LastName == "publisher" {
				publisherRepo := persistence.MakePublisherRepo(db)
				publisher, err := publisherRepo.FindByName(params.Username)
				if (publisher == nil || err != nil){
					app.CreatePublisher(getSessionToken(request), params.Username, db)
				}
				
			}
			responseJson, err = json.Marshal(response)
		}
		// type Account struct {
		// 	Id        AccountId
		// 	Email     string `validate:"email,required"`
		// 	Username  string `validate:"required"`
		// 	FirstName string `validate:"required"`
		// 	LastName  string `validate:"required"`
		// 	Password  string `validate:"required,max=72"`
		// }

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
