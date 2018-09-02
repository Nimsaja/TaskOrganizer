package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/Nimsaja/TaskOrganizer/task"
	"github.com/Nimsaja/TaskOrganizer/taskview"
	"github.com/gorilla/mux"
)

var (
	server = httptest.NewServer(handler())
)

func init() {
	os.Setenv("RUN_IN_CLOUD", "NotSet")
}

func TestTaskList(t *testing.T) {
	r := httptest.NewRequest("GET", "http://localhost:8080/tasks", nil)
	w := httptest.NewRecorder()
	taskList(w, r)

	// check status code
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status ok (200), but is: %v", resp.StatusCode)
	}

	// results
	result := make([]taskview.TaskView, 0)
	body, _ := ioutil.ReadAll(resp.Body)
	err := json.Unmarshal(body, &result)
	if err != nil {
		t.Errorf("No err expected: %v", err)
	}

	// items to check against
	o := task.New("TestOwner")

	// check size of tasks
	if len(o.Tasks) != len(result) {
		t.Errorf("task size expected: %v and get: %v", len(o.Tasks), len(result))
	}

	// check entry 6
	if o.Tasks[6].Name != result[6].Name {
		t.Errorf("task 6 expected: %v and get: %v", o.Tasks[6], result[6])
	}

	// check next month of entry 5
	nextMonth := task.GetNextMonthForTask(o.Tasks[5], int8(time.Now().Month()))
	if nextMonth != result[5].Next {
		t.Errorf("next month for task 5 expected: %v and get: %v", nextMonth, result[5].Next)
	}

}

func TestMonthTasks(t *testing.T) {
	r := httptest.NewRequest("GET", "http://localhost:8080/tasks", nil)
	r = mux.SetURLVars(r, map[string]string{"m": "6"})
	w := httptest.NewRecorder()
	monthTasks(w, r)

	// check status code
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status ok (200), but is: %v", resp.StatusCode)
	}

	// results
	tasks := make([]task.Task, 0)
	body, _ := ioutil.ReadAll(resp.Body)
	err := json.Unmarshal(body, &tasks)
	if err != nil {
		t.Errorf("No err expected: %v", err)
	}

	// check size of tasks
	if len(tasks) != 4 {
		t.Errorf("task size expected: %v and get: %v", 4, len(tasks))
	}
}
