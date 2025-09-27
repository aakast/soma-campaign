package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

type facebookPost struct {
	ID           string    `json:"id"`
	Message      string    `json:"message,omitempty"`
	CreatedTime  time.Time `json:"created_time"`
	PermalinkURL string    `json:"permalink_url,omitempty"`
	FullPicture  string    `json:"full_picture,omitempty"`
}

type facebookFeedResponse struct {
	Posts          []facebookPost `json:"posts"`
	Cached         bool           `json:"cached"`
	Source         string         `json:"source,omitempty"`
	FallbackReason string         `json:"fallback_reason,omitempty"`
}

type graphPostsResponse struct {
	Data []struct {
		ID           string `json:"id"`
		Message      string `json:"message"`
		CreatedTime  string `json:"created_time"`
		PermalinkURL string `json:"permalink_url"`
		FullPicture  string `json:"full_picture"`
	} `json:"data"`
}

var fbCache struct {
	mu        sync.Mutex
	data      facebookFeedResponse
	expiresAt time.Time
}

func FacebookFeed(c *fiber.Ctx) error {
	// Serve from cache if valid
	fbCache.mu.Lock()
	if time.Now().Before(fbCache.expiresAt) {
		cached := fbCache.data
		cached.Cached = true
		fbCache.mu.Unlock()
		return c.JSON(cached)
	}
	fbCache.mu.Unlock()

	pageID := os.Getenv("FACEBOOK_PAGE_ID")
	if pageID == "" {
		pageID = os.Getenv("FACEBOOK_PAGE_USERNAME")
	}
	if pageID == "" {
		pageID = "SomamayelRV"
	}

	accessToken := os.Getenv("FACEBOOK_ACCESS_TOKEN")
	if accessToken == "" {
		// Graceful fallback if no token configured
		resp := facebookFeedResponse{
			Posts:          []facebookPost{},
			Cached:         false,
			Source:         "fallback",
			FallbackReason: "missing_access_token",
		}
		// Short cache to avoid hammering
		fbCache.mu.Lock()
		fbCache.data = resp
		fbCache.expiresAt = time.Now().Add(5 * time.Minute)
		fbCache.mu.Unlock()
		return c.JSON(resp)
	}

	graphVersion := os.Getenv("FACEBOOK_GRAPH_VERSION")
	if graphVersion == "" {
		graphVersion = "v18.0"
	}

	limit := os.Getenv("FACEBOOK_FEED_LIMIT")
	if limit == "" {
		limit = "5"
	}

	// Build Graph API URL
	fields := "message,created_time,permalink_url,full_picture"
	url := fmt.Sprintf("https://graph.facebook.com/%s/%s/posts?fields=%s&limit=%s&access_token=%s",
		graphVersion, pageID, fields, limit, accessToken,
	)

	client := &http.Client{Timeout: 10 * time.Second}
    res, err := client.Get(url)
    if err != nil {
        return apiFallbackWithCache(c, err)
	}
	defer res.Body.Close()

    if res.StatusCode < 200 || res.StatusCode >= 300 {
        return apiFallbackWithCache(c, fmt.Errorf("facebook api status %d", res.StatusCode))
	}

	var gp graphPostsResponse
    if err := json.NewDecoder(res.Body).Decode(&gp); err != nil {
        return apiFallbackWithCache(c, err)
	}

	posts := make([]facebookPost, 0, len(gp.Data))
	for _, p := range gp.Data {
		// Parse created_time (RFC3339)
		createdAt, err := time.Parse(time.RFC3339, p.CreatedTime)
		if err != nil {
			// If parse fails, keep zero time
		}
		posts = append(posts, facebookPost{
			ID:           p.ID,
			Message:      p.Message,
			CreatedTime:  createdAt,
			PermalinkURL: p.PermalinkURL,
			FullPicture:  p.FullPicture,
		})
	}

	resp := facebookFeedResponse{
		Posts:  posts,
		Cached: false,
		Source: "facebook_graph",
	}

	// Cache result
	fbCache.mu.Lock()
	fbCache.data = resp
	fbCache.expiresAt = time.Now().Add(5 * time.Minute)
	fbCache.mu.Unlock()

	return c.JSON(resp)
}

func apiFallbackWithCache(c *fiber.Ctx, err error) error {
	resp := facebookFeedResponse{
		Posts:          []facebookPost{},
		Cached:         false,
		Source:         "fallback",
		FallbackReason: err.Error(),
	}
	fbCache.mu.Lock()
	fbCache.data = resp
	fbCache.expiresAt = time.Now().Add(2 * time.Minute)
	fbCache.mu.Unlock()
    return c.JSON(resp)
}

