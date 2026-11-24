package utils

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"
)

var AllowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".webp": true,
}

func UploadFile(file io.Reader, OriginalFileName string, path string) (string, error) {
	ext := strings.ToLower(filepath.Ext(OriginalFileName))
	if !AllowedExtensions[ext] {
		return "", fmt.Errorf("недопустимый формат файла: %s", ext)
	}
}
