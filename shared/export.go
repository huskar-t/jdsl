package main

import "C"
import (
	"encoding/json"
	"github.com/huskar-t/jsontemplate/logic"
)

type CParseResult struct {
	Result *logic.ParseResult `json:"result"`
	Sql    *logic.CreatSql    `json:"sql"`
	Error  string             `json:"error"`
}

//export ParseJson
func ParseJson(config string) *C.char {
	var result CParseResult
	var parseConfig logic.Config
	err := json.Unmarshal([]byte(config), &parseConfig)
	if err != nil {
		result.Error = err.Error()
		d, _ := json.Marshal(result)
		return C.CString(string(d))
	}
	parseResult, createSql, err := logic.ParseJson(&parseConfig)
	if err != nil {
		result.Error = err.Error()
		d, _ := json.Marshal(result)
		return C.CString(string(d))
	}
	result.Result = parseResult
	result.Sql = createSql
	s, err := json.Marshal(result)
	if err != nil {
		d, _ := json.Marshal(CParseResult{Error: err.Error()})
		return C.CString(string(d))
	}
	return C.CString(string(s))
}

type CReadResult struct {
	SqlList []string `json:"sqlList"`
	Error   string   `json:"error"`
}

//export Read
func Read(parseJsonResult string, msg string, extra string) *C.char {
	var result CReadResult
	var parseResult logic.ParseResult
	err := json.Unmarshal([]byte(parseJsonResult), &parseResult)
	if err != nil {
		result.Error = err.Error()
		d, _ := json.Marshal(result)
		return C.CString(string(d))
	}
	var extraMsg map[string]string
	if extra != "" {
		err = json.Unmarshal([]byte(extra), &extraMsg)
		if err != nil {
			result.Error = err.Error()
			d, _ := json.Marshal(result)
			return C.CString(string(d))
		}
	}
	sqlList, err := logic.Read(&parseResult, msg, extraMsg)
	if err != nil {
		result.Error = err.Error()
		d, _ := json.Marshal(result)
		return C.CString(string(d))
	}
	result.SqlList = sqlList
	s, err := json.Marshal(result)
	if err != nil {
		d, _ := json.Marshal(CParseResult{Error: err.Error()})
		return C.CString(string(d))
	}
	return C.CString(string(s))
}

func main() {

}
