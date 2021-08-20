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
	if cache.CheckTemporaryCache(t.Project + t.Name) {
		return errors.New("Pool is exist")
	}
	cache := cache.Cache{Name: t.Project + t.Name, BufferLen: t.BufferLen, WorkersCount: t.Workers, Life: 24 * time.Hour, Clean: 0}
	cache.TemporaryInit()
	return nil
}
func (p TemporaryPool) GetValue() (string, error) {
	temporaryCache, err := cache.GetChan(p.Project + p.Name)
	if err != nil {
		return "", errors.New("TemporaryPool - GetValue error " + err.Error())
	}
	return <-temporaryCache, nil
}

func (p TemporaryPool) Delete() (err error) {
	cache, err := cache.GetTemporaryCache(p.Project + p.Name)
	if err != nil {
		return errors.New("PersistedPool Delete - GetPersistedCache error " + err.Error())
	}
	err = cache.TemporaryDelete()
	if err != nil {
		return errors.New("PersistedPool Delete - Delete error " + err.Error())
	}
	return nil
}
