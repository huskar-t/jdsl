package logic

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/huskar-t/jsontemplate/scout"
	"strings"
)

func ParseJson(c *Config) (*ParseResult, *CreatSql, error) {
	var result = &ParseResult{
		Value:      c.Value,
		DBName:     c.DBName,
		TableName:  c.TableName,
		STableName: c.STableName,
	}

	template := c.PayloadTemplate
	var tagsStr []string
	var tagsTemplateStr []string
	var tagsTemplateVerMap = map[string]bool{}
	for tagName, tagTemplate := range c.TagsMap {
		tagsStr = append(tagsStr, tagName+" binary(64)")
		tagsTemplateStr = append(tagsTemplateStr, tagTemplate)
		if isTemplate(tagTemplate) {
			tagsTemplateVerMap[tagTemplate] = true
		}
	}
	createDB := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s KEEP %d", c.DBName, c.Keep)
	superFloat := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s.%s_double (ts timestamp, value Double) TAGS (%s)`, c.DBName, c.STableName, strings.Join(tagsStr, ","))
	superString := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s.%s_string (ts timestamp, value NCHAR(64)) TAGS (%s)`, c.DBName, c.STableName, strings.Join(tagsStr, ","))

	var createSql = &CreatSql{
		DBSql:          createDB,
		TableFloatSql:  superFloat,
		TableStringSql: superString,
	}
	result.TagsTemplateStr = tagsTemplateStr
	var jsonResult interface{}
	err := json.Unmarshal([]byte(template), &jsonResult)
	if err != nil {
		return nil, nil, err
	}
	var tempChars []byte
	var templateVers []string
	var start = false
	for i := 0; i < len(template); i++ {
		if start {
			tempChars = append(tempChars, template[i])
		}
		if template[i] == '$' {
			if !start {
				start = true
				tempChars = append(tempChars, template[i])
			} else {
				start = false
				templateVers = append(templateVers, string(tempChars))
				tempChars = nil
			}
		}
	}
	if start {
		return nil, nil, errors.New("template error with wrong variable")
	}
	tagsTemplatesMap := map[string][]string{}
	var ns []scout.Found
	for _, s := range templateVers {
		s2 := scout.New(s, jsonResult)
		parseResult, err := s2.DoSearch()
		if err != nil {
			return nil, nil, err
		}
		if len(parseResult) != 1 {
			return nil, nil, errors.New("repetitive definition")
		}
		pr := parseResult[0]
		if tagsTemplateVerMap[pr.Name] {
			tagsTemplatesMap[pr.Name] = getTemplate(pr.Path)
		}
		ns = append(ns, pr)
	}
	result.NS = ns
	tempRelation := map[string]bool{}
	for tag, tagRelation := range tagsTemplatesMap {
		tempRelation[tag] = true
		for _, relation := range tagRelation {
			tempRelation[relation] = true
		}
	}

	return result, createSql, err
}
