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
	Done  bool   `json:"done"`
}

// GetDefaultList a list of task to start with
func GetDefaultList() []Task {
	return []Task{
		Task{Name: "Living Room", Descr: "clean", Freq: 3, Start: 1, Done: false},
		Task{Name: "Sleeping Room", Descr: "hoover the floor", Freq: 2, Start: 2, Done: false},
		Task{Name: "Basement", Descr: "kill spiders", Freq: 6, Start: 3, Done: false},
		Task{Name: "Bathroom", Descr: "clean the shower", Freq: 3, Start: 4, Done: false},
		Task{Name: "Kitchen", Descr: "clean cupboards", Freq: 4, Start: 5, Done: false},
		Task{Name: "Corridor", Descr: "hoover", Freq: 2, Start: 6, Done: false},
		Task{Name: "Bed", Descr: "change the sheets", Freq: 2, Start: 7, Done: false},
		Task{Name: "Kitchen", Descr: "clean the fridge", Freq: 12, Start: 8, Done: false},
		Task{Name: "Living Room", Descr: "clean the cupboards", Freq: 3, Start: 9, Done: false},
		Task{Name: "Guest Room", Descr: "clean", Freq: 4, Start: 10, Done: false},
	}
}

// GetNextMonthForTask ...
func GetNextMonthForTask(t Task, m int8) int8 {
	next := (float64(m) - float64(t.Start)) / float64(t.Freq)
	if next < 0 {
		next = -next
	}
	next = math.Ceil(next)*float64(t.Freq) + float64(t.Start)
	return int8(int(next) % 11)
}
