package store

type Store interface {
	Car() CarRepository
	People() PeopleRepository
	DBController() DBControllerRepository
}
