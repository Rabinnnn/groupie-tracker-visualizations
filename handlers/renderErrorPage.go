package handlers

import (
	"net/http"
	"strconv"
	"html/template"
	"groupie-tracker/api"
)


func renderErrorPage(w http.ResponseWriter, errorText string, statusCode int){
	w.WriteHeader(statusCode)

	content := api.ErrorContent{
		Message : errorText,
		Code : strconv.Itoa(statusCode),
	}

	temp, err := template.ParseFiles("static/templates/errorPage.html")
	if err != nil{
		http.Error(w, "Internal Server Error!", http.StatusInternalServerError)
		return
	}

	err = temp.Execute(w, content)
	if err != nil{
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

}