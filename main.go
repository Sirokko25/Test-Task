package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"

	"testwork/handlers"
	"testwork/storage"
	"testwork/structs"

)

func main() {

	db, err := storage.OpenDb()
	if err != nil {
		log.Fatalf("Ошибка доступа к базе данных: %v", err)
	}

	handlers := handlers.Handlers{db}

	r := chi.NewRouter()
	r.Post("/api/task/add", handlers.AddNote())
	r.Get("/api/tasks", handlers.GetNotes())

	log.Printf("Сервер слушает порт %s", structs.Port)

	if err := http.ListenAndServe(":"+structs.Port, r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
