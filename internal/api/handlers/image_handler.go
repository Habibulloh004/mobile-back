package handlers

import (
	"mobilka/internal/service"
	"mobilka/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// ImageHandler handles image requests
type ImageHandler struct {
	imageService *service.ImageService
}

// NewImageHandler creates a new image handler
func NewImageHandler(imageService *service.ImageService) *ImageHandler {
	return &ImageHandler{
		imageService: imageService,
	}
}

// Upload handles uploading an image
func (h *ImageHandler) Upload(c *fiber.Ctx) error {
	// Get the uploaded file
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "No image provided",
		})
	}

	// Check file size
	if file.Size > utils.MaxImageSize {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Image size exceeds the maximum allowed size",
		})
	}

	// Save the image
	filename, err := h.imageService.SaveImage(file)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Failed to upload image",
		})
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": utils.StatusSuccess,
		"data": fiber.Map{
			"filename": filename,
			"url":      filename,
		},
	})
}

// Get handles retrieving an image
func (h *ImageHandler) Get(c *fiber.Ctx) error {
	// Get the filename from URL
	filename := c.Params("filename")
	if filename == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "No filename provided",
		})
	}

	// Get the image path
	imagePath := h.imageService.GetImagePath(filename)

	// Return the image
	return c.SendFile(imagePath)
}

// Delete handles deleting an image
func (h *ImageHandler) Delete(c *fiber.Ctx) error {
	// Get the filename from URL
	filename := c.Params("filename")
	if filename == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "No filename provided",
		})
	}

	// Delete the image
	err := h.imageService.DeleteImage(filename)
	if err != nil {
		if err == utils.ErrResourceNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  utils.StatusError,
				"message": "Image not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  utils.StatusError,
			"message": "Failed to delete image",
		})
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  utils.StatusSuccess,
		"message": "Image deleted successfully",
	})
}
