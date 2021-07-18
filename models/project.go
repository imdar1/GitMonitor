package models

type Project struct {
	ProjectId         int64
	ProjectDir        string
	DefaultBranchName string
	DefaultRemoteName string
}
