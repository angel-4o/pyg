package persistence

import (
	"database/sql"

	"pyg.com/api/app/data"
	apperrors "pyg.com/api/app/errors"
	"pyg.com/api/domain"

)

type gameRepo struct {
	db *sql.DB
}

func MakeGameRepo(db *sql.DB) *gameRepo {
	return &gameRepo{
		db: db,
	}
}

func (repo *gameRepo) Create(game *domain.Game) (domain.GameId, error) {
	err := repo.db.QueryRow(
		`insert into games (developer_id, name, description, genre, platform)
		 values ($1, $2, $3, $4, $5)
		 returning id`,
		game.DeveloperId,
		game.Name,
		game.Description,
		domain.GenreToString(game.Genre),
		domain.PlatformToString(game.Platform),
	).Scan(
		&game.Id,
	)

	if err != nil {
		return domain.GameId(0), err
	}

	return game.Id, nil
}

func (repo *gameRepo) Get(gameId domain.GameId) (*domain.Game, error) {
	game := domain.Game{}

	var genre, platform string

	err := repo.db.QueryRow(`
		select id, developer_id, name, description, genre, platform
		from games
		where id = $1`,
		gameId,
	).Scan(
		&game.Id,
		&game.DeveloperId,
		&game.Name,
		&game.Description,
		&genre,
		&platform,
	)

	if err != nil {
		return &domain.Game{}, err
	}

	game.Genre, err = domain.GenreFromString(genre)
	if err != nil {
		return &domain.Game{}, err
	}

	game.Platform, err = domain.PlatformFromString(platform)
	if err != nil {
		return &domain.Game{}, err
	}

	return &game, nil
}

func (repo *gameRepo) GetGames(pageToken data.PageToken, genre *domain.Genre, platform *domain.Platform) ([]domain.Game, data.PageToken, error) {
	rows, _ := repo.db.Query(`
		select id, developer_id, name, description, genre, platform
		from games
		where id > $1
		limit 30`,
		pageToken,
	)

	games := make([]domain.Game, 0)

	for rows.Next() {
		var game domain.Game
		var genre, platform string

		err := rows.Scan(
			&game.Id,
			&game.DeveloperId,
			&game.Name,
			&game.Description,
			&genre,
			&platform,
		)

		if err != nil {
			return nil, pageToken, apperrors.ErrNotFound
		}
		// var repo := persistence.MakeAccountRepo(db)
		games = append(games, game)
	}

	nextPageToken := pageToken
	if count := len(games); count > 0 {
		nextPageToken = data.PageToken(games[count-1].Id)
	}

	return games, nextPageToken, nil
}


func (repo *gameRepo) GetGames2(pageToken data.PageToken, genre *domain.Genre, platform *domain.Platform, dp *developerRepo) ([]domain.Game, data.PageToken, error) {
	rows, _ := repo.db.Query(`
		select id, developer_id, name, description, genre, platform
		from games
		where id > $1
		limit 10000`,
		pageToken,
	)

	games := make([]domain.Game, 0)

	for rows.Next() {
		var game domain.Game
		var genreString, platformString string

		err := rows.Scan(
			&game.Id,
			&game.DeveloperId,
			&game.Name,
			&game.Description,
			&genreString,
			&platformString,
		)

		if err != nil {
			return nil, pageToken, apperrors.ErrNotFound
		}

		dev, err := dp.FindById(game.DeveloperId)
		if (err != nil){
			game.DeveloperName = "Error"
		}else {
			game.DeveloperName = dev.Name;
		}
		
		g, err := domain.GenreFromString(genreString)
		if (err == nil){
			game.Genre = g
		}
		p, err := domain.PlatformFromString(platformString)
		if (err == nil){
			game.Platform = p
		}
		// test := *platform == nil ? "" : ""
		// addPlatform := false
		// if (platform == nil || game.Platform == *platform){
		// 	addPlatform = true
		// }
		// addGenre := false
		// if (genre == nil || game.Genre == *genre){
		// 	addGenre = true;
		// }
		// if (addPlatform && addGenre){
		// 	games = append(games, game)
		// }

		
		if (platform == nil && genre == nil){
			games = append(games, game)
		} else if (platform == nil && game.Genre == *genre){
			games = append(games, game)
		} else if (genre == nil && game.Platform == *platform){
			games = append(games, game)
		} else if (genre != nil && platform != nil && game.Platform == *platform && game.Genre == *genre){
			games = append(games, game)
		}
	}

	nextPageToken := pageToken
	if count := len(games); count > 0 {
		nextPageToken = data.PageToken(games[count-1].Id)
	}

	return games, nextPageToken, nil
}



