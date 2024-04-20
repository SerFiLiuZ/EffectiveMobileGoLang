package main

import (
	"log"
	"net/http"
	"os"

	"github.com/SerFiLiuZ/EffectiveMobileGoLang/internal/config"
	"github.com/SerFiLiuZ/EffectiveMobileGoLang/internal/database"
	"github.com/SerFiLiuZ/EffectiveMobileGoLang/internal/handlers"
	"github.com/SerFiLiuZ/EffectiveMobileGoLang/internal/utils"
	"github.com/gorilla/mux"
)

func main() {
	err := config.LoadEnv()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	logger := utils.NewLogger()

	db, err := database.Connect()
	if err != nil {
		logger.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	carHandler := handlers.NewCarHandler(db, logger)

	router := mux.NewRouter()

	// Получение данных с фильтрацией по всем полям и пагинацией
	router.HandleFunc("/cars", carHandler.GetCars).Methods("GET")

	// Удаление по идентификатору
	router.HandleFunc("/cars/{id}", carHandler.DeleteCar).Methods("DELETE")

	// Изменение одного или нескольких полей по идентификатору
	router.HandleFunc("/cars/{id}", carHandler.UpdateCar).Methods("PUT")

	// Добавление новых автомобилей
	router.HandleFunc("/cars", carHandler.AddCar).Methods("POST")

	// Получение данных об одном автомобиле по идентификатору
	router.HandleFunc("/cars/{id}", carHandler.GetCarByID).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Infof("Server is starting on port %s", port)

	err = http.ListenAndServe(":"+port, router)
	if err != nil {
		logger.Fatalf("Error starting server: %v", err)
	}
}
