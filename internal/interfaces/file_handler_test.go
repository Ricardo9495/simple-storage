package interfaces

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"mymodule/internal/domain/entity"
	"mymodule/internal/domain/repository"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetAllFile_failure(t *testing.T) {
	fileApp.GetAllFileFn = func(ctx context.Context) ([]entity.File, error) {
		return []entity.File{}, errors.New("Errors!")
	}

	req, err := http.NewRequest(http.MethodGet, "/files", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	r := gin.Default()
	r.GET("/files", f.GetAllFile)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusInternalServerError)
}

func TestUploadFile_success(t *testing.T) {
	fileApp.FindByNameFn = func(ctx context.Context, name string) (*entity.File, error) {
		return nil, repository.ErrNotFound
	}

	fileApp.CreateFn = func(ctx context.Context, file entity.File) (*entity.File, error) {
		file.ID = 1
		return &file, nil
	}
	fileUpload.UploadFileFn = func(file multipart.File, fileName string) (string, error) {
		return "dummy_url", nil
	}

	body := strings.NewReader("file content")
	req, err := http.NewRequest(http.MethodPost, "/files", body)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	req.Header.Set("Content-Type", "multipart/form-data")
	req.Header.Set("filename", "test.txt")

	r := gin.Default()
	r.POST("/files", f.UploadFile)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestUploadFile_BadRequest(t *testing.T) {
	body := strings.NewReader("file content")
	req, err := http.NewRequest(http.MethodPost, "/files", body)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}

	r := gin.Default()
	r.POST("/files", f.UploadFile)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestUploadFile_Duplicate(t *testing.T) {
	fileApp.FindByNameFn = func(ctx context.Context, name string) (*entity.File, error) {
		return &entity.File{
			ID:        1,
			Name:      "test.txt",
			Path:      "/tmp/test.txt",
			Size:      1024,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}, nil
	}
	req := prepareRequest(t)

	r := gin.Default()
	r.POST("/files", f.UploadFile)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestUploadFile_ValidateError(t *testing.T) {
	fileApp.FindByNameFn = func(ctx context.Context, name string) (*entity.File, error) {
		return nil, repository.ErrNotFound
	}

	body := strings.NewReader("file content")
	req, err := http.NewRequest(http.MethodPost, "/files", body)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	req.Header.Set("Content-Type", "multipart/form-data")
	req.Header.Set("filename", "test.txt")

	r := gin.Default()
	r.POST("/files", f.UploadFile)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestUploadFile_CreateError(t *testing.T) {
	fileApp.FindByNameFn = func(ctx context.Context, name string) (*entity.File, error) {
		return nil, repository.ErrNotFound
	}

	fileApp.CreateFn = func(ctx context.Context, file entity.File) (*entity.File, error) {
		return nil, errors.New("error db")
	}
	req := prepareRequest(t)

	r := gin.Default()
	r.POST("/files", f.UploadFile)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestUploadFile_UploadError(t *testing.T) {
	fileApp.FindByNameFn = func(ctx context.Context, name string) (*entity.File, error) {
		return nil, repository.ErrNotFound
	}

	fileApp.CreateFn = func(ctx context.Context, file entity.File) (*entity.File, error) {
		file.ID = 1
		return &file, nil
	}

	fileUpload.UploadFileFn = func(file multipart.File, fileName string) (string, error) {
		return "", errors.New("upload error")
	}
	req := prepareRequest(t)

	r := gin.Default()
	r.POST("/files", f.UploadFile)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestUploadFile_FindError(t *testing.T) {
	fileApp.FindByNameFn = func(ctx context.Context, name string) (*entity.File, error) {
		return nil, errors.New("some error")
	}

	req := prepareRequest(t)

	r := gin.Default()
	r.POST("/files", f.UploadFile)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func prepareRequest(t *testing.T) *http.Request {
	// File to upload
	filePath := "./sample.txt"
	fieldName := "file" // same as form field name expected by server

	// Open file
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Create a buffer to write our form data into
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Create form file field
	part, err := writer.CreateFormFile(fieldName, file.Name())
	if err != nil {
		panic(err)
	}

	// Copy file content to the part
	_, err = io.Copy(part, file)
	if err != nil {
		panic(err)
	}

	// Optionally add other form fields
	_ = writer.WriteField("user_id", "42")

	// Close the writer to finalize the form
	writer.Close()
	req, err := http.NewRequest(http.MethodPost, "/files", &requestBody)
	if err != nil {
		panic(err)
	}

	// Set the correct content type for multipart
	req.Header.Set("Content-Type", writer.FormDataContentType())

	return req
}

func TestDeleteFile_Success(t *testing.T) {
	fileApp.FindByNameFn = func(ctx context.Context, name string) (*entity.File, error) {
		return &entity.File{}, nil
	}
	fileApp.DeleteByNameFn = func(ctx context.Context, name string) error {
		return nil
	}
	req, err := http.NewRequest(http.MethodDelete, "/files/test.txt", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	r := gin.Default()
	r.DELETE("/files/:filename", f.DeleteFile)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestDeleteFile_Failure(t *testing.T) {
	fileApp.FindByNameFn = func(ctx context.Context, name string) (*entity.File, error) {
		return &entity.File{}, nil
	}
	fileApp.DeleteByNameFn = func(ctx context.Context, name string) error {
		return errors.New("error to delete")
	}
	req, err := http.NewRequest(http.MethodDelete, "/files/test.txt", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	r := gin.Default()
	r.DELETE("/files/:filename", f.DeleteFile)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestDeleteFile_FindByNameError(t *testing.T) {
	fileApp.FindByNameFn = func(ctx context.Context, name string) (*entity.File, error) {
		return nil, errors.New("find error")
	}
	req, err := http.NewRequest(http.MethodDelete, "/files/test.txt", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	r := gin.Default()
	r.DELETE("/files/:filename", f.DeleteFile)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestDeleteFile_NotFound(t *testing.T) {
	fileApp.FindByNameFn = func(ctx context.Context, name string) (*entity.File, error) {
		return nil, repository.ErrNotFound
	}
	req, err := http.NewRequest(http.MethodDelete, "/files/test.txt", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	r := gin.Default()
	r.DELETE("/files/:filename", f.DeleteFile)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestDeleteFile_DeleteError(t *testing.T) {
	fileApp.FindByNameFn = func(ctx context.Context, name string) (*entity.File, error) {
		return &entity.File{}, nil
	}
	fileUpload.DeleteFileFn = func(url string) error {
		return errors.New("delete error")
	}
	req, err := http.NewRequest(http.MethodDelete, "/files/test.txt", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	r := gin.Default()
	r.DELETE("/files/:filename", f.DeleteFile)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestDeleteFile_DeleteByNameError(t *testing.T) {
	fileApp.FindByNameFn = func(ctx context.Context, name string) (*entity.File, error) {
		return &entity.File{}, nil
	}
	fileUpload.DeleteFileFn = func(url string) error {
		return nil
	}
	fileApp.DeleteByNameFn = func(ctx context.Context, name string) error {
		return errors.New("delete by name error")
	}
	req, err := http.NewRequest(http.MethodDelete, "/files/test.txt", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	r := gin.Default()
	r.DELETE("/files/:filename", f.DeleteFile)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

// func TestDeleteFile_BadRequest(t *testing.T) {
// 	req, err := http.NewRequest(http.MethodDelete, "/files/", nil)
// 	if err != nil {
// 		t.Errorf("this is the error: %v\n", err)
// 	}
// 	r := gin.Default()
// 	r.DELETE("/files/:filename", f.DeleteFile)
// 	rr := httptest.NewRecorder()
// 	r.ServeHTTP(rr, req)
// 	assert.Equal(t, http.StatusBadRequest, rr.Code)
// }

func TestGetAllFile_success(t *testing.T) {
	fileApp.GetAllFileFn = func(ctx context.Context) ([]entity.File, error) {
		return []entity.File{
			{
				ID:        1,
				Name:      "test.txt",
				Path:      "/tmp/test.txt",
				Size:      1024,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				ID:        2,
				Name:      "test2.txt",
				Path:      "/tmp/test2.txt",
				Size:      2048,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}, nil
	}

	req, err := http.NewRequest(http.MethodGet, "/files", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	r := gin.Default()
	r.GET("/files", f.GetAllFile)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	var file []entity.File
	err = json.Unmarshal(rr.Body.Bytes(), &file)
	if err != nil {
		t.Errorf("cannot unmarshal response: %v\n", err)
	}
	assert.Equal(t, rr.Code, 200)
	assert.EqualValues(t, len(file), 2)
}
