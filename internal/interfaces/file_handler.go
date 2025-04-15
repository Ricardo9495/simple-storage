package interfaces

import (
	"errors"
	"fmt"
	"mymodule/internal/application"
	"mymodule/internal/config"
	"mymodule/internal/domain/entity"
	"mymodule/internal/domain/repository"
	"mymodule/internal/interfaces/fileupload"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FileHandler struct {
	fileApp    application.FileAppInterface
	fileUpload fileupload.FileUploadInterface
	config     *config.Config
}

// File constructor
func NewFileHandler(fileApp application.FileAppInterface, fileUpload fileupload.FileUploadInterface, config *config.Config) *FileHandler {
	return &FileHandler{
		fileApp:    fileApp,
		fileUpload: fileUpload,
		config:     config,
	}
}

func (fi *FileHandler) GetAllFile(gc *gin.Context) {
	files, err := fi.fileApp.GetAllFile(gc.Request.Context())
	if err != nil {
		gc.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	gc.JSON(http.StatusOK, files)
}

func (fi *FileHandler) UploadFile(gc *gin.Context) {
	// Get the file name from the request.
	fileName := gc.Param("file_name")
	if fileName == "" {
		fmt.Println("[Error] UploadFile: File name is required")
		gc.JSON(http.StatusBadRequest, gin.H{"error": "File name is required"})
		return
	}

	// Get the file from the request.
	file, header, err := gc.Request.FormFile("file")
	if err != nil {
		fmt.Println("[Error] UploadFile: Failed to get file from request")
		gc.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get file from request"})
		return
	}
	defer file.Close()

	// 1. Check duplication
	fileExist, err := fi.fileApp.FindByName(gc.Request.Context(), fileName)
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		fmt.Println("[Error] UploadFile - FindByName:", err.Error())
		gc.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if fileExist != nil {
		fmt.Println("[Error] UploadFile: File Exists!")
		gc.JSON(http.StatusInternalServerError, gin.H{"error": "File Exists!"})
		return
	}

	// 2. Validate file
	newFile := &entity.File{
		Name: fileName,
		Size: header.Size,
	}

	err = newFile.Validate(fi.config)
	if err != nil {
		fmt.Println("[Error] UploadFile - Validate:", err.Error())
		gc.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 3. Upload file
	filePath, err := fi.fileUpload.UploadFile(file, fileName)
	if err != nil {
		fmt.Println("[Error] UploadFile - UploadFile:", err.Error())
		gc.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	newFile.Path = filePath

	// 4. Store in DB
	newFile, err = fi.fileApp.Create(gc.Request.Context(), *newFile)
	if err != nil {
		fmt.Println("[Error] UploadFile - Create Metadata:", err.Error())
		gc.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	gc.JSON(http.StatusOK, gin.H{
		"file": newFile,
	})
}

func (fi *FileHandler) DeleteFile(gc *gin.Context) {
	// 1. Get Name from query
	fileName := gc.Param("file_name")
	if fileName == "" {
		fmt.Println("[Error] DeleteFile: File name is required")
		gc.JSON(http.StatusBadRequest, gin.H{"error": "File name is required"})
		return
	}

	// 2. Find file in the database
	_, err := fi.fileApp.FindByName(gc.Request.Context(), fileName)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			fmt.Println("[Error] DeleteFile - FindByName: File not found")
			gc.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
			return
		}
		fmt.Println("[Error] DeleteFile - FindByName:", err.Error())
		gc.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 3. Delete file from storage
	err = fi.fileUpload.DeleteFile(fileName)
	if err != nil {
		fmt.Println("[Error] DeleteFile - DeleteFile:", err.Error())
		gc.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 4. Delete file from database
	err = fi.fileApp.DeleteByName(gc.Request.Context(), fileName)
	if err != nil {
		fmt.Println("[Error] DeleteFile - DeleteByName:", err.Error())
		gc.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	gc.JSON(http.StatusOK, gin.H{"message": "File deleted successfully"})
}
