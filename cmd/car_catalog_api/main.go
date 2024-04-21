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
	defer db.Db.Close()

	err = database.RollbackMigrations(db)
	if err != nil {
		logger.Fatalf("Error rollback migrations: %v", err)
	}

	err = database.ApplyMigrations(db)
	if err != nil {
		logger.Fatalf("Error applying migrations: %v", err)
	}

	err = database.InsertData(db)
	if err != nil {
		logger.Fatalf("Error inserdata: %v", err)
	}

	carHandler := handlers.NewCarHandler(db, logger)

	router := mux.NewRouter()

	// Получение данных об одном автомобиле по идентификатору
	router.HandleFunc("/cars", carHandler.GetCar).Methods("GET")

	// Удаление по идентификатору
	router.HandleFunc("/cars/{id}", carHandler.DeleteCar).Methods("DELETE")

	// Изменение одного или нескольких полей по идентификатору
	router.HandleFunc("/cars/{id}", carHandler.UpdateCar).Methods("PUT")

	// Добавление новых автомобилей
	router.HandleFunc("/cars", carHandler.AddCar).Methods("POST")

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
