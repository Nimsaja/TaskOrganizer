package task

import (
	"testing"
)

func testNew(tasks []Task) *Organizer {
	o := New("For Tests")
	o.SetTasks(tasks)
	return o
}

func TestGetNextMonthForTask(t *testing.T) {
	var tests = []struct {
		m        int8  // input: month
		task     *Task // input: Task
		expected int8  // expected result: next month
	}{
		{6, &Task{Name: "Bathroom", Descr: "clean the shower", Freq: 3, Start: 4}, 7},
		{2, &Task{Name: "Bathroom", Descr: "clean the shower", Freq: 3, Start: 4}, 4},
		{2, &Task{Name: "Kitchen", Descr: "clean the fridge", Freq: 4, Start: 8}, 4},
		{6, &Task{Name: "Kitchen", Descr: "clean the fridge", Freq: 4, Start: 8}, 8},
	}

	for _, tt := range tests {
		next := GetNextMonthForTask(*tt.task, tt.m)
		if next != tt.expected {
			t.Errorf("next month for task %v expected: %v and get: %v", tt.task, tt.expected, next)
		}
	}
}

func TestCalculateMonthList(t *testing.T) {
	tasks := []Task{
		Task{Name: "Living Room", Descr: "clean", Freq: 2, Start: 0},
		Task{Name: "Bathroom", Descr: "clean the shower", Freq: 3, Start: 0},
		Task{Name: "Work1", Descr: "something todo", Freq: 4, Start: 0},
		Task{Name: "Work2", Descr: "something", Freq: 6, Start: 0},
		Task{Name: "Living Room", Descr: "clean the cupboards", Freq: 3, Start: 9},
		Task{Name: "Guest Room", Descr: "clean", Freq: 4, Start: 10},
	}
	o := testNew(tasks)
	o.SetTasks(tasks)

	var tests = []struct {
		m        int8 // input: month
		expected int  // expected result: length of task list for month m
	}{
		{0, 5},
		{1, 0},
		{2, 2},
		{6, 5},
	}

	for _, tt := range tests {
		list := o.GetTaskForMonth(tt.m)
		if len(list) != tt.expected {
			t.Errorf("should find %v tasks for month %v and get: %v", tt.expected, tt.m, len(list))
		}
	}
}
