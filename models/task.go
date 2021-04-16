package models

type Task struct {
	TaskId        int
	ProjectId     int
	BranchId      int
	Name          string
	AssigneeName  string
	AssigneeEmail string
	TaskStatus    int
	StartDate     int64
	EndDate       int64
}
