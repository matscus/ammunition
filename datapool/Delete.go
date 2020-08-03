package datapool

import (
	"github.com/allegro/bigcache"
	"github.com/matscus/ammunition/database"
)

func (d Datapool) Delete() (err error) {
	ch, ok := chanMap.Load(d.ProjectName + d.ScriptName)
	if ok {
		close(ch.(chan string))
	}
	cache, ok := cacheMap.Load(d.ProjectName + d.ScriptName)
	if ok {
		cache.(*bigcache.BigCache).Close()
	}
	return database.DropPool(d.ProjectName, d.ScriptName)
}
