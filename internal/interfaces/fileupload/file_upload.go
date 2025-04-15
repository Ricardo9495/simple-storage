package fileupload

import (
	"fmt"
	"io"
	"mime/multipart"
	"mymodule/internal/config"
	"os"
)

func NewFileUpload(config *config.Config) *fileUpload {
	return &fileUpload{config}
}

type fileUpload struct {
	config *config.Config
}

type FileUploadInterface interface {
	UploadFile(file multipart.File, name string) (string, error)
	DeleteFile(name string) error
}

var _ FileUploadInterface = &fileUpload{}

func (fu *fileUpload) UploadFile(file multipart.File, name string) (string, error) {
	// Create a new file in the "tmp" directory.
	//TODO: fix ./tmp to constant
	fmt.Println("[DEBUG] TMP", fu.config.StorageDir)
	out, err := os.Create(fu.config.StorageDir + "/" + name)
	if err != nil {
		fmt.Println("[Error] Failed To Create File")
		return "", err
	}
	defer out.Close()
	// Copy the file content to the new file.
	_, err = io.Copy(out, file)
	if err != nil {
		return "", err
	}

	return out.Name(), nil
}

func (fu *fileUpload) DeleteFile(name string) error {
	err := os.Remove(fu.config.StorageDir + "/" + name)
	if err != nil {
		fmt.Println("[DEBUG] Failed to delete file")
		return err
	}

	return nil
}
