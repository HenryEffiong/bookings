package forms

type errors map[string][]string

func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

func (e errors) Get(field string) string {
	error := e[field]
	if len(error) == 0 {
		return ""
	}
	return error[0]
}
