package logic

import (
	"errors"
	"fmt"
	"github.com/huskar-t/jsontemplate/scout"
	"github.com/tidwall/gjson"
	"reflect"
	"strings"
)

func Read(parseResult *ParseResult, msg string, extra map[string]string) ([]string, error) {
	ns := parseResult.NS
	if len(ns) == 0 {
		return nil, errors.New("no node to collect")
	}
	tmp := map[string]interface{}{}
	var ret = map[string][]*Value{}
	r := gjson.Parse(msg)
I:
	for _, n := range ns {
		var keyBuilder [][]string
		var tags []string
		var tagIndex = map[int]bool{}
		for index, key := range n.Path {
			if key[0] == '$' && key[len(key)-1] == '$' {
				tagIndex[index] = true
				tags = append(tags, key)
				templateValues, ok := tmp[key]
				if !ok {
					break I
					//return nil, errors.New("not found templateValues")
				}
				if templateValues == nil {
					break I
				}
				switch templateValues.(type) {
				case string:
					keyBuilder = append(keyBuilder, []string{templateValues.(string)})
				case []interface{}:
					var tmpTv []string
					for _, tv := range templateValues.([]interface{}) {
						tmpTv = append(tmpTv, tv.(string))
					}
					keyBuilder = append(keyBuilder, tmpTv)
				}
			} else {
				keyBuilder = append(keyBuilder, []string{key})
			}
		}

		var keys [][]string
		for i := 0; i < len(keyBuilder); i++ {
			itemLen := len(keyBuilder[i])
			keyBuilderItem := make([]string, itemLen)
			copy(keyBuilderItem, keyBuilder[i])
			if itemLen > 1 {
				if len(keys) == 0 {
					for _, s := range keyBuilderItem {
						keys = append(keys, []string{s})
					}
				} else {
					var keysCopy = make([][]string, len(keys))
					copy(keysCopy, keys)
					keys = nil
					copyLen := len(keysCopy)
					keys = make([][]string, copyLen*itemLen)
					for j := 0; j < copyLen; j++ {
						for k := 0; k < itemLen; k++ {
							keys[j*itemLen+k] = make([]string, len(keysCopy[j])+1)
							for ind := 0; ind < len(keysCopy[j]); ind++ {
								keys[j*itemLen+k][ind] = keysCopy[j][ind]
							}
							keys[j*itemLen+k][len(keys[j*itemLen+k])-1] = keyBuilderItem[k]
						}
					}
				}
			} else if itemLen == 1 {
				if len(keys) == 0 {
					keys = append(keys, []string{keyBuilderItem[0]})
				} else {
					for l := 0; l < len(keys); l++ {
						keys[l] = append(keys[l], keyBuilderItem[0])
					}
				}
			}
		}

		var tt []interface{}
		arrayCount := 0

		//把数组打平
		for _, keyType := range n.Types {
			if keyType == scout.TypeArray {
				arrayCount += 1
			}
		}
		if len(keys) == 0 {
			keys = append(keys, []string{"*"})
		}
		for _, key := range keys {
			var item = &Value{
				Name: n.Name,
				Tags: tags,
			}
			var gKey []string
			for index, keyType := range n.Types {
				switch keyType {
				case scout.TypeObj:
					if tagIndex[index] {
						item.TagsValue = append(item.TagsValue, key[index])
					}
					gKey = append(gKey, key[index])
				case scout.TypeArray:
					gKey = append(gKey, n.Path[index])
				case scout.TypeSingle:
					if len(n.Path) != 1 || len(n.Types) != 1 {
						return nil, errors.New("template error")
					}
					gKey = append(gKey, n.Path[index])
				}
			}
			if gKey != nil && gKey[len(gKey)-1] == "#" {
				if n.Type == scout.TypeKey {

					gKey[len(gKey)-1] = "#(*)#"
				} else {
					gKey = gKey[:len(gKey)-1]
					if len(gKey) == 0 {
						gKey = append(gKey, "*")
					}
				}
			} else if gKey == nil && len(key) > 0 && key[0] == "*" {
				gKey = append(gKey, "*")
			}
			gKeyStr := strings.Join(gKey, ".")
			//fmt.Println(n.Name, gKeyStr)
			var s gjson.Result
			if gKeyStr == "*" {
				s = r
			} else {
				s = r.Get(gKeyStr)
			}
			if n.Type == scout.TypeKey {
				for i := 0; i < arrayCount-1; i++ {
					k := s.Array()
					if len(k) == 0 {
						break
					}
					s = k[0]
				}
				if s.IsArray() && len(s.Array()) != 0 {
					s.ForEach(func(_, value gjson.Result) bool {
						for k := range value.Map() {
							tt = append(tt, k)
							item.Values = append(item.Values, k)
						}
						return true
					})
					ret[n.Name] = append(ret[n.Name], item)
				}
				if s.IsObject() {
					for k := range s.Map() {
						tt = append(tt, k)
						item.Values = append(item.Values, k)
						ret[n.Name] = append(ret[n.Name], item)
					}
				}
			} else {
				reduce := arrayCount - 1
				for reduce > 0 {
					va := s.Array()
					if len(va) == 0 {
						break
					}
					if len(va) == 2 {
						break
					}
					s = va[0]
					reduce -= 1
				}
				if s.Exists() {
					if s.IsArray() {
						if len(s.Array()) > 0 {
							v := s.Value()
							tt = append(tt, v)
							item.Values = v.([]interface{})
							ret[n.Name] = append(ret[n.Name], item)
						}
					} else {
						v := s.Value()
						tt = append(tt, v)
						item.Values = append(item.Values, v)
						ret[n.Name] = append(ret[n.Name], item)
					}
				}
			}
		}
		if tt != nil {
			if n.Type == scout.TypeKey {
				//做一次去重
				tmpMap := map[string]bool{}
				for _, ttItem := range tt {
					tmpMap[ttItem.(string)] = true
				}
				tt = make([]interface{}, len(tmpMap))
				index := 0
				for s, _ := range tmpMap {
					tt[index] = s
					index += 1
				}

			}
			tmp[n.Name] = tt
		}
	}
	sqlList, err := genSQL(parseResult, ret, extra)
	if err != nil {
		return nil, err
	}
	return sqlList, nil
}

func genSQL(c *ParseResult, values map[string][]*Value, extraValue map[string]string) ([]string, error) {
	var sqlList []string
	if isTemplate(c.Value) {
		searchTag := make([]string, len(c.TagsTemplateStr))
		copy(searchTag, c.TagsTemplateStr)
		if isTemplate(c.TableName) {
			if strings.Contains(c.TableName, c.Value) {
				return nil, errors.New("can not use value in table name")
			}
			searchTag = append(searchTag, c.TableName)
		}
		insertFTemplate := fmt.Sprintf(`IMPORT INTO %s.t_%s USING %s.%s_double TAGS ("%s") VALUES (now,%%v)`, c.DBName, c.TableName, c.DBName, c.STableName, strings.Join(c.TagsTemplateStr, "\",\""))
		insertSTemplate := fmt.Sprintf(`IMPORT INTO %s.t_%s USING %s.%s_string TAGS ("%s") VALUES (now,'%%v')`, c.DBName, c.TableName, c.DBName, c.STableName, strings.Join(c.TagsTemplateStr, "\",\""))
		results := values[c.Value]
	I:
		for _, result := range results {
			for index, rv := range result.Values {
				valueType := reflect.TypeOf(rv).String()
				resultValueMap := map[string]interface{}{result.Name: rv}
				temp := map[string]interface{}{}
				tempTagValue := map[string][]*Value{}
				for i, tag := range result.Tags {
					temp[tag] = result.TagsValue[i]
				}

				for _, tag := range searchTag {
					if isTemplate(tag) {
						if extraValue != nil {
							_, exist := extraValue[tag]
							if exist {
								continue
							}
						}
						tagResult, ok := values[tag]
						if !ok {
							return nil, errors.New("tags not found error")
						}
						if len(tagResult) < 1 {
							return nil, errors.New("tags element is nil (inner error)")
						}
						for _, v := range tagResult {
							if len(v.Values) == 0 {
								return nil, errors.New("tags element Value is nil (inner error)")
							}
							related := true
							for i, t := range v.Tags {
								tagValue, ok := temp[t]
								if ok {
									if tagValue != v.TagsValue[i] {
										related = false
										break
									}
								}
							}
							if related {
								tempTagValue[v.Name] = append(tempTagValue[v.Name], v)
							}
						}
					}
				}
				var tagValues []map[string]string
				if extraValue != nil {
					tagValues = append(tagValues, extraValue)
				} else {
					tagValues = append(tagValues, map[string]string{})
				}

				for i := 0; i < len(result.Tags); i++ {
					tagValues[0][result.Tags[i]] = result.TagsValue[i].(string)
				}
				for _, s := range searchTag {
					if _, ok := tagValues[0][s]; ok {
						continue
					}
					if isTemplate(s) {
						values, ok := tempTagValue[s]
						if !ok {
							break I
						}
						for _, v := range values {
							if len(v.Values) == len(result.Values) {
								//	相等认为是一对一关系
								if len(tagValues) == 0 {
									tagValues = append(tagValues, map[string]string{v.Name: v.Values[index].(string)})
								} else {
									for _, tagValue := range tagValues {
										tagValue[v.Name] = v.Values[index].(string)
									}
								}
							} else if len(v.Values) == 1 {
								//数量为1认为是全部都加
								if len(tagValues) == 0 {
									tagValues = append(tagValues, map[string]string{v.Name: v.Values[0].(string)})
								} else {
									for i := range tagValues {
										tagValues[i][v.Name] = v.Values[0].(string)
									}
								}
							} else {
								//数量不相等直接往上加
								if len(tagValues) == 0 {
									if v.Values != nil {
										for _, vv := range v.Values {
											tagValues = append(tagValues, map[string]string{v.Name: vv.(string)})
										}
									}
								} else {
									cpTagValues := make([]map[string]string, len(tagValues))
									copy(cpTagValues, tagValues)
									tagValues = nil
									tagValues = make([]map[string]string, len(cpTagValues)*len(v.Values))
									for i := 0; i < len(cpTagValues); i++ {
										for j := 0; j < len(v.Values); j++ {
											tagValues[i*len(v.Values)+j] = cpTagValues[i]
											tagValues[i*len(v.Values)+j][v.Name] = v.Values[j].(string)
										}
									}

								}
							}
						}
					} else {
						if len(tagValues) == 0 {
							tagValues = append(tagValues, map[string]string{s: s})
						} else {
							for i := range tagValues {
								tagValues[i][s] = s
							}
						}
					}
				}
				var templateSql string
				for _, vs := range tagValues {
					if valueType == "string" {
						templateSql = insertSTemplate
					} else if strings.HasPrefix(valueType, "int") || strings.HasPrefix(valueType, "float") {
						templateSql = insertFTemplate
					} else if valueType == "bool" {
						templateSql = insertFTemplate
					} else {
						return nil, errors.New("can not parse value type:" + valueType)
					}
					variables := getTemplateStr(templateSql)
					cpTemplate := templateSql
					for _, variable := range variables {
						if variable == c.Value {
							return nil, errors.New("can not using value as tag:" + variable)
						}
						templateValue, ok := vs[variable]
						if !ok {
							return nil, errors.New("get template value not found:" + variable)
						}

						cpTemplate = strings.ReplaceAll(cpTemplate, variable, templateValue)
					}
					pointValue := resultValueMap[c.Value]
					pv := pointValue
					if bv, isBool := pointValue.(bool); isBool {
						if bv {
							pv = 1
						} else {
							pv = 0
						}
					}
					sql := fmt.Sprintf(cpTemplate, pv)
					sqlList = append(sqlList, sql)
				}
			}
		}
	} else {
		return nil, errors.New("value must be a template variable")
	}
	return sqlList, nil
}
