package database

//CreateTablePool -  table from pool
func CreateTablePool(scheme string, table string) error {
	_, err := DB.Exec("CREATE TABLE IF NOT EXISTS "+scheme+"."+table+" (id serial NOT NULL PRIMARY KEY,pool json NOT null)", scheme, table)
	if err != nil {
		return err
	}
	return nil
}
