package utils

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
)

func (h Handler) UploadingUserPhoto(userId int, header *multipart.FileHeader) error {
	photoPath, err := UploadFile(header, h.config.ImgDefaultPath, h.config.ImgMaxUploadMb*1024*1024)
	if err != nil {
		return fmt.Errorf("Не удалось загрузить фото: %w", err)
	}

	query := `UPDATE Users SET PhotoURL = $1 WHERE UserId = $2`

	_, err = h.db.Exec(query, photoPath, userId)
	if err != nil {

		fullPath := filepath.Join(h.config.ImgDefaultPath, photoPath)
		_ = os.Remove(fullPath)
		return fmt.Errorf("Ошибка БД: %w", err)
	}
	return nil
}
