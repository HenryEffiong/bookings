package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	isValid := form.Valid()

	if !isValid {
		t.Error("Got invalid when should have been valid")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("Form shows valid when required fields are missing")
	}
	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "b")
	postedData.Add("c", "c")

	r, _ = http.NewRequest("POST", "/whatever", nil)

	r.PostForm = postedData

	form = New(r.PostForm)

	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("Form shows invalid when required fields are present")
	}
}

func TestForm_IsEmail(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	if form.IsEmail("email") {
		t.Error("Validated an invalid email")
	}

	postedData := url.Values{}
	postedData.Add("email", "henry.hogan2012@gmail.com")

	r, _ = http.NewRequest("POST", "/whatever", nil)
	r.PostForm = postedData
	form = New(r.PostForm)

	if !form.IsEmail("email") {
		t.Error("Expected a valid email, but got an invalid one")
	}
}

func TestForm_MinLength(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	if !form.MinLength("email", 0, r) {
		t.Errorf("Expected %d length but got %d length characters", 0, len(form.Get("email")))
	}

	postedData := url.Values{}
	postedData.Add("email", "henry.hogan2012@gmail.com")

	r, _ = http.NewRequest("POST", "/whatever", nil)
	r.PostForm = postedData
	form = New(r.PostForm)

	if !form.MinLength("email", 3, r) {
		t.Error(form.Get("email"))
		t.Error("Expected email to be more than minimum length but is not")
	}

}
