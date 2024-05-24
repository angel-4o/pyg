package app

import (
	"database/sql"

	"pyg.com/api/app/data"
	apperrors "pyg.com/api/app/errors"
	"pyg.com/api/domain"
	"pyg.com/api/persistence"
)

func CreateGame(sessionToken data.SessionToken, developerId domain.DeveloperId, name, description string, genre domain.Genre, platform domain.Platform, db *sql.DB) (domain.GameId, error) {
	sessionDataSource := persistence.MakeSessionDataSource(db)
	session, err := sessionDataSource.FindSessionByToken(sessionToken)
	if err != nil {
		return domain.GameId(0), apperrors.ErrUnauthenticated
	}

	developerRepo := persistence.MakeDeveloperRepo(db)

	//persistence.MakeAccountRepo()
	//account_repo
	//developerRepo.find§
	developer, err := developerRepo.FindById(developerId)
	if err != nil {
		return domain.GameId(0), err
	}

	game, err := domain.CreateGame(session.AccountId, developer, name, description, genre, platform)
	if err != nil {
		return domain.GameId(0), err
	}

	gameRepo := persistence.MakeGameRepo(db)
	gameId, err := gameRepo.Create(game)
	if err != nil {
		return domain.GameId(0), err
	}

	return domain.GameId(gameId), nil
}

func CreateGame2(sessionToken data.SessionToken, developerName string, name, description string, genre domain.Genre, platform domain.Platform, db *sql.DB) (domain.GameId, error) {
	sessionDataSource := persistence.MakeSessionDataSource(db)
	session, err := sessionDataSource.FindSessionByToken(sessionToken)
	if err != nil {
		return domain.GameId(0), apperrors.ErrUnauthenticated
	}

	developerRepo := persistence.MakeDeveloperRepo(db)

	//persistence.MakeAccountRepo()
	//account_repo
	//developerRepo.find§
	developer, err := developerRepo.FindByName(developerName)
	if err != nil {
		return domain.GameId(0), err
	}

	game, err := domain.CreateGame(session.AccountId, developer, name, description, genre, platform)
	if err != nil {
		return domain.GameId(0), err
	}

	gameRepo := persistence.MakeGameRepo(db)
	gameId, err := gameRepo.Create(game)
	if err != nil {
		return domain.GameId(0), err
	}

	return domain.GameId(gameId), nil
}

func GetGame(gameId domain.GameId, db *sql.DB) (*domain.Game, error) {
	gameRepo := persistence.MakeGameRepo(db)
	return gameRepo.Get(gameId)
}

func GetGames(pageToken data.PageToken, genre *domain.Genre, platform *domain.Platform, db *sql.DB) ([]domain.Game, data.PageToken, error) {
	gameRepo := persistence.MakeGameRepo(db)
	return gameRepo.GetGames(pageToken, genre, platform)
}
