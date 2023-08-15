package render

import (
	"net/http"
	"testing"

	"github.com/henryeffiong/bookings/internal/models"
)

func TestAddDefaultData(t *testing.T) {
	var td models.TemplateData
	r, err := getSession()

	if err != nil {
		t.Error(err)
	}

	session.Put(r.Context(), "flash", "123")
	result := AddDefaultData(&td, r)

	if result.FlashMessage != "123" {
		t.Error(err)
	}
}

func TestRenderTemplate(t *testing.T) {
	pathToTemplates = "./../../templates"
	templateCache, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
	app.TemplateCache = templateCache

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	var ww MyWriter

	err = RenderTemplate(&ww, "home.page.tmpl", &models.TemplateData{}, r)
	if err != nil {
		t.Error("error writing template to browser", err)
	}

	err = RenderTemplate(&ww, "non-existent.page.tmpl", &models.TemplateData{}, r)
	if err == nil {
		t.Error("rendered template that does not exist")
	}
}

func getSession() (*http.Request, error) {
	r, err := http.NewRequest("GET", "/", nil)

	if err != nil {
		return nil, err
	}

	ctx := r.Context()

	ctx, _ = session.Load(ctx, r.Header.Get("X-Session"))
	r = r.WithContext(ctx)

	return r, nil
}

func TestNewTemplates(t *testing.T) {
	NewTemplate(app)
}

func TestCreateTemplateCache(t *testing.T) {
	pathToTemplates = "./../../templates"

	_, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
}
