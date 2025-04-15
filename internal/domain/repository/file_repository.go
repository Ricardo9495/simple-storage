package repository

import (
	"context"
	"errors"
	"mymodule/internal/domain/entity"
)

var (
	ErrNotFound = errors.New("file not found")
)

type FileRepository interface {
	GetAllFile(ctx context.Context) ([]entity.File, error)
	FindByName(ctx context.Context, name string) (*entity.File, error)
	Create(ctx context.Context, file entity.File) (*entity.File, error)
	DeleteByName(ctx context.Context, name string) error
}
