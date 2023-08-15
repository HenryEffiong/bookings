package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/henryeffiong/bookings/internal/config"
	"github.com/henryeffiong/bookings/internal/handlers"
	"github.com/henryeffiong/bookings/internal/models"
	"github.com/henryeffiong/bookings/internal/render"
)

const portNumber = ":8080"

var appConfig config.AppConfig
var session *scs.SessionManager

func main() {
	errr := run()

	if errr != nil {
		log.Fatal(errr)
	}

	fmt.Printf("Starting application on port %s", portNumber)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&appConfig),
	}
	errr = srv.ListenAndServe()
	log.Fatal(errr)
}

func run() error {
	appConfig.InProduction = false
	gob.Register(models.Reservation{})
	session = scs.New()
	session.Lifetime = 24 * time.Hour

	session.Cookie.Persist = false
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = appConfig.InProduction

	appConfig.Session = session

	templateCache, errr := render.CreateTemplateCache()
	if errr != nil {
		fmt.Println("Error creating template cache: ", errr)
		log.Fatal(errr)
		return errr
	}

	appConfig.TemplateCache = templateCache
	appConfig.UseCache = false

	repo := handlers.NewRepo(&appConfig)
	handlers.NewHandlers(repo)

	render.NewTemplate(&appConfig)

	return nil
}
