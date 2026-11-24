package utils

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"uuid"

	"github.com/disintegration/imaging"
	"github.com/gabriel-vasile/mimetype"
)

var allowedMimeTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/webp": true,
}

func UploadFile(header *multipart.FileHeader, uploadDir string, MaxLoad int64) (string, error) {
	//сделать гибче!!
	if header.Size > MaxLoad {
		return "", fmt.Errorf("Файл слишком большой")
	}
	src, err := header.Open()
	if err != nil {
		return "", fmt.Errorf("не удалось открыть файл: %w", err)
	}
	defer src.Close()

	data, err := io.ReadAll(io.LimitReader(src, MaxLoad+1))
	if err != nil {
		return "", fmt.Errorf("ошибка чтения файла: %w", err)
	}

	if int64(len(data)) > MaxLoad { // ← вот этого не хватает
		return "", fmt.Errorf("файл превышает лимит %d байт", MaxLoad)
	}

	mtype := mimetype.Detect(data)
	if !allowedMimeTypes[mtype.String()] {
		return "", fmt.Errorf("недопустимый тип файла: %s", mtype.String())
	}

	if _, err := imaging.Decode(bytes.NewReader(data)); err != nil {
		return "", fmt.Errorf("недопустимый тип файла")
	}

	//Директория с безопасными правами
	if err := os.MkdirAll(uploadDir, 0o755); err != nil {
		return "", fmt.Errorf("не удалось создать директорию: %w", err)
	}

	newFileName := uuid.New().String() + mtype.Extension()
	fullPath := filepath.Join(uploadDir, newFileName)

	tmp, err := os.CreateTemp(uploadDir, "upload-*")
	if err != nil {
		return "", fmt.Errorf("не удалось создать временный файл: %w", err)
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)

	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return "", fmt.Errorf("ошибка записи: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return "", fmt.Errorf("ошибка закрытия: %w", err)
	}

	if err := os.Chmod(tmpName, 0o644); err != nil {
		return "", fmt.Errorf("ошибка прав: %w", err)
	}

	if err := os.Rename(tmpName, fullPath); err != nil {
		return "", fmt.Errorf("ошибка сохранения: %w", err)
	}

	return newFileName, nil
}
