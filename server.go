package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"google.golang.org/appengine"
)

var path = "output.txt"
var inCloud bool

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
	// inCloud, _ := strconv.ParseBool(os.Getenv("RUN_IN_CLOUD"))
	inCloud = true

	http.Handle("/", handler())

	if inCloud {
		appengine.Main()
	} else {
		log.Println("Init is ready and start the server on: http://localhost:8080")
		log.Fatalln(http.ListenAndServe(":8080", nil))
	}
}

func taskList(w http.ResponseWriter, r *http.Request) {
	text := r.FormValue("text")

	if inCloud {
		if len(text) == 0 {
			text = "Try to change the url to .../tasks?text=MyText"
		}

		b, err := json.Marshal(text)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "%s", string(b))
	} else {

		appendToOutputFile(text)

		readFile(w)
	}
}

func appendToOutputFile(t string) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Printf("Error %s ", err)
		panic(err)
	}

	defer f.Close()

	fmt.Fprintln(f, t)
}

func readFile(w http.ResponseWriter) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0600)
	if err != nil {
		log.Printf("Error %s ", err)
		panic(err)
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		b, err := json.Marshal(scanner.Text())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "%s\n", string(b))
	}
}
