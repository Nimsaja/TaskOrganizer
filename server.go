package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Nimsaja/TaskOrganizer/task"
	"github.com/gorilla/mux"
	"google.golang.org/appengine"
)

var path = "output.txt"
var inCloud bool
var lastCalledMonth = -1
var thisCalledMonth = int(time.Now().Month())

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
	urlWithMonth := "/tasks/{m:[0-9]+}"
	router.HandleFunc(urlWithMonth, monthTasks).Methods("GET")
	router.HandleFunc(urlWithMonth, func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "invalid method: "+r.Method, http.StatusBadRequest)
	}).Methods("DELETE", "PATH", "COPY", "HEAD", "LINK", "UNLINK", "PURGE", "LOCK", "UNLOCK", "VIEW", "PROPFIND")

	return corsAndOptionHandler(router)
}

func main() {
	inCloud, _ = strconv.ParseBool(os.Getenv("RUN_IN_CLOUD"))

	http.Handle("/", handler())

	appengine.Main()
}

func taskList(w http.ResponseWriter, r *http.Request) {
	//Try to read in stored task, if not possible (because the app is opened
	//for the first time e.g.) read in the default task list
	s := r.URL.Query().Get("owner")
	o := task.New(s)

	//if currentMonth is different to the saved actMonth recalculate the nextMonth property
	if thisCalledMonth != lastCalledMonth {
		o.RecalculateNextMonthProp(thisCalledMonth)
		lastCalledMonth = thisCalledMonth

		for t := range o.Tasks {
			fmt.Println(t)
		}

		writeOutAsJSON(w, o.Tasks)
	} else {
		writeOutAsJSON(w, o.Tasks)
	}
}

func monthTasks(w http.ResponseWriter, r *http.Request) {
	s := r.URL.Query().Get("owner")

	//for now take default list TODO: create hashmap ownerName - organizer sync!
	o := task.New(s)

	vars := mux.Vars(r)
	varMonth := vars["m"]

	month, err := strconv.Atoi(varMonth)
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid month: %v in params: %v", vars, month), http.StatusBadRequest)
		return
	}

	tasks := o.GetTaskForMonth(month)

	writeOutAsJSON(w, tasks)
}

func writeOutAsJSON(w http.ResponseWriter, v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%s\n", string(b))
}
