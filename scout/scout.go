package scout

import (
	"fmt"
	"strconv"
	"strings"
)

// Scout is the main type that handles the traversing of JSON paths.
type Scout struct {
	data       interface{}
	found      []Found
	lookingFor string
}

const (
	TypeKey = iota + 1
	TypeValue
)
const (
	TypeObj = iota + 1
	TypeArray
	TypeSingle
)

type PathItem struct {
	path string
	Type int
}
type Found struct {
	Path  []string `json:"path"`
	Types []int    `json:"types"`
	Type  int      `json:"type"`
	Name  string   `json:"name"`
}

// DoSearch searches the parsed JSON and returns an array of strings that correspond to
// found locations.
func (s *Scout) DoSearch() ([]Found, error) {
	var foundTypes []int
	var foundPath []string
	switch val := s.data.(type) {
	case []interface{}:
		s.parseArray(val, foundPath, foundTypes)
	case map[string]interface{}:
		s.parseMap(val, foundPath, foundTypes)
	//case string:
	//case int,int32,int64:
	//case float32,float64:
	//case bool:
	default:
		return []Found{
			{
				Path:  []string{"*"},
				Types: []int{TypeSingle},
				Type:  TypeValue,
				Name:  s.lookingFor,
			},
		}, nil
	}
	for i, _ := range s.found {
		s.found[i].Name = s.lookingFor
	}
	return s.found, nil
}

// New instantiates and returns a Scout
func New(lookingFor string, target interface{}) Scout {
	return Scout{
		data:       target,
		lookingFor: lookingFor,
	}
}

func (s *Scout) parseMap(data map[string]interface{}, path []string, types []int) {
	cpType := make([]int, len(types))
	copy(cpType, types)
	types = append(types, TypeObj)
	nextTypes := make([]int, len(types))
	copy(nextTypes, types)
	for key, val := range data {
		if fmt.Sprintf("%v", key) == s.lookingFor {
			s.found = append(s.found, Found{
				Path:  path,
				Type:  TypeKey,
				Types: cpType,
			})
		}
		nextPath := make([]string, len(path))
		copy(nextPath, path)
		nextPath = append(nextPath, strings.ReplaceAll(key, ".", "\\."))

		switch valData := val.(type) {
		case map[string]interface{}:
			s.parseMap(val.(map[string]interface{}), nextPath, nextTypes)
		case []interface{}:
			s.parseArray(val.([]interface{}), nextPath, nextTypes)
		default:
			if fmt.Sprintf("%v", valData) == s.lookingFor {
				s.found = append(s.found, Found{
					Path:  nextPath,
					Type:  TypeValue,
					Types: types,
				})
			}
		}
	}
	return
}

func (s *Scout) parseArray(anArray []interface{}, path []string, types []int) {
	types = append(types, TypeArray)
	nextTypes := make([]int, len(types))
	copy(nextTypes, types)
	getAll := true
	if len(anArray) > 1 {
		getAll = false
	}
	for idx, val := range anArray {
		arrayPath := strconv.Itoa(idx)
		if getAll {
			arrayPath = "#"
		}
		nextPath := make([]string, len(path))
		copy(nextPath, path)
		nextPath = append(nextPath, arrayPath)

		switch valData := val.(type) {
		case map[string]interface{}:
			s.parseMap(val.(map[string]interface{}), nextPath, nextTypes)
		case []interface{}:
			s.parseArray(val.([]interface{}), nextPath, nextTypes)
		default:
			if fmt.Sprintf("%v", valData) == s.lookingFor {
				s.found = append(s.found, Found{
					Path:  nextPath,
					Type:  TypeValue,
					Types: types,
				})
			}
		}
	}
}
