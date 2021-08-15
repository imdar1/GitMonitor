package models

type Project struct {
	ProjectId         int64
	ProjectDir        string
	ProjectStartDate  int64
	ProjectEndDate    int64
	DefaultBranchName string
	DefaultRemoteName string
	IsFirstTime       bool
}
