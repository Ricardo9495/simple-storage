package interfaces

import (
	"mymodule/internal/config"
	"mymodule/internal/interfaces/mock"
)

var (
	fileApp    mock.FileAppInterface
	fileUpload mock.FileUploadInterface

	f = NewFileHandler(&fileApp, &fileUpload, mockConfig)
)

var mockConfig = &config.Config{
	ApiHost:     "localhost",
	ApiPort:     8080,
	MaxFileSize: 10 << 20, // 10 MB
	StorageDir:  "./uploads",
	PG:          pgConfig,
}

var pgConfig = config.PG{
	Host:     "localhost",
	Port:     int(5432),
	User:     "myuser",
	Password: "mypassword",
	DBName:   "postgres",
}
