package task

import (
	"math"
)

// Task the struct
type Task struct {
	Name  string `json:"name"`
	Descr string `json:"desc"`
	Freq  int8   `json:"freq"`
	Start int8   `json:"start"`
	Next  int8   `json:"next"`
}

var tasks = []Task{
	Task{Name: "Living Room", Descr: "clean", Freq: 3, Start: 1},
	Task{Name: "Sleeping Room", Descr: "hoover the floor", Freq: 2, Start: 2},
	Task{Name: "Basement", Descr: "kill spiders", Freq: 6, Start: 3},
	Task{Name: "Bathroom", Descr: "clean the shower", Freq: 3, Start: 4},
	Task{Name: "Kitchen", Descr: "clean cupboards", Freq: 4, Start: 5},
	Task{Name: "Corridor", Descr: "hoover", Freq: 2, Start: 6},
	Task{Name: "Bed", Descr: "change the sheets", Freq: 2, Start: 7},
	Task{Name: "Kitchen", Descr: "clean the fridge", Freq: 12, Start: 8},
	Task{Name: "Living Room", Descr: "clean the cupboards", Freq: 3, Start: 9},
	Task{Name: "Guest Room", Descr: "clean", Freq: 4, Start: 10},
}

// DefaultTasks default tasks to begin with and to test against
var DefaultTasks = tasks

var monthList = make(map[int][]Task)

// SetTasksList override default task list
func SetTasksList(list []Task) {
	if list == nil {
		tasks = DefaultTasks
	} else {
		tasks = list
	}
}

// GetDefaultList a list of task to start with
func GetDefaultList() []Task {
	return tasks
}

// GetNextMonthForTask ...
func GetNextMonthForTask(t Task, m int) int {
	next := (float64(m+1) - float64(t.Start+1)) / float64(t.Freq)
	next = math.Ceil(next)*float64(t.Freq) + float64(t.Start+1)
	nextInt := int(next)%12 - 1

	if nextInt < 0 {
		nextInt = nextInt + 12
	}
	return nextInt
}

// RecalculateNextMonthProp ...
func RecalculateNextMonthProp(tasks []*Task, m int) {
	for _, el := range tasks {
		el.Next = int8(GetNextMonthForTask(*el, m))
	}
}

// GetTaskForMonth gets or calculates the list of tasks for month @m
func GetTaskForMonth(m int) []Task {
	_, exists := monthList[m]

	if !exists {
		//calculate task list, save and return
		l := make([]Task, 0)
		for _, el := range tasks {
			if GetNextMonthForTask(el, m) == m {
				l = append(l, el)
			}
		}
		monthList[m] = l

		//need to remove old task lists - everything that is older than 3 month
		//to keep it simple just calculate the difference of the two months
		if len(monthList) > 3 {
			for key := range monthList {
				if math.Abs(float64(key-m)) > 3 {
					delete(monthList, key)
				}
			}
		}
	}

	return monthList[m]
}
