package database

func GetAllPoolName() ([]ProjectScheme, error) {
	rows, err := DB.Query("select distinct project,scriptname from system.tDatapools")
	if err != nil {
		return nil, err
	}
	res := make([]ProjectScheme, 0)
	for rows.Next() {
		var ps ProjectScheme
		err := rows.Scan(&ps.ProjectName, &ps.ScriptName)
		if err != nil {
			return nil, err
		}
		res = append(res, ps)
	}
	return res, nil
}
