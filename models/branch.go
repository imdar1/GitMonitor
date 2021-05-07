package models

type Branch struct {
	BranchId      int
	ProjectId     int
	Name          string
	IsMergeTarget bool
	IsDeleted     bool
}
