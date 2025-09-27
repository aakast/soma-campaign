package main

import (
    "context"
    "log"
    "os"
    "sort"
    "strings"
    "time"

    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/compress"
    "github.com/gofiber/fiber/v2/middleware/basicauth"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/gofiber/fiber/v2/middleware/logger"
    "github.com/gofiber/fiber/v2/middleware/recover"
    "github.com/gofiber/template/html/v2"
    "github.com/joho/godotenv"

    "github.com/gofiber/adaptor/v2"
    "github.com/dracory/entitystore"
    "github.com/gouniverse/cms"
    _ "modernc.org/sqlite"
    "net/http"
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

    // Register shortcode to read settings inside CMS pages: [setting key="..." default="..."]
    cmsInstance.ShortcodeAdd(SettingShortcode{cms: cmsInstance})

    // Routes powered by CMS data, rendered with existing templates
    app.Get("/", func(c *fiber.Ctx) error {
        posts := getLatestCmsPosts(cmsInstance, 3)

        heroTitle := getSettingString(cmsInstance, "home_hero_title", "Sammen skaber vi et grønnere og mere inkluderende Fredensborg")
        heroSubtitle := getSettingString(cmsInstance, "home_hero_subtitle", "Stem på Soma Mayel - Radikale Venstre")
        heroVideoURL := getSettingString(cmsInstance, "home_hero_video_url", "/static/videos/hero-video.mp4")
        facebookURL := getSettingString(cmsInstance, "facebook_page_url", "https://www.facebook.com/somamayel")

        return c.Render("home", fiber.Map{
            "Title":    getSettingString(cmsInstance, "site_title", "Soma Mayel - Spidskandidat for Radikale Venstre i Fredensborg"),
            "Posts":    posts,
            "Featured": filterFeatured(posts),
            "Hero": fiber.Map{
                "VideoURL": heroVideoURL,
                "Title":    heroTitle,
                "Subtitle": heroSubtitle,
            },
            "FacebookPageURL": facebookURL,
        })
    })

    app.Get("/nyheder", func(c *fiber.Ctx) error {
        posts := getAllCmsPosts(cmsInstance)
        return c.Render("news", fiber.Map{
            "Title": getSettingString(cmsInstance, "news_title", "Nyheder og Blog"),
            "Posts": posts,
        })
    })

    app.Get("/blog/:slug", func(c *fiber.Ctx) error {
        slug := c.Params("slug")
        post := getCmsPostBySlug(cmsInstance, slug)
        if post == nil {
            return c.Status(404).Render("404", fiber.Map{
                "Title": "Side ikke fundet",
            })
        }
        return c.Render("blog-post", fiber.Map{
            "Title": post.Title,
            "Post":  post,
        })
    })

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

    // Legacy JSON admin API removed in favor of CMS

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

// == CMS-backed helpers =========================================================

type CmsPost struct {
    ID         string
    Title      string
    Slug       string
    Content    string
    Excerpt    string
    Author     string
    Date       time.Time
    Image      string
    Tags       []string
    IsFeatured bool
}

func getAllCmsPosts(c *cms.Cms) []CmsPost {
    if c == nil || c.EntityStore == nil {
        return []CmsPost{}
    }
    list, err := c.EntityStore.EntityList(entitystore.EntityQueryOptions{EntityType: "post"})
    if err != nil {
        return []CmsPost{}
    }
    posts := make([]CmsPost, 0, len(list))
    for _, e := range list {
        title, _ := e.GetString("title", "")
        slug, _ := e.GetString("slug", "")
        excerpt, _ := e.GetString("excerpt", "")
        content, _ := e.GetString("content", "")
        author, _ := e.GetString("author", "")
        image, _ := e.GetString("image", "")
        dateStr, _ := e.GetString("date", "")
        tagsStr, _ := e.GetString("tags", "")
        featuredStr, _ := e.GetString("is_featured", "false")

        t := time.Now()
        if dateStr != "" {
            if tt, err := time.Parse(time.RFC3339, dateStr); err == nil {
                t = tt
            } else if tt2, err2 := time.Parse("2006-01-02", dateStr); err2 == nil {
                t = tt2
            }
        }
        tags := []string{}
        if tagsStr != "" {
            for _, s := range strings.Split(tagsStr, ",") {
                s = strings.TrimSpace(s)
                if s != "" {
                    tags = append(tags, s)
                }
            }
        }
        isFeatured := strings.EqualFold(strings.TrimSpace(featuredStr), "true")

        posts = append(posts, CmsPost{
            ID:         e.ID(),
            Title:      title,
            Slug:       slug,
            Content:    content,
            Excerpt:    excerpt,
            Author:     author,
            Date:       t,
            Image:      image,
            Tags:       tags,
            IsFeatured: isFeatured,
        })
    }
    sort.Slice(posts, func(i, j int) bool { return posts[i].Date.After(posts[j].Date) })
    return posts
}

func getLatestCmsPosts(c *cms.Cms, limit int) []CmsPost {
    all := getAllCmsPosts(c)
    if limit > 0 && len(all) > limit {
        return all[:limit]
    }
    return all
}

func getCmsPostBySlug(c *cms.Cms, slug string) *CmsPost {
    list := getAllCmsPosts(c)
    for i := range list {
        if list[i].Slug == slug {
            return &list[i]
        }
    }
    return nil
}

func filterFeatured(list []CmsPost) []CmsPost {
    featured := []CmsPost{}
    for _, p := range list {
        if p.IsFeatured {
            featured = append(featured, p)
            if len(featured) >= 3 {
                break
            }
        }
    }
    return featured
}

func cmsPostQueryOptions() (opts interface{}) {
    // The entity store uses its own query options type, but EntityList with an empty
    // options yields all entities of a type when we ensure later filtering. To keep
    // the code simple, we fetch all posts and sort.
    return struct{}{}
}

func getSettingString(c *cms.Cms, key string, def string) string {
    if c == nil || c.SettingStore == nil {
        return def
    }
    val, err := c.SettingStore.Get(context.Background(), key, def)
    if err != nil || val == "" {
        return def
    }
    return val
}

// == Shortcodes ================================================================

type SettingShortcode struct{ cms *cms.Cms }

func (s SettingShortcode) Alias() string       { return "setting" }
func (s SettingShortcode) Description() string { return "Render a setting value by key" }
func (s SettingShortcode) Render(r *http.Request, raw string, m map[string]string) string {
    key := strings.TrimSpace(m["key"])
    def := m["default"]
    if key == "" || s.cms == nil || s.cms.SettingStore == nil {
        return def
    }
    v, err := s.cms.SettingStore.Get(r.Context(), key, def)
    if err != nil {
        return def
    }
    return v
}