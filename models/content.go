package models

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Content struct {
	ID         string                 `json:"id"`
	Collection string                 `json:"collection"`
	Data       map[string]interface{} `json:"data"`
}

func SaveContent(content Content) error {
	contentDir := filepath.Join("./content", content.Collection)
	
	// Create directory if it doesn't exist
	if err := os.MkdirAll(contentDir, 0755); err != nil {
		return err
	}
	
	// Marshal content to JSON
	data, err := json.MarshalIndent(content, "", "  ")
	if err != nil {
		return err
	}
	
	// Write to file
	filename := filepath.Join(contentDir, content.ID+".json")
	return ioutil.WriteFile(filename, data, 0644)
}

func GetContentByCollection(collection string) ([]Content, error) {
	contentDir := filepath.Join("./content", collection)
	
	if _, err := os.Stat(contentDir); os.IsNotExist(err) {
		return []Content{}, nil
	}
	
	files, err := ioutil.ReadDir(contentDir)
	if err != nil {
		return nil, err
	}
	
	var contents []Content
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			data, err := ioutil.ReadFile(filepath.Join(contentDir, file.Name()))
			if err != nil {
				continue
			}
			
			var content Content
			if err := json.Unmarshal(data, &content); err != nil {
				continue
			}
			
			contents = append(contents, content)
		}
	}
	
	return contents, nil
}