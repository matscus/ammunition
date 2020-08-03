package database

import "errors"

//CreateScheme - create scheme and table from pool
func CreateScheme(scheme string) error {
	_, err := DB.Exec("CREATE SCHEMA IF NOT EXISTS " + scheme + " AUTHORIZATION postgres")
	if err != nil {
		return errors.New("Create scheme error: " + err.Error())
	}
	return nil
}
