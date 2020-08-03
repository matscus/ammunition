package datapool

import (
	"mime/multipart"

	"github.com/matscus/ammunition/database"
	"github.com/matscus/ammunition/parser"
)

func (d Datapool) Update(file *multipart.File) (err error) {
	err = database.ClearTablePool(d.ProjectName, d.ScriptName)
	if err != nil {
		return err
	}
	jsonSlice, err := parser.CSVToJSON(*file)
	if err != nil {
		return err
	}
	strs := make([]string, 0)
	for _, v := range jsonSlice {
		strs = append(strs, string(v))
	}
	return database.InsertMultiValuePool(d.ProjectName, d.ScriptName, strs)
}
