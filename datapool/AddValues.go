package datapool

import (
	"mime/multipart"

	"github.com/matscus/ammunition/database"
	"github.com/matscus/ammunition/parser"
)

func (d Datapool) Add(file *multipart.File) (err error) {
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
