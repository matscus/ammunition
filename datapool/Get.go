package datapool

import "errors"

func (d Datapool) Get() (string, error) {
	ch, ok := chanMap.Load(d.ProjectName + d.ScriptName)
	if ok {
		res := <-ch.(chan string)
		return res, nil
	}
	return "", errors.New("chanel is empty")
}
