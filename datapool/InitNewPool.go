package datapool

import (
	"log"
	"strconv"
	"time"

	"github.com/allegro/bigcache"
	"github.com/matscus/ammunition/database"
)

func (d Datapool) InitNewPool() (err error) {
	data, err := database.GetPool(d.ProjectName, d.ScriptName)
	if err != nil {
		println(err.Error())
		return err
	}
	cacheConfig := bigcache.DefaultConfig(5 * time.Minute)
	cacheConfig.CleanWindow = 0 * time.Minute
	cache, _ := bigcache.NewBigCache(cacheConfig)
	cacheMap.Store(d.ProjectName+d.ScriptName, cache)
	defer cache.Close()
	for k, v := range *data {
		cache.Set(strconv.Itoa(k), []byte(v))
	}
	ch := make(chan string, 100)
	go func() {
		for {
			for i := 0; i < cache.Len(); i++ {
				d, err := cache.Get(strconv.Itoa(i))
				if err != nil {
					log.Println(err)
				}
				ch <- string(d)
			}
		}

	}()
	chanMap.Store(d.ProjectName+d.ScriptName, ch)
	return nil
}
