package pool

import (
	"errors"
	"time"

	"github.com/matscus/ammunition/cache"
)

type TemporaryPool struct {
	Project   string        `json:"project"`
	Name      string        `json:"script"`
	Retantion time.Duration `json:"retantion"`
	BufferLen int           `json:"bufferlen,omitempty"`
	Workers   int           `json:"workers,omitempty"`
}

func (t TemporaryPool) Create() (err error) {
	if cache.CheckCache(t.Project + t.Name) {
		return errors.New("Pool is exist")
	} else {

	}
	cache, err := cache.New(t.Project+t.Name, t.BufferLen, t.Workers, t.Retantion, 1*time.Minute)
	if err != nil {
		return errors.New("Func Create - CreatePersistedCache error " + err.Error())
	}
	cache.Init(strs)
	// cache, err := cache.CreateDefaultCache(p.Project+p.Script, p.BufferLen, p.WorkersCount)
	// if err != nil {
	// 	return err
	// }
	// cache.Init(strs)
	// return scheme.InsertMultiValues(strs)
	return nil
}

func (p TemporaryPool) CheckPool() (ok bool) {
	// jsonSlice, err := parser.CSVToJSON(*file)
	// if err != nil {
	// 	return err
	// }
	// strs := make([]string, 0)
	// for _, v := range jsonSlice {
	// 	strs = append(strs, string(v))
	// }
	// scheme := database.PoolScheme{Project: p.Project, Script: p.Script}
	// err = newScheme(scheme)
	// if err != nil {
	// 	return err
	// }
	// cache, err := cache.CreateDefaultCache(p.Project+p.Script, p.BufferLen, p.WorkersCount)
	// if err != nil {
	// 	return err
	// }
	// cache.Init(strs)
	// return scheme.InsertMultiValues(strs)
	return false
}
