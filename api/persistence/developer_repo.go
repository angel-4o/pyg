package persistence

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"pyg.com/api/app/data"
	apperrors "pyg.com/api/app/errors"
	"pyg.com/api/domain"
)

type developerRepo struct {
	db *sql.DB
}

func MakeDeveloperRepo(db *sql.DB) *developerRepo {
	return &developerRepo{db}
}

func (repo *developerRepo) GetDevelopers(pageToken data.PageToken, pageSize int) ([]domain.Developer, data.PageToken, error) {
	rows, _ := repo.db.Query(
		`select id, name, created_by, members
		 from developers
		 where id > $1
		 limit $2`,
		pageToken,
		pageSize,
	)

	developers := make([]domain.Developer, 0)

	for rows.Next() {
		var developer domain.Developer
		var membersJson []byte

		err := rows.Scan(
			&developer.Id,
			&developer.Name,
			&developer.CreatedBy,
			&membersJson,
		)

		if err != nil {
			return nil, pageToken, fmt.Errorf("%w: %w", apperrors.ErrDeveloperNotFound, err)
		}

		err = json.Unmarshal(membersJson, &developer.Members)
		if err != nil {
			return nil, pageToken, fmt.Errorf("%w: %w", apperrors.ErrCannotUnmarshalMembers, err)
		}
		// developer.Name
		developers = append(developers, developer)
	}

	nextPageToken := pageToken
	if count := len(developers); count > 0 {
		nextPageToken = data.PageToken(developers[count-1].Id)
	}

	return developers, nextPageToken, nil
}

func (repo *developerRepo) Create(developer *domain.Developer) (domain.DeveloperId, error) {
	membersJson, err := json.Marshal(developer.Members)
	if err != nil {
		return domain.DeveloperId(0), fmt.Errorf("%w: %w", apperrors.ErrCannotMarshalMembers, err)
	}

	err = repo.db.QueryRow(
		`insert into developers (name, created_by, members)
		 values ($1, $2, $3)
		 returning id`,
		developer.Name,
		developer.CreatedBy,
		membersJson,
	).Scan(
		&developer.Id,
	)

	if err != nil {
		return domain.DeveloperId(0), fmt.Errorf("%w: %w", apperrors.ErrCreateDeveloperFailed, err)
	}

	return developer.Id, nil
}

func (repo *developerRepo) Update(developer domain.Developer) error {
	if !developer.IsValid() {
		return apperrors.ErrDeveloperValidationFailed
	}

	membersJson, err := json.Marshal(developer.Members)
	if err != nil {
		return fmt.Errorf("%w: %w", apperrors.ErrCannotMarshalMembers, err)
	}

	_, err = repo.db.Exec(
		`update developers
		 set (name, created_by, members) = ($2, $3, $4)
		 where id = $1`,
		developer.Id,
		developer.Name,
		developer.CreatedBy,
		membersJson,
	)

	if err != nil {
		return fmt.Errorf("%w: %w", apperrors.ErrUpdateDeveloperFailed, err)
	}

	return nil
}

func (repo *developerRepo) FindById(developerId domain.DeveloperId) (*domain.Developer, error) {
	var developer domain.Developer
	var membersJson []byte

	err := repo.db.QueryRow(
		`select id, name, created_by, members
		 from developers
		 where id = $1`,
		developerId,
	).Scan(
		&developer.Id,
		&developer.Name,
		&developer.CreatedBy,
		&membersJson,
	)

	if err != nil {
		return nil, fmt.Errorf("%w: %w", apperrors.ErrDeveloperNotFound, err)
	}

	err = json.Unmarshal(membersJson, &developer.Members)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", apperrors.ErrCannotUnmarshalMembers, err)
	}

	return &developer, nil
}

func (repo *developerRepo) FindByName(developerName string) (*domain.Developer, error) {
	var developer domain.Developer
	var membersJson []byte

	err := repo.db.QueryRow(
		`SELECT id, name, created_by, members
		 FROM developers
		 WHERE name = $1`,
		developerName,
	).Scan(
		&developer.Id,
		&developer.Name,
		&developer.CreatedBy,
		&membersJson,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%w: %s", apperrors.ErrDeveloperNotFound, developerName)
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	err = json.Unmarshal(membersJson, &developer.Members)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", apperrors.ErrCannotUnmarshalMembers, err)
	}

	return &developer, nil
}
