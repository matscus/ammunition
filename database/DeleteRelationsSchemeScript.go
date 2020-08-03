package database

import "errors"

func DeleteRelationsSchemeScript(scheme string, table string) (err error) {
	_, err = DB.Exec("DELETE system.tDatapools  WHERE project=$1 and scriptname=$2", scheme, table)
	if err != nil {
		return errors.New("Delete Relations Scheme and Script error: " + err.Error())
	}
	return nil
}
