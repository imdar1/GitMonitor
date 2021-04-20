package db

import (
	"database/sql"
	"fmt"
	"gitmonitor/models"
	"gitmonitor/services/utils"
	"sort"
)

func (db *DBConfig) GetBranchesData(projectId int64) ([]models.Branch, error) {
	var branches []models.Branch
	query := fmt.Sprintf("SELECT * FROM branch WHERE project_id=%d;", projectId)
	rows, err := db.Driver.Query(query)

	if err != nil && err != sql.ErrNoRows {
		return branches, err
	}

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
	return branches, nil
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
	query := fmt.Sprintf("SELECT branch_id FROM branch WHERE name='%s' LIMIT 1;", branchName)
	rows := db.Driver.QueryRow(query)
	err := rows.Scan(&branchId)
	utils.CheckErr(err)
	return branchId
}

func getBranchesName(branches []models.Branch) []string {
	branchName := []string{}
	for _, v := range branches {
		branchName = append(branchName, v.Name)
	}
	return branchName
}

func (db *DBConfig) SyncBranches(projectId int64, branches []string) error {
	sort.Strings(branches)

	branchesModel, err := db.GetBranchesData(projectId)
	if err != nil {
		return err
	}

	// Check if a branch exists in DB also exists in remote branches if not delete the record
	delTx, err := db.Driver.Begin()
	if err != nil {
		return err
	}
	for _, v := range branchesModel {
		isExist := utils.IsExist(v.Name, branches)
		if !isExist {
			deleteQuery := fmt.Sprintf("DELETE FROM branch WHERE name='%s';", v.Name)
			_, err := delTx.Exec(deleteQuery)
			if err != nil {
				delTx.Rollback()
				return err
			}
		}
	}
	err = delTx.Commit()
	if err != nil {
		return err
	}

	branchesModelList := getBranchesName(branchesModel)
	isDefault := 0
	insTx, err := db.Driver.Begin()
	if err != nil {
		return err
	}
	for _, v := range branches {
		// insert if not exist
		isExist := utils.IsExist(v, branchesModelList)
		if !isExist {
			insertQuery := fmt.Sprintf(
				"INSERT INTO branch(project_id, name, is_default) VALUES(%d, '%s', %d); ",
				projectId,
				v,
				isDefault,
			)

			_, err := insTx.Exec(insertQuery)
			if err != nil {
				insTx.Rollback()
				return err
			}
		}
	}
	err = insTx.Commit()
	if err != nil {
		return err
	}

	return nil
}
