package database

func InsertSingleValuePool(scheme string, table string, data string) error {
	_, err := DB.Exec("INSERT INTO "+scheme+"."+table+" (pool) VALUES($1)", data)
	if err != nil {
		return err
	}
	return nil
}
