package datapool

import (
	"errors"
	"mime/multipart"

	"github.com/allegro/bigcache"
	"github.com/matscus/ammunition/database"
	"github.com/matscus/ammunition/parser"
)

type Datapool struct {
	ScriptName  string `json:"scriptname"`
	ProjectName string `json:"projectname"`
	UniqueData  bool   `json:"uniquedata"`
}

func (d Datapool) New(file *multipart.File) (err error) {
	jsonSlice, err := parser.CSVToJSON(*file)
	if err != nil {
		return err
	}
	strs := make([]string, 0)
	for _, v := range jsonSlice {
		strs = append(strs, string(v))
	}
	scheme := database.DatabaseScheme{ProjectName: d.ProjectName, ScriptName: d.ScriptName}
	err = scheme.AddRelationsSchemeScript()
	if err != nil {
		return err
	}
	err = scheme.CreateScheme()
	if err != nil {
		return err
	}
	err = scheme.CreateTable()
	if err != nil {
		return err
	}
	initCashe(d.ProjectName+d.ScriptName, &strs)
	return scheme.InsertMultiValues(strs)
}
func (d Datapool) Update(file *multipart.File) (err error) {
	scheme := database.DatabaseScheme{ProjectName: d.ProjectName, ScriptName: d.ScriptName}
	err = scheme.ClearTable()
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
	err = reInitCashe(d.ProjectName+d.ScriptName, &strs)
	if err != nil {
		return err
	}
	return scheme.InsertMultiValues(strs)
}

func (d Datapool) AddValues(file *multipart.File) (err error) {
	jsonSlice, err := parser.CSVToJSON(*file)
	if err != nil {
		return err
	}
	strs := make([]string, 0)
	for _, v := range jsonSlice {
		strs = append(strs, string(v))
	}
	err = addValueInCashe(d.ProjectName+d.ScriptName, &strs)
	if err != nil {
		return err
	}

	return database.DatabaseScheme{ProjectName: d.ProjectName, ScriptName: d.ScriptName}.InsertMultiValues(strs)
}

func (d Datapool) Delete() (err error) {
	ch, ok := chanMap.Load(d.ProjectName + d.ScriptName)
	if ok {
		close(ch.(chan string))
	}
	cache, ok := cacheMap.Load(d.ProjectName + d.ScriptName)
	if ok {
		cache.(*bigcache.BigCache).Close()
	}
	return database.DatabaseScheme{ProjectName: d.ProjectName, ScriptName: d.ScriptName}.DropTable()
}

//InitPoolFromDB - Datapool initialization function.
//gets all data from the database, based on the project name and script name fields,
//and initializes the data cache and the upload channel.
func (d Datapool) InitPoolFromDB() (err error) {
	data, err := database.DatabaseScheme{ProjectName: d.ProjectName, ScriptName: d.ScriptName}.GetPool()
	if err != nil {
		println(err.Error())
		return err
	}
	initCashe(d.ProjectName+d.ScriptName, data)
	return nil
}

func (d Datapool) GetValue() (string, error) {
	ch, ok := chanMap.Load(d.ProjectName + d.ScriptName)
	if ok {
		res := <-ch.(chan string)
		return res, nil
	}
	return "", errors.New("chanel is empty")
}
