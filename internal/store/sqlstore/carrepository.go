package sqlstore

import (
	"database/sql"
	"fmt"

	"github.com/SerFiLiuZ/EffectiveMobileGoLang/internal/models"
	"github.com/SerFiLiuZ/EffectiveMobileGoLang/internal/utils"
	"github.com/pkg/errors"
)

type CarRepository struct {
	store *Store
}

func (r *CarRepository) InitCarRepository() bool {
	return r != nil
}

func (r *CarRepository) GetCarByRegNum(regNum string) (*models.Car, error) {
	logger := utils.NewLogger()
	logger.EnableDebug()
	logger.DisableDebug()

	var car models.Car

	query := `
	SELECT regNum, mark, model, year, people.name, people.surname, people.patronymic
		FROM car
		JOIN people ON car.owner_name = people.name AND car.owner_surname = people.surname
		WHERE car.regNum = $1;
	`
	logger.Debugf("r.store.db: %v", r.store.db)
	row := r.store.db.QueryRow(query, regNum)

	err := row.Scan(&car.RegNum, &car.Mark, &car.Model, &car.Year, &car.Owner.Name, &car.Owner.Surname, &car.Owner.Patronymic)
	if err == sql.ErrNoRows {
		return nil, errors.New("car not found")
	} else if err != nil {
		return nil, err
	}

	return &car, nil
}

func (r *CarRepository) DeleteCarByRegNum(regNum string) error {
	var count int
	query := "SELECT COUNT(regNum) FROM car WHERE regNum = $1"
	row := r.store.db.QueryRow(query, regNum)
	err := row.Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New("car not found")
	}

	query = "DELETE FROM car WHERE regNum = $1"
	_, err = r.store.db.Exec(query, regNum)
	if err != nil {
		return err
	}

	return nil
}

func (r *CarRepository) UpdateCarByRegNum(regNum, mark, model string, year int, owner models.People) error {
	var count int
	query := "SELECT COUNT(regNum) FROM car WHERE regNum = $1"
	row := r.store.db.QueryRow(query, regNum)
	err := row.Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New("car not found")
	}

	if mark == "" && model == "" && year == 0 && (owner == models.People{}) {
		return errors.New("no data to update")
	}

	query = "UPDATE car SET "

	if mark != "" {
		query += fmt.Sprintf("mark = '%v',", mark)
	}

	if model != "" {
		query += fmt.Sprintf("model = '%v',", model)
	}

	if year != 0 {
		query += fmt.Sprintf("year = '%v',", year)
	}

	query = query[:len(query)-1]

	query += " WHERE regNum = $1"

	_, err = r.store.db.Exec(query, regNum)
	if err != nil {
		return err
	}

	return nil
}

func (r *CarRepository) AddCar(newCar models.Car) error {
	logger := utils.NewLogger()
	logger.EnableDebug()
	logger.DisableDebug()

	logger.Debugf("Start AddCar")
	logger.Debugf("Form req")

	logger.Debugf("newCar: %v", newCar)

	logger.Debugf("CarRepository: %v", r)

	existingOwner, err := r.store.People().GetOwnerByName(newCar.Owner.Name, newCar.Owner.Surname, newCar.Owner.Patronymic)

	if err != nil {
		return err
	}

	logger.Debugf("Req complite")

	if existingOwner == nil {
		newOwner := models.People{
			Name:       newCar.Owner.Name,
			Surname:    newCar.Owner.Surname,
			Patronymic: newCar.Owner.Patronymic,
		}
		err := r.store.People().AddOwner(newOwner)
		if err != nil {
			return err
		}
	}

	query := `
		INSERT INTO car (regNum, mark, model, year, owner_name, owner_surname)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err = r.store.db.Exec(query, newCar.RegNum, newCar.Mark, newCar.Model, newCar.Year, newCar.Owner.Name, newCar.Owner.Surname)
	if err != nil {
		return err
	}

	return nil
}
