package models

type Branch struct {
	BranchId  int
	ProjectId int
	Name      string
	IsDefault bool
	IsDeleted bool
}
