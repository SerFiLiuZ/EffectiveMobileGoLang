package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/SerFiLiuZ/EffectiveMobileGoLang/internal/database"
	"github.com/SerFiLiuZ/EffectiveMobileGoLang/internal/models"
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
	regNum := r.URL.Query().Get("regNum")
	if regNum == "" {
		h.Logger.Errorf("Parameter 'regNum' is required")
		http.Error(w, "Parameter 'regNum' is required", http.StatusBadRequest)
		return
	}

	h.Logger.Debugf("GetCar: regNum: %v", regNum)

	car, err := h.DB.GetCarByRegNum(regNum)
	if err != nil {
		h.Logger.Errorf("Failed to fetch car information: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	var requestBody struct {
		RegNums         []string `json:"regNums"`
		Mark            []string `json:"mark"`
		Model           []string `json:"model"`
		Year            []int    `json:"year"`
		OwnerName       []string `json:"ownerName"`
		OwnerSurname    []string `json:"ownerSurname"`
		OwnerPatronymic []string `json:"ownerPatronymic"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		h.Logger.Errorf("Failed to decode JSON: %v", err)
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}

	h.Logger.Debugf("AddCars: requestBody: %v", requestBody)

	lengths := []int{
		len(requestBody.RegNums),
		len(requestBody.Mark),
		len(requestBody.Model),
		len(requestBody.Year),
		len(requestBody.OwnerName),
		len(requestBody.OwnerSurname),
		len(requestBody.OwnerPatronymic),
	}

	for i := 1; i < len(lengths); i++ {
		if lengths[i] != lengths[0] {
			err := errors.New("data is incomplete")
			h.Logger.Errorf("Failed to add car: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	for i := 0; i < len(requestBody.RegNums); i++ {
		if requestBody.RegNums[i] == "" ||
			requestBody.Mark[i] == "" ||
			requestBody.Model[i] == "" ||
			requestBody.Year[i] == 0 ||
			requestBody.OwnerName[i] == "" ||
			requestBody.OwnerSurname[i] == "" {
			err := errors.New("data is incomplete")
			h.Logger.Errorf("Failed to add car: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		newCar := models.Car{
			RegNum:          requestBody.RegNums[i],
			Mark:            requestBody.Mark[i],
			Model:           requestBody.Model[i],
			Year:            requestBody.Year[i],
			OwnerName:       requestBody.OwnerName[i],
			OwnerSurname:    requestBody.OwnerSurname[i],
			OwnerPatronymic: requestBody.OwnerPatronymic[i],
		}

		err := h.DB.AddCar(newCar)
		if err != nil {
			h.Logger.Errorf("Failed to add car: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		h.Logger.Infof("Car add successfully: %s", requestBody.RegNums[i])
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Cars added successfully"))
}

func (h *CarHandler) UpdateCar(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		RegNum          string `json:"regNum"`
		Mark            string `json:"mark"`
		Model           string `json:"model"`
		Year            int    `json:"year"`
		OwnerName       string `json:"ownerName"`
		OwnerSurname    string `json:"ownerSurname"`
		OwnerPatronymic string `json:"ownerPatronymic"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)

	h.Logger.Debugf("UpdateCar: requestBody: %v", requestBody)

	requestData := map[string]interface{}{
		"mark":            requestBody.Mark,
		"model":           requestBody.Model,
		"year":            requestBody.Year,
		"ownerName":       requestBody.OwnerName,
		"ownerSurname":    requestBody.OwnerSurname,
		"ownerPatronymic": requestBody.OwnerPatronymic,
	}

	for key, value := range requestData {
		if value == nil || value == "" || value == 0 {
			delete(requestData, key)
		}
	}

	if err != nil {
		h.Logger.Errorf("Failed to decode JSON: %v", err)
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}

	if requestBody.RegNum == "" {
		h.Logger.Errorf("Parameter 'regNum' is required")
		http.Error(w, "Parameter 'regNum' is required", http.StatusBadRequest)
		return
	}

	err = h.DB.UpdateCarByRegNum(requestBody.RegNum, requestData)
	if err != nil {
		h.Logger.Errorf("Failed to update car: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.Logger.Infof("Car updated successfully: %s", requestBody.RegNum)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Car updated successfully"))
}

func (h *CarHandler) DeleteCar(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		RegNum string `json:"regNum"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)

	h.Logger.Debugf("UpdateCar: requestBody: %v", requestBody)

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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.Logger.Infof("Car deleted successfully: %s", regNum)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Car deleted successfully"))
}
