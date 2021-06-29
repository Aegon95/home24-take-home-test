package internal

import (
	forms "github.com/Aegon95/home24-webscraper/pkg/form"
	models2 "github.com/Aegon95/home24-webscraper/pkg/models"
	"html/template"
	"path"
	"path/filepath"
)

// TemplateData contains any dynamic data that we want to pass to our HTML templates.
type TemplateData struct {
	CurrentYear int
	Flash       string
	Stats       *models2.WebStats
	Form        *forms.Form
}

// NewTemplateCache acts as a cache for all the html templates in the directory
func NewTemplateCache(dir string) (map[string]*template.Template, error) {

	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(path.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {

		name := filepath.Base(page)

		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
