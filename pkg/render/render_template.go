package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/henryeffiong/bookings/pkg/config"
	"github.com/henryeffiong/bookings/pkg/models"
)

var appConfigFromConfigFile *config.AppConfig

var functions = template.FuncMap{}

func NewTemplates(appConfig *config.AppConfig) {
	appConfigFromConfigFile = appConfig
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {

	return td
}

func RenderTemplate(w http.ResponseWriter, tmpl string, templateData *models.TemplateData) {

	var templateCache map[string]*template.Template
	if appConfigFromConfigFile.UseCache {
		templateCache = appConfigFromConfigFile.TemplateCache
	} else {
		var err error
		templateCache, err = CreateTemplateCache()
		if err != nil {
			log.Fatal(err)
			// fmt.Println("Error getting template cache:", err)
		}
	}

	myTemplateFromCache, okay := templateCache[tmpl]
	if !okay {
		log.Fatal("Could not get template from app config")

	}

	myBytesBuffer := new(bytes.Buffer)

	templateData = AddDefaultData(templateData)

	_ = myTemplateFromCache.Execute(myBytesBuffer, templateData)

	_, err := myBytesBuffer.WriteTo(w)
	if err != nil {
		fmt.Println("error parsing template:", err)
	}

	// parsedTemplate, _ := template.ParseFiles("./templates/" + tmpl)
	// err = parsedTemplate.Execute(w, nil)
	// if err != nil {
	// 	fmt.Println("error parsing template:", err)
	// 	return
	// }

	// template.ParseFiles
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		templateSet, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}
		matches, err := filepath.Glob("./*templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			templateSet, err = templateSet.ParseGlob("./*templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}
		myCache[name] = templateSet
	}
	return myCache, nil
}
