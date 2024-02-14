package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

func SavePhoto(c *gin.Context, file *multipart.FileHeader, userID string) (string, error) {
	assetsDir := os.Getenv("PHOTO_DIRECTORY")
	domain := os.Getenv("DOMAIN_NAME")
	if _, err := os.Stat(assetsDir); os.IsNotExist(err) {
		if err := os.Mkdir(assetsDir, 0755); err != nil {
			return "", err
		}
	}

	fileExt := filepath.Ext(file.Filename)
	currentTimestamp := time.Now().UnixNano() / int64(time.Millisecond)

	filename := fmt.Sprintf("%d_%s%s", currentTimestamp, userID, fileExt)

	dst := filepath.Join(assetsDir, filename)
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return "", err
	}
	if err := c.SaveUploadedFile(file, dst); err != nil {
		return "", err
	}

	photoPath := fmt.Sprintf("%s/images/%s", domain, filename)
	return photoPath, nil
}
