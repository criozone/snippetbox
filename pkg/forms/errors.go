package forms

type errors map[string][]string

func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

func (e errors) Get(field string) string {
	s := e[field]
	if len(s) == 0 {
		return ""
	}

	return s[0]
}
