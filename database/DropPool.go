package database

func DropPool(scheme string, table string) error {
	_, err := DB.Exec("DROP TABLE " + scheme + "." + table)
	if err != nil {
		return err
	}
	return nil
}
