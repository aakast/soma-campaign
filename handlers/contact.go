package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func Contact(c *fiber.Ctx) error {
	return c.Render("contact", fiber.Map{
		"Title": "Kontakt Soma Mayel",
		"ContactInfo": fiber.Map{
			"Email":    "soma@radikale-fredensborg.dk",
			"Phone":    "+45 XX XX XX XX",
			"Facebook": "https://www.facebook.com/somamayel",
			"Address":  "Kokkedal, Fredensborg Kommune",
		},
	})
}