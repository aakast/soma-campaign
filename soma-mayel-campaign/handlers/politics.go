package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func Politics(c *fiber.Ctx) error {
	return c.Render("politics", fiber.Map{
		"Title": "Politik og Mærkesager",
		"PolicyAreas": []fiber.Map{
			{
				"Title":       "Børn og Uddannelse",
				"Icon":        "child_care",
				"Description": "Alle børn fortjener en god start på livet",
				"Points": []string{
					"REELLE minimumsnormeringer i daginstitutioner",
					"Mental trivsel på skoleskemaet",
					"Kvalitetsuddannelse for alle uanset baggrund",
				},
			},
			{
				"Title":       "Unge og Integration",
				"Icon":        "diversity",
				"Description": "Fredensborg skal være attraktiv for unge",
				"Points": []string{
					"Større ungdomsrepræsentation i politik",
					"Gøre kommunen attraktiv for unge at bosætte sig",
					"Bryde barrierer for unge med forskellig baggrund",
				},
			},
			{
				"Title":       "Lighed og Inklusion",
				"Icon":        "equality",
				"Description": "En kommune for alle",
				"Points": []string{
					"Lige muligheder uanset baggrund",
					"Bedre retssikkerhed for alle borgere",
					"Borgerambassadør til at hjælpe med systemnavigation",
				},
			},
			{
				"Title":       "Bæredygtighed og Miljø",
				"Icon":        "nature",
				"Description": "Grøn omstilling med ansvar",
				"Points": []string{
					"Bevæge Fredensborg i en mere bæredygtig retning",
					"Balance mellem økonomi og miljøansvar",
					"Grønne initiativer i alle kommunale områder",
				},
			},
			{
				"Title":       "Social Retfærdighed",
				"Icon":        "balance",
				"Description": "Kamp mod ulighed",
				"Points": []string{
					"Et samfund baseret på mindre uretfærdighed",
					"Politik der samler frem for at splitte",
					"Fokus på de mest udsatte borgere",
				},
			},
		},
	})
}