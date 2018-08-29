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
	url = "/tasks/{m:[0-9]+}"
	router.HandleFunc(url, monthTasks).Methods("GET")
	router.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
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
	//Try to read in stored task, if not possible (because the app is openen
	//for the first time e.g.) read in the default task list
	tasks := task.GetDefaultList()

	//if currentMonth is different to the saved actMonth recalculate the nextMonth property
	if thisCalledMonth != lastCalledMonth {
		var tasksPointers []*task.Task
		for i := 0; i < len(tasks); i++ {
			tasksPointers = append(tasksPointers, &tasks[i])
		}
		task.RecalculateNextMonthProp(tasksPointers, thisCalledMonth)

		lastCalledMonth = thisCalledMonth

		writeOutAsJSON(w, tasksPointers)
	} else {
		writeOutAsJSON(w, tasks)
	}
}

func monthTasks(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	varMonth := vars["m"]
	month, err := strconv.Atoi(varMonth)
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid id: %v in params: %v", vars, month), http.StatusBadRequest)
		return
	}

	// task.GetTaskForMonth(month)
}

func writeOutAsJSON(w http.ResponseWriter, v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%s\n", string(b))
}
