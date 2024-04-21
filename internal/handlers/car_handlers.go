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
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	regNum := r.Form.Get("regNum")

	if regNum == "" {
		http.Error(w, "Parameter 'regNum' is required", http.StatusBadRequest)
		return
	}

	car, err := h.DB.GetCarByRegNum(regNum)
	if err != nil {
		http.Error(w, "Failed to fetch car information", http.StatusInternalServerError)
		return
	}

	jsonBytes, err := json.Marshal(car)
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *CarHandler) AddCar(w http.ResponseWriter, r *http.Request) {
	// Реализация метода добавления нового автомобиля
}

func (h *CarHandler) UpdateCar(w http.ResponseWriter, r *http.Request) {
	// Реализация метода обновления данных автомобиля
}

func (h *CarHandler) DeleteCar(w http.ResponseWriter, r *http.Request) {
	// Реализация метода удаления автомобиля
}
