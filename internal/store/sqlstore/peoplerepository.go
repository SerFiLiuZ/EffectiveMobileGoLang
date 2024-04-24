package sqlstore

import (
	"database/sql"

	"github.com/SerFiLiuZ/EffectiveMobileGoLang/internal/models"
	"github.com/SerFiLiuZ/EffectiveMobileGoLang/internal/utils"
	"github.com/pkg/errors"
)

type PeopleRepository struct {
	store *Store
}

func (r *PeopleRepository) InitPeopleRepository() bool {
	return r != nil
}

func (r *PeopleRepository) GetOwnerByName(name, surname, patronymic string) (*models.People, error) {
	logger := utils.NewLogger()
	logger.EnableDebug()
	logger.DisableDebug()

	logger.Debugf("Start GetOwnerByName")

	logger.Debugf("r: %v", r)
	logger.Debugf("r.store: %v", r.store)
	logger.Debugf("r.store.carRepository: %v", r.store.carRepository)
	logger.Debugf("r.store.peopleRepository: %v", r.store.peopleRepository)
	logger.Debugf("r.store.dbcontrollerRepository: %v", r.store.dbcontrollerRepository)
	logger.Debugf("r.store.db: %v", r.store.db)

	query := `
		SELECT name, surname, patronymic
		FROM people
		WHERE name = $1 AND surname = $2 AND patronymic = $3
	`

	logger.Debugf("Start req")
	logger.Debugf("Params: name - %v : surname - %v : patronymic - %v", name, surname, patronymic)

	var owner models.People
	logger.Debugf("owner: %v", owner)

	if r.store.db == nil {
		logger.Debugf("this")
	} else {
		logger.Debugf("thisssssssssssss")
	}

	logger.Debugf("Start QueryRow")

	err := r.store.db.QueryRow(query, name, surname, patronymic).Scan(&owner.Name, &owner.Surname, &owner.Patronymic)
	logger.Debugf("Err: %v", err)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Debugf("Err: ErrNoRows")
			return nil, nil
		}
		return nil, errors.Wrap(err, "failed to get owner by name")
	}

	logger.Debugf("Req complite")

	return &owner, nil
}

func (r *PeopleRepository) AddOwner(owner models.People) error {
	logger := utils.NewLogger()
	logger.EnableDebug()
	logger.DisableDebug()
	logger.Debugf("AddOwner : r.store.db: %v", r.store.db)
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
