package database

import "errors"

func AddRelationsSchemeScript(scheme string, table string) (err error) {
	_, err = DB.Exec("INSERT INTO system.tDatapools (project,scriptname) VALUES($1,$2)", scheme, table)
	if err != nil {
		return errors.New("ADD Relations Scheme and Script error: " + err.Error())
	}
	return nil
}
