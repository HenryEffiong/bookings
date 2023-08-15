package models

import "github.com/henryeffiong/bookings/internal/forms"

type TemplateData struct {
	StringMap    map[string]string
	IntMap       map[string]int
	FloatMap     map[string]float32
	Data         map[string]interface{}
	CSRF         string
	FlashMessage string
	Warning      string
	Error        string
	Form         *forms.Form
}
