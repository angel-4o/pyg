package infra

import (
	"database/sql"
	"net/http"
)

func NewRootHandler(db *sql.DB) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("POST /v1/register", contentTypeMiddleware(Register(db)))
	mux.Handle("POST /v1/login", contentTypeMiddleware(Login(db)))
	mux.Handle("POST /v1/login2", contentTypeMiddleware(Login2(db)))

	// we are not going to use these endpoints
	mux.Handle("POST /v1/developer", contentTypeMiddleware(CreateDeveloper(db)))
	mux.Handle("PUT /v1/developer", contentTypeMiddleware(UpdateDeveloper(db)))
	mux.Handle("GET /v1/developer/{id}", contentTypeMiddleware(GetDeveloper(db)))
	mux.Handle("GET /v1/developers", contentTypeMiddleware(GetDevelopers(db)))
	mux.Handle("POST /v1/publisher", contentTypeMiddleware(CreatePublisher(db)))
	mux.Handle("PUT /v1/publisher", contentTypeMiddleware(UpdatePublisher(db)))
	mux.Handle("GET /v1/publisher/{id}", contentTypeMiddleware(GetPublisher(db)))
	mux.Handle("GET /v1/publishers", contentTypeMiddleware(GetPublisher(db)))
	// until here

	mux.Handle("POST /v1/game", contentTypeMiddleware(CreateGame(db)))
	mux.Handle("GET /v1/game/{id}", contentTypeMiddleware(GetGame(db)))
	mux.Handle("GET /v1/games", contentTypeMiddleware(GetGames(db)))
	mux.Handle("GET /v1/allGames", contentTypeMiddleware(GetGames2(db)))

	/*
		/publishers
		/publisher/{id}
		/developer/{id}
		/developers

		/games/genre/...
		/games (filter)
		/game/{id}
		/developer/requests

		/games/published
		/games/requested
		/games/saved


		/developer:send_message
	*/

	return mux
}
