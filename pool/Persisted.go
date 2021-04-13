package pool

import (
	"mime/multipart"

	"github.com/matscus/ammunition/cache"
	"github.com/matscus/ammunition/database"
	"github.com/matscus/ammunition/parser"
)

var (
	databaseSchemes []database.PoolScheme
)

type PersistedPool struct {
	Project      string `json:"project"`
	Script       string `json:"script"`
	BufferLen    int    `json:"bufferlen,omitempty"`
	WorkersCount int    `json:"workerscount,omitempty"`
}

func InitAllPersistedPools() (err error) {
	databaseSchemes, err = database.GetAllName()
	if err != nil {
		return err
	}
	for _, v := range databaseSchemes {
		go PersistedPool{Project: v.Project, Script: v.Script}.InitPoolFromDB()
	}
	return nil
}

func (p PersistedPool) Create(file *multipart.File) (err error) {
	jsonSlice, err := parser.CSVToJSON(*file)
	if err != nil {
		return err
	}
	strs := make([]string, 0)
	for _, v := range jsonSlice {
		strs = append(strs, string(v))
	}
	scheme := database.PoolScheme{Project: p.Project, Script: p.Script}
	err = newScheme(scheme)
	if err != nil {
		return err
	}
	cache, err := cache.CreateDefaultCache(p.Project+p.Script, p.BufferLen, p.WorkersCount)
	if err != nil {
		return err
	}
	cache.Init(strs)
	return scheme.InsertMultiValues(strs)
}
func (p PersistedPool) Update(file *multipart.File) (err error) {
	scheme := database.PoolScheme{Project: p.Project, Script: p.Script}
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
	cache, err := cache.GetCache(p.Project + p.Script)
	if err != nil {
		return err
	}
	err = cache.ReInit(strs)
	if err != nil {
		return err
	}
	return scheme.InsertMultiValues(strs)
}

func (p PersistedPool) AddValues(file *multipart.File) (err error) {
	jsonSlice, err := parser.CSVToJSON(*file)
	if err != nil {
		return err
	}
	strs := make([]string, 0)
	for _, v := range jsonSlice {
		strs = append(strs, string(v))
	}
	cache, err := cache.GetCache(p.Project + p.Script)
	if err != nil {
		return err
	}
	cache.AddValues(strs)
	return database.PoolScheme{Project: p.Project, Script: p.Script}.InsertMultiValues(strs)
}

func (p PersistedPool) Delete() (err error) {
	defer func() {
		recover()
	}()
	cache, err := cache.GetCache(p.Project + p.Script)
	if err != nil {
		return err
	}
	cache.Delete()
	return database.PoolScheme{Project: p.Project, Script: p.Script}.DropTable()
}

//InitPoolFromDB - Datapool initialization function.
//gets all data from the database, based on the project name and script name fields,
//and initializes the data cache and the upload channel.
func (p PersistedPool) InitPoolFromDB() (err error) {
	data, err := database.PoolScheme{Project: p.Project, Script: p.Script}.GetPool()
	if err != nil {
		println(err.Error())
		return err
	}
	cache, err := cache.CreateDefaultCache(p.Project+p.Script, p.BufferLen, p.WorkersCount)
	if err != nil {
		return err
	}
	cache.Init(data)
	return nil
}

func (p PersistedPool) GetValue() (string, error) {
	cache, err := cache.GetCache(p.Project + p.Script)
	if err != nil {
		return "", err
	}
	return <-cache.CH, nil
}

func newScheme(ps database.PoolScheme) (err error) {
	err = ps.AddRelationsSchemeScript()
	if err != nil {
		return err
	}
	err = ps.CreateScheme()
	if err != nil {
		return err
	}
	err = ps.CreateTable()
	if err != nil {
		return err
	}
	return nil
}
