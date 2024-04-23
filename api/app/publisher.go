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

func GetPublishers(pageToken data.PageToken, pageSize int, db *sql.DB) ([]domain.Publisher, data.PageToken, error) {
	publisherRepo := persistence.MakePublisherRepo(db)
	return publisherRepo.GetPublishers(pageToken, pageSize)
}

func GetPublisher(publisherId domain.PublisherId, db *sql.DB) (*domain.Publisher, error) {
	publisherRepo := persistence.MakePublisherRepo(db)
	return publisherRepo.FindById(publisherId)
}

func CreatePublisher(sessionToken data.SessionToken, name string, db *sql.DB) (domain.PublisherId, error) {
	sessionDataSource := persistence.MakeSessionDataSource(db)
	session, err := sessionDataSource.FindSessionByToken(sessionToken)
	if err != nil {
		return domain.PublisherId(0), err
	}

	publisherRepo := persistence.MakePublisherRepo(db)
	publisherId, err := publisherRepo.Create(domain.MakePublisher(name, session.AccountId))
	if err != nil {
		return domain.PublisherId(0), err
	}

	return publisherId, nil
}

func UpdatePublisher(sessionToken data.SessionToken, publisherId domain.PublisherId, name *string, members domain.MemberList, db *sql.DB) error {
	publisherRepo := persistence.MakePublisherRepo(db)
	publisher, err := publisherRepo.FindById(publisherId)
	if err != nil {
		return fmt.Errorf("%w: %w", apperrors.ErrPublisherNotFound, err)
	}

	if !auth.HasAdminRole(sessionToken, publisher.Members, db) {
		return fmt.Errorf("%w: %w", apperrors.ErrUpdateUnauthorized, err)
	}

	if name != nil {
		publisher.Name = *name
	}
	if members != nil {
		publisher.Members = members
	}

	err = publisherRepo.Update(*publisher)
	if err != nil {
		return fmt.Errorf("%w: %w", apperrors.ErrUpdatePublisherFailed, err)
	}

	return nil
}
