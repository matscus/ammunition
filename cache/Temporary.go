package cache

import "errors"

func (c Cache) TemporaryInit() {
	c.CH = make(chan string, c.BufferLen)
	ChanMap.Store(c.Name, c.CH)
	for i := 0; i < c.WorkersCount; i++ {
		go c.RunWorker()
	}
}

func CheckTemporaryCache(name string) bool {
	_, ok := TemporaryCacheMap.Load(name)
	if ok {
		return true
	}
	return false
}

func (c Cache) TemporaryDelete() error {
	close(c.CH)
	TemporaryCacheMap.Delete(c.Name)
	ChanMap.Delete(c.CH)
	return c.BigCache.Close()
}

func GetTemporaryCache(name string) (Cache, error) {
	cache, ok := TemporaryCacheMap.Load(name)
	if ok {
		return cache.(Cache), nil
	}
	return cache.(Cache), errors.New("Cache not found")
}
