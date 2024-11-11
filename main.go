package main

import (
	"fmt"
	"groupie-tracker/fileio"
	"groupie-tracker/handlers"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
)

func main() {
	// configure file logging to temporary application logger file
	{
		logFilePath := path.Join(os.TempDir(), strconv.Itoa(os.Getpid()), "groupie-logger.log")
		logger, err := os.OpenFile(logFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			log.Printf("failed to setup file logging: logging to stderr instead: %v\n", err)
		}
		log.SetOutput(logger)
		defer fileio.Close(logger)
	}

	http.HandleFunc("/", handlers.IndexHandler)
	http.HandleFunc("/details", handlers.DetailsHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
