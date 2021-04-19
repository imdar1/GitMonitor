package db

import (
	"fmt"
	"gitmonitor/models"
	"gitmonitor/services/utils"
)

func (db *DBConfig) GetBranchesData(projectId int64) []models.Branch {
	var branches []models.Branch
	query := fmt.Sprintf("SELECT * FROM branch WHERE project_id=%d;", projectId)
	rows, err := db.Driver.Query(query)
	utils.CheckErr(err)

	if err == nil {
		for rows.Next() {
			var branch models.Branch
			err = rows.Scan(
				&branch.BranchId,
				&branch.ProjectId,
				&branch.Name,
				&branch.IsDefault,
			)
			utils.CheckErr(err)
			branches = append(branches, branch)
		}
		rows.Close()
	}
	return branches
}

func (db *DBConfig) GetBranchById(branchId int) models.Branch {
	var branch models.Branch
	var isDefault int

	query := fmt.Sprintf("SELECT * FROM branch WHERE branch_id='%d' LIMIT 1;", branchId)
	rows := db.Driver.QueryRow(query)
	err := rows.Scan(&branch.BranchId, &branch.ProjectId, &branch.Name, &isDefault)
	branch.IsDefault = isDefault == 1
	utils.CheckErr(err)

	return branch
}

func (db *DBConfig) GetBranchIdByName(branchName string) int {
	var branchId int
	query := fmt.Sprintf("SELECT branch_id FROM branch WHERE branch_name='%s' LIMIT 1;", branchName)
	rows := db.Driver.QueryRow(query)
	err := rows.Scan(&branchId)
	utils.CheckErr(err)
	return branchId
}

func (db *DBConfig) SyncBranches(branches []string) error {
	return nil
}
