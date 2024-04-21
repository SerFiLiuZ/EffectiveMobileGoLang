package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/SerFiLiuZ/EffectiveMobileGoLang/internal/config"
	"github.com/SerFiLiuZ/EffectiveMobileGoLang/internal/models"
	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
)

type DB struct {
	Db *sql.DB
}

func Connect() (*DB, error) {
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

	log.Println("Connected to database")

	return &DB{Db: dbConn}, nil
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

	log.Println("Migrations applied successfully")
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

	log.Println("Migrations rolled back successfully")
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

	log.Println("People data inserted successfully")

	_, err = db.Db.Exec(`INSERT INTO car (regNum, mark, model, year, owner_name, owner_surname, owner_patronymic) VALUES
						('X123XX150', 'Lada', 'Vesta', 2002, 'Иван', 'Иванов', 'Иванович'),
						('Y456YY200', 'Toyota', 'Camry', 2015, 'Петр', 'Петров', 'Петрович'),
						('Z789ZZ250', 'BMW', 'X5', 2019, 'Анна', 'Сидорова', 'Ивановна')`)
	if err != nil {
		return err
	}

	log.Println("Car data inserted successfully")

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

	row := db.Db.QueryRow("SELECT regNum, mark, model, year, owner_name, owner_surname, owner_patronymic FROM car WHERE regNum = $1", regNum)
	err := row.Scan(&car.RegNum, &car.Mark, &car.Model, &car.Year, &car.OwnerName, &car.OwnerSurname, &car.OwnerPatronymic)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &car, nil
}

func (db *DB) DeleteCarByRegNum(regNum string) error {
	query := "DELETE FROM car WHERE regNum = $1"

	_, err := db.Db.Exec(query, regNum)
	if err != nil {
		return err
	}

	return nil
}
