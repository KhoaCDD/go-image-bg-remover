package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"

	services "go-image-bg-remover/services"
)

func UploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("Upload file error: %s", err.Error()))
		return
	}

	// Check if the file is an image
	fileHeader := make([]byte, 512)
	fileContent, err := file.Open()
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Open file error: %s", err.Error()))
		return
	}
	defer fileContent.Close()

	_, err = fileContent.Read(fileHeader)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Read file error: %s", err.Error()))
		return
	}

	fileType := http.DetectContentType(fileHeader)
	if !strings.HasPrefix(fileType, "image/") {
		c.String(http.StatusBadRequest, "The uploaded file is not an image")
		return
	}

	// Get the string from the form
	clientID := c.PostForm("clientID")
	if clientID == "" {
		c.String(http.StatusBadRequest, "clientID is required")
		return
	}

	UploadPath := os.Getenv("UPLOAD_DIR")

	fileName := strings.TrimSuffix(file.Filename, filepath.Ext(file.Filename))

	randomFileName := fmt.Sprintf("%s_%s%s", fileName, services.GenerateRandomString(16), filepath.Ext(file.Filename))
	dst := filepath.Join(UploadPath, randomFileName)

	// Save the file to the server static directory
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Upload file error: %s", err.Error()))
		return
	}

	// start process image
	go services.ProcessedImage(randomFileName, clientID)

	c.String(http.StatusOK, fmt.Sprintf("File %s uploaded successfully.", file.Filename))
}

// Download handler that reads the validated file path from the context
func DownloadImage(c *gin.Context) {
	// Retrieve the validated file path from the context
	filePath, exists := c.Get("filePath")
	if !exists {
		c.JSON(500, gin.H{"error": "filePath not found in context"})
		return
	}

	// Serve the file
	filePathStr, _ := filePath.(string)
	c.File("static/" + filePathStr)
}
