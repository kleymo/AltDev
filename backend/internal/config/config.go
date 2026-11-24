package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	PostgresUser         string
	PostgresPassword     string
	PostgresPort         string
	PostgresHost         string
	Postgres_db          string
	JwtSecret            string
	BackendPort          string
	DefaultAdminLogin    string
	DefaultAdminPassword string
	ImgMaxUploadMb       int64
	DocMaxUploadMb       int64
	ImgDefaultPath       string
	DocDefaultPath       string
}

func getEnvInt64(key string) int64 {
	val := os.Getenv(key)
	if val == "" {
		return 0
	}
	i, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0
	}
	return i
}

func Load() (*Config, error) {
	// Загружаем .env
	if err := godotenv.Load("../.env"); err != nil {
		return nil, fmt.Errorf("ERROR loading .env: %w", err)
	}

	return &Config{
		PostgresUser:         os.Getenv("postgres_user"),
		PostgresPassword:     os.Getenv("postgres_password"),
		PostgresPort:         os.Getenv("postgres_port"),
		PostgresHost:         os.Getenv("postgres_host"),
		Postgres_db:          os.Getenv("postgres_db"),
		JwtSecret:            os.Getenv("jwt_secret"),
		BackendPort:          os.Getenv("backend_port"),
		DefaultAdminLogin:    os.Getenv("default_admin_login"),
		DefaultAdminPassword: os.Getenv("default_admin_password"),

		ImgMaxUploadMb: getEnvInt64("img_max_upload_mb"),
		DocMaxUploadMb: getEnvInt64("doc_max_upload_mb"),

		ImgDefaultPath: os.Getenv("upload_img_dir"),
		DocDefaultPath: os.Getenv("upload_doc_dir"),
	}, nil
}
