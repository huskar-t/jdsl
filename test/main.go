package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/huskar-t/jsontemplate/logic"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var templatePath = flag.String("p", ".", "template path")
	flag.Parse()
	basePath, err := filepath.Abs(*templatePath)
	//basePath := "./example/3"
	if err != nil {
		panic(err)
	}
	templateFilePath := filepath.Join(basePath, "template.json")
	msgPath := filepath.Join(basePath, "msg.json")
	outPath := filepath.Join(basePath, "result.sql")
	t, err := os.Open(templateFilePath)
	if err != nil {
		panic(err)
	}
	m, err := os.Open(msgPath)
	if err != nil {
		panic(err)
	}
	r := io.Reader(t)
	msg, err := ioutil.ReadAll(m)
	if err != nil {
		panic(err)
	}
	c := logic.Config{}
	if err = json.NewDecoder(r).Decode(&c); err != nil {
		panic(err)
	}
	parseResult, createSql, err := logic.ParseJson(&c)
	if err != nil {
		panic(err)
	}
	sqls, err := logic.Read(parseResult, string(msg), nil)
	if err != nil {
		panic(err)
	}
	result := []string{createSql.DBSql, createSql.TableFloatSql, createSql.TableStringSql}
	result = append(result, sqls...)
	out := strings.Join(result, "\n")
	err = ioutil.WriteFile(outPath, []byte(out), os.ModePerm)
	if err != nil {
		panic(err)
	}
	fmt.Println(createSql.DBSql)
	fmt.Println(createSql.TableFloatSql)
	fmt.Println(createSql.TableStringSql)
	for _, sql := range sqls {
		fmt.Println(sql)
	}
}
