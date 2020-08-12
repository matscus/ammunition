package datapool

import (
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/allegro/bigcache"
)

func initCashe(name string, data *[]string) {
	cacheConfig := bigcache.DefaultConfig(5 * time.Minute)
	cacheConfig.CleanWindow = 0 * time.Minute
	cache, _ := bigcache.NewBigCache(cacheConfig)
	cacheMap.Store(name, cache)
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
	chanMap.Store(name, ch)
}

func reInitCashe(name string, data *[]string) (err error) {
	c, ok := cacheMap.Load(name)
	if ok {
		cache := c.(*bigcache.BigCache)
		err = cache.Reset()
		if err != nil {
			return err
		}
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
		chanMap.Store(name, ch)
	} else {
		return errors.New("Cache not found")
	}
	return nil
}

func addValueInCashe(name string, data *[]string) (err error) {
	c, ok := cacheMap.Load(name)
	if ok {
		cache := c.(*bigcache.BigCache)
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
		chanMap.Store(name, ch)
	} else {
		return errors.New("Cache not found")
	}
	return nil
}
