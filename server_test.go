package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Nimsaja/TaskOrganizer/task"
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
	tasks := make([]task.Task, 0)
	body, _ := ioutil.ReadAll(resp.Body)
	err := json.Unmarshal(body, &tasks)
	if err != nil {
		t.Errorf("No err expected: %v", err)
	}

	// items to check against
	testTasks := task.GetDefaultList()

	// check size of tasks
	if len(testTasks) != len(tasks) {
		t.Errorf("task size expected: %v and get: %v", len(testTasks), len(tasks))
	}

	// check entry 6
	if testTasks[6] != tasks[6] {
		t.Errorf("task 6 expected: %v and get: %v", testTasks[6], tasks[6])
	}

}

func TestRecalcOfNextMonthIfRunTheFirstTime(t *testing.T) {
	r := httptest.NewRequest("GET", "http://localhost:8080/tasks", nil)
	w := httptest.NewRecorder()

	newTasks := make([]task.Task, 0)
	tk := task.Task{Name: "Test", Descr: "every month", Freq: 1, Start: 1, Done: false}
	newTasks = append(newTasks, tk)

	task.SetTasksList(newTasks)

	//should calculate
	thisCalledMonth = 2

	taskList(w, r)

	// check status code
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status ok (200), but is: %v", resp.StatusCode)
	}

	// results
	resultTasks := make([]task.Task, 0)
	body, _ := ioutil.ReadAll(resp.Body)
	err := json.Unmarshal(body, &resultTasks)
	if err != nil {
		t.Errorf("No err expected: %v", err)
	}

	if resultTasks[0].Next != 2 {
		t.Errorf("next month expected: %v and get: %v", 2, resultTasks[0].Next)
	}

	if lastCalledMonth != 2 {
		t.Errorf("last Called Month value in server.go expected: %v and get: %v", 2, lastCalledMonth)
	}
}
func TestRecalcOfNextMonthIfRunTheSecondTime(t *testing.T) {
	r := httptest.NewRequest("GET", "http://localhost:8080/tasks", nil)
	w := httptest.NewRecorder()

	newTasks := make([]task.Task, 0)
	tk := task.Task{Name: "Test", Descr: "every month", Freq: 1, Start: 1, Done: false}
	newTasks = append(newTasks, tk)

	task.SetTasksList(newTasks)

	//should not calculate
	thisCalledMonth = 2
	lastCalledMonth = 2

	taskList(w, r)

	// check status code
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status ok (200), but is: %v", resp.StatusCode)
	}

	// results
	resultTasks := make([]task.Task, 0)
	body, _ := ioutil.ReadAll(resp.Body)
	err := json.Unmarshal(body, &resultTasks)
	if err != nil {
		t.Errorf("No err expected: %v", err)
	}

	if resultTasks[0].Next != 0 {
		t.Errorf("next month expected: %v and get: %v", 0, resultTasks[0].Next)
	}
}

func TestRecalcOfNextMonthIfMonthChanged(t *testing.T) {
	r := httptest.NewRequest("GET", "http://localhost:8080/tasks", nil)
	w := httptest.NewRecorder()

	newTasks := make([]task.Task, 0)
	tk := task.Task{Name: "Test", Descr: "every month", Freq: 1, Start: 1, Done: false}
	newTasks = append(newTasks, tk)

	task.SetTasksList(newTasks)

	//should calculate
	thisCalledMonth = 7
	lastCalledMonth = 3

	taskList(w, r)

	// check status code
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status ok (200), but is: %v", resp.StatusCode)
	}

	// results
	resultTasks := make([]task.Task, 0)
	body, _ := ioutil.ReadAll(resp.Body)
	err := json.Unmarshal(body, &resultTasks)
	if err != nil {
		t.Errorf("No err expected: %v", err)
	}

	if resultTasks[0].Next != 7 {
		t.Errorf("next month expected: %v and get: %v", 7, resultTasks[0].Next)
	}

	if lastCalledMonth != 7 {
		t.Errorf("last Called Month value in server.go expected: %v and get: %v", 7, lastCalledMonth)
	}
}
