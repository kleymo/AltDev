package utils

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
	"github.com/gabriel-vasile/mimetype"
	"github.com/google/uuid"
)

var allowedMimeTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/webp": true,
}

func UploadFile(header *multipart.FileHeader, uploadDir string, maxLoad int64) (string, error) {
	if header.Size > maxLoad {
		return "", fmt.Errorf("файл превышает лимит %d байт", maxLoad)
	}

	src, err := header.Open()
	if err != nil {
		return "", fmt.Errorf("не удалось открыть файл: %w", err)
	}
	defer src.Close()

	data, err := io.ReadAll(io.LimitReader(src, maxLoad))
	if err != nil {
		return "", fmt.Errorf("ошибка чтения файла: %w", err)
	}

	mtype := mimetype.Detect(data)
	if !allowedMimeTypes[mtype.String()] {
		return "", fmt.Errorf("недопустимый тип файла: %s", mtype.String())
	}

	if _, err := imaging.Decode(bytes.NewReader(data)); err != nil {
		return "", fmt.Errorf("файл поврежден или не является изображением")
	}

	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", fmt.Errorf("не удалось создать директорию: %w", err)
	}

	ext := ""
	switch mtype.String() {
	case "image/jpeg":
		ext = ".jpg"
	case "image/png":
		ext = ".png"
	case "image/webp":
		ext = ".webp"
	default:
		ext = filepath.Ext(header.Filename)
	}

	newFileName := uuid.New().String() + ext
	fullPath := filepath.Join(uploadDir, newFileName)

	err = os.WriteFile(fullPath, data, 0644)
	if err != nil {
		return "", fmt.Errorf("ошибка при сохранении файла: %w", err)
	}

	return newFileName, nil
}
