package main

import (
	"flag"
	"fmt"
	"groupie-tracker/fileio"
	"groupie-tracker/handlers"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

var port = flag.Int("P", 8080, "port to listen on")

func main() {
	// configure file logging to temporary application logger file
	{
		logFilePath := path.Join(os.TempDir(), fmt.Sprintf("%d-groupie-logger.log", os.Getpid()))
		logger, err := os.OpenFile(logFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			log.Printf("failed to setup file logging: logging to stderr instead: %v\n", err)
		}
		log.SetOutput(logger)
		defer fileio.Close(logger)
	}

	http.HandleFunc("/", handlers.IndexHandler)
	http.HandleFunc("/details", handlers.DetailsHandler)
	// Browsers ping for the /favicon.ico icon, redirect to the respective static file
	http.Handle("/favicon.ico", http.RedirectHandler("/static/images/favicon.svg", http.StatusMovedPermanently))
	// Server static files from the static dir, but, ensure not to expose the directory entries
	staticDirFileServer := http.StripPrefix("/static/", http.FileServer(http.Dir("static")))
	http.HandleFunc(
		"/static/", func(w http.ResponseWriter, r *http.Request) {
			// clean to remove trailing slash in path, so that the
			//paths `/static` and `/static/` both translate to `/static`
			reqPath := filepath.Clean(r.URL.Path)
			switch reqPath {
			case "/static", "/static/css", "/static/images":
				handlers.RenderErrorPage(w, "Bad Request", http.StatusBadRequest)
				return
			}
			staticDirFileServer.ServeHTTP(w, r)
		},
	)

	servePort := fmt.Sprintf(":%d", *port)
	fmt.Printf("Server running at http://localhost%s\n", servePort)
	log.Fatal(http.ListenAndServe(servePort, nil))
}
