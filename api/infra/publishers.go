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

func GetPublisher(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		id, err := strconv.Atoi(request.PathValue("id"))
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		publisher, err := app.GetPublisher(domain.PublisherId(id), db)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		responseBody, err := json.Marshal(publisher)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		writer.Write(responseBody)
	})
}

func GetPublishers(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		var params struct {
			PageToken int64 `json:"pageToken"`
		}
		err := parseBody(request, &params)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		publishersRepo := persistence.MakePublisherRepo(db)
		publishers, nextPageToken, err := publishersRepo.GetPublishers(data.PageToken(params.PageToken), pageSize)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		var response = struct {
			PageToken  int64              `json:"pageToken"`
			Publishers []domain.Publisher `json:"publishers"`
		}{
			PageToken:  int64(nextPageToken),
			Publishers: publishers,
		}

		responseBody, err := json.Marshal(response)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		writer.Write(responseBody)
	})
}

func CreatePublisher(db *sql.DB) http.Handler {
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

		publisherId, err := app.CreatePublisher(sessionToken, params.Name, db)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		var response = struct {
			Id int64 `json:"id"`
		}{
			Id: int64(publisherId),
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

func UpdatePublisher(db *sql.DB) http.Handler {
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

		err = app.UpdatePublisher(sessionToken, domain.PublisherId(params.Id), params.Name, params.MemberList, db)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		writer.WriteHeader(http.StatusOK)
	})
}
