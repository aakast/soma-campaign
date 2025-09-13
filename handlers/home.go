package handlers

import (
	"soma-mayel-campaign/models"

	"github.com/gofiber/fiber/v2"
)

func Home(c *fiber.Ctx) error {
	// Get latest blog posts
	posts := models.GetLatestPosts(3)

	// Get featured content
	featured := models.GetFeaturedContent()

	return c.Render("home", fiber.Map{
		"Title":    "Soma Mayel - Spidskandidat for Radikale Venstre i Fredensborg",
		"Posts":    posts,
		"Featured": featured,
		"Hero": fiber.Map{
			"VideoURL": "/static/videos/hero-video.mp4",
			"Title":    "Sammen skaber vi et grønnere og mere inkluderende Fredensborg",
			"Subtitle": "Stem på Soma Mayel - Radikale Venstre",
		},
	})
}