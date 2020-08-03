package datapool

import (
	"errors"
	"log"
	"strconv"

	"github.com/allegro/bigcache"
	"github.com/matscus/ammunition/database"
)

func (d Datapool) ReInitPool() (err error) {
	c, ok := cacheMap.Load(d.ProjectName + d.ScriptName)
	if ok {
		data, err := database.GetPool(d.ProjectName, d.ScriptName)
		if err != nil {
			log.Println(err.Error())
			return err
		}
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
		chanMap.Store(d.ProjectName+d.ScriptName, ch)
	} else {
		return errors.New("Cache not found")
	}
	return nil
}
