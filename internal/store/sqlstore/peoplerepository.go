package sqlstore

import (
	"database/sql"

	"github.com/SerFiLiuZ/EffectiveMobileGoLang/internal/models"
	"github.com/pkg/errors"
)

type PeopleRepository struct {
	store *Store
}

func (r *PeopleRepository) GetOwnerByName(name, surname, patronymic string) (*models.People, error) {
	query := `
		SELECT name, surname, patronymic
		FROM people
		WHERE name = $1 AND surname = $2 AND patronymic = $3
		LIMIT 1
	`

	var owner models.People
	err := r.store.db.QueryRow(query, name, surname, patronymic).Scan(&owner.Name, &owner.Surname, &owner.Patronymic)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "failed to get owner by name")
	}

	return &owner, nil
}

func (r *PeopleRepository) AddOwner(owner models.People) error {
	query := `
        INSERT INTO people (name, surname, patronymic)
        VALUES ($1, $2, $3)
    `

	_, err := r.store.db.Exec(query, owner.Name, owner.Surname, owner.Patronymic)
	if err != nil {
		return errors.Wrap(err, "failed to add owner")
	}

	return nil
}
