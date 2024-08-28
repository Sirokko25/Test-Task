package structs

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
)

const Port = "8080"

type Request struct {
	Note string `json:"note"`
}

type YandexSpeller struct {
	Code int64    `json:"code"`
	Pos  int64    `json:"pos"`
	Row  int64    `json:"row"`
	Col  int64    `json:"col"`
	Len  int64    `json:"len"`
	Word string   `json:"word"`
	S    []string `json:"s"`
}

func (y YandexSpeller) CheckYandexSpeller(note string) (string, error) {
	url := "https://speller.yandex.net/services/spellservice.json/checkText?text="
	noteslice := strings.Split(note, " ")
	for i, v := range noteslice {
		if (len(noteslice) - i) == 1 {
			url += v
		} else {
			url += (v + "+")
		}
	}
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Ошибка доступа к YandexSpeller")
		return "", errors.New("Ошибка доступа к YandexSpeller")
	}
	defer resp.Body.Close()

	var resultWord []YandexSpeller
	err = json.NewDecoder(resp.Body).Decode(&resultWord)
	if err != nil {
		log.Println("Ошибка работы с JSON")
		return "", errors.New("Ошибка работы с JSON")
	}
	if len(resultWord) != 0 {
		for _, wordjson := range resultWord {
			for p, wordnote := range noteslice {
				if wordjson.Word == wordnote {
					noteslice[p] = wordjson.S[0]
				}
			}
		}
	}
	note = ""
	for i, v := range noteslice {
		if (len(noteslice) - i) == 1 {
			note += v
		} else {
			note += (v + " ")
		}
	}
	return note, nil
}
