package database

func GetAllPoolName() ([]DatabaseScheme, error) {
	rows, err := DB.Query("select distinct project,scriptname from system.tDatapools")
	if err != nil {
		return nil, err
	}
	res := make([]DatabaseScheme, 0)
	for rows.Next() {
		var ps DatabaseScheme
		err := rows.Scan(&ps.ProjectName, &ps.ScriptName)
		if err != nil {
			return nil, err
		}
		res = append(res, ps)
	}
	return res, nil
}
