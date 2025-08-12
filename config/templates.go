package config

import (
	"html/template"
	"log"
)

var Templates *template.Template

func LoadTemplates() {
	var err error
	Templates, err = template.ParseGlob("internal/view/*.html")
	if err != nil {
		log.Fatalf("Error parsing templates: %v", err)
	}
}
