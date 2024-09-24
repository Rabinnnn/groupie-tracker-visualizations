package handlers

import (
	"groupie-tracker/api"
	"net/http"
	"text/template"
)

func DetailsHandler(w http.ResponseWriter, r *http.Request){
	if r.Method != "GET" {
		renderErrorPage(w, "Method Not Allowed!", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	isValid := CheckId(id)

	if !isValid{
		renderErrorPage(w, "Not Found!", http.StatusNotFound)
		return
	}

	data, err := api.GetAllDetails(id)
	if err != nil{
		renderErrorPage(w, "Internal Server Error!", http.StatusInternalServerError)
		return
	}

	temp, err := template.ParseFiles("static/templates/details.html")
	if err !=nil {
		renderErrorPage(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = temp.Execute(w, data)
	if err != nil{
		renderErrorPage(w, "Internal Server error", http.StatusInternalServerError)
		return
	}

}