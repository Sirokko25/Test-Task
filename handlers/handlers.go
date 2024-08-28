package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"testwork/storage"
	"testwork/structs"

)

type Handlers struct {
	TaskStorage storage.DB
}

func (h *Handlers) AddNote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		var request structs.Request
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			log.Print("Ошибка десериализации JSON")
			response := map[string]interface{}{
				"error": err,
			}
			json.NewEncoder(w).Encode(response)
			return
		}
		if request.Note == "" {
			log.Print("Не указаны данные")
			response := map[string]interface{}{
				"error": "Не указаны данные",
			}
			json.NewEncoder(w).Encode(response)
			return
		}

		login, password, ok := r.BasicAuth()
		if !ok {
			log.Print("Ошибка авторизации")
			response := map[string]interface{}{
				"error": "Отсутствует обязательный хэдер авторизации",
			}
			json.NewEncoder(w).Encode(response)
			return
		}

		idUser, errstr := h.TaskStorage.Indetification(login, password)
		if errstr != "" {
			log.Print("Ошибка авторизации")
			response := map[string]interface{}{
				"error": errstr,
			}
			json.NewEncoder(w).Encode(response)
			return
		}
		//тут интеграция яндекс Спеллер
		var Speller structs.YandexSpeller
		refreshNote, err := Speller.CheckYandexSpeller(request.Note)
		if err != nil {
			log.Print("Ошибка при обращении к YandexSpeller")
			response := map[string]interface{}{
				"error": err,
			}
			json.NewEncoder(w).Encode(response)
			return
		}

		err = h.TaskStorage.AddNoteToDatabase(idUser, refreshNote)
		if err != nil {
			log.Print("Ошибка десериализации JSON")
			response := map[string]interface{}{
				"error": err,
			}
			json.NewEncoder(w).Encode(response)
			return
		}

		log.Printf("Добавление заметки для пользователя %s прошло успешно", login)
		response := map[string]interface{}{
			"ok": "Заметка была добавлена",
		}
		json.NewEncoder(w).Encode(response)
	}
}

func (h *Handlers) GetNotes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		var response interface{}
		login, password, ok := r.BasicAuth()
		if !ok {
			log.Print("Ошибка авторизации")
			response := map[string]interface{}{
				"error": "Отсутствует обязательный хэдер авторизации",
			}
			json.NewEncoder(w).Encode(response)
			return
		}
		idUser, errstr := h.TaskStorage.Indetification(login, password)
		if errstr != "" {
			log.Print("Ошибка авторизации")
			response := map[string]interface{}{
				"error": errstr,
			}
			json.NewEncoder(w).Encode(response)
			return
		}
		notes, err := h.TaskStorage.ReturnNotesFromDB(idUser)
		if err != nil {
			log.Print("Ошибка прочтения заметок")
			response = map[string]interface{}{
				"error": err,
			}
			json.NewEncoder(w).Encode(response)
			return
		}
		log.Printf("Заметки пользователя %s получены", login)
		json.NewEncoder(w).Encode(notes)
	}
}
