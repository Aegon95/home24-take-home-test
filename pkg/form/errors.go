package forms

// errors type, holds the validation errors
// The name of the form field will be used as the key in this map.
type errors map[string][]string

// Add adds new error messages for a given field to the map.
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

// Get retrieves error messages for a given field to the map.
func (e errors) Get(field string) string {
	es := e[field]
	if len(es) == 0 {
		return ""
	}
	return es[0]
}
