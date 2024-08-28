package storage

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	conn *sql.DB
}

func OpenDb() (DB, error) {
	db, err := sql.Open("sqlite3", "table.db")
	if err != nil {
		log.Fatal(err)
		return DB{conn: nil}, err
	}
	return DB{conn: db}, nil
}

func (db *DB) Indetification(login string, password string) (int64, string) {
	passwordCheck := ""
	var idUser int64 = 0
	query := `SELECT userPassword, id FROM Users WHERE userLogin = ?`
	err := db.conn.QueryRow(query, login).Scan(&passwordCheck, &idUser)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, "Пользователь не найден"
		} else {
			return 0, "Ошибка выполнения запроса"
		}
	}
	result := sha256.Sum256([]byte(password))
	passwordEntry := hex.EncodeToString(result[:])
	if passwordEntry != passwordCheck {
		return 0, "Неправильный пароль"
	}
	return idUser, ""
}

func (db *DB) AddNoteToDatabase(id int64, note string) error {
	query := `INSERT INTO UserNotes (idUser, note) VALUES (?, ?)`
	_, err := db.conn.Exec(query, id, note)
	if err != nil {
		return errors.New("Ошибка добавления задачи")
	}
	return nil
}

func (db *DB) ReturnNotesFromDB(id int64) ([]string, error) {
	rows, err := db.conn.Query(`SELECT note FROM UserNotes WHERE idUser = :id`, sql.Named("id", id))
	if err != nil {
		return []string{}, errors.New("Ошибка выполнения запроса")
	}
	defer rows.Close()

	notes := make([]string, 0, 0)

	for rows.Next() {
		note := ""
		if err := rows.Scan(&note); err != nil {
			return []string{}, errors.New("Ошибка чтения строки: ")
		}
		notes = append(notes, note)
	}

	if err := rows.Err(); err != nil {
		return []string{}, errors.New("Ошибка обработки результата: ")
	}
	return notes, nil
}
