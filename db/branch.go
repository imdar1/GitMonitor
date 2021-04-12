package db

import (
	"fmt"
	"gitmonitor/models"
	"gitmonitor/services"
)

func (db *DBConfig) GetBranchesData(projectId int64) []models.Branch {
	var branches []models.Branch
	query := fmt.Sprintf("SELECT * FROM branch WHERE project_id=%d;", projectId)
	rows, err := db.Driver.Query(query)
	services.CheckErr(err)

	if rows != nil {
		for rows.Next() {
			var branch models.Branch
			err = rows.Scan(
				&branch.BranchId,
				&branch.ProjectId,
				&branch.Name,
				&branch.IsDefault,
			)
			services.CheckErr(err)
			branches = append(branches, branch)
		}
		rows.Close()
	}
	return branches
}
