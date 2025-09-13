package handlers

import (
	"soma-mayel-campaign/models"

	"github.com/gofiber/fiber/v2"
)

func News(c *fiber.Ctx) error {
	posts := models.GetAllPosts()
	
	return c.Render("news", fiber.Map{
		"Title": "Nyheder og Blog",
		"Posts": posts,
	})
}

func BlogPost(c *fiber.Ctx) error {
	slug := c.Params("slug")
	post := models.GetPostBySlug(slug)
	
	if post == nil {
		return c.Status(404).Render("404", fiber.Map{
			"Title": "Side ikke fundet",
		})
	}
	
	return c.Render("blog-post", fiber.Map{
		"Title": post.Title,
		"Post":  post,
	})
}