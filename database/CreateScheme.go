package database

//CreateScheme - create scheme and table from pool
func CreateScheme(scheme string) error {
	_, err := DB.Exec("CREATE SCHEMA IF NOT EXISTS " + scheme + " AUTHORIZATION postgres;")
	if err != nil {
		return err
	}
	return nil
}
