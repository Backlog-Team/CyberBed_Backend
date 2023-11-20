package coder

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func extractImageExtension(dataURI string) (string, error) {
	parts := strings.Split(dataURI, ";")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid data URI format")
	}

	mimeTypePart := parts[0]
	mimeTypeParts := strings.Split(mimeTypePart, "/")
	if len(mimeTypeParts) != 2 {
		return "", fmt.Errorf("invalid MIME type format")
	}

	imageType := mimeTypeParts[1]
	return imageType, nil
}

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

	imageType, err := extractImageExtension(parts[0])
	if err != nil {
		return "", err
	}

	uuid := uuid.New()

	imagePath := fmt.Sprintf("%s.%s", uuid.String(), imageType)
	if err = os.WriteFile(imagePath, decoded, 0644); err != nil {
		return "", err
	}

	return imagePath, nil
}

func EncodeToBase64(content, extension string) (string, error) {
	var imageType string
	switch extension {
	case ".png":
		imageType = "image/png"
	case ".jpg", ".jpeg", ".jfif":
		imageType = "image/jpeg"
	default:
		return "", fmt.Errorf("unsupported image format: %s", extension)
	}

	encoded := base64.StdEncoding.EncodeToString([]byte(content))
	encodedData := fmt.Sprintf("data:%s;base64,%s", imageType, encoded)
	return encodedData, nil
}
