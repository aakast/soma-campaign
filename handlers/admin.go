package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"soma-mayel-campaign/models"

	"github.com/gofiber/fiber/v2"
)

// AdminPage renders the admin UI
func AdminPage(c *fiber.Ctx) error {
	return c.Render("admin", fiber.Map{
		"Title": "Admin - Content Manager",
	})
}

// AdminListPosts returns all posts as JSON
func AdminListPosts(c *fiber.Ctx) error {
	posts := models.GetAllPosts()
	return c.JSON(posts)
}

// AdminGetPost returns a single post by id
func AdminGetPost(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "missing id"})
	}
	post := models.GetPostByID(id)
	if post == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "post not found"})
	}
	return c.JSON(post)
}

type upsertPostRequest struct {
	ID         string   `json:"id"`
	Title      string   `json:"title"`
	Slug       string   `json:"slug"`
	Content    string   `json:"content"`
	Excerpt    string   `json:"excerpt"`
	Author     string   `json:"author"`
	Date       string   `json:"date"`
	Image      string   `json:"image"`
	Tags       []string `json:"tags"`
	IsFeatured bool     `json:"is_featured"`
}

// AdminUpsertPost creates or updates a post
func AdminUpsertPost(c *fiber.Ctx) error {
	var req upsertPostRequest
	if err := json.Unmarshal(c.Body(), &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}

	// Basic validation
	if strings.TrimSpace(req.Title) == "" || strings.TrimSpace(req.Slug) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "title and slug are required"})
	}

	// Parse date (accept YYYY-MM-DD or RFC3339)
	var parsedDate time.Time
	var err error
	if req.Date == "" {
		parsedDate = time.Now()
	} else {
		parsedDate, err = time.Parse(time.RFC3339, req.Date)
		if err != nil {
			parsedDate, err = time.Parse("2006-01-02", req.Date)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid date format"})
			}
		}
	}

	post := models.Post{
		ID:         req.ID,
		Title:      strings.TrimSpace(req.Title),
		Slug:       strings.TrimSpace(req.Slug),
		Content:    req.Content,
		Excerpt:    req.Excerpt,
		Author:     req.Author,
		Date:       parsedDate,
		Image:      req.Image,
		Tags:       req.Tags,
		IsFeatured: req.IsFeatured,
	}

	if err := models.SavePost(&post); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to save post"})
	}

	return c.JSON(post)
}

// AdminDeletePost deletes a post by id
func AdminDeletePost(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "missing id"})
	}
	if err := models.DeletePost(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to delete post"})
	}
	return c.JSON(fiber.Map{"success": true})
}

// AdminUpload handles image uploads and returns a public URL
func AdminUpload(c *fiber.Ctx) error {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "file is required"})
	}

	file, err := fileHeader.Open()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "unable to open file"})
	}
	defer file.Close()

	url, err := saveUpload(fileHeader.Filename, file)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to save upload"})
	}

	return c.JSON(fiber.Map{"url": url})
}

func saveUpload(originalName string, src multipart.File) (string, error) {
	// Ensure uploads directory exists
	uploadDir := filepath.Join("./static", "images", "uploads")
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", err
	}

	// Create unique filename
	ext := filepath.Ext(originalName)
	base := strings.TrimSuffix(filepath.Base(originalName), ext)
	if base == "" {
		base = "upload"
	}
	filename := fmt.Sprintf("%s-%d%s", slugify(base), time.Now().UnixNano(), ext)
	dstPath := filepath.Join(uploadDir, filename)

	dst, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", err
	}

	// Public URL
	return "/static/images/uploads/" + filename, nil
}

func slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "_", "-")
	// Remove characters that are not letters, numbers, or dash
	var b strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			b.WriteRune(r)
		}
	}
	res := b.String()
	if res == "" {
		return "file"
	}
	return res
}

