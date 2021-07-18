package models

func GetBranchesName(branches []Branch) []string {
	var branchesName []string
	for _, v := range branches {
		branchesName = append(branchesName, v.Name)
	}
	return branchesName
}

func GetBranchName(branchId int, branches []Branch) string {
	name := ""
	for _, v := range branches {
		if v.BranchId == branchId {
			name = v.Name
			break
		}
	}
	return name
}
