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

func CreateGame(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		sessionToken := getSessionToken(request)
		if !sessionToken.IsValid() {
			http.Error(writer, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		// errr := request.ParseMultipartForm(10 << 20)
		var params struct {
			DeveloperName string
			Name          string
			Description   string
			Genre         string
			Platform      string
		}
		err := parseBody(request, &params)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		genre, err := domain.GenreFromString(params.Genre)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		platform, err := domain.PlatformFromString(params.Platform)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		// domain.Developer.Id // todo use developer id
		gameId, err := app.CreateGame2(sessionToken, params.DeveloperName,
			params.Name, params.Description, genre, platform, db)

		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		var response = struct {
			Id int64 `json:"id"`
		}{
			Id: int64(gameId),
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

func GetGame(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		id, err := strconv.Atoi(request.PathValue("id"))
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		gameRepo := persistence.MakeGameRepo(db)
		game, err := gameRepo.Get(domain.GameId(id))
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		responseBody, err := json.Marshal(game)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		writer.Write(responseBody)
	})
}

func GetGames(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		var params struct {
			PageToken int64            `json:"pageToken"`
			Genre     *domain.Genre    `json:"genre"`
			Platform  *domain.Platform `json:"platform"`
		}
		err := parseBody(request, &params)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		gameRepo := persistence.MakeGameRepo(db)
		developerRepo := persistence.MakeDeveloperRepo(db)
		games, nextPageToken, err := gameRepo.GetGames2(data.PageToken(params.PageToken), params.Genre, params.Platform, developerRepo)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		var response = struct {
			PageToken int64         `json:"pageToken"`
			Games     []domain.Game `json:"developers"`
		}{
			PageToken: int64(nextPageToken),
			Games:     games,
		}

		responseBody, err := json.Marshal(response)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		writer.Write(responseBody)
	})
}

func GetGames2(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		var params struct {
			PageToken int64            `json:"pageToken"`
			Genre     *domain.Genre    `json:"genre"`
			Platform  *domain.Platform `json:"platform"`
		}
		err := parseBody(request, &params)
		// if err != nil {
		// 	http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		// 	return
		// }

		gameRepo := persistence.MakeGameRepo(db)
		developerRepo := persistence.MakeDeveloperRepo(db)
		games, nextPageToken, err := gameRepo.GetGames2(data.PageToken(params.PageToken), params.Genre, params.Platform, developerRepo)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		var response = struct {
			PageToken int64         `json:"pageToken"`
			Games     []domain.Game `json:"developers"`
		}{
			PageToken: int64(nextPageToken),
			Games:     games,
		}

		responseBody, err := json.Marshal(response)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		writer.Write(responseBody)
	})
}
