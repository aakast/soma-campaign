package models

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// GetPostByID returns a post by its ID by scanning JSON files
func GetPostByID(id string) *Post {
	posts := GetAllPosts()
	for _, post := range posts {
		if post.ID == id {
			return &post
		}
	}
	return nil
}

// SavePost writes or updates a post JSON file under content/posts
func SavePost(post *Post) error {
	if strings.TrimSpace(post.ID) == "" {
		candidate := strings.TrimSpace(post.Slug)
		if candidate == "" {
			candidate = strings.ReplaceAll(strings.ToLower(post.Title), " ", "-")
		}
		post.ID = sanitizeID(candidate)
	}

	contentDir := "./content/posts"
	if err := os.MkdirAll(contentDir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(post, "", "  ")
	if err != nil {
		return err
	}

	filename := filepath.Join(contentDir, post.ID+".json")
	return ioutil.WriteFile(filename, data, 0644)
}

// DeletePost removes a post JSON file by ID
func DeletePost(id string) error {
	contentDir := "./content/posts"
	filename := filepath.Join(contentDir, id+".json")
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil
	}
	return os.Remove(filename)
}

func sanitizeID(s string) string {
	s = strings.ToLower(s)
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "_", "-")
	var b strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			b.WriteRune(r)
		}
	}
	res := b.String()
	if res == "" {
		return "post"
	}
	return res
}

