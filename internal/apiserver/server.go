package apiserver

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/SerFiLiuZ/EffectiveMobileGoLang/internal/models"
	"github.com/SerFiLiuZ/EffectiveMobileGoLang/internal/store"
	"github.com/SerFiLiuZ/EffectiveMobileGoLang/internal/utils"
	"github.com/gorilla/mux"
)

type server struct {
	router *mux.Router
	logger *utils.Logger
	store  store.Store
}

func newServer(store store.Store, logger *utils.Logger) *server {
	s := &server{
		router: mux.NewRouter(),
		logger: logger,
		store:  store,
	}

	s.configureRouter()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.HandleFunc("/info", s.getCar()).Methods("GET")
	s.router.HandleFunc("/add", s.addCars()).Methods("POST")
	s.router.HandleFunc("/update", s.updateCar()).Methods("PUT")
	s.router.HandleFunc("/del", s.deleteCar()).Methods("DELETE")
}

func (s *server) getCar() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		regNum := r.URL.Query().Get("regNum")
		if regNum == "" {
			s.logger.Errorf("Parameter 'regNum' is required")
			http.Error(w, "Parameter 'regNum' is required", http.StatusBadRequest)
			return
		}

		s.logger.Debugf("GetCar: regNum: %v", regNum)

		car, err := s.store.Car().GetCarByRegNum(regNum)
		if err != nil {
			s.logger.Errorf("Failed to fetch car information: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jsonBytes, err := json.Marshal(car)
		if err != nil {
			s.logger.Errorf("Failed to marshal JSON: %v", err)
			http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
			return
		}

		s.logger.Infof("Get car info successfully: %s", regNum)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonBytes)
	}
}

func (s *server) addCars() http.HandlerFunc {
	type request struct {
		RegNums []string        `json:"regNums"`
		Mark    []string        `json:"mark"`
		Model   []string        `json:"model"`
		Year    []int           `json:"year"`
		Owner   []models.People `json:"owner"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.logger.Errorf("Failed to decode JSON: %v", err)
			http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
			return
		}

		s.logger.Debugf("AddCars: requestBody: %v", req)

		lengths := []int{
			len(req.RegNums),
			len(req.Mark),
			len(req.Model),
			len(req.Year),
			len(req.Owner),
		}

		for i := 1; i < len(lengths); i++ {
			if lengths[i] != lengths[0] {
				err := errors.New("data is incomplete")
				s.logger.Errorf("Failed to add car: %v", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		for i := 0; i < len(req.RegNums); i++ {
			if req.RegNums[i] == "" ||
				req.Mark[i] == "" ||
				req.Model[i] == "" ||
				req.Year[i] == 0 {
				err := errors.New("data is incomplete")
				s.logger.Errorf("Failed to add car: %v", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		for _, owner := range req.Owner {
			if owner.Name == "" || owner.Surname == "" {
				err := errors.New("data is incomplete")
				s.logger.Errorf("Failed to add car: %v", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		for i := 0; i < len(req.RegNums); i++ {
			newCar := models.Car{
				RegNum: req.RegNums[i],
				Mark:   req.Mark[i],
				Model:  req.Model[i],
				Year:   req.Year[i],
				Owner:  req.Owner[i],
			}

			err := s.store.Car().AddCar(newCar)
			if err != nil {
				s.logger.Errorf("Failed to add car: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			s.logger.Infof("Car add successfully: %s", req.RegNums[i])
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Cars added successfully"))
	}
}

func (s *server) updateCar() http.HandlerFunc {
	type request struct {
		RegNum string        `json:"regNum"`
		Mark   string        `json:"mark"`
		Model  string        `json:"model"`
		Year   int           `json:"year"`
		Owner  models.People `json:"owner"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.logger.Errorf("Failed to decode JSON: %v", err)
			http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
			return
		}

		s.logger.Debugf("UpdateCar: requestBody: %v", req)

		if req.RegNum == "" {
			s.logger.Errorf("Parameter 'regNum' is required")
			http.Error(w, "Parameter 'regNum' is required", http.StatusBadRequest)
			return
		}

		err := s.store.Car().UpdateCarByRegNum(
			req.RegNum,
			req.Mark,
			req.Model,
			req.Year,
			req.Owner,
		)

		if err != nil {
			s.logger.Errorf("Failed to update car: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		s.logger.Infof("Car updated successfully: %s", req.RegNum)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Car updated successfully"))
	}
}

func (s *server) deleteCar() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		regNum := r.URL.Query().Get("regNum")
		if regNum == "" {
			s.logger.Errorf("Parameter 'regNum' is required")
			http.Error(w, "Parameter 'regNum' is required", http.StatusBadRequest)
			return
		}

		s.logger.Debugf("GetCar: regNum: %v", regNum)

		if regNum == "" {
			s.logger.Errorf("Parameter 'regNum' is required")
			http.Error(w, "Parameter 'regNum' is required", http.StatusBadRequest)
			return
		}

		err := s.store.Car().DeleteCarByRegNum(regNum)
		if err != nil {
			s.logger.Errorf("Failed to delete car: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		s.logger.Infof("Car deleted successfully: %s", regNum)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Car deleted successfully"))
	}
}
