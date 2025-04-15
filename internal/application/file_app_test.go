package application

import (
	"context"
	"errors"
	"mymodule/internal/domain/entity"
	"testing"

	"github.com/stretchr/testify/assert"
)

type fileRepositoryMock struct{}

var (
	getAllFile   func(ctx context.Context) ([]entity.File, error)
	findByName   func(ctx context.Context, name string) (*entity.File, error)
	create       func(ctx context.Context, file entity.File) (*entity.File, error)
	deleteByName func(ctx context.Context, name string) error
)

func (m *fileRepositoryMock) GetAllFile(ctx context.Context) ([]entity.File, error) {
	return getAllFile(ctx)
}

func (m *fileRepositoryMock) FindByName(ctx context.Context, name string) (*entity.File, error) {
	return findByName(ctx, name)
}

func (m *fileRepositoryMock) DeleteByName(ctx context.Context, name string) error {
	return deleteByName(ctx, name)
}

func (m *fileRepositoryMock) Create(ctx context.Context, file entity.File) (*entity.File, error) {
	return create(ctx, file)
}

var fileAppMock FileAppInterface = &fileRepositoryMock{}

func TestGetAllFile_success(t *testing.T) {
	getAllFile = func(ctx context.Context) ([]entity.File, error) {
		return []entity.File{
			{
				ID:   1,
				Name: "file1",
				Path: "/path/to/file1",
				Size: 1024,
			},
			{
				ID:   2,
				Name: "file2",
				Path: "/path/to/file2",
				Size: 2048,
			},
		}, nil
	}

	f, err := fileAppMock.GetAllFile(context.Background())
	assert.Nil(t, err)
	assert.EqualValues(t, len(f), 2)
}

func TestGetAllFile_failure(t *testing.T) {
	getAllFile = func(_ context.Context) ([]entity.File, error) {
		return []entity.File{}, errors.New("Errors!")
	}

	f, err := fileAppMock.GetAllFile(context.Background())
	assert.EqualError(t, err, "Errors!")
	assert.EqualValues(t, len(f), 0)
}

func TestFindByName_Success(t *testing.T) {
	findByName = func(ctx context.Context, name string) (*entity.File, error) {
		return &entity.File{
			ID:   1,
			Name: "testfile",
			Path: "/path/to/testfile",
			Size: 1024,
		}, nil
	}

	file, err := fileAppMock.FindByName(context.Background(), "testfile")

	assert.Nil(t, err)
	assert.Equal(t, "testfile", file.Name)
	assert.Equal(t, "/path/to/testfile", file.Path)
	assert.Equal(t, int64(1024), file.Size)
}

func TestFindByName_Failure(t *testing.T) {
	findByName = func(ctx context.Context, name string) (*entity.File, error) {
		return nil, errors.New("not found")
	}
	file, err := fileAppMock.FindByName(context.Background(), "not_exist")
	assert.EqualError(t, err, "not found")
	assert.Nil(t, file)
}

func TestCreate_Success(t *testing.T) {
	create = func(ctx context.Context, file entity.File) (*entity.File, error) {
		file.ID = 1
		return &file, nil
	}
	fileToCreate := entity.File{
		Name: "newfile",
		Path: "/path/to/newfile",
		Size: 2048,
	}

	file, err := fileAppMock.Create(context.Background(), fileToCreate)
	assert.Nil(t, err)
	assert.Equal(t, file.ID, int64(1))
	assert.Equal(t, file.Name, "newfile")
	assert.Equal(t, file.Size, int64(2048))
}

func TestCreate_Failure(t *testing.T) {
	create = func(ctx context.Context, file entity.File) (*entity.File, error) {
		return nil, errors.New("Failed to create")
	}
	fileToCreate := entity.File{
		Name: "newfile",
		Path: "/path/to/newfile",
		Size: 2048,
	}
	file, err := fileAppMock.Create(context.Background(), fileToCreate)
	assert.EqualError(t, err, "Failed to create")
	assert.Nil(t, file)
}

func TestDeleteByName_Success(t *testing.T) {
	deleteByName = func(ctx context.Context, name string) error {
		return nil
	}

	err := fileAppMock.DeleteByName(context.Background(), "testfile")

	assert.Nil(t, err)
}

func TestDeleteByName_Failure(t *testing.T) {
	deleteByName = func(ctx context.Context, name string) error {
		return errors.New("Failed to delete")
	}

	err := fileAppMock.DeleteByName(context.Background(), "not_exist")

	assert.EqualError(t, err, "Failed to delete")

}
