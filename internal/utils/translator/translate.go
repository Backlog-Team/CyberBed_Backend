package translator

import (
	"net/http"

	translator "github.com/turk/free-google-translate"
)

func Translate(text string) (string, error) {
	if text == "" {
		return "", nil
	}

	client := http.Client{}
	t := translator.NewTranslator(&client)
	result, err := t.Translate(text, "en", "ru")
	if err != nil {
		return "", err
	}

	return result, nil
}
