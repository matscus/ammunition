package database

import "errors"

func ClearTablePool(scheme string, table string) error {
	_, err := DB.Exec("DELETE FROM  " + scheme + "." + table)
	if err != nil {
		return errors.New("Delete pool error: " + err.Error())
	}
	return nil
}
