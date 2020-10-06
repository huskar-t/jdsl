package logic

import "github.com/huskar-t/jsontemplate/scout"

type Config struct {
	DBName          string            `json:"dbName"`
	STableName      string            `json:"sTableName"`
	TableName       string            `json:"tableName"`
	Keep            int               `json:"keep"`
	TagsMap         map[string]string `json:"tagsMap"`
	PayloadTemplate string            `json:"payloadTemplate"`
	Value           string            `json:"value"`
}

type Value struct {
	Name      string
	Tags      []string
	TagsValue []interface{}
	Values    []interface{}
}

type ParseResult struct {
	Value           string        `json:"value"`
	TagsTemplateStr []string      `json:"tagsTemplateStr"`
	NS              []scout.Found `json:"ns"`
	DBName          string        `json:"dbName"`
	TableName       string        `json:"tableName"`
	STableName      string        `json:"sTableName"`
}

type CreatSql struct {
	DBSql          string `json:"dbSql"`
	TableFloatSql  string `json:"tableFloatSql"`
	TableStringSql string `json:"tableStringSql"`
}

func isTemplate(item string) bool {
	return item[0] == '$' && item[len(item)-1] == '$'
}

func getTemplate(path []string) []string {
	var templateList []string
	for i := 0; i < len(path); i++ {
		if isTemplate(path[i]) {
			templateList = append(templateList, path[i])
		}
	}
	return templateList
}

func getTemplateStr(str string) []string {
	var tempChars []byte
	var templateVers []string
	var start = false
	for i := 0; i < len(str); i++ {
		if start {
			tempChars = append(tempChars, str[i])
		}
		if str[i] == '$' {
			if !start {
				start = true
				tempChars = append(tempChars, str[i])
			} else {
				start = false
				templateVers = append(templateVers, string(tempChars))
				tempChars = nil
			}
		}
	}
	return templateVers
}
