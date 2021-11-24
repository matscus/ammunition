package database

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

var (
	DB     *sqlx.DB
	scheme = `
	CREATE SCHEMA IF NOT EXISTS system AUTHORIZATION postgres;

	CREATE TABLE IF NOT EXISTS system.tDatapools (
		id SERIAL PRIMARY key,
		project VARCHAR NOT NULL,
		name  VARCHAR NOT NULL,
		bufferlen  SMALLINT NOT NULL,
		workers  SMALLINT NOT NULL,
		CONSTRAINT progectScripts UNIQUE (project, name)
	);

	CREATE EXTENSION IF NOT EXISTS pg_stat_statements;
	`
)

//PoolScheme  - type from data struct
type PoolScheme struct {
	Project   string `json:"project,omitempty" db:"project"`
	Name      string `json:"name ,omitempty" db:"name"`
	BufferLen int    `json:"bufferlen ,omitempty" db:"bufferlen"`
	Workers   int    `json:"workers ,omitempty" db:"workers"`
}

func InitDB(connStr string) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("InitDB must exec panic recover ", err)
		}
	}()
	DB, err = sqlx.Connect("postgres", connStr)
	if err != nil {
		return errors.New("Database connection error " + err.Error())
	}
	go func() {
		for {
			err := DB.Ping()
			if err != nil {
				log.Error("Database ping error ", err)
			}
			time.Sleep(10 * time.Second)
		}
	}()
	DB.MustExec(scheme)
	return nil
}

//GetAllPool - get all pools data
func GetAllPools() (pools []PoolScheme, err error) {
	return pools, DB.Select(&pools, "SELECT distinct project, name, bufferlen, workers FROM system.tDatapools")
}

//CreateScheme - create scheme and table from pool
func (ds PoolScheme) CreateScheme() sql.Result {
	defer func() {
		if err := recover(); err != nil {
			log.Error("Func CreateScheme recover panic ", err)
		}
	}()
	create := "CREATE SCHEMA IF NOT EXISTS " + ds.Project + " AUTHORIZATION postgres"
	result := DB.MustExec(create)
	return result
}

//CreateTable -  table from pool
func (ds PoolScheme) CreateTable() error {
	create := "CREATE TABLE IF NOT EXISTS " + ds.Project + "." + ds.Name + " (id serial NOT NULL PRIMARY KEY,pool jsonb NOT null)"
	_, err := DB.Exec(create)
	if err != nil {
		return errors.New("Func CreateTable exec error: " + err.Error())
	}
	return nil
}

//ClearTable - delete all values from table
func (ds PoolScheme) ClearTable() error {
	delete := "DELETE FROM " + ds.Project + "." + ds.Name
	_, err := DB.Exec(delete)
	if err != nil {
		return errors.New("Func ClearTable exec error: " + err.Error())
	}
	return nil
}

//DropTable - drop table
func (ds PoolScheme) DropTable() error {
	drop := "DROP TABLE " + ds.Project + "." + ds.Name
	_, err := DB.Exec(drop)
	if err != nil {
		return errors.New("Func DropTable exec error: " + err.Error())
	}
	return nil
}

//AddRelationsSchemeScript
func (ds PoolScheme) AddRelationsSchemeScript() (err error) {
	_, err = DB.NamedExec("INSERT INTO system.tDatapools (project, name, bufferlen, workers) VALUES(:project, :name, :bufferlen, :workers)", ds)
	if err != nil {
		return errors.New("Func AddRelationsSchemeScript named exec error: " + err.Error())
	}
	return nil
}

// DeleteRelationsSchemeScript
func (ds PoolScheme) DeleteRelationsSchemeScript() (err error) {
	_, err = DB.Exec("DELETE from system.tDatapools  WHERE project=$1 and name=$2", ds.Project, ds.Name)
	if err != nil {
		return errors.New("Delete Relations Scheme Script exec error: " + err.Error())
	}
	return nil
}

// DeleteRelationsSchemeScript
func (ds PoolScheme) DropScheme() (err error) {
	_, err = DB.Exec("DROP SCHEMA " + ds.Project)
	if err != nil {
		return errors.New("DROP SCHEMA exec error: " + err.Error())
	}
	return nil
}

func (ds PoolScheme) InsertSingleValuePool(data string) sql.Result {
	defer func() {
		if err := recover(); err != nil {
			log.Error("Insert Single Value Pool recover panic ", err)
		}
	}()
	var builder strings.Builder
	builder.WriteString("INSERT INTO ")
	builder.WriteString(ds.Project)
	builder.WriteString(".")
	builder.WriteString(" (pool) VALUES($1)")
	result := DB.MustExec(builder.String(), data)
	return result
}

//InsertMultiValues - multi values insert from
func (ds PoolScheme) InsertMultiValues(data []string) error {
	var builder strings.Builder
	builder.WriteString("INSERT INTO ")
	builder.WriteString(ds.Project)
	builder.WriteString(".")
	builder.WriteString(ds.Name)
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
		return errors.New("Insert Multi Values error: " + err.Error())
	}
	return nil
}

//GetPool - get once pool
func (ds PoolScheme) GetPool() ([]string, error) {
	res := make([]string, 0)
	query := "SELECT pool FROM " + ds.Project + "." + ds.Name
	rows, err := DB.Query(query)
	if err != nil {
		return nil, errors.New("Func GetPool query error: " + err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var str string
		err := rows.Scan(&str)
		if err != nil {
			return nil, errors.New("Func GetPool scan error: " + err.Error())
		}
		res = append(res, str)
	}
	return res, nil
}
