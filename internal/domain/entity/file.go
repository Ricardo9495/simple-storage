package entity

import (
	"errors"
	"mymodule/internal/config"
	"time"
)

var (
	MaxSize         = int64(500)
	ErrNameRequired = errors.New("Name is required")
	ErrFileEmpty    = errors.New("File is empty")
	ErrFileTooLarge = errors.New("File is too large")
)

type File struct {
	ID        int64
	Name      string
	Path      string
	Size      int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (f *File) Validate(cf *config.Config) error {
	if f.Name == "" {
		return ErrNameRequired
	}
	if f.Size == 0 {
		return ErrFileEmpty
	}
	if f.Size >= cf.MaxFileSize {
		return ErrFileTooLarge
	}
	return nil
}
