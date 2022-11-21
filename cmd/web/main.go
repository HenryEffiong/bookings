package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/henryeffiong/bookings/pkg/config"
	"github.com/henryeffiong/bookings/pkg/handlers"
	"github.com/henryeffiong/bookings/pkg/render"
)

const portName = ":8080"

var appConfig config.AppConfig

var session *scs.SessionManager

func main() {

	appConfig.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = appConfig.InProduction

	appConfig.Session = session

	myTemplateFromCache, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	appConfig.TemplateCache = myTemplateFromCache
	appConfig.UseCache = false

	repo := handlers.NewRepo(&appConfig)
	handlers.NewHandlers(repo)
	render.NewTemplates(&appConfig)
	// 	name := "Henry"
	// 	response := fmt.Sprintf("Hello World from %s", name)
	// 	fmt.Println(response)

	handler := http.HandlerFunc(handlers.HandleRequest)

	http.Handle("/img", handler)

	// http.HandleFunc("/", handlers.Repo.Home)
	// http.HandleFunc("/about", handlers.Repo.About)
	// fmt.Printf("Starting server on port %s", portName)
	// http.ListenAndServe(":8080", nil)

	fmt.Printf("Starting server on port %s", portName)
	server := &http.Server{
		Addr:    portName,
		Handler: routes(&appConfig),
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}

// go run ./cmd/web/ .

// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 	n, err := fmt.Fprintf(w, "My First Golang Server")
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	fmt.Printf("Number of bytes: %d", n)

// })

// =================================Serving a page on the browser=======================================

// Home page handler. A handler requires a
// func Home(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "This is my Home page.")
// }

// func About(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "This is my About page and 2 + 2 is: %d", sum(2, 2))
// }

// func sum(x, y int) int {
// 	return x + y
// }

// func main() {
// 	// 	name := "Henry"
// 	// 	response := fmt.Sprintf("Hello World from %s", name)
// 	// 	fmt.Println(response)

// 	http.HandleFunc("/", Home)
// 	http.HandleFunc("/about", About)

// 	fmt.Printf("Starting server on port %s", portName)
// 	http.ListenAndServe(":8080", nil)

// }
