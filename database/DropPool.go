package database

import "errors"

func DropPool(scheme string, table string) error {
	_, err := DB.Exec("DROP TABLE " + scheme + "." + table)
	if err != nil {
		return errors.New("Drop table error: " + err.Error())
	}
	return nil
}
