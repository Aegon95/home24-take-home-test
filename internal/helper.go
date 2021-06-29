package internal

import (
	"bytes"
	"fmt"
	"go.uber.org/zap"
	"html/template"
	"net/http"
	"runtime/debug"
	"time"
)

type Helper interface {
	Render(w http.ResponseWriter, r *http.Request, name string, td *TemplateData)
	ServerError(w http.ResponseWriter, err error)
	ClientError(w http.ResponseWriter, status int)
}

func NewHelper(logger *zap.SugaredLogger, templateCache map[string]*template.Template) Helper {
	return &helper{
		logger,
		templateCache,
	}
}

type helper struct {
	logger        *zap.SugaredLogger
	templateCache map[string]*template.Template
}

// Render function reads the template name from cache and writes it in responsewritter
func (h *helper) Render(w http.ResponseWriter, r *http.Request, name string, td *TemplateData) {

	// Retrieve the appropriate template set from the cache based on the page name
	ts, ok := h.templateCache[name]
	if !ok {
		h.logger.Errorf("The template %s does not exist", name)
		h.ServerError(w, fmt.Errorf("The template %s does not exist", name))
		return
	}

	buf := new(bytes.Buffer)

	// Execute the template set, passing the dynamic data with the current
	err := ts.Execute(buf, addDefaultData(td, r))
	if err != nil {
		h.logger.Errorf("Error occurred while executing the template - %s", err.Error())
		h.ServerError(w, err)
		return
	}

	// Write the contents of the buffer to the http.ResponseWriter
	buf.WriteTo(w)
}

// ServerError helper writes an error message and stack trace to stdout, sends a generic 500 Internal Server Error response to the user
func (h *helper) ServerError(w http.ResponseWriter, err error) {

	h.logger.Debugf("%s\n%s", err.Error(), debug.Stack())

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// ClientError helper sends a specific status code and corresponding description to the user.
func (h *helper) ClientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

//Helpers

// addDefaultData helper adds default data which are need by templates
func addDefaultData(td *TemplateData, r *http.Request) *TemplateData {
	if td == nil {
		td = &TemplateData{}
	}
	td.CurrentYear = time.Now().Year()

	return td
}
