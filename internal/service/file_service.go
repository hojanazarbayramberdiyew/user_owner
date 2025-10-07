package service

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type FileService interface {
	UploadUserLogo(ctx context.Context, userID uuid.UUID, fileHeader *multipart.FileHeader) (string, error)
}

type fileService struct {
	uploadDir    string
	allowedTypes map[string]bool
}

func NewFileService(uploadDir string) FileService {
	return &fileService{
		uploadDir: uploadDir,
		allowedTypes: map[string]bool{
			"image/jpeg": true,
			"image/png":  true,
			"image/gif":  true,
			"image/webp": true,
		},
	}
}

func (s *fileService) UploadUserLogo(ctx context.Context, userID uuid.UUID, fileHeader *multipart.FileHeader) (string, error) {

	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("fayl acylmady: %v", err)
	}

	defer file.Close()

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return "", fmt.Errorf("MIME tipi barlap bolmady: %v", err)
	}

	mimeType := http.DetectContentType(buffer)

	if !s.allowedTypes[mimeType] {
		return "", fmt.Errorf("rugsat berilmedik fayl gornusi :%s", mimeType)
	}

	fileExt := filepath.Ext(fileHeader.Filename)
	fileName := fmt.Sprintf("%s%s", userID.String(), fileExt)
	filePath := filepath.Join(s.uploadDir, "users", fileName)

	err = os.MkdirAll(filepath.Dir(filePath), 0755)
	if err != nil {
		return "", fmt.Errorf("directory doredilmedi: %v", err)
	}

	outFile, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("file doredilmedi %v", err)
	}

	defer outFile.Close()

	_, err = file.Seek(0, 0)
	if err != nil {
		return "", fmt.Errorf("seek hatasy: %v", err)
	}

	_, err = io.Copy(outFile, file)
	if err != nil {
		return "", fmt.Errorf("file yazylmady: %v", err)
	}

	return filePath, nil

}
