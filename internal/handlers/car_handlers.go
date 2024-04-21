package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/SerFiLiuZ/EffectiveMobileGoLang/internal/database"
	"github.com/SerFiLiuZ/EffectiveMobileGoLang/internal/utils"
)

type CarHandler struct {
	DB     *database.DB
	Logger *utils.Logger
}

func NewCarHandler(db *database.DB, logger *utils.Logger) *CarHandler {
	return &CarHandler{
		DB:     db,
		Logger: logger,
	}
}

func (h *CarHandler) GetCar(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		RegNum string `json:"regNum"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		h.Logger.Errorf("Failed to decode JSON: %v", err)
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}

	regNum := requestBody.RegNum

	if regNum == "" {
		h.Logger.Errorf("Parameter 'regNum' is required")
		http.Error(w, "Parameter 'regNum' is required", http.StatusBadRequest)
		return
	}

	car, err := h.DB.GetCarByRegNum(regNum)
	if err != nil {
		h.Logger.Errorf("Failed to fetch car information: %v", err)
		http.Error(w, "Failed to fetch car information", http.StatusInternalServerError)
		return
	}

	jsonBytes, err := json.Marshal(car)
	if err != nil {
		h.Logger.Errorf("Failed to marshal JSON: %v", err)
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}

	h.Logger.Infof("Get car info successfully: %s", regNum)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *CarHandler) AddCars(w http.ResponseWriter, r *http.Request) {
	// Реализация метода добавления нового автомобиля или автомоблей
}

func (h *CarHandler) UpdateCar(w http.ResponseWriter, r *http.Request) {
	// Реализация метода обновления данных автомобиля
}

func (h *CarHandler) DeleteCar(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		RegNum string `json:"regNum"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		h.Logger.Errorf("Failed to decode JSON: %v", err)
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}

	regNum := requestBody.RegNum

	if regNum == "" {
		h.Logger.Errorf("Parameter 'regNum' is required")
		http.Error(w, "Parameter 'regNum' is required", http.StatusBadRequest)
		return
	}

	err = h.DB.DeleteCarByRegNum(regNum)
	if err != nil {
		h.Logger.Errorf("Failed to delete car: %v", err)
		http.Error(w, "Failed to delete car", http.StatusInternalServerError)
		return
	}

	h.Logger.Infof("Car deleted successfully: %s", regNum)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Car deleted successfully"))
}
