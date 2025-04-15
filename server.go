// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Hello is a simple hello, world demonstration web server.
//
// It serves version information on /version and answers
// any other request like /name by saying "Hello, name!".
//
// See golang.org/x/example/outyet for a more sophisticated server.
package main

import (
	"fmt"
	"log"
	"mymodule/internal/config"
	"mymodule/internal/infrastructure/repository/postgres"
	"mymodule/internal/interfaces"
	"mymodule/internal/interfaces/fileupload"
	"mymodule/internal/interfaces/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	//To load our environmental variables.
	if err := godotenv.Load(); err != nil {
		log.Println("no env gotten")
	}
}

func main() {
	config := config.NewConfig()
	pg := postgres.NewOrGetSingleton(config)

	fileUpload := fileupload.NewFileUpload(config)

	fileRepository := postgres.NewFileRepository(pg)
	files := interfaces.NewFileHandler(fileRepository, fileUpload, config)

	addr := fmt.Sprintf("%s:%d", config.ApiHost, config.ApiPort)
	fmt.Println("addr", addr)

	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	r.GET("/", ping)
	r.GET("/files", files.GetAllFile)
	r.POST("/files/:file_name", middleware.MaxSizeAllowed(8192000), files.UploadFile)
	r.DELETE("/files/:file_name", files.DeleteFile)
	log.Fatal(r.Run(fmt.Sprintf(":%d", config.ApiPort)))

	log.Printf("serving http://%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func ping(r *gin.Context) {
	r.JSON(http.StatusOK, "Hello")
}
