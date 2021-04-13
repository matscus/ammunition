package database

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

var (
	DB     *sqlx.DB
	schema = `
	CREATE SCHEMA IF NOT EXISTS system AUTHORIZATION;

	CREATE TABLE IF NOT EXISTS system.tDatapools (
		id SERIAL PRIMARY key,
		project varchar NOT NULL,
		script  varchar NOT NULL,
		CONSTRAINT progectScripts UNIQUE (project, script)
	);

	CREATE EXTENSION IF NOT EXISTS pg_stat_statements;
	`
)

//PoolScheme  - type from data struct
type PoolScheme struct {
	Project string `json:"project,omitempty" db:"project"`
	Script  string `json:"script ,omitempty" db:"script"`
}

func InitDB(connStr string) (err error) {
	DB, err = sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Printf("Database connection error %s\n", err)
		return err
	}
	log.Printf("Database connection completed")
	go func() {
		for {
			err := DB.Ping()
			if err != nil {
				log.Printf("database ping error %s\n", err)
			}
			time.Sleep(5 * time.Second)
		}
	}()
	DB.MustExec(schema)
	log.Println("schema init complited")
	return nil
}

//GetAllNames - get all pools name
func GetAllName() (pools []PoolScheme, err error) {
	return pools, DB.Select(&pools, "SELECT distinct project,script FROM system.tDatapools")
}

//CreateScheme - create scheme and table from pool
func (ds PoolScheme) CreateScheme() error {
	_, err := DB.NamedExec("CREATE SCHEMA IF NOT EXISTS :project AUTHORIZATION postgres", ds)
	if err != nil {
		return errors.New("Create scheme error: " + err.Error())
	}
	return nil
}

//CreateTable -  table from pool
func (ds PoolScheme) CreateTable() error {
	_, err := DB.Exec("CREATE TABLE IF NOT EXISTS $1.$2 (id serial NOT NULL PRIMARY KEY,pool json NOT null)", ds.Project, ds.Script)
	if err != nil {
		return errors.New("Create table error: " + err.Error())
	}
	return nil
}

//ClearTable - delete all values from table
func (ds PoolScheme) ClearTable() error {
	_, err := DB.Exec("DELETE FROM $1.$2", ds.Project, ds.Script)
	if err != nil {
		return errors.New("Delete pool error: " + err.Error())
	}
	return nil
}

//DropTable - drop table
func (ds PoolScheme) DropTable() error {
	_, err := DB.Exec("DROP TABLE $1.$2", ds.Project, ds.Script)
	if err != nil {
		return errors.New("Drop table error: " + err.Error())
	}
	return nil
}

//AddRelationsSchemeScript
func (ds PoolScheme) AddRelationsSchemeScript() (err error) {
	_, err = DB.NamedExec("INSERT INTO system.tDatapools (project,script) VALUES(:project,:script)", ds)
	if err != nil {
		return errors.New("ADD Relations Scheme and Script error: " + err.Error())
	}
	return nil
}

// DeleteRelationsSchemeScript
func (ds PoolScheme) DeleteRelationsSchemeScript() (err error) {
	_, err = DB.Exec("DELETE system.tDatapools  WHERE project=$1 and script=$2", ds.Project, ds.Script)
	if err != nil {
		return errors.New("Delete Relations Scheme and Script error: " + err.Error())
	}
	return nil
}

func (ds PoolScheme) InsertSingleValuePool(data string) sql.Result {
	var builder strings.Builder
	builder.WriteString("INSERT INTO ")
	builder.WriteString(ds.Project)
	builder.WriteString(".")
	builder.WriteString(" (pool) VALUES($1)")
	return DB.MustExec(builder.String(), data)
}

//InsertMultiValues - multi values insert from
func (ds PoolScheme) InsertMultiValues(data []string) error {
	var builder strings.Builder
	builder.WriteString("INSERT INTO ")
	builder.WriteString(ds.Project)
	builder.WriteString(".")
	builder.WriteString(ds.Script)
	builder.WriteString(" (pool) VALUES ")
	l := len(data)
	for i := 0; i < l; i++ {
		builder.WriteString("('")
		builder.WriteString(data[i])
		builder.WriteString("')")
		if i < l-1 {
			builder.WriteString(",")
		}
	}
	_, err := DB.Exec(builder.String())
	if err != nil {
		return errors.New("Insert Multi Value Pool error: " + err.Error())
	}
	return nil
}

//GetPool - get once pool
func (ds PoolScheme) GetPool() ([]string, error) {
	rows, err := DB.Query(`SELECT pool FROM $1.$2`, ds.Project, ds.Script)
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
	return res, nil
}
