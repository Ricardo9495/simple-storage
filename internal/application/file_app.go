package application

import (
	"context"
	"mymodule/internal/domain/entity"
	"mymodule/internal/domain/repository"
)

type fileApp struct {
	fr repository.FileRepository
}

var _ FileAppInterface = &fileApp{}

type FileAppInterface interface {
	GetAllFile(ctx context.Context) ([]entity.File, error)
	FindByName(ctx context.Context, name string) (*entity.File, error)
	Create(ctx context.Context, file entity.File) (*entity.File, error)
	DeleteByName(ctx context.Context, name string) error
}

func (f *fileApp) GetAllFile(ctx context.Context) ([]entity.File, error) {
	return f.fr.GetAllFile(ctx)
}

func (f *fileApp) FindByName(ctx context.Context, name string) (*entity.File, error) {
	return f.fr.FindByName(ctx, name)
}

func (f *fileApp) Create(ctx context.Context, file entity.File) (*entity.File, error) {
	return f.fr.Create(ctx, file)
}

func (f *fileApp) DeleteByName(ctx context.Context, name string) error {
	return f.fr.DeleteByName(ctx, name)
}
