package database

import "errors"

func InsertSingleValuePool(scheme string, table string, data string) error {
	_, err := DB.Exec("INSERT INTO "+scheme+"."+table+" (pool) VALUES($1)", data)
	if err != nil {
		return errors.New("Insert Single Value Pool error: " + err.Error())
	}
	return nil
}
