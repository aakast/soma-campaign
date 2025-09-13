package handlers

import (
	"encoding/json"
	"soma-mayel-campaign/models"

	"github.com/gofiber/fiber/v2"
)

// TinaContentAPI handles content updates from TinaCMS
func TinaContentAPI(c *fiber.Ctx) error {
	var content models.Content
	if err := json.Unmarshal(c.Body(), &content); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid content format",
		})
	}

	// Save content to database/file system
	if err := models.SaveContent(content); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to save content",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Content saved successfully",
	})
}

// TinaGetContent retrieves content for TinaCMS
func TinaGetContent(c *fiber.Ctx) error {
	collection := c.Params("collection")
	
	content, err := models.GetContentByCollection(collection)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Content not found",
		})
	}

	return c.JSON(content)
}