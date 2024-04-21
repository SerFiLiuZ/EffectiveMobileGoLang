package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/SerFiLiuZ/EffectiveMobileGoLang/internal/config"
	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
)

func Connect() (*sql.DB, error) {
	conf := config.GetConfig()
	connStr := getConnStrForConnectDB(conf)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	log.Println("Connected to database")

	return db, nil
}

func ApplyMigrations(db *sql.DB) error {
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

func RollbackMigrations(db *sql.DB) error {
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

	// Откатываем миграции
	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	log.Println("Migrations rolled back successfully")
	return nil
}

func InsertData(db *sql.DB) error {
	_, err := db.Exec(`INSERT INTO people (name, surname, patronymic) VALUES
						('Иван', 'Иванов', 'Иванович'),
						('Петр', 'Петров', 'Петрович'),
						('Анна', 'Сидорова', 'Ивановна')`)
	if err != nil {
		return err
	}

	log.Println("People data inserted successfully")

	_, err = db.Exec(`INSERT INTO car (regNum, mark, model, year, owner_name, owner_surname, owner_patronymic) VALUES
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
