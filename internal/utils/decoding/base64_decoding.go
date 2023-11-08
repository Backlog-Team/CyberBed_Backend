package decoding

import (
	"encoding/base64"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func DecodeBase64(encryptedData string) (string, error) {
	parts := strings.Split(encryptedData, ",")
	if len(parts) != 2 {
		return "", errors.New("incorrect data to encode")
	}

	imageData := parts[1]
	decoded, err := base64.StdEncoding.DecodeString(imageData)
	if err != nil {
		return "", err
	}

	uuid := uuid.New()
	if err = os.WriteFile(uuid.String()+".png", decoded, 0644); err != nil {
		return "", err
	}

	return uuid.String() + ".png", nil
}
