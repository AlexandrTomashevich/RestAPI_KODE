package utils

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const endpoint = "https://speller.yandex.net/services/spellservice.json/checkText"

type SpellError struct {
	Word string   `json:"word"`
	S    []string `json:"s"`
}

func CheckSpelling(text string) (string, error) {
	resp, err := http.Get(endpoint + "?text=" + url.QueryEscape(text) + "&options=512")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("Speller service returned non-OK status")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var errors []SpellError
	if err := json.Unmarshal(body, &errors); err != nil {
		return "", err
	}

	for _, spellError := range errors {
		if len(spellError.S) > 0 {
			log.Printf("Corrected word '%s' to '%s'. Updated text: %s", spellError.Word, spellError.S[0], text)
			text = strings.Replace(text, spellError.Word, spellError.S[0], -1)
		}
	}
	return text, nil
}
