package forms

import (
	"net/url"
	"strings"
)

// Form contains a url.Values object to hold the form data and an Errors field to hold any validation errors
type Form struct {
	url.Values
	Errors errors
}

// New Creates new form
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Required field checks if a field in form is empty or not
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

// IsValidUrl checks if a field containing URL is valid or not
func (f *Form) IsValidUrl(field string) {
	value := strings.TrimSpace(f.Get(field))
	u, err := url.ParseRequestURI(value)
	if err != nil || u.Scheme == "" || u.Host == "" {
		f.Errors.Add(field, "URL is invalid")
	}
}

// Valid function returns if errors are present in form or not.
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
