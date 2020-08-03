package database

func AddRelationsSchemeScript(scheme string, table string) (err error) {
	_, err = DB.Exec("INSERT INTO system.tDatapools (project,scriptname) VALUES($1,$1)", scheme, table)
	return err
}
