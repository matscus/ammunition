package database

import "errors"

//CreateTablePool -  table from pool
func CreateTablePool(scheme string, table string) error {
	_, err := DB.Exec("CREATE TABLE IF NOT EXISTS " + scheme + "." + table + " (id serial NOT NULL PRIMARY KEY,pool json NOT null)")
	if err != nil {
		return errors.New("Create table error: " + err.Error())
	}
	return nil
}
