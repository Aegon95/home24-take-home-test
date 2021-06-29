package main

import (
	"github.com/Aegon95/home24-webscraper/internal"
	"github.com/Aegon95/home24-webscraper/internal/server"
	logger "github.com/Aegon95/home24-webscraper/pkg"
	"net/http"
	"time"
)

func main() {

	// Initialize zap logger
	log := logger.InitLogger()

	// Initialize a new template cache...
	templateCache, err := internal.NewTemplateCache("./ui/html/")
	if err != nil {
		log.Error(err)
	}

	// Initialize a new instance of Application containing the dependencies.
	app := &server.Application{
		Logger:        log,
		TemplateCache: templateCache,
	}

	svr := http.Server{
		Handler:      app.Routes(),
		Addr:         ":3000",
		IdleTimeout:  time.Minute,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// start the server
	log.Fatal(svr.ListenAndServe())
}
