package rendering

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"path"
)

type HtmlRenderer struct {
	templatesDir string
	cache        map[string]*template.Template
}

func NewHtmlRenderer(templatesDir string) *HtmlRenderer {
	return &HtmlRenderer{
		templatesDir: templatesDir,
		cache:        make(map[string]*template.Template),
	}
}

func (r *HtmlRenderer) Render(w http.ResponseWriter, name string, data interface{}) error {
	tmpl, ok := r.cache[name]
	var err error
	if !ok {
		tmplPath := path.Join(r.templatesDir, name)
		tmpl, err = template.ParseFiles(tmplPath)
		if err != nil {
			return err
		}
		r.cache[name] = tmpl
	}

	if tmpl == nil {
		return errors.New("template not found")
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		return err
	}

	return nil
}

func (r *HtmlRenderer) RenderWithLayout(w http.ResponseWriter, layoutName, contentName string, data interface{}) error {
	// Generate cache key for the combined template
	cacheKey := layoutName + ":" + contentName

	// Check if the combined template is already cached
	tmpl, ok := r.cache[cacheKey]
	var err error

	if !ok {
		// Get paths for both templates
		layoutPath := path.Join(r.templatesDir, layoutName)
		contentPath := path.Join(r.templatesDir, contentName)

		// First parse the layout template
		tmpl, err = template.ParseFiles(layoutPath, contentPath)
		if err != nil {
			return fmt.Errorf("error parsing layout template: %w", err)
		}

		// Cache the combined template
		r.cache[cacheKey] = tmpl
	}

	if tmpl == nil {
		return errors.New("template not found")
	}

	// Execute the template, which will automatically handle the {{ .Content }} block
	err = tmpl.ExecuteTemplate(w, "layout", data)
	if err != nil {
		return fmt.Errorf("error executing template: %w", err)
	}

	return nil
}
