package datapool

import (
	"log"
	"sync"

	"github.com/matscus/ammunition/database"
)

var (
	cacheMap      sync.Map
	chanMap       sync.Map
	projectScheme []database.ProjectScheme
	chMap         = make(map[string]*chan string)
)

func init() {
	err := getProjectScheme()
	if err != nil {
		log.Println(err)
	}
	for _, v := range projectScheme {
		go Datapool{ProjectName: v.ProjectName, ScriptName: v.ScriptName}.InitNewPool()
	}
}
func getProjectScheme() error {
	var err error
	projectScheme, err = database.GetAllPoolName()
	if err != nil {
		return err
	}
	return nil
}
