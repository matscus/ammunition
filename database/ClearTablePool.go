package database

func ClearTablePool(scheme string, table string) error {
	_, err := DB.Exec("DELETE FROM  " + scheme + "." + table)
	if err != nil {
		return err
	}
	return nil
}
