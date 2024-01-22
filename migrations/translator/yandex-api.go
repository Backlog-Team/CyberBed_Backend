package translator

import (
	"context"
	"net/http"
	"strings"

	"github.com/carlmjohnson/requests"
	"github.com/pkg/errors"
)

const URL = "https://translate.api.cloud.yandex.net/translate/v2/translate"

type Message struct {
	SourceLanguageCode string   `json:"sourceLanguageCode"`
	TargetLanguageCode string   `json:"targetLanguageCode"`
	Format             string   `json:"format"`
	Texts              []string `json:"texts"`
	FolderId           string   `json:"folderId"`
	Model              string   `json:"model"`
	GlossaryConfig     struct {
		GlossaryData struct {
			GlossaryPairs []struct {
				SourceText     string `json:"sourceText"`
				TranslatedText string `json:"translatedText"`
			} `json:"glossaryPairs"`
		} `json:"glossaryData"`
	} `json:"glossaryConfig"`
	Speller bool `json:"speller"`
}

func NewMsg(msg string) Message {
	return Message{
		SourceLanguageCode: "en",
		TargetLanguageCode: "ru",
		Format:             "PLAIN_TEXT",
		Texts: []string{
			msg,
		},
	}
}

type Response struct {
	Translations []struct {
		Text                 string `json:"text"`
		DetectedLanguageCode string `json:"detectedLanguageCode"`
	} `json:"translations"`
}

func (r Response) ToStringArr() []string {
	res := make([]string, 0, len(r.Translations))
	for _, translation := range r.Translations {
		res = append(res, translation.Text)
	}

	return res
}

func Translate(ctx context.Context, message, apiKey string) (string, error) {
	var resp Response

	if message == "" {
		return "", nil
	}

	if err := requests.
		URL(URL).
		Method(http.MethodPost).
		BodyJSON(NewMsg(message)).
		Header("Content-Type", "application/json").
		Header("Authorization", "Api-Key").
		ToJSON(&resp).
		Fetch(ctx); err != nil {
		return "", errors.Wrapf(err, "failed to translate %s", message)
	}

	return strings.Join(resp.ToStringArr(), "\n"), nil
}
