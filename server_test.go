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
