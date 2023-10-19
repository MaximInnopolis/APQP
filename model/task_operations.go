package model

func CalculateCurrentIter(task *Task) float64 {
	return task.CurrentIter + task.D
}
