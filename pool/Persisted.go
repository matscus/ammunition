package pool

import (
	"errors"
	"mime/multipart"
	"time"

	"github.com/matscus/ammunition/cache"
	"github.com/matscus/ammunition/database"
	"github.com/matscus/ammunition/parser"
)

var (
	databaseSchemes []database.PoolScheme
)

type PersistedPool struct {
	Project   string `json:"project"`
	Name      string `json:"name"`
	BufferLen int    `json:"bufferlen,omitempty"`
	Workers   int    `json:"workers,omitempty"`
}

func InitAllPersistedPools() (err error) {
	databaseSchemes, err = database.GetAllPools()
	if err != nil {
		return errors.New("InitAllPersistedPools - GetAllPoolserror: " + err.Error())
	}
	for _, v := range databaseSchemes {
		go PersistedPool{Project: v.Project, Name: v.Name, BufferLen: v.BufferLen, Workers: v.Workers}.InitPoolFromDB()
	}
	return nil
}

func (p PersistedPool) Create(file *multipart.File) (err error) {
	jsonSlice, err := parser.CSVToJSON(*file)
	if err != nil {
		return errors.New("PersistedPool Create - CSVToJSON error " + err.Error())
	}
	strs := make([]string, 0)
	for _, v := range jsonSlice {
		strs = append(strs, string(v))
	}
	scheme := database.PoolScheme{Project: p.Project, Name: p.Name, BufferLen: p.BufferLen, Workers: p.Workers}
	err = newScheme(scheme)
	if err != nil {
		return errors.New("PersistedPool Create - newScheme error " + err.Error())
	}
	cache := cache.Cache{Name: p.Project + p.Name, BufferLen: p.BufferLen, WorkersCount: p.Workers, Life: 24 * time.Hour, Clean: 0}
	err = cache.PersistedInit(strs)
	if err != nil {
		return errors.New("PersistedPool Create - CreatePersistedCache error " + err.Error())
	}
	err = scheme.InsertMultiValues(strs)
	if err != nil {
		return errors.New("PersistedPool Create - InsertMultiValues error " + err.Error())
	}
	return nil
}
func (p PersistedPool) Update(file *multipart.File) (err error) {
	scheme := database.PoolScheme{Project: p.Project, Name: p.Name}
	jsonSlice, err := parser.CSVToJSON(*file)
	if err != nil {
		return errors.New("PersistedPool Update - CSVToJSON error " + err.Error())
	}
	strs := make([]string, 0)
	for _, v := range jsonSlice {
		strs = append(strs, string(v))
	}
	oldCache, err := cache.GetPersistedCache(p.Project + p.Name)
	if err != nil {
		return errors.New("PersistedPool Update - GetPersistedCache error " + err.Error())
	}
	err = oldCache.PersistedDelete()
	if err != nil {
		return errors.New("PersistedPool Update - Delete error " + err.Error())
	}
	cache := cache.Cache{Name: p.Project + p.Name, BufferLen: p.BufferLen, WorkersCount: p.Workers, Life: 24 * time.Hour, Clean: 0}
	err = cache.PersistedInit(strs)
	if err != nil {
		return errors.New("PersistedPool Update - CreatePersistedCache error " + err.Error())
	}
	err = scheme.ClearTable()
	if err != nil {
		return errors.New("PersistedPool Update - ClearTable error " + err.Error())
	}
	err = scheme.InsertMultiValues(strs)
	if err != nil {
		return errors.New("PersistedPool Update - InsertMultiValues error " + err.Error())
	}
	return nil
}

func (p PersistedPool) AddValues(file *multipart.File) (err error) {
	jsonSlice, err := parser.CSVToJSON(*file)
	if err != nil {
		return errors.New("PersistedPoolAddValues - CSVToJSON error " + err.Error())
	}
	strs := make([]string, 0)
	for _, v := range jsonSlice {
		strs = append(strs, string(v))
	}
	cache, err := cache.GetPersistedCache(p.Project + p.Name)
	if err != nil {
		return errors.New("PersistedPool AddValues - GetPersistedCache error " + err.Error())
	}
	cache.SetValues(strs)
	err = database.PoolScheme{Project: p.Project, Name: p.Name}.InsertMultiValues(strs)
	if err != nil {
		return errors.New("PersistedPool AddValues - InsertMultiValues error " + err.Error())
	}
	return nil
}

func (p PersistedPool) Delete() (err error) {
	cache, err := cache.GetPersistedCache(p.Project + p.Name)
	if err != nil {
		return errors.New("PersistedPool Delete - GetPersistedCache error " + err.Error())
	}
	err = cache.PersistedDelete()
	if err != nil {
		return errors.New("PersistedPool Delete - Delete error " + err.Error())
	}
	pool := database.PoolScheme{Project: p.Project, Name: p.Name}
	err = pool.DropTable()
	if err != nil {
		return errors.New("PersistedPoolDelete - DropTable error " + err.Error())
	}
	err = pool.DeleteRelationsSchemeScript()
	if err != nil {
		return errors.New("PersistedPool Delete - DeleteRelationsSchemeScript error " + err.Error())
	}
	return nil
}

//InitPoolFromDB - Datapool initialization function.
//gets all data from the database, based on the project name and script name fields,
//and initializes the data cache and the upload channel.
func (p PersistedPool) InitPoolFromDB() (err error) {
	data, err := database.PoolScheme{Project: p.Project, Name: p.Name}.GetPool()
	if err != nil {
		return errors.New("PersistedPool InitPoolFromDB - GetPool error " + err.Error())
	}
	cache := cache.Cache{Name: p.Project + p.Name, BufferLen: p.BufferLen, WorkersCount: p.Workers, Life: 24 * time.Hour, Clean: 0}
	err = cache.PersistedInit(data)
	if err != nil {
		return errors.New("PersistedPool InitPoolFromDB - PersistedInit error " + err.Error())
	}
	return nil
}

func (p PersistedPool) GetValue() (string, error) {
	persistedCache, err := cache.GetChan(p.Project + p.Name)
	if err != nil {
		return "", errors.New("PersistedPool - GetValue error " + err.Error())
	}
	return <-persistedCache, nil
}

func newScheme(ps database.PoolScheme) (err error) {
	defer func() {
		recover()
	}()
	err = ps.AddRelationsSchemeScript()
	if err != nil {
		return errors.New("AddRelationsSchemeScript error " + err.Error())
	}
	ps.CreateScheme()
	err = ps.CreateTable()
	if err != nil {
		return errors.New("CreateTable error " + err.Error())
	}
	return nil
}
