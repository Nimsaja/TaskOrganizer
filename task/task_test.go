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

func TestRecalculateNextMonthProp(t *testing.T) {
	tasks := []Task{
		Task{Name: "Living Room", Descr: "clean", Freq: 3, Start: 1, Done: false},
		Task{Name: "Bathroom", Descr: "clean the shower", Freq: 3, Start: 4, Done: false},
		Task{Name: "DecWork", Descr: "something for the last month", Freq: 12, Start: 11, Done: false},
		Task{Name: "JanWork", Descr: "something for the first month", Freq: 12, Start: 0, Done: false},
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
