package datapool

import (
	"mime/multipart"

	"github.com/matscus/ammunition/database"
	"github.com/matscus/ammunition/parser"
)

func (d Datapool) New(file *multipart.File) (err error) {
	jsonSlice, err := parser.CSVToJSON(*file)
	if err != nil {
		return err
	}
	strs := make([]string, 0)
	for _, v := range jsonSlice {
		strs = append(strs, string(v))
	}
	err = database.AddRelationsSchemeScript(d.ProjectName, d.ScriptName)
	if err != nil {
		return err
	}
	err = database.CreateScheme(d.ProjectName)
	if err != nil {
		return err
	}
	err = database.CreateTablePool(d.ProjectName, d.ScriptName)
	if err != nil {
		return err
	}
	return database.InsertMultiValuePool(d.ProjectName, d.ScriptName, strs)
}
