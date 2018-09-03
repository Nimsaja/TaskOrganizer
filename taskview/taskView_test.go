package taskview

import (
	"encoding/json"
	"testing"

	"github.com/Nimsaja/TaskOrganizer/task"
)

func TestConvertTask2TaskView(t *testing.T) {
	var tests = []struct {
		task     task.Task // input: Task
		expected int8      // expected result: next month
	}{
		{task.Task{Name: "Living Room", Descr: "clean", Freq: 3, Start: 1}, 7},
		{task.Task{Name: "Bathroom", Descr: "clean the shower", Freq: 3, Start: 4}, 7},
		{task.Task{Name: "DecWork", Descr: "something for the last month", Freq: 12, Start: 11}, 11},
		{task.Task{Name: "JanWork", Descr: "something for the first month", Freq: 12, Start: 0}, 0},
	}

	actMonth = 5

	var tv TaskView
	for _, tt := range tests {
		tv = convertTask2TaskView(tt.task)
		if tv.Next != tt.expected {
			t.Errorf("next month for task %v expected: %v and get: %v", tv, tt.expected, tv.Next)
		}
	}
}

func TestConvertTasks2TaskViews(t *testing.T) {
	tasks := []task.Task{
		task.Task{Name: "Living Room", Descr: "clean", Freq: 3, Start: 1},
		task.Task{Name: "Bathroom", Descr: "clean the shower", Freq: 3, Start: 4},
		task.Task{Name: "DecWork", Descr: "something for the last month", Freq: 12, Start: 11},
		task.Task{Name: "JanWork", Descr: "something for the first month", Freq: 12, Start: 0},
	}

	actMonth = 5

	tests := [4]int8{7, 7, 11, 0}

	taskViews := ConvertTasks2TaskViews(tasks)

	for i, tt := range tests {
		if taskViews[i].Next != tt {
			t.Errorf("next month for task %v expected: %v and get: %v", taskViews[i], tt, taskViews[i].Next)
		}
	}
}

func TestJson(t *testing.T) {
	tasks := []task.Task{
		task.Task{Name: "Living Room", Descr: "clean", Freq: 3, Start: 1},
		task.Task{Name: "Bathroom", Descr: "clean the shower", Freq: 3, Start: 4},
		task.Task{Name: "DecWork", Descr: "something for the last month", Freq: 12, Start: 11},
		task.Task{Name: "JanWork", Descr: "something for the first month", Freq: 12, Start: 0},
	}
	taskViews := ConvertTasks2TaskViews(tasks)

	j, err := json.Marshal(taskViews)
	if err != nil {
		t.Errorf("no err expected: %v", err)
	}

	taskViewsAfter := make([]TaskView, 0)
	err = json.Unmarshal(j, &taskViewsAfter)
	if err != nil {
		t.Errorf("no err expected: %v", err)
	}
	if len(taskViews) != len(taskViewsAfter) {
		t.Errorf("different length: %v != %v", len(taskViews), len(taskViewsAfter))
	}
	tva := taskViewsAfter[0]
	if taskViews[0].Name != tva.Name {
		t.Errorf("different name: %v != %v", taskViews[0].Name, tva.Name)
	}
}
