package task

import (
	"testing"
)

func TestGetNextMonthForTask(t *testing.T) {
	task := Task{Name: "Bathroom", Descr: "clean the shower", Freq: 3, Start: 4, Done: false}
	next := GetNextMonthForTask(task, 6)

	if next != 7 {
		t.Errorf("next month for task %v expected: %v and get: %v", task, 7, next)
	}

	task = Task{Name: "Kitchen", Descr: "clean the fridge", Freq: 4, Start: 8, Done: false}
	next = GetNextMonthForTask(task, 2)

	if next != 5 {
		t.Errorf("next month for task %v expected: %v and get: %v", task, 5, next)
	}
}
