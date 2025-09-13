package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func About(c *fiber.Ctx) error {
	return c.Render("about", fiber.Map{
		"Title": "Om Soma Mayel",
		"Timeline": []fiber.Map{
			{
				"Year":        "1993",
				"Title":       "Født i Afghanistan",
				"Description": "Født i Mazar-e Sharif, Afghanistan den 24. juni",
			},
			{
				"Year":        "2001",
				"Title":       "Ankomst til Danmark",
				"Description": "Kom til Danmark som 7-årig flygtning med sin familie",
			},
			{
				"Year":        "2015",
				"Title":       "Dansk statsborgerskab",
				"Description": "Modtog dansk pas og startede karriere i Udlændingestyrelsen",
			},
			{
				"Year":        "2019",
				"Title":       "Cand.jur.",
				"Description": "Afsluttede juridisk kandidatgrad ved Københavns Universitet",
			},
			{
				"Year":        "2021",
				"Title":       "Byrådsmedlem",
				"Description": "Valgt til Fredensborg Byråd med 811 personlige stemmer",
			},
			{
				"Year":        "2024",
				"Title":       "Spidskandidat",
				"Description": "Blev spidskandidat for Radikale Venstre i Fredensborg Kommune",
			},
		},
		"Roles": []fiber.Map{
			{"Title": "Formand for Fritids- og Idrætsudvalget", "Icon": "sports"},
			{"Title": "Medlem af Økonomiudvalget", "Icon": "economy"},
			{"Title": "Formand for Handicaprådet", "Icon": "accessibility"},
			{"Title": "Underviser ved Københavns Universitet", "Icon": "education"},
		},
	})
}