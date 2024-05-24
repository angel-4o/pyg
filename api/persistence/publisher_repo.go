package persistence

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"pyg.com/api/app/data"
	apperrors "pyg.com/api/app/errors"
	"pyg.com/api/domain"
)

type publisherRepo struct {
	db *sql.DB
}

func MakePublisherRepo(db *sql.DB) *publisherRepo {
	return &publisherRepo{db}
}

func (repo *publisherRepo) GetPublishers(pageToken data.PageToken, pageSize int) ([]domain.Publisher, data.PageToken, error) {
	rows, _ := repo.db.Query(
		`select id, name, created_by, members
		 from publishers
		 where id > $1
		 limit $2`,
		pageToken,
		pageSize,
	)

	publishers := make([]domain.Publisher, 0)

	for rows.Next() {
		var publisher domain.Publisher
		var membersJson []byte

		err := rows.Scan(
			&publisher.Id,
			&publisher.Name,
			&publisher.CreatedBy,
			&membersJson,
		)

		if err != nil {
			return nil, pageToken, fmt.Errorf("%w: %w", apperrors.ErrPublisherNotFound, err)
		}

		err = json.Unmarshal(membersJson, &publisher.Members)
		if err != nil {
			return nil, pageToken, fmt.Errorf("%w: %w", apperrors.ErrCannotUnmarshalMembers, err)
		}

		publishers = append(publishers, publisher)
	}

	nextPageToken := pageToken
	if count := len(publishers); count > 0 {
		nextPageToken = data.PageToken(publishers[count-1].Id)
	}

	return publishers, nextPageToken, nil
}

func (repo *publisherRepo) Create(publisher *domain.Publisher) (domain.PublisherId, error) {
	membersJson, err := json.Marshal(publisher.Members)
	if err != nil {
		return domain.PublisherId(0), fmt.Errorf("%w: %w", apperrors.ErrCannotMarshalMembers, err)
	}

	err = repo.db.QueryRow(
		`insert into publishers (name, created_by, members)
		 values ($1, $2, $3)
		 returning id`,
		publisher.Name,
		publisher.CreatedBy,
		membersJson,
	).Scan(
		&publisher.Id,
	)

	if err != nil {
		return domain.PublisherId(0), fmt.Errorf("%w: %w", apperrors.ErrCreatePublisherFailed, err)
	}

	return publisher.Id, nil
}

func (repo *publisherRepo) Update(publisher domain.Publisher) error {
	if !publisher.IsValid() {
		return apperrors.ErrPublisherValidationFailed
	}

	membersJson, err := json.Marshal(publisher.Members)
	if err != nil {
		return fmt.Errorf("%w: %w", apperrors.ErrCannotMarshalMembers, err)
	}

	_, err = repo.db.Exec(
		`update publishers
		 set (name, created_by, members) = ($2, $3, $4)
		 where id = $1`,
		publisher.Id,
		publisher.Name,
		publisher.CreatedBy,
		membersJson,
	)

	if err != nil {
		return fmt.Errorf("%w: %w", apperrors.ErrUpdatePublisherFailed, err)
	}

	return nil
}

func (repo *publisherRepo) FindById(publisherId domain.PublisherId) (*domain.Publisher, error) {
	var publisher domain.Publisher
	var membersJson []byte

	err := repo.db.QueryRow(
		`select id, name, created_by, members
		 from publishers
		 where id = $1`,
		publisherId,
	).Scan(
		&publisher.Id,
		&publisher.Name,
		&publisher.CreatedBy,
		&membersJson,
	)

	if err != nil {
		return nil, fmt.Errorf("%w: %w", apperrors.ErrPublisherNotFound, err)
	}

	err = json.Unmarshal(membersJson, &publisher.Members)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", apperrors.ErrCannotUnmarshalMembers, err)
	}

	return &publisher, nil
}

func (repo *publisherRepo) FindByName(publisherName string) (*domain.Publisher, error) {
	var publisher domain.Publisher
	var membersJson []byte

	err := repo.db.QueryRow(
		`select id, name, created_by, members
		 from publishers
		 where name = $1`,
		 publisherName,
	).Scan(
		&publisher.Id,
		&publisher.Name,
		&publisher.CreatedBy,
		&membersJson,
	)

	if err != nil {
		return nil, fmt.Errorf("%w: %w", apperrors.ErrPublisherNotFound, err)
	}

	err = json.Unmarshal(membersJson, &publisher.Members)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", apperrors.ErrCannotUnmarshalMembers, err)
	}

	return &publisher, nil
}