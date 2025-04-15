package mock

import (
	"mime/multipart"
)

type FileUploadInterface struct {
	UploadFileFn func(file multipart.File, name string) (string, error)
	DeleteFileFn func(name string) error
}

func (fu *FileUploadInterface) UploadFile(file multipart.File, name string) (string, error) {
	return fu.UploadFileFn(file, name)
}

func (fu *FileUploadInterface) DeleteFile(name string) error {
	return fu.DeleteFileFn(name)
}
