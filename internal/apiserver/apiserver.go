package apiserver

import (
	"database/sql"
	"errors"
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

	logger.Debugf("db: %v", db)

	store := sqlstore.New(db)

	logger.Debugf("store: %v", store)

	srv := newServer(store, logger)

	err = InitRepositoris(store, logger)
	if err != nil {
		return err
	}

	logger.Infof("All repositoris init")

	logger.Debugf("srv: %v", srv)

	logger.Debugf("srv.store: %v", srv.store)

	err = store.DBController().RollbackMigrations(db, config.DBMigrationsdir, config.DatabaseURL)
	if err != nil {
		return err
	}

	logger.Infof("Rollback migrations complite")

	err = store.DBController().ApplyMigrations(db, config.DBMigrationsdir, config.DatabaseURL)
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

	logger.Debugf("srv.store: %v", srv.store)

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

func InitRepositoris(store *sqlstore.Store, logger *utils.Logger) error {
	if store.People().InitPeopleRepository() {
		logger.Debugf("store.People() - init")
	} else {
		return errors.New("PeopleRepository dont init")
	}
	if store.Car().InitCarRepository() {
		logger.Debugf("store.Car() - init")
	} else {
		return errors.New("CarRepository dont init")
	}
	if store.DBController().InitDBControllerRepository() {
		logger.Debugf("store.DBController() - init")
	} else {
		return errors.New("DBControllerRepository dont init")
	}
	return nil
}
