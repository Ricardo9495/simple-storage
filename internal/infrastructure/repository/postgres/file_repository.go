package postgres

import (
	"context"
	"fmt"
	"mymodule/internal/domain/entity"
	"mymodule/internal/domain/repository"

	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
)

type FileRepository struct {
	*Postgres
}

func NewFileRepository(pg *Postgres) *FileRepository {
	return &FileRepository{pg}
}

var _ repository.FileRepository = &FileRepository{}

func (f *FileRepository) Create(ctx context.Context, file entity.File) (*entity.File, error) {
	sql, args, err := f.Builder.Insert("file").
		Columns("name", "path", "size").
		Values(file.Name, file.Path, file.Size).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("FileRepository - Create - f.Builder: %w", err)
	}
	err = f.Pool.QueryRow(ctx, sql, args...).Scan(&file.ID)
	if err != nil {
		return nil, fmt.Errorf("FileRepository - Create - row.Scan: %w", err)
	}

	return &file, nil
}

func (f *FileRepository) GetAllFile(ctx context.Context) ([]entity.File, error) {
	sql, _, err := f.Builder.Select("*").From("file").ToSql()
	if err != nil {
		return nil, fmt.Errorf("FileRepository - GetAllFile - r.Builder: %w", err)
	}

	rows, err := f.Pool.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("FileRepository - GetAllFile - r.Pool.Query: %w", err)
	}
	defer rows.Close()
	var files []entity.File
	for rows.Next() {
		e := entity.File{}

		err = rows.Scan(&e.ID, &e.Name, &e.Path, &e.Size, &e.CreatedAt, &e.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("FileRepository - GetAllFile - rows.Scan: %w", err)
		}

		files = append(files, e)
	}

	return files, nil
}

func (f *FileRepository) FindByName(ctx context.Context, name string) (*entity.File, error) {
	sql, args, err := f.Builder.Select("*").From("file").Where(squirrel.Eq{"name": name}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("FileRepository - FindByName - f.Builder: %w", err)
	}
	row := f.Pool.QueryRow(ctx, sql, args...)
	file := &entity.File{}
	err = row.Scan(&file.ID, &file.Name, &file.Path, &file.Size, &file.CreatedAt, &file.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("FileRepository - FindByName - row.Scan: %w", err)
	}

	return file, nil
}

func (f *FileRepository) DeleteByName(ctx context.Context, name string) error {
	sql, args, err := f.Builder.Delete("file").Where(squirrel.Eq{"name": name}).ToSql()
	if err != nil {
		return fmt.Errorf("FileRepository - DeleteByName - f.Builder: %w", err)
	}
	_, err = f.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("FileRepository - DeleteByName - f.Pool.Exec: %w", err)
	}

	return nil
}
