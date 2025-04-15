package mock

import (
	"context"
	"mymodule/internal/domain/entity"
)

type FileAppInterface struct {
	GetAllFileFn   func(ctx context.Context) ([]entity.File, error)
	FindByNameFn   func(ctx context.Context, name string) (*entity.File, error)
	CreateFn       func(ctx context.Context, file entity.File) (*entity.File, error)
	DeleteByNameFn func(ctx context.Context, name string) error
}

// Create implements application.FileAppInterface.
func (f *FileAppInterface) Create(ctx context.Context, file entity.File) (*entity.File, error) {
	return f.CreateFn(ctx, file)
}

func (f *FileAppInterface) GetAllFile(ctx context.Context) ([]entity.File, error) {
	return f.GetAllFileFn(ctx)
}

func (f *FileAppInterface) FindByName(ctx context.Context, name string) (*entity.File, error) {
	return f.FindByNameFn(ctx, name)
}

func (f *FileAppInterface) DeleteByName(ctx context.Context, name string) error {
	return f.DeleteByNameFn(ctx, name)
}
