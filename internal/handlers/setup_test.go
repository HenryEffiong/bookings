package handlers

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/henryeffiong/bookings/internal/config"
	"github.com/henryeffiong/bookings/internal/models"
	"github.com/henryeffiong/bookings/internal/render"
	"github.com/justinas/nosurf"
)

var appConfig config.AppConfig
var session *scs.SessionManager
var pathToTemplates = "./../../templates"
var functions = template.FuncMap{}

func getRoutes() http.Handler {
	appConfig.InProduction = false
	gob.Register(models.Reservation{})
	session = scs.New()
	session.Lifetime = 24 * time.Hour

	session.Cookie.Persist = false
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = appConfig.InProduction

	appConfig.Session = session

	templateCache, errr := CreateTestTemplateCache()
	if errr != nil {
		fmt.Println("Error creating template cache: ", errr)
		log.Fatal(errr)

	}

	appConfig.TemplateCache = templateCache
	appConfig.UseCache = true

	repo := NewRepo(&appConfig)
	NewHandlers(repo)

	render.NewTemplate(&appConfig)

	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	// mux.Use(WriteToConsole)
	// mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", Repo.Home)
	mux.Get("/about", Repo.About)
	mux.Get("/generals", Repo.Generals)
	mux.Get("/majors", Repo.Majors)
	mux.Get("/search-availability", Repo.Availability)
	mux.Post("/search-availability", Repo.PostAvailability)
	mux.Post("/search-availability-json", Repo.AvailabilityJSON)
	mux.Get("/contact", Repo.Contact)
	mux.Get("/make-reservation", Repo.Reservation)
	mux.Post("/make-reservation", Repo.PostReservation)
	mux.Get("/reservation-summary", Repo.ReservationSummary)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return mux

}

func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   appConfig.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

// Loads and saves the session on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

func CreateTestTemplateCache() (map[string]*template.Template, error) {
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
