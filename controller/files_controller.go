package controller

import (
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func generateUniqueFileName(file *multipart.FileHeader) string {
	timestamp := time.Now().Unix() // Get the current Unix timestamp in seconds
	fileExtension := filepath.Ext(file.Filename)
	uniqueFileName := strconv.FormatInt(timestamp, 10) + fileExtension
	return uniqueFileName
}

func uploadFile(file *multipart.FileHeader) (string, error) {
	dirPath := os.Getenv("FILES_LOCATION")

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Generate a unique file name
	fileName := generateUniqueFileName(file)

	// Create the destination file
	dst, err := os.Create(dirPath + fileName)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// Copy the uploaded file to the destination
	_, err = io.Copy(dst, src)
	if err != nil {
		return "", err
	}
	return fileName, nil
}

func getFilePath(filename string) (string, error) {
	dirPath := os.Getenv("FILES_LOCATION")

	// Check if the file exists
	filePath := filepath.Join(dirPath, filename) // Assuming the files are stored in the "/app/data/files" directory
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return "", err // If the file doesn't exist, set the file URL to an empty string
	}

	return filePath, nil
}

func getFileURL(c *gin.Context, filename string) (string, error) {
	hostname := c.Request.Host
	_, err := getFilePath(filename)
	if err != nil {
		return "", err
	}
	return hostname + "/data/files/" + filename, nil

}

func GetFile(c *gin.Context) {
	filename := c.Param("filename")

	filePath, err := getFilePath(filename)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "file not found",
		})
		return
	}

	// Determine the content type based on the file extension
	contentType := mime.TypeByExtension(filepath.Ext(filename))
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// Set the appropriate headers
	c.Header("Content-Type", contentType)
	c.Header("Content-Disposition", "inline; filename="+filename)

	// Serve the file
	c.File(filePath)

}
