package database

func GetPool(scheme string, table string) (*[]string, error) {
	rows, err := DB.Query("SELECT pool FROM " + scheme + "." + table)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := make([]string, 0, 0)
	for rows.Next() {
		var str string
		err := rows.Scan(&str)
		if err != nil {
			return nil, err
		}
		res = append(res, str)
	}
	return &res, nil
}
