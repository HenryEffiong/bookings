package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/henryeffiong/bookings/internal/config"
	"github.com/henryeffiong/bookings/internal/forms"
	"github.com/henryeffiong/bookings/internal/models"
	"github.com/henryeffiong/bookings/internal/render"
)

var Repo *Respository

type Respository struct {
	App *config.AppConfig
}

func NewRepo(pointerToAppConfig *config.AppConfig) *Respository {
	return &Respository{
		App: pointerToAppConfig,
	}
}

func NewHandlers(pointerToRepository *Respository) {
	Repo = pointerToRepository
}

func (pointerToRepository *Respository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	pointerToRepository.App.Session.Put(r.Context(), "remote_ip", remoteIP)
	render.RenderTemplate(w, "home.page.tmpl", &models.TemplateData{}, r)
}

func (pointerToRepository *Respository) About(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello World, New Data"

	// pointerToRepository.App.Session

	remoteIP := pointerToRepository.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP
	render.RenderTemplate(w, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	}, r)

}

func (pointerToRepository *Respository) Reservation(w http.ResponseWriter, r *http.Request) {
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	render.RenderTemplate(w, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	}, r)
}

func (pointerToRepository *Respository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Phone:     r.Form.Get("phone"),
		Email:     r.Form.Get("email"),
	}

	form := forms.New(r.PostForm)
	// form.Has("first_name", r)
	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3, r)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.RenderTemplate(w, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		}, r)
	}

	pointerToRepository.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "reservation-summary", http.StatusSeeOther)
}

func (pointerToRepository *Respository) Generals(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "generals.page.tmpl", &models.TemplateData{}, r)
}

func (pointerToRepository *Respository) Majors(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "majors.page.tmpl", &models.TemplateData{}, r)
}

func (pointerToRepository *Respository) Availability(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "search-availability.page.tmpl", &models.TemplateData{}, r)
}

func (pointerToRepository *Respository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")
	w.Write([]byte(fmt.Sprintf("Posted to search-availabilty %s & %s", start, end)))
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

func (pointerToRepository *Respository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	data := jsonResponse{
		OK:      true,
		Message: "Available",
	}
	response, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)

}

func (pointerToRepository *Respository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "contact.page.tmpl", &models.TemplateData{}, r)
}

func (pointerToRepository *Respository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, okay := pointerToRepository.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !okay {
		pointerToRepository.App.Session.Put(r.Context(), "error", "Unable to retrieve reservation.")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	pointerToRepository.App.Session.Remove(r.Context(), "reservation")
	data := make(map[string]interface{})
	data["reservation"] = reservation
	render.RenderTemplate(w, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
	}, r)
}
