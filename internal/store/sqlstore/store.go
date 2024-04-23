package sqlstore

import (
	"database/sql"

	"github.com/SerFiLiuZ/EffectiveMobileGoLang/internal/store"
)

type Store struct {
	db                     *sql.DB
	carRepository          *CarRepository
	peopleRepository       *PeopleRepository
	dbcontrollerRepository *DBControllerRepository
}

func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) Car() store.CarRepository {
	if s.carRepository != nil {
		return s.carRepository
	}

	s.carRepository = &CarRepository{
		store: s,
	}

	return s.carRepository
}

func (s *Store) People() store.PeopleRepository {
	if s.carRepository != nil {
		return s.peopleRepository
	}

	s.peopleRepository = &PeopleRepository{
		store: s,
	}

	return s.peopleRepository
}

func (s *Store) DBController() store.DBControllerRepository {
	if s.dbcontrollerRepository != nil {
		return s.dbcontrollerRepository
	}

	s.dbcontrollerRepository = &DBControllerRepository{
		store: s,
	}

	return s.dbcontrollerRepository
}
