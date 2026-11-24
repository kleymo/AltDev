package utils

import (
	"altayCTF/internal/config"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"uuid"
)

var AllowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".webp": true,
}

func GetImgDirAndMaxUpload() (string, string, error) {
	cfg, err := config.Load()
	if err != nil {
		return "", "", fmt.Errorf("Не удалось получить конфиг:%s", err)
	}
	return cfg.ImgDefaultPath, cfg.ImgMaxUploadMb, nil
}

func UploadFile(file io.Reader, OriginalFileName string, path string) (string, error) {
	uploadDir, uploadMax, err := GetImgDirAndMaxUpload()
	if err != nil {
		log.Println("Error in ENV: ", err)
		return "", err
	}

	//сделать гибче!!

	ext := strings.ToLower(filepath.Ext(OriginalFileName))
	if !AllowedExtensions[ext] {
		return "", fmt.Errorf("недопустимый формат файла: %s", ext)
	}
	newFileName := uuid.New().String() + ext

	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err := os.MkdirAll(uploadDir, os.ModePerm)
		if err != nil {
			return "", fmt.Errorf("не удалось создать директорию: %w", err)
		}
	}

	fullPath := filepath.Join(uploadDir, newFileName)

	out, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("не удалось создать файл на диске: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		return "", fmt.Errorf("ошибка при записи данных файла: %w", err)
	}
	return newFileName, nil
}
