package database

import (
	"database/sql"
	"fmt"

	"github.com/SerFiLiuZ/EffectiveMobileGoLang/internal/config"
	"github.com/SerFiLiuZ/EffectiveMobileGoLang/internal/models"
	"github.com/SerFiLiuZ/EffectiveMobileGoLang/internal/utils"
	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/pkg/errors"
)

type DB struct {
	Db     *sql.DB
	Logger *utils.Logger
}

func Connect(logger *utils.Logger) (*DB, error) {
	conf := config.GetConfig()
	connStr := getConnStrForConnectDB(conf)

	dbConn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = dbConn.Ping()
	if err != nil {
		return nil, err
	}

	logger.Infof("Connected to database")

	return &DB{Db: dbConn, Logger: logger}, nil
}

func ApplyMigrations(db *DB) error {
	migrationsDir := "C:/Users/serfi/Desktop/Work/GoLang Project/v2/EffectiveMobileGoLang/migrations"

	conf := config.GetConfig()
	connStr := getConnStrForMigrations(conf)

	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsDir),
		connStr,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	db.Logger.Infof("Migrations applied successfully")
	return nil
}

func RollbackMigrations(db *DB) error {
	migrationsDir := "C:/Users/serfi/Desktop/Work/GoLang Project/v2/EffectiveMobileGoLang/migrations"

	conf := config.GetConfig()
	connStr := getConnStrForMigrations(conf)

	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsDir),
		connStr,
	)
	if err != nil {
		return err
	}

	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	db.Logger.Infof("Migrations rolled back successfully")
	return nil
}

func InsertData(db *DB) error {
	_, err := db.Db.Exec(`INSERT INTO people (name, surname, patronymic) VALUES
						('Иван', 'Иванов', 'Иванович'),
						('Петр', 'Петров', 'Петрович'),
						('Анна', 'Сидорова', 'Ивановна')`)
	if err != nil {
		return err
	}

	db.Logger.Debugf("People data inserted successfully")

	_, err = db.Db.Exec(`INSERT INTO car (regNum, mark, model, year, owner_name, owner_surname) VALUES
						('X123XX150', 'Lada', 'Vesta', 2002, 'Иван', 'Иванов'),
						('Y456YY200', 'Toyota', 'Camry', 2015, 'Петр', 'Петров'),
						('Z789ZZ250', 'BMW', 'X5', 2019, 'Анна', 'Сидорова')`)
	if err != nil {
		return err
	}

	db.Logger.Debugf("Car data inserted successfully")

	return nil
}

// Вероятно кривая реализация
func getConnStrForConnectDB(conf *config.Config) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		conf.DBHost, conf.DBPort, conf.DBUser, conf.DBPassword, conf.DBName)
}

// Вероятно кривая реализация
func getConnStrForMigrations(conf *config.Config) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		conf.DBUser, conf.DBPassword, conf.DBHost, conf.DBPort, conf.DBName)
}

func (db *DB) GetCarByRegNum(regNum string) (*models.Car, error) {
	var car models.Car

	query := `
	SELECT regNum, mark, model, year, people.name, people.surname, people.patronymic
		FROM car
		JOIN people ON car.owner_name = people.name AND car.owner_surname = people.surname
		WHERE car.regNum = $1;
	`

	row := db.Db.QueryRow(query, regNum)

	err := row.Scan(&car.RegNum, &car.Mark, &car.Model, &car.Year, &car.Owner.Name, &car.Owner.Surname, &car.Owner.Patronymic)
	if err == sql.ErrNoRows {
		return nil, errors.New("car not found")
	} else if err != nil {
		return nil, err
	}

	return &car, nil
}

func (db *DB) DeleteCarByRegNum(regNum string) error {
	var count int
	query := "SELECT COUNT(regNum) FROM car WHERE regNum = $1"
	row := db.Db.QueryRow(query, regNum)
	err := row.Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New("car not found")
	}

	query = "DELETE FROM car WHERE regNum = $1"
	_, err = db.Db.Exec(query, regNum)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) UpdateCarByRegNum(regNum, mark, model string, year int, owner models.People) error {
	var count int
	query := "SELECT COUNT(regNum) FROM car WHERE regNum = $1"
	row := db.Db.QueryRow(query, regNum)
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

	_, err = db.Db.Exec(query, regNum)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) AddCar(newCar models.Car) error {
	existingOwner, err := db.GetOwnerByName(newCar.Owner.Name, newCar.Owner.Surname, newCar.Owner.Patronymic)
	if err != nil {
		return err
	}

	if existingOwner == nil {
		newOwner := models.People{
			Name:       newCar.Owner.Name,
			Surname:    newCar.Owner.Surname,
			Patronymic: newCar.Owner.Patronymic,
		}

		err := db.AddOwner(newOwner)
		if err != nil {
			return err
		}
	}

	query := `
		INSERT INTO car (regNum, mark, model, year, owner_name, owner_surname, owner_patronymic)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err = db.Db.Exec(query, newCar.RegNum, newCar.Mark, newCar.Model, newCar.Year, newCar.Owner.Name, newCar.Owner.Surname, newCar.Owner.Patronymic)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) GetOwnerByName(name, surname, patronymic string) (*models.People, error) {
	query := `
		SELECT name, surname, patronymic
		FROM people
		WHERE name = $1 AND surname = $2 AND patronymic = $3
		LIMIT 1
	`

	var owner models.People
	err := db.Db.QueryRow(query, name, surname, patronymic).Scan(&owner.Name, &owner.Surname, &owner.Patronymic)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "failed to get owner by name")
	}

	return &owner, nil
}

func (db *DB) AddOwner(owner models.People) error {
	query := `
        INSERT INTO people (name, surname, patronymic)
        VALUES ($1, $2, $3)
    `

	_, err := db.Db.Exec(query, owner.Name, owner.Surname, owner.Patronymic)
	if err != nil {
		return errors.Wrap(err, "failed to add owner")
	}

	return nil
}
