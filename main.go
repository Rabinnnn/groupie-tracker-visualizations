package main

import(
	"fmt"
	"net/http"
	"groupie-tracker/handlers"
	"log"
)

func main(){
	http.HandleFunc("/", handlers.IndexHandler)
	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}