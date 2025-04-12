package service

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"mobilka/internal/utils"

	"github.com/google/uuid"
)

// ImageService handles image operations
type ImageService struct {
	uploadPath string
}

// NewImageService creates a new image service
func NewImageService(uploadPath string) *ImageService {
	return &ImageService{
		uploadPath: uploadPath,
	}
}

// SaveImage saves an uploaded image to the local file system
func (s *ImageService) SaveImage(file *multipart.FileHeader) (string, error) {
	// Get the file extension
	ext := filepath.Ext(file.Filename)
	origName := strings.TrimSuffix(filepath.Base(file.Filename), ext)

	// Generate a unique ID
	uniqueID := uuid.New().String()

	// Create the final filename with the format: filename!_unique_id.extension
	filename := fmt.Sprintf("%s!_%s%s", origName, uniqueID, ext)

	// Create the full path
	fullPath := filepath.Join(s.uploadPath, filename)

	// Open the source file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	// Create the destination file
	dst, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	// Copy the contents
	if _, err = io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("failed to copy file contents: %w", err)
	}

	return filename, nil
}

// DeleteImage deletes an image from the local file system
func (s *ImageService) DeleteImage(filename string) error {
	fullPath := filepath.Join(s.uploadPath, filename)

	// Check if file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return utils.ErrResourceNotFound
	}

	// Delete the file
	if err := os.Remove(fullPath); err != nil {
		return fmt.Errorf("failed to delete image: %w", err)
	}

	return nil
}

// GetImagePath returns the full path to an image
func (s *ImageService) GetImagePath(filename string) string {
	return filepath.Join(s.uploadPath, filename)
}
