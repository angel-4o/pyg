package infra

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"pyg.com/api/app"
	"pyg.com/api/app/data"
	"pyg.com/api/domain"
	"pyg.com/api/persistence"
)

const (
	pageSize = 30
)

func GetDeveloper(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		id, err := strconv.Atoi(request.PathValue("id"))
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		developer, err := app.GetDeveloper(domain.DeveloperId(id), db)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		responseBody, err := json.Marshal(developer)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		writer.Write(responseBody)
	})
}

func GetDevelopers(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		var params struct {
			PageToken int64 `json:"pageToken"`
		}
		err := parseBody(request, &params)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		developersRepo := persistence.MakeDeveloperRepo(db)
		developers, nextPageToken, err := developersRepo.GetDevelopers(data.PageToken(params.PageToken), pageSize)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		var response = struct {
			PageToken  int64              `json:"pageToken"`
			Developers []domain.Developer `json:"developers"`
		}{
			PageToken:  int64(nextPageToken),
			Developers: developers,
		}

		responseBody, err := json.Marshal(response)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		writer.Write(responseBody)
	})
}

func CreateDeveloper(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		sessionToken := getSessionToken(request)
		if !sessionToken.IsValid() {
			http.Error(writer, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		var params struct {
			Name string `json:"name"`
		}
		err := parseBody(request, &params)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		developerId, err := app.CreateDeveloper(sessionToken, params.Name, db)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		var response = struct {
			Id int64 `json:"id"`
		}{
			Id: int64(developerId),
		}

		responseBody, err := json.Marshal(response)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		writer.WriteHeader(http.StatusCreated)
		writer.Write(responseBody)
	})
}

func UpdateDeveloper(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		sessionToken := getSessionToken(request)
		if !sessionToken.IsValid() {
			http.Error(writer, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		var params struct {
			Id         int64             `json:"id"`
			Name       *string           `json:"name"`
			MemberList domain.MemberList `json:"memberList"`
		}
		err := parseBody(request, &params)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		err = app.UpdateDeveloper(sessionToken, domain.DeveloperId(params.Id), params.Name, params.MemberList, db)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		writer.WriteHeader(http.StatusOK)
	})
}
