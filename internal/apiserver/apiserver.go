package apiserver

import (
	"database/sql"
	"net/http"

	"github.com/SerFiLiuZ/EffectiveMobileGoLang/internal/store/sqlstore"
	"github.com/SerFiLiuZ/EffectiveMobileGoLang/internal/utils"

	"github.com/gorilla/handlers"

	_ "github.com/golang-migrate/migrate/database/postgres"
)

func Start(config *Config, logger *utils.Logger) error {
	db, err := Connect(config.DatabaseURL)
	if err != nil {
		return err
	}
	defer db.Close()

	logger.Infof("Connected to database")

	store := sqlstore.New(db)
	srv := newServer(store, logger)

	migrationsDir := "C:/Users/serfi/Desktop/Work/GoLang Project/v2/EffectiveMobileGoLang/migrations"

	err = store.DBController().RollbackMigrations(db, migrationsDir, config.DatabaseURL)
	if err != nil {
		return err
	}

	logger.Infof("Rollback migrations complite")

	err = store.DBController().ApplyMigrations(db, migrationsDir, config.DatabaseURL)
	if err != nil {
		return err
	}

	logger.Infof("Apply migrations complite")

	err = store.DBController().InsertTestData()
	if err != nil {
		return err
	}

	logger.Infof("Insert test data complite")

	logger.Infof("Server started on port %s", config.Port)

	return http.ListenAndServe(config.Port,
		handlers.CORS(
			handlers.AllowedOrigins([]string{"*"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
			handlers.AllowedHeaders([]string{"Content-Type", "Authorization", "X-Requested-With", "Cookie"}),
			handlers.ExposedHeaders([]string{"Set-Cookie"}),
			handlers.AllowCredentials(),
		)(srv))
}

func Connect(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
