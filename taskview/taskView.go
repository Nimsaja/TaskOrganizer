package taskview

import (
	"time"

	"github.com/Nimsaja/TaskOrganizer/task"
)

// TaskView the task struct for the frontend
type TaskView struct {
	task.Task
	Next int8 `json:"next"`
}

var actMonth = int8(time.Now().Month())

func convertTask2TaskView(t task.Task) TaskView {
	return TaskView{
		Task: task.Task{
			Name:  t.Name,
			Descr: t.Descr,
			Freq:  t.Freq,
			Start: t.Start,
		},
		Next: task.GetNextMonthForTask(t, actMonth),
	}
}

// ConvertTasks2TaskViews converts the backend task list to the frontend taskview list
func ConvertTasks2TaskViews(tl []task.Task) []TaskView {
	var tvl = make([]TaskView, len(tl))
	for i, t := range tl {
		tvl[i] = convertTask2TaskView(t)
	}

	return tvl
}
