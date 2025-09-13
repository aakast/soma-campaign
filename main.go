package main

import (
	"log"
	"os"
	"soma-mayel-campaign/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	godotenv.Load()

	// Create template engine
	engine := html.New("./templates", ".html")
	engine.Reload(true)
	engine.Debug(true)

	// Create fiber app with template engine
	app := fiber.New(fiber.Config{
		Views:       engine,
		ViewsLayout: "layouts/main",
	})

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(compress.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// Static files
	app.Static("/static", "./static")
	app.Static("/content", "./content")

	// Routes
	app.Get("/", handlers.Home)
	app.Get("/om-soma", handlers.About)
	app.Get("/politik", handlers.Politics)
	app.Get("/nyheder", handlers.News)
	app.Get("/kontakt", handlers.Contact)
	app.Get("/blog/:slug", handlers.BlogPost)

	// TinaCMS API routes
	app.Post("/api/tina/content", handlers.TinaContentAPI)
	app.Get("/api/tina/content/:collection", handlers.TinaGetContent)

	// Admin authentication (Basic Auth)
	adminUser := os.Getenv("ADMIN_USERNAME")
	adminPass := os.Getenv("ADMIN_PASSWORD")
	if adminUser == "" || adminPass == "" {
		adminUser = "admin"
		adminPass = "admin123"
	}
	adminAuth := basicauth.New(basicauth.Config{
		Users: map[string]string{
			adminUser: adminPass,
		},
		Realm: "Restricted",
	})

	// Admin UI
	adminUI := app.Group("/admin", adminAuth)
	adminUI.Get("/", handlers.AdminPage)

	// Admin API
	adminAPI := app.Group("/api/admin", adminAuth)
	adminAPI.Get("/posts", handlers.AdminListPosts)
	adminAPI.Get("/posts/:id", handlers.AdminGetPost)
	adminAPI.Post("/posts", handlers.AdminUpsertPost)
	adminAPI.Delete("/posts/:id", handlers.AdminDeletePost)
	adminAPI.Post("/upload", handlers.AdminUpload)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Server starting on http://localhost:%s", port)
	log.Fatal(app.Listen(":" + port))
}