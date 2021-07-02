package pool

type TemporaryPool struct {
	Project   string      `json:"project"`
	Script    string      `json:"script"`
	Operation string      `json:"operation"`
	Data      interface{} `json:"data"`
}

func (p TemporaryPool) Create() (err error) {
	// jsonSlice, err := parser.CSVToJSON(*file)
	// if err != nil {
	// 	return err
	// }
	// strs := make([]string, 0)
	// for _, v := range jsonSlice {
	// 	strs = append(strs, string(v))
	// }
	// scheme := database.PoolScheme{Project: p.Project, Script: p.Script}
	// err = newScheme(scheme)
	// if err != nil {
	// 	return err
	// }
	// cache, err := cache.CreateDefaultCache(p.Project+p.Script, p.BufferLen, p.WorkersCount)
	// if err != nil {
	// 	return err
	// }
	// cache.Init(strs)
	// return scheme.InsertMultiValues(strs)
	return nil
}

func (p TemporaryPool) CheckPool() (ok bool) {
	// jsonSlice, err := parser.CSVToJSON(*file)
	// if err != nil {
	// 	return err
	// }
	// strs := make([]string, 0)
	// for _, v := range jsonSlice {
	// 	strs = append(strs, string(v))
	// }
	// scheme := database.PoolScheme{Project: p.Project, Script: p.Script}
	// err = newScheme(scheme)
	// if err != nil {
	// 	return err
	// }
	// cache, err := cache.CreateDefaultCache(p.Project+p.Script, p.BufferLen, p.WorkersCount)
	// if err != nil {
	// 	return err
	// }
	// cache.Init(strs)
	// return scheme.InsertMultiValues(strs)
	return false
}
