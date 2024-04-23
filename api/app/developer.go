package app

import (
	"database/sql"
	"fmt"

	"pyg.com/api/app/auth"
	"pyg.com/api/app/data"
	apperrors "pyg.com/api/app/errors"
	"pyg.com/api/domain"
	"pyg.com/api/persistence"
)

func GetDevelopers(pageToken data.PageToken, pageSize int, db *sql.DB) ([]domain.Developer, data.PageToken, error) {
	developerRepo := persistence.MakeDeveloperRepo(db)
	return developerRepo.GetDevelopers(pageToken, pageSize)
}

func GetDeveloper(developerId domain.DeveloperId, db *sql.DB) (*domain.Developer, error) {
	developerRepo := persistence.MakeDeveloperRepo(db)
	return developerRepo.FindById(developerId)
}

func CreateDeveloper(sessionToken data.SessionToken, name string, db *sql.DB) (domain.DeveloperId, error) {
	sessionDataSource := persistence.MakeSessionDataSource(db)
	session, err := sessionDataSource.FindSessionByToken(sessionToken)
	if err != nil {
		return domain.DeveloperId(0), err
	}

	developerRepo := persistence.MakeDeveloperRepo(db)
	developerId, err := developerRepo.Create(domain.MakeDeveloper(name, session.AccountId))
	if err != nil {
		return domain.DeveloperId(0), err
	}

	return developerId, nil
}

func UpdateDeveloper(sessionToken data.SessionToken, developerId domain.DeveloperId, name *string, members domain.MemberList, db *sql.DB) error {
	developerRepo := persistence.MakeDeveloperRepo(db)
	developer, err := developerRepo.FindById(developerId)
	if err != nil {
		return fmt.Errorf("%w: %w", apperrors.ErrDeveloperNotFound, err)
	}

	if !auth.HasAdminRole(sessionToken, developer.Members, db) {
		return fmt.Errorf("%w: %w", apperrors.ErrUpdateUnauthorized, err)
	}

	if name != nil {
		developer.Name = *name
	}
	if members != nil {
		developer.Members = members
	}

	err = developerRepo.Update(*developer)
	if err != nil {
		return fmt.Errorf("%w: %w", apperrors.ErrUpdateDeveloperFailed, err)
	}

	return nil
}
