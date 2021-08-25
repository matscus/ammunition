package cache

import (
	"errors"
	"mime/multipart"
	"strconv"
	"sync"
	"time"

	"github.com/allegro/bigcache"
	"github.com/matscus/ammunition/config"
	"github.com/matscus/ammunition/database"
	"github.com/matscus/ammunition/metrics"
	"github.com/matscus/ammunition/parser"
	log "github.com/sirupsen/logrus"
)

var (
	persistedCacheMap sync.Map
	chanMap           sync.Map
	databaseSchemes   []database.PoolScheme
)

type PersistedPool struct {
	Project   string `json:"project"`
	Name      string `json:"name"`
	BufferLen int    `json:"bufferlen,omitempty"`
	Workers   int    `json:"workers,omitempty"`
}

func init() {
	go getPersistCacheMetrics()
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
	cache := Cache{Name: p.Project + p.Name, BufferLen: p.BufferLen, WorkersCount: p.Workers, Life: 24 * time.Hour, Clean: 0}
	err = cache.persistedInit(strs)
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
	oldCache, err := getPersistedCache(p.Project + p.Name)
	if err != nil {
		return errors.New("PersistedPool Update - GetPersistedCache error " + err.Error())
	}
	err = oldCache.persistedDelete()
	if err != nil {
		return errors.New("PersistedPool Update - Delete error " + err.Error())
	}
	cache := Cache{Name: p.Project + p.Name, BufferLen: p.BufferLen, WorkersCount: p.Workers, Life: 24 * time.Hour, Clean: 0}
	err = cache.persistedInit(strs)
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
	cache, err := getPersistedCache(p.Project + p.Name)
	if err != nil {
		return errors.New("PersistedPool AddValues - GetPersistedCache error " + err.Error())
	}
	cache.setValues(strs)
	err = database.PoolScheme{Project: p.Project, Name: p.Name}.InsertMultiValues(strs)
	if err != nil {
		return errors.New("PersistedPool AddValues - InsertMultiValues error " + err.Error())
	}
	return nil
}

func (p PersistedPool) Delete() (err error) {
	cache, err := getPersistedCache(p.Project + p.Name)
	if err != nil {
		return errors.New("PersistedPool Delete - GetPersistedCache error " + err.Error())
	}
	err = cache.persistedDelete()
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
	cache := Cache{Name: p.Project + p.Name, BufferLen: p.BufferLen, WorkersCount: p.Workers, Life: 24 * time.Hour, Clean: 0}
	err = cache.persistedInit(data)
	if err != nil {
		return errors.New("PersistedPool InitPoolFromDB - PersistedInit error " + err.Error())
	}
	return nil
}

func (p PersistedPool) GetValue() (string, error) {
	persistedCache, err := getChan(p.Project + p.Name)
	if err != nil {
		return "", errors.New("PersistedPool - GetValue error " + err.Error())
	}
	return <-persistedCache, nil
}

func (c Cache) persistedInit(data []string) (err error) {
	config := config.DefaultConfig
	config.LifeWindow = c.Life
	config.CleanWindow = c.Clean
	c.BigCache, err = bigcache.NewBigCache(config)
	if err != nil {
		return err
	}
	persistedCacheMap.Store(c.Name, c)
	for k, v := range data {
		c.BigCache.Set(strconv.Itoa(k), []byte(v))
	}
	c.CH = make(chan string, c.BufferLen)
	chanMap.Store(c.Name, c.CH)
	for i := 0; i < c.WorkersCount; i++ {
		go c.runWorker()
	}
	return nil
}

func (c Cache) persistedDelete() error {
	close(c.CH)
	persistedCacheMap.Delete(c.Name)
	chanMap.Delete(c.CH)
	return c.BigCache.Close()
}

func (c Cache) setValues(data []string) {
	for k, v := range data {
		c.BigCache.Set(strconv.Itoa(k), []byte(v))
	}
}

func (c Cache) runWorker() {
	defer func() {
		recover()
	}()
	metrics.WorkerCount.WithLabelValues(c.Name).Inc()
	for {
		for i := 0; i < c.BigCache.Len(); i++ {
			start := time.Now()
			d, err := c.BigCache.Get(strconv.Itoa(i))
			if err != nil {
				log.Println("Worker get values error: ", err)
			}
			c.CH <- string(d)
			metrics.WorkerDuration.WithLabelValues(c.Name).Observe(float64(time.Since(start).Milliseconds()))
		}
	}
}

func getChan(name string) (ch chan string, err error) {
	tempChan, ok := chanMap.Load(name)
	if ok {
		return tempChan.(chan string), nil
	}
	return tempChan.(chan string), errors.New("Chan not found")
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

func checkPersistedCache(name string) bool {
	_, ok := persistedCacheMap.Load(name)
	if ok {
		return true
	}
	return false
}

func getPersistedCache(name string) (Cache, error) {
	cache, ok := persistedCacheMap.Load(name)
	if ok {
		return cache.(Cache), nil
	}
	return cache.(Cache), errors.New("Cache not found")
}

func getPersistCacheMetrics() {
	defer func() {
		recover()
	}()
	var i float64
	for {
		persistedCacheMap.Range(func(k, v interface{}) bool {
			metrics.CacheLen.WithLabelValues("persist", k.(string)).Set(float64(v.(Cache).BigCache.Len()))
			metrics.CacheCap.WithLabelValues("persist", k.(string)).Set(float64(v.(Cache).BigCache.Capacity()))
			i++
			return true
		})
		metrics.CacheCount.WithLabelValues("persist", "test").Set(i)
		i = 0
		time.Sleep(10 * time.Second)
	}
}
