package server

import (
	"github.com/Aegon95/home24-webscraper/internal"
	"github.com/Aegon95/home24-webscraper/internal/app"
	"github.com/go-chi/chi/v5"
	"net/http"
)

// Routes Setups the Router and the middlewares
func (a *Application) Routes() http.Handler {

	r := chi.NewRouter()

	helperService := internal.NewHelper(a.Logger, a.TemplateCache)
	middlewareManager := NewMiddlewareManager(a.Logger, helperService)

	// middlewares
	r.Use(middlewareManager.RecoverPanic)
	r.Use(middlewareManager.LogRequest)
	r.Use(middlewareManager.SecureHeaders)

	analyzerServer := app.NewAnalyzeWebService(http.DefaultClient, a.Logger)
	mainHandler := app.NewMainHandler(a.Logger, analyzerServer, helperService, &internal.TemplateData{})
	// Routes
	r.Get("/", mainHandler.Home)
	r.Post("/submit", mainHandler.Submit)
	r.Get("/results", mainHandler.Results)

	// Static file server
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	r.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return r
}
