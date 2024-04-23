package store

import (
	"database/sql"

	"github.com/SerFiLiuZ/EffectiveMobileGoLang/internal/models"
)

type CarRepository interface {
	GetCarByRegNum(regNum string) (*models.Car, error)
	DeleteCarByRegNum(regNum string) error
	UpdateCarByRegNum(regNum, mark, model string, year int, owner models.People) error
	AddCar(newCar models.Car) error
}

type PeopleRepository interface {
	GetOwnerByName(name, surname, patronymic string) (*models.People, error)
	AddOwner(owner models.People) error
}

type DBControllerRepository interface {
	ApplyMigrations(db *sql.DB, migrationsDir, connStr string) error
	RollbackMigrations(db *sql.DB, migrationsDir, connStr string) error
	InsertTestData() error
}
