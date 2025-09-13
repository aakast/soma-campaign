package models

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type Post struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Slug        string    `json:"slug"`
	Content     string    `json:"content"`
	Excerpt     string    `json:"excerpt"`
	Author      string    `json:"author"`
	Date        time.Time `json:"date"`
	Image       string    `json:"image"`
	Tags        []string  `json:"tags"`
	IsFeatured  bool      `json:"is_featured"`
}

func GetLatestPosts(limit int) []Post {
	posts := GetAllPosts()
	if len(posts) > limit {
		return posts[:limit]
	}
	return posts
}

func GetAllPosts() []Post {
	var posts []Post
	
	contentDir := "./content/posts"
	if _, err := os.Stat(contentDir); os.IsNotExist(err) {
		// Return sample posts if directory doesn't exist
		return getSamplePosts()
	}
	
	files, err := ioutil.ReadDir(contentDir)
	if err != nil {
		return getSamplePosts()
	}
	
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".json") {
			data, err := ioutil.ReadFile(filepath.Join(contentDir, file.Name()))
			if err != nil {
				continue
			}
			
			var post Post
			if err := json.Unmarshal(data, &post); err != nil {
				continue
			}
			
			posts = append(posts, post)
		}
	}
	
	// Sort by date, newest first
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Date.After(posts[j].Date)
	})
	
	if len(posts) == 0 {
		return getSamplePosts()
	}
	
	return posts
}

func GetPostBySlug(slug string) *Post {
	posts := GetAllPosts()
	for _, post := range posts {
		if post.Slug == slug {
			return &post
		}
	}
	return nil
}

func GetFeaturedContent() []Post {
	posts := GetAllPosts()
	var featured []Post
	
	for _, post := range posts {
		if post.IsFeatured {
			featured = append(featured, post)
			if len(featured) >= 3 {
				break
			}
		}
	}
	
	return featured
}

func getSamplePosts() []Post {
	return []Post{
		{
			ID:         "1",
			Title:      "Sammen skaber vi fremtidens Fredensborg",
			Slug:       "sammen-skaber-vi-fremtidens-fredensborg",
			Content:    "Som jeres spidskandidat for Radikale Venstre går jeg til valg på at skabe et grønnere, mere inkluderende og stærkere Fredensborg Kommune...",
			Excerpt:    "Min vision for Fredensborg Kommune bygger på bæredygtighed, inklusion og fællesskab.",
			Author:     "Soma Mayel",
			Date:       time.Now().AddDate(0, 0, -1),
			Image:      "/static/images/campaign-launch.jpg",
			Tags:       []string{"Vision", "Fredensborg", "Valg2025"},
			IsFeatured: true,
		},
		{
			ID:         "2",
			Title:      "Børn fortjener bedre - reelle minimumsnormeringer nu",
			Slug:       "born-fortjener-bedre",
			Content:    "Det er ikke nok at have minimumsnormeringer på papir. Vi skal sikre, at de også bliver til virkelighed i vores daginstitutioner...",
			Excerpt:    "Alle børn fortjener en god start på livet med kvalitet i daginstitutionerne.",
			Author:     "Soma Mayel",
			Date:       time.Now().AddDate(0, 0, -3),
			Image:      "/static/images/children-education.jpg",
			Tags:       []string{"Børn", "Uddannelse", "Politik"},
			IsFeatured: true,
		},
		{
			ID:         "3",
			Title:      "Fra flygtning til folkevalgt - min historie",
			Slug:       "fra-flygtning-til-folkevalgt",
			Content:    "I 2001 kom jeg til Danmark som 7-årig flygtning fra Afghanistan. I dag er jeg byrådsmedlem og kæmper for at gøre vores kommune til et sted, hvor alle kan trives...",
			Excerpt:    "Min personlige rejse har givet mig en unik forståelse for vigtigheden af inklusion.",
			Author:     "Soma Mayel",
			Date:       time.Now().AddDate(0, 0, -7),
			Image:      "/static/images/soma-story.jpg",
			Tags:       []string{"Personligt", "Integration", "Historie"},
			IsFeatured: false,
		},
	}
}