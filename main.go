package main

import (
	"flag"
	"fmt"
	"groupie-tracker/fileio"
	"groupie-tracker/filter"
	"groupie-tracker/handlers"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
)

var port = flag.Int("P", 8080, "port to listen on")
var open = flag.Bool("O", false, "whether to open page in default browser")

// openBrowser function opens a URL in the default web browser based on the operating
// system that the code is running on. It handles Linux, Windows,and macOS platforms.
// It takes a single parameter which is a string representing the URL to open.
func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll, FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	if err != nil {
		log.Printf("Failed to open browser:%v", err)
	}
}

func main() {
	// parse the defined command-line flags
	flag.Parse()
	// configure file logging to temporary application logger file
	{
		logFilePath := path.Join(os.TempDir(), fmt.Sprintf("%d-groupie-logger.log", os.Getpid()))
		logger, err := os.OpenFile(logFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			log.Printf("failed to setup file logging: logging to stderr instead: %v\n", err)
		}
		mw := io.MultiWriter(os.Stdout, logger)
		log.SetOutput(mw)
		defer fileio.Close(logger)
	}

	http.HandleFunc("/", handlers.IndexHandler)
	http.HandleFunc("/details", handlers.DetailsHandler)
	http.HandleFunc("/search-suggestions", handlers.SearchHandler)
	http.HandleFunc("/filter", handlers.Filter)
	http.HandleFunc("/api/filter", filter.API)

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
			case "/static", "/static/css", "/static/fonts", "/static/gifs", "/static/images", "/static/js":
				handlers.RenderErrorPage(w, "Bad Request", http.StatusBadRequest)
				return
			}
			staticDirFileServer.ServeHTTP(w, r)
		},
	)

	servePort := fmt.Sprintf(":%d", *port)
	url := fmt.Sprintf("http://localhost%s\n", servePort)
	fmt.Printf("Server running at %s\n", url)

	if *open {
		openBrowser(url)
	}

	log.Fatal(http.ListenAndServe(servePort, nil))
}
