package database

import (
	"database/sql"
	"fmt"

	"github.com/SerFiLiuZ/EffectiveMobileGoLang/internal/config"
	"github.com/SerFiLiuZ/EffectiveMobileGoLang/internal/utils"
	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
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
