package domain

import (
	apperrors "pyg.com/api/app/errors"
)

type GameId int64

type Game struct {
	Id          GameId
	DeveloperId DeveloperId
	DeveloperName DeveloperName
	Name        string
	Description string
	Genre       Genre
	Platform    Platform
}

func CreateGame(accountId AccountId, developer *Developer, name, description string, genre Genre, platform Platform) (*Game, error) {
	if !developer.IsAdmin(accountId) {
		return nil, apperrors.ErrUnauthorized
	}

	game := &Game{
		DeveloperId: developer.Id,
		Name:        name,
		Description: description,
		Genre:       genre,
		Platform:    platform,
	}

	return game, nil
}
