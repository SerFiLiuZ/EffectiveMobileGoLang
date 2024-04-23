package sqlstore

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
)

type DBControllerRepository struct {
	store *Store
}

func (r *DBControllerRepository) ApplyMigrations(db *sql.DB, migrationsDir, connStr string) error {
	// migrationsDir := "C:/Users/serfi/Desktop/Work/GoLang Project/v2/EffectiveMobileGoLang/migrations"

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

	return nil
}

func (r *DBControllerRepository) RollbackMigrations(db *sql.DB, migrationsDir, connStr string) error {
	// migrationsDir := "C:/Users/serfi/Desktop/Work/GoLang Project/v2/EffectiveMobileGoLang/migrations"

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

	return nil
}

func (r *DBControllerRepository) InsertTestData() error {
	_, err := r.store.db.Exec(`INSERT INTO people (name, surname, patronymic) VALUES
						('Иван', 'Иванов', 'Иванович'),
						('Петр', 'Петров', 'Петрович'),
						('Анна', 'Сидорова', 'Ивановна')`)
	if err != nil {
		return err
	}

	_, err = r.store.db.Exec(`INSERT INTO car (regNum, mark, model, year, owner_name, owner_surname) VALUES
						('X123XX150', 'Lada', 'Vesta', 2002, 'Иван', 'Иванов'),
						('Y456YY200', 'Toyota', 'Camry', 2015, 'Петр', 'Петров'),
						('Z789ZZ250', 'BMW', 'X5', 2019, 'Анна', 'Сидорова')`)
	if err != nil {
		return err
	}

	return nil
}
