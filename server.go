package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// handle CORS and the OPION method
func corsAndOptionHandler(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
		} else {
			h.ServeHTTP(w, r)
		}
	}
}

// create all used Handler
func handler() http.Handler {
	router := mux.NewRouter()

	url := "/tasks"
	router.HandleFunc(url, taskList).Methods("GET")
	router.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "invalid method: "+r.Method, http.StatusBadRequest)
	}).Methods("DELETE", "PATH", "COPY", "HEAD", "LINK", "UNLINK", "PURGE", "LOCK", "UNLOCK", "VIEW", "PROPFIND")

	return corsAndOptionHandler(router)
}

func main() {
	// http.HandleFunc("/tasks", taskList)
	http.Handle("/", handler())
	log.Fatalln(http.ListenAndServe(":8080", nil))
	log.Println("Init is ready and start the server on: http://localhost:8080")
}

func taskList(w http.ResponseWriter, r *http.Request) {
	log.Println("taskList...")

	text := r.FormValue("text")

	if len(text) == 0 {
		text = "This is my first task entry! :-D"
	}

	b, err := json.Marshal(text)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%s", string(b))
}
