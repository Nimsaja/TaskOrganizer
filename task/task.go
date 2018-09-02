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

// Organizer the struct
type Organizer struct {
	Owner     string  `json:"owner"`
	Tasks     []*Task `json:"tasks"`
	monthList map[int][]*Task
}

// New instance of Organizer
func New(owner string) *Organizer {
	return &Organizer{
		Owner:     owner,
		Tasks:     GetDefaultTasks(),
		monthList: make(map[int][]*Task),
	}
}

// SetTasks override default task list
func (o *Organizer) SetTasks(tasks []*Task) {
	o.Tasks = tasks
}

// GetDefaultTasks a list of task to start with
func GetDefaultTasks() []*Task {
	return []*Task{
		&Task{Name: "Living Room", Descr: "clean", Freq: 3, Start: 1},
		&Task{Name: "Sleeping Room", Descr: "hoover the floor", Freq: 2, Start: 2},
		&Task{Name: "Basement", Descr: "kill spiders", Freq: 6, Start: 3},
		&Task{Name: "Bathroom", Descr: "clean the shower", Freq: 3, Start: 4},
		&Task{Name: "Kitchen", Descr: "clean cupboards", Freq: 4, Start: 5},
		&Task{Name: "Corridor", Descr: "hoover", Freq: 2, Start: 6},
		&Task{Name: "Bed", Descr: "change the sheets", Freq: 2, Start: 7},
		&Task{Name: "Kitchen", Descr: "clean the fridge", Freq: 12, Start: 8},
		&Task{Name: "Living Room", Descr: "clean the cupboards", Freq: 3, Start: 9},
		&Task{Name: "Guest Room", Descr: "clean", Freq: 4, Start: 10},
	}
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
func (o *Organizer) RecalculateNextMonthProp(m int) {
	for _, el := range o.Tasks {
		el.Next = int8(GetNextMonthForTask(*el, m))
	}
}

// GetTaskForMonth gets or calculates the list of tasks for month @m
func (o *Organizer) GetTaskForMonth(m int) []*Task {
	_, exists := o.monthList[m]

	if !exists {
		//calculate task list, save and return
		l := make([]*Task, 0)
		for _, el := range o.Tasks {
			if GetNextMonthForTask(*el, m) == m {
				l = append(l, el)
			}
		}
		o.monthList[m] = l

		//need to remove old task lists - everything that is older than 3 month
		//to keep it simple just calculate the difference of the two months
		if len(o.monthList) > 3 {
			for key := range o.monthList {
				if math.Abs(float64(key-m)) > 3 {
					delete(o.monthList, key)
				}
			}
		}
	}

	return o.monthList[m]
}
