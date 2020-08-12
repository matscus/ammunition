package datapool

import (
	"log"
	"sync"

	"github.com/matscus/ammunition/database"
)

var (
	cacheMap        sync.Map
	chanMap         sync.Map
	databaseSchemes []database.DatabaseScheme
	chMap           = make(map[string]*chan string)
)

func init() {
	err := getProjectScheme()
	if err != nil {
		log.Println(err)
	}
	for _, v := range databaseSchemes {
		go Datapool{ProjectName: v.ProjectName, ScriptName: v.ScriptName}.InitPoolFromDB()
	}
}
func getProjectScheme() error {
	var err error
	databaseSchemes, err = database.GetAllPoolName()
	if err != nil {
		return err
	}
	return nil
}
