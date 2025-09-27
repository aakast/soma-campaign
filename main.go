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

    "github.com/gofiber/adaptor/v2"
    "github.com/gouniverse/cms"
    _ "modernc.org/sqlite"
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

    // CMS initialization (gouniverse/cms)
    dbDriver := os.Getenv("DB_DRIVER")
    if dbDriver == "" {
        dbDriver = "sqlite"
    }
    dbDsn := os.Getenv("DB_DSN")
    if dbDsn == "" {
        dbDsn = "file:data/cms.db?_pragma=journal_mode(WAL)&_pragma=busy_timeout=5000"
    }

    cmsInstance, errCms := cms.NewCms(cms.Config{
        DbDriver:            dbDriver,
        DbDsn:               dbDsn,
        Prefix:              "cms_",
        EntitiesAutomigrate: true,
        BlocksEnable:        true,
        MenusEnable:         true,
        PagesEnable:         true,
        TemplatesEnable:     true,
        WidgetsEnable:       true,
        SettingsEnable:      true,
        SettingsAutomigrate: true,
        CacheEnable:         true,
        CacheAutomigrate:    true,
        LogsEnable:          true,
        LogsAutomigrate:     true,
        SessionEnable:       true,
        SessionAutomigrate:  true,
        CustomEntityList: []cms.CustomEntityStructure{
            {
                Type:      "post",
                TypeLabel: "Blog Post",
                AttributeList: []cms.CustomAttributeStructure{
                    {Name: "title", Type: "string", FormControlLabel: "Title", FormControlType: "input"},
                    {Name: "slug", Type: "string", FormControlLabel: "Slug", FormControlType: "input"},
                    {Name: "excerpt", Type: "string", FormControlLabel: "Excerpt", FormControlType: "textarea"},
                    {Name: "content", Type: "string", FormControlLabel: "Content", FormControlType: "textarea"},
                    {Name: "author", Type: "string", FormControlLabel: "Author", FormControlType: "input"},
                    {Name: "image", Type: "string", FormControlLabel: "Image URL", FormControlType: "input"},
                    {Name: "date", Type: "string", FormControlLabel: "Date", FormControlType: "input"},
                    {Name: "tags", Type: "string", FormControlLabel: "Tags (comma-separated)", FormControlType: "input"},
                    {Name: "is_featured", Type: "string", FormControlLabel: "Featured (true/false)", FormControlType: "input"},
                },
            },
        },
    })
    if errCms != nil {
        log.Fatalf("Failed to initialize CMS: %v", errCms)
    }

	// Routes
	app.Get("/", handlers.Home)
	app.Get("/om-soma", handlers.About)
	app.Get("/politik", handlers.Politics)
	app.Get("/nyheder", handlers.News)
	app.Get("/kontakt", handlers.Contact)
	app.Get("/blog/:slug", handlers.BlogPost)
	app.Get("/api/facebook/feed", handlers.FacebookFeed)

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
    // Redirect legacy admin root to CMS admin
    adminUI.Get("/", func(c *fiber.Ctx) error { return c.Redirect("/admin/cms", fiber.StatusFound) })
    // Mount CMS Admin under /admin/cms
    adminUI.All("/cms", adaptor.HTTPHandlerFunc(cmsInstance.Router))
    adminUI.All("/cms/*", adaptor.HTTPHandlerFunc(cmsInstance.Router))

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
    // Catch-all to CMS frontend (after all other routes and statics)
    app.All("/*", adaptor.HTTPHandlerFunc(cmsInstance.FrontendHandler))
    log.Fatal(app.Listen(":" + port))
}