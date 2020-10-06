package main

import "C"
import (
	"encoding/json"
	"fmt"
	"github.com/huskar-t/jsontemplate/logic"
	"time"
)

type KeysNode struct {
	name     string
	children []*KeysNode
	parent   *KeysNode
}

const a = `{"$ver$":"$value$"}`
const b = `{"temperature":"15","humidity":"17"}`

func main() {
	start := time.Now()
	c := &logic.Config{
		DBName:     "db",
		STableName: "stb",
		TableName:  "zone",
		Keep:       365,
		TagsMap: map[string]string{
			"type": "$ver$",
		},
		PayloadTemplate: a,
		Value:           "$value$",
	}
	d, _ := json.Marshal(c)
	fmt.Println(string(d))

	result, createSqls, err := logic.ParseJson(c)
	if err != nil {
		fmt.Println(err)
	}
	_ = createSqls
	fmt.Println(createSqls.DBSql)
	fmt.Println(createSqls.TableFloatSql)
	fmt.Println(createSqls.TableStringSql)
	r, err := logic.Read(result, b, nil)
	if err != nil {
		fmt.Println(err)
	}
	for _, i2 := range r {
		fmt.Println(i2)
	}
	fmt.Println(r)
	fmt.Println(d[0])
	fmt.Println(time.Now().Sub(start))
}
