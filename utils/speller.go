package utils

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type SpellError struct {
	Word    string `json:"word"`
	Suggest string `json:"s"`
}

func CheckSpelling(text string) ([]SpellError, error) {
	endpoint := "https://speller.yandex.net/services/spellservice.json/checkText"
	resp, err := http.PostForm(endpoint, url.Values{"text": {text}})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var errors []SpellError
	if err := json.NewDecoder(resp.Body).Decode(&errors); err != nil {
		return nil, err
	}

	return errors, nil
}
