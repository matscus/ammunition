package database

func DeleteRelationsSchemeScript(scheme string, table string) (err error) {
	_, err = DB.Exec("DELETE system.tDatapools  WHERE project=$1 and scriptname=$2", scheme, table)
	return err
}
