package task

import (
	"testing"
)

func TestGetNextMonthForTask(t *testing.T) {
	var tests = []struct {
		m        int  // input: month
		task     Task // input: Task
		expected int  // expected result: next month
	}{
		{6, Task{Name: "Bathroom", Descr: "clean the shower", Freq: 3, Start: 4}, 7},
		{2, Task{Name: "Bathroom", Descr: "clean the shower", Freq: 3, Start: 4}, 4},
		{2, Task{Name: "Kitchen", Descr: "clean the fridge", Freq: 4, Start: 8}, 4},
		{6, Task{Name: "Kitchen", Descr: "clean the fridge", Freq: 4, Start: 8}, 8},
	}

	for _, tt := range tests {
		next := GetNextMonthForTask(tt.task, tt.m)
		if next != tt.expected {
			t.Errorf("next month for task %v expected: %v and get: %v", tt.task, tt.expected, next)
		}
	}
}

func TestRecalculateNextMonthProp(t *testing.T) {
	tasks := []Task{
		Task{Name: "Living Room", Descr: "clean", Freq: 3, Start: 1},
		Task{Name: "Bathroom", Descr: "clean the shower", Freq: 3, Start: 4},
		Task{Name: "DecWork", Descr: "something for the last month", Freq: 12, Start: 11},
		Task{Name: "JanWork", Descr: "something for the first month", Freq: 12, Start: 0},
	}
	var tasksPointers []*Task
	for i := 0; i < len(tasks); i++ {
		tasksPointers = append(tasksPointers, &tasks[i])
	}

	tests := [4]int8{7, 7, 11, 0}

	RecalculateNextMonthProp(tasksPointers, 5)

	for i, tt := range tests {
		if tasks[i].Next != tt {
			t.Errorf("next month for task %v expected: %v and get: %v", tasks[0], tt, tasks[0].Next)
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

	SetTasksList(tasks)

	if len(monthList) != 0 {
		t.Errorf("monthList should be empty when called the first time! Has length %v", len(monthList))
	}

	var tests = []struct {
		m        int // input: month
		expected int // expected result: length of task list for month m
	}{
		{0, 5},
		{1, 0},
		{2, 2},
		{6, 5},
	}

	for _, tt := range tests {
		list := GetTaskForMonth(tt.m)
		if len(list) != tt.expected {
			t.Errorf("should find %v tasks for month %v and get: %v", tt.expected, tt.m, len(list))
		}
	}

	//should remove list for month 0, 1 and 2
	//TODO: split this test here when owner functionality is implemented!!
	var testLists = []struct {
		m        int    // input: month
		expected bool   // expected result: if list exists
		msg      string // error message
	}{
		{0, false, "list for month %v should be not available!"},
		{1, false, "list for month %v should be not available!"},
		{2, false, "list for month %v should be not available!"},
		{6, true, "list for month %v should be available!"},
	}

	for _, tt := range testLists {
		_, ok := monthList[tt.m]
		if ok != tt.expected {
			t.Errorf(tt.msg, tt.m)
		}
	}
}
