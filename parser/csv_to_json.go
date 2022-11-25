package parser

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"mime/multipart"
	"strconv"
	"strings"
)

//CSVToJSON - func from parse csv file and convert data to Json struct
//return slice byte or error
func CSVToJSON(file multipart.File) (rs [][]byte, err error) {
	reader := csv.NewReader(file)
	content, err := reader.ReadAll()
	if len(content) < 1 {
		return nil, errors.New("File is empty")
	}
	csvHeaders := make([]string, 0, len(content[0]))
	for _, h := range content[0] {
		csvHeaders = append(csvHeaders, h)
	}
	content = content[1:]
	var buffer bytes.Buffer
	for _, d := range content {
		buffer.WriteString("{")
		for j, y := range d {
			buffer.WriteString(`"` + csvHeaders[j] + `":`)
			_, fErr := strconv.ParseFloat(y, 32)
			_, bErr := strconv.ParseBool(y)
			if fErr == nil {
				buffer.WriteString(y)
			} else if bErr == nil {
				buffer.WriteString(strings.ToLower(y))
			} else {
				buffer.WriteString((`"` + y + `"`))
			}
			if j < len(d)-1 {
				buffer.WriteString(",")
			}
		}
		buffer.WriteString("}")
		rawMessage := json.RawMessage(buffer.String())
		x, _ := json.MarshalIndent(rawMessage, "", "  ")
		rs = append(rs, x)
		buffer.Reset()
	}
	return rs, nil
}
