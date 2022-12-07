package cache

import (
	"bytes"
	"fmt"
	"sync"
	"time"

	"ammunition/config"
	"ammunition/metrics"

	"github.com/allegro/bigcache"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

var (
	temporaryCaches sync.Map
)

type Temporary struct {
	Name         string `json:"name"`
	Worker       int    `json:"worker"`
	LiveDuration int64  `json:"live_duration"`
	Chans        []struct {
		Name      string `json:"name"`
		BufferLen int    `json:"buffer_len"`
	} `json:"chans"`
	ChanMap    sync.Map
	BigCache   *bigcache.BigCache
	Context    context.Context
	CancelFunc context.CancelFunc
}

func (t *Temporary) New() (err error) {
	c := bigcache.DefaultConfig(time.Duration(config.Config.Temporary.LifeWindow) * time.Minute)
	c.CleanWindow = time.Duration(config.Config.Temporary.CleanWindow) * time.Minute
	c.HardMaxCacheSize = config.Config.Temporary.HardMaxCacheSize
	c.MaxEntrySize = config.Config.Temporary.MaxEntrySize
	c.Shards = config.Config.Temporary.Shards
	c.Verbose = config.Config.Temporary.Verbose
	t.BigCache, err = bigcache.NewBigCache(c)
	if err != nil {
		return err
	}
	if t.LiveDuration == 0 {
		t.LiveDuration = 1
	}
	t.Context, t.CancelFunc = context.WithTimeout(context.Background(), time.Minute*time.Duration(t.LiveDuration))
	temporaryCaches.Store(t.Name, t)
	for _, v := range t.Chans {
		t.ChanMap.Store(v.Name, make(chan []byte, v.BufferLen))
	}
	for i := 0; i < t.Worker; i++ {
		go temporaryWorker(t.Context, t.Name)
	}
	go getTemporaryCacheMetrics(t)
	go cleaner(t)
	log.Infof("Init temorary cache %s is completed, cache lives in %d minutes", t.Name, t.LiveDuration)
	return err
}

func SetTemporaryValue(cacheName string, queue string, key string, values []byte) error {
	tempCache, ok := temporaryCaches.Load(cacheName)
	if !ok {
		return fmt.Errorf("Cache %s not found", cacheName)
	}
	buf := bytes.Buffer{}
	buf.Write(make([]byte, (16 - len(queue))))
	buf.Write([]byte(queue))
	buf.Write(values)
	return tempCache.(*Temporary).BigCache.Set(key, buf.Bytes())
}

func GetTemporaryValue(cacheName string, queue string) []byte {
	tempCache, ok := temporaryCaches.Load(cacheName)
	if !ok {
		return []byte("{\"Message\":\"Cache not found\"}") //"Cache %s not found", cacheName
	}
	temporaryChan, ok := tempCache.(*Temporary).ChanMap.Load(queue)
	if !ok {
		return []byte("{\"Message\":\"Chan not found\"}")
	}
	select {
	case res, ok := <-temporaryChan.(chan []byte):
		if ok {
			return res
		} else {
			return []byte("{\"Message\":\"Chan is close\"}")
		}
	default:
		return []byte("{\"Message\":\"Chan is empty\"}")
	}
}

func DeleteTemporaryCache(cacheName string) error {
	tempCache, ok := temporaryCaches.Load(cacheName)
	if !ok {
		return errors.Errorf("Cache %s not found", cacheName)
	}
	tempCache.(*Temporary).CancelFunc()
	err := tempCache.(*Temporary).BigCache.Close()
	if err != nil {
		return err
	}
	tempCache.(*Temporary).ChanMap.Range(func(key, value interface{}) bool {
		close(value.(chan []byte))
		return true
	})
	temporaryCaches.Delete(cacheName)
	return nil
}

func cleaner(t *Temporary) {
	for {
		select {
		case <-t.Context.Done():
			t.ChanMap.Range(func(key, value interface{}) bool {
				_, ok := <-value.(chan []byte)
				if ok {
					close(value.(chan []byte))
				}
				return true
			})
			temporaryCaches.Delete(t.Name)
		default:
			time.Sleep(1 * time.Second)
		}
	}
}

func temporaryWorker(ctx context.Context, cacheName string) {
	log.Infof("Start worker from cache %s", cacheName)
	defer func() {
		if err := recover(); err != nil {
			log.Error("Temporary worker recover panic ", err)
		}
		select {
		case <-ctx.Done():
			return
		default:
			go temporaryWorker(ctx, cacheName)
		}
	}()
	tempCache, ok := temporaryCaches.Load(cacheName)
	if !ok {
		log.Panic("Cache", cacheName, "not found ")
	}
	var firstBytes string
	for {
		select {
		case <-ctx.Done():
			log.Printf("End worker from %s", cacheName)
			return
		default:
			iterator := tempCache.(*Temporary).BigCache.Iterator()
			start := time.Now()
			for iterator.SetNext() {
				entry, err := iterator.Value()
				if err != nil {
					log.Error("Worker iterarion ", err)
				} else {
					firstBytes = string(bytes.Trim(entry.Value()[0:16], "\x00"))
					ch, ok := tempCache.(*Temporary).ChanMap.Load(firstBytes)
					if ok {
						ch.(chan []byte) <- entry.Value()[16:]
						tempCache.(*Temporary).BigCache.Delete(entry.Key())
					} else {
						log.Panic("Worker panir: chan ", firstBytes, " not found")
					}
					metrics.WorkerDuration.WithLabelValues(cacheName).Observe(float64(time.Since(start).Milliseconds()))
				}
			}
		}
	}
}

func getTemporaryCacheMetrics(t *Temporary) {
	log.Infof("Temporary metrics init from cache %s is completed", t.Name)
	metrics.WorkerCount.WithLabelValues(t.Name).Inc()
	metrics.CacheCount.WithLabelValues("in-memory", t.Name).Set(1)
	for {
		select {
		case <-t.Context.Done():
			log.Printf("End worker from %s", t.Name)
			return
		default:
			metrics.CacheLen.WithLabelValues("in-memory", t.Name).Set(float64(t.BigCache.Len()))
			metrics.CacheCap.WithLabelValues("in-memory", t.Name).Set(float64(t.BigCache.Capacity()))
			time.Sleep(10 * time.Second)
		}
	}
}
