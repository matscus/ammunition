package database

import (
	"bytes"
	"errors"
)

func InsertMultiValuePool(scheme string, table string, data []string) error {
	var buf bytes.Buffer
	buf.WriteString("INSERT INTO " + scheme + "." + table + " (pool) VALUES ")
	l := len(data)
	for i := 0; i < l; i++ {
		buf.WriteString("('" + data[i] + "')")
		if i < l-1 {
			buf.WriteString(",")
		}
	}
	_, err := DB.Exec(buf.String())
	if err != nil {
		return errors.New("Insert Multi Value Pool error: " + err.Error())
	}
	return nil
}
