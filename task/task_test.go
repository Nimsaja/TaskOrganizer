package task

import (
	"testing"
)

func TestGetNextMonthForTask(t *testing.T) {
	task := Task{Name: "Bathroom", Descr: "clean the shower", Freq: 3, Start: 4}
	next := GetNextMonthForTask(task, 6)

	if next != 7 {
		t.Errorf("next month for task %v expected: %v and get: %v", task, 7, next)
	}

	task = Task{Name: "Kitchen", Descr: "clean the fridge", Freq: 4, Start: 8}
	next = GetNextMonthForTask(task, 2)

	if next != 4 {
		t.Errorf("next month for task %v expected: %v and get: %v", task, 4, next)
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
	RecalculateNextMonthProp(tasksPointers, 5)

	if tasks[0].Next != 7 {
		t.Errorf("next month for task %v expected: %v and get: %v", tasks[0], 7, tasks[0].Next)
	}
	if tasks[1].Next != 7 {
		t.Errorf("next month for task %v expected: %v and get: %v", tasks[1], 7, tasks[1].Next)
	}
	if tasks[2].Next != 11 {
		t.Errorf("next month for task %v expected: %v and get: %v", tasks[2], 11, tasks[2].Next)
	}
	if tasks[3].Next != 0 {
		t.Errorf("next month for task %v expected: %v and get: %v", tasks[3], 0, tasks[3].Next)
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

	l := GetTaskForMonth(0)

	if len(l) != 5 {
		t.Errorf("should find %v tasks for month %v and get: %v", 5, 0, len(l))
	}

	l = GetTaskForMonth(1)

	if len(l) != 0 {
		t.Errorf("should find %v tasks for month %v and get: %v", 0, 1, len(l))
	}

	l = GetTaskForMonth(2)

	if len(l) != 2 {
		t.Errorf("should find %v tasks for month %v and get: %v", 2, 2, len(l))
	}

	l = GetTaskForMonth(6)

	if len(l) != 5 {
		t.Errorf("should find %v tasks for month %v and get: %v", 5, 6, len(l))
	}

	//should remove list for month 0, 1 and 2
	_, ok := monthList[0]

	if ok {
		t.Errorf("list for month %v should be not available!", 0)
	}

	_, ok = monthList[1]

	if ok {
		t.Errorf("list for month %v should be not available!", 1)
	}

	_, ok = monthList[2]

	if ok {
		t.Errorf("list for month %v should be not available!", 2)
	}

	//list for month 6 should still be there
	_, ok = monthList[6]

	if !ok {
		t.Errorf("list for month %v should be available!", 6)
	}

}
