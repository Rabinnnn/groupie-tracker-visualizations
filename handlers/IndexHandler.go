package handlers

import (
	"html/template"
	"net/http"
	"groupie-tracker/api"
)

// handle requests for index page (home page)
func IndexHandler(w http.ResponseWriter, r *http.Request){
	if r.Method != "GET"{
		renderErrorPage(w, "Method Not Allowed!", http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/"{
		renderErrorPage(w, "Page Not Found!", http.StatusNotFound)
		return
	}

	artists, err := api.GetArtists()
	if err != nil{
		renderErrorPage(w, "Internal Server Error!", http.StatusInternalServerError)
		return
	}

	temp, err := template.ParseFiles("static/templates/index.html")
	if err != nil{
		renderErrorPage(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = temp.Execute(w, artists)
	if err != nil{
		renderErrorPage(w, "Internal Server Error", http.StatusInternalServerError)
		return 
	}
}

