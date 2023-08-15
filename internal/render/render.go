package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/henryeffiong/bookings/internal/config"
	"github.com/henryeffiong/bookings/internal/models"
	"github.com/justinas/nosurf"
)

var app *config.AppConfig
var pathToTemplates = "./templates"

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.FlashMessage = app.Session.PopString(r.Context(), "flash")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.Error = app.Session.PopString(r.Context(), "error")

	td.CSRF = nosurf.Token(r)
	return td
}

func NewTemplate(pointerToAppConfig *config.AppConfig) {
	app = pointerToAppConfig
}

var functions = template.FuncMap{}

func RenderTemplate(w http.ResponseWriter, tmpl string, templateData *models.TemplateData, r *http.Request) error {
	var templateCache map[string]*template.Template
	if app.UseCache {
		templateCache = app.TemplateCache
	} else {
		templateCache, _ = CreateTemplateCache()
	}

	template, okay := templateCache[tmpl]
	if !okay {
		fmt.Println("Not okay")
		return errors.New("could not get template from template cache")
	}
	templateData = AddDefaultData(templateData, r)
	buf := new(bytes.Buffer)

	_ = template.Execute(buf, templateData)
	_, errr := buf.WriteTo(w)
	if errr != nil {
		fmt.Println("Error writing to buffer: ", errr)
		return errr
	}

	// parsedTemplate, _ := template.ParseFiles("./templates/" + tmpl)
	// err := parsedTemplate.Execute(w, nil)
	// if err != nil {
	// 	fmt.Println("Error parsingg template:", err)
	// }
	return nil
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// Get all the page files in the template folder i.e everything with .page.tmpl
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))

	if err != nil {
		fmt.Println("Error getting templates: ", err)
	}

	for _, page := range pages {
		// Get the name of the file e.g. "about.page.tmpl"
		name := filepath.Base(page)

		// Create the pointer to a template with the name and associate it with the page
		templateSet, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			fmt.Println("Error parsing templateSet: ", err)
		}

		matches, errr := filepath.Glob(fmt.Sprintf("%s/*layout.tmpl", pathToTemplates))
		if errr != nil {
			fmt.Println("Error parsing matches: ", errr)
		}

		// Check if we have any layout files and parse with the layout
		if len(matches) > 0 {
			templateSet, err = templateSet.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				fmt.Println("Error parsing matches: ", err)
			}
		}
		myCache[name] = templateSet
	}

	// fmt.Println("myCache: ", myCache)

	return myCache, nil
}
