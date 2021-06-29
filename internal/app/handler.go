package app

import (
	"github.com/Aegon95/home24-webscraper/internal"
	forms "github.com/Aegon95/home24-webscraper/pkg/form"
	"go.uber.org/zap"
	"net/http"
)

type MainHandler interface {
	Home(w http.ResponseWriter, r *http.Request)
	Results(w http.ResponseWriter, r *http.Request)
	Submit(w http.ResponseWriter, r *http.Request)
}

func NewMainHandler(logger *zap.SugaredLogger, analyzerService AnalyzeWebService, helper internal.Helper, data *internal.TemplateData) MainHandler {
	return &handler{
		logger,
		analyzerService,
		helper,
		data,
	}
}

type handler struct {
	logger          *zap.SugaredLogger
	analyzerService AnalyzeWebService
	helper          internal.Helper
	data            *internal.TemplateData
}

// home handler runs when route "/" is hit
func (h *handler) Home(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("Entering Home handler")
	// Use the new render helper.
	h.helper.Render(w, r, "home.page.tmpl", &internal.TemplateData{
		Form: forms.New(nil),
	})
}

// Results handler runs when route "/results" is hit
func (h *handler) Results(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("Entering results handler")
	h.helper.Render(w, r, "show.page.tmpl", h.data)
}

// Submit handler runs when route "/submit" is hit
func (h *handler) Submit(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("Entering submit handler")
	err := r.ParseForm()
	if err != nil {
		h.logger.Errorf("Error occurred while parsing the forms - %s", err.Error())
		h.helper.ClientError(w, http.StatusBadRequest)
		return
	}

	// Validate the form contents
	form := forms.New(r.PostForm)
	form.Required("url")
	form.IsValidUrl("url")

	// If there are any errors, return the home page with errors
	if !form.Valid() {
		h.helper.Render(w, r, "home.page.tmpl", &internal.TemplateData{Form: form})
		return
	}
	//temporary cache for form data
	h.data = nil

	// scraping the URl and getting the results
	url := form.Get("url")
	webCount, err := h.analyzerService.Scraper(url)

	// If there are any errors, return the home page with errors
	if err != nil {
		h.logger.Errorf("Error occurred while scraping the website - %s", err.Error())
		form.Errors.Add("url", err.Error())
		h.data = &internal.TemplateData{
			Form:  form,
			Stats: webCount,
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	//temporary cache for form data
	h.data = &internal.TemplateData{
		Form:  form,
		Stats: webCount,
	}

	http.Redirect(w, r, "/results", http.StatusSeeOther)

}
