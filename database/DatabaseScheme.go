package database

import (
	"bytes"
	"errors"
	"os"
)

//DatabaseScheme  - type from data struct
type DatabaseScheme struct {
	ProjectName string
	ScriptName  string
}

//CreateScheme - create scheme and table from pool
func (ds DatabaseScheme) CreateScheme() error {
	_, err := DB.Exec("CREATE SCHEMA IF NOT EXISTS " + ds.ProjectName + " AUTHORIZATION " + os.Getenv("POSTGRESUSER"))
	if err != nil {
		return errors.New("Create scheme error: " + err.Error())
	}
	return nil
}

//CreateTable -  table from pool
func (ds DatabaseScheme) CreateTable() error {
	_, err := DB.Exec("CREATE TABLE IF NOT EXISTS " + ds.ProjectName + "." + ds.ScriptName + " (id serial NOT NULL PRIMARY KEY,pool json NOT null)")
	if err != nil {
		return errors.New("Create table error: " + err.Error())
	}
	return nil
}

func (ds DatabaseScheme) ClearTable() error {
	_, err := DB.Exec("DELETE FROM  " + ds.ProjectName + "." + ds.ScriptName)
	if err != nil {
		return errors.New("Delete pool error: " + err.Error())
	}
	return nil
}
func (ds DatabaseScheme) DropTable() error {
	_, err := DB.Exec("DROP TABLE " + ds.ProjectName + "." + ds.ScriptName)
	if err != nil {
		return errors.New("Drop table error: " + err.Error())
	}
	return nil
}

func (ds DatabaseScheme) AddRelationsSchemeScript() (err error) {
	_, err = DB.Exec("INSERT INTO system.tDatapools (project,scriptname) VALUES($1,$2)", ds.ProjectName, ds.ScriptName)
	if err != nil {
		return errors.New("ADD Relations Scheme and Script error: " + err.Error())
	}
	return nil
}

func (ds DatabaseScheme) DeleteRelationsSchemeScript() (err error) {
	_, err = DB.Exec("DELETE system.tDatapools  WHERE project=$1 and scriptname=$2", ds.ProjectName, ds.ScriptName)
	if err != nil {
		return errors.New("Delete Relations Scheme and Script error: " + err.Error())
	}
	return nil
}

func (ds DatabaseScheme) InsertSingleValuePool(data string) error {
	_, err := DB.Exec("INSERT INTO "+ds.ProjectName+"."+ds.ScriptName+" (pool) VALUES($1)", data)
	if err != nil {
		return errors.New("Insert Single Value Pool error: " + err.Error())
	}
	return nil
}

func (ds DatabaseScheme) InsertMultiValues(data []string) error {
	var buf bytes.Buffer
	buf.WriteString("INSERT INTO " + ds.ProjectName + "." + ds.ScriptName + " (pool) VALUES ")
	l := len(data)
	for i := 0; i < l; i++ {
		buf.WriteString("('" + data[i] + "')")
		if i < l-1 {
			buf.WriteString(",")
		}
	}
	_, err := DB.Exec(buf.String())
	if err != nil {
		return errors.New("Insert Multi Value Pool error: " + err.Error())
	}
	return nil
}

func (ds DatabaseScheme) GetPool() (*[]string, error) {
	rows, err := DB.Query("SELECT pool FROM " + ds.ProjectName + "." + ds.ScriptName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := make([]string, 0, 0)
	for rows.Next() {
		var str string
		err := rows.Scan(&str)
		if err != nil {
			return nil, err
		}
		res = append(res, str)
	}
	return &res, nil
}
