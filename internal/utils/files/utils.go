package fileUtils

import (
	"mime/multipart"
	"net/http"
)

func GetMimeType(fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return "", err
	}

	file.Seek(0, 0)

	mimeType := http.DetectContentType(buffer)
	return mimeType, nil
}
