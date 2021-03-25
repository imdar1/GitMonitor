package models

type TaskStat int

type Task struct {
	TaskId        int
	ProjectId     int
	BranchId      int
	Name          string
	AssigneeName  string
	AssigneeEmail string
	TaskStatus    TaskStat
	StartDate     int
	EndDate       int
}
