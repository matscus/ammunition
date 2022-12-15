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
		t.LiveDuration = 10
	}
	t.Context, t.CancelFunc = context.WithTimeout(context.Background(), time.Minute*time.Duration(t.LiveDuration))
	temporaryCaches.Store(t.Name, t)
	for _, v := range t.Chans {
		t.ChanMap.Store(v.Name, make(chan []byte, v.BufferLen))
	}
	for i := 0; i < t.Worker; i++ {
		go temporaryWorker(t)
	}
	temporaryCaches.Range(func(key, value interface{}) bool {
		return true
	})
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

func GetTemporaryValue(cacheName string, queue string) (res []byte,err error) {
	tempCache, ok := temporaryCaches.Load(cacheName)
	if !ok {
		return nil, errors.New("cache not found")
	}
	temporaryChan, ok := tempCache.(*Temporary).ChanMap.Load(queue)
	if !ok {
		return nil, errors.New("chan not found")
	}
	select {
	case res, ok := <-temporaryChan.(chan []byte):
		if ok {
			return res, nil
		} else {
			return nil, errors.New("chan is close")
		}
	default:
		return nil, errors.New("chan is empty")
	}
}

func DeleteTemporaryCache(cacheName string) error {
	tempCache, ok := temporaryCaches.Load(cacheName)
	if !ok {
		return errors.Errorf("Cache %s not found", cacheName)
	}
	tempCache.(*Temporary).CancelFunc()
	return nil
}

func cleaner(t *Temporary) {
	for {
		select {
		case <-t.Context.Done():
			temporaryCaches.Delete(t.Name)
			t.ChanMap.Range(func(key, value interface{}) bool {
				close(value.(chan []byte))
				return true
			})
			return
		default:
			time.Sleep(1 * time.Second)
		}
	}
}

func temporaryWorker(t *Temporary) {
	log.Infof("Start worker from cache %s", t.Name)
	defer func() {
		if err := recover(); err != nil {
			log.Error("Temporary worker recover panic ", err)
		}
		select {
		case <-t.Context.Done():
			return
		default:
			go temporaryWorker(t)
		}
	}()
	tempCache, ok := temporaryCaches.Load(t.Name)
	if !ok {
		log.Panic("Worker not found  cache ", t.Name)
	}
	var firstBytes string
	for {
		select {
		case <-t.Context.Done():
			log.Printf("End worker from %s", t.Name)
			metrics.WorkerDuration.DeleteLabelValues(t.Name)
			return
		default:
			iterator := tempCache.(*Temporary).BigCache.Iterator()
			start := time.Now()
			for iterator.SetNext() {
				entry, err := iterator.Value()
				if err != nil {
					log.Errorf("Worker: iterarion %s", err)
				} else {
					firstBytes = string(bytes.Trim(entry.Value()[0:16], "\x00"))
					ch, ok := tempCache.(*Temporary).ChanMap.Load(firstBytes)
					if ok {
						ch.(chan []byte) <- entry.Value()[16:]
						tempCache.(*Temporary).BigCache.Delete(entry.Key())
					} else {
						log.Warnf("Worker: chan %s not found", firstBytes)
					}
					metrics.WorkerDuration.WithLabelValues(t.Name).Observe(float64(time.Since(start).Milliseconds()))
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
			metrics.WorkerCount.DeleteLabelValues(t.Name)
			metrics.CacheCount.DeleteLabelValues("in-memory", t.Name)
			metrics.CacheLen.DeleteLabelValues("in-memory", t.Name)
			metrics.CacheCap.DeleteLabelValues("in-memory", t.Name)
			t.ChanMap.Range(func(key, value interface{}) bool {
				metrics.ChanLen.DeleteLabelValues(t.Name,key.(string))
				return true
			})
			return
		default:
			metrics.CacheLen.WithLabelValues("in-memory", t.Name).Set(float64(t.BigCache.Len()))
			metrics.CacheCap.WithLabelValues("in-memory", t.Name).Set(float64(t.BigCache.Capacity()))
			t.ChanMap.Range(func(key, value interface{}) bool {
				metrics.ChanLen.WithLabelValues(t.Name,key.(string)).Set(float64(len(value.(chan []byte))))
				return true
			})
			time.Sleep(10 * time.Second)
		}
	}
}
