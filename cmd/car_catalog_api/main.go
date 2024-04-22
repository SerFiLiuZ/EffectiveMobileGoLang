package main

import (
	"net/http"
	"os"

	"github.com/SerFiLiuZ/EffectiveMobileGoLang/internal/config"
	"github.com/SerFiLiuZ/EffectiveMobileGoLang/internal/database"
	"github.com/SerFiLiuZ/EffectiveMobileGoLang/internal/handlers"
	"github.com/SerFiLiuZ/EffectiveMobileGoLang/internal/utils"
	"github.com/gorilla/mux"
)

func main() {
	logger := utils.NewLogger()
	logger.EnableDebug()

	err := config.LoadEnv(logger)
	if err != nil {
		logger.Fatal("Error loading .env file: %v", err)
	}

	db, err := database.Connect(logger)
	if err != nil {
		logger.Fatal("Error connecting to database: %v", err)
	}
	defer db.Db.Close()

	err = database.RollbackMigrations(db)
	if err != nil {
		logger.Fatal("Error rollback migrations: %v", err)
	}

	err = database.ApplyMigrations(db)
	if err != nil {
		logger.Fatal("Error applying migrations: %v", err)
	}

	err = database.InsertData(db)
	if err != nil {
		logger.Fatal("Error inserdata: %v", err)
	}

	carHandler := handlers.NewCarHandler(db, logger)

	router := mux.NewRouter()

	// Получение данных об одном автомобиле по идентификатору
	router.HandleFunc("/info", carHandler.GetCar).Methods("GET")

	// Удаление по идентификатору
	router.HandleFunc("/del", carHandler.DeleteCar).Methods("DELETE")

	// Изменение одного или нескольких полей по идентификатору
	router.HandleFunc("/update", carHandler.UpdateCar).Methods("PUT")

	// Добавление новых автомобилей
	router.HandleFunc("/add", carHandler.AddCars).Methods("POST")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Infof("Server is starting on port %s", port)

	err = http.ListenAndServe(":"+port, router)
	if err != nil {
		logger.Fatal("Error starting server: %v", err)
	}
}
