package handlers

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/henryeffiong/bookings/pkg/config"
	"github.com/henryeffiong/bookings/pkg/models"
	"github.com/henryeffiong/bookings/pkg/render"
)

// Repo is the repository used by the handlers
var Repo *Repository

type Repository struct {
	App *config.AppConfig
}

// Creates a new repository
func NewRepo(appConfig *config.AppConfig) *Repository {
	return &Repository{
		App: appConfig,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home page handler. A handler requires a
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)
	render.RenderTemplate(w, "home.page.tmpl", &models.TemplateData{})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello from the other side"

	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	render.RenderTemplate(w, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

func HandleRequest(w http.ResponseWriter, r *http.Request) {

	buf, err := ioutil.ReadFile("sid.png")

	if err != nil {

		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Disposition", `attachment;filename="sid.png"`)

	w.Write(buf)
}
