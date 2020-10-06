# json解析模板 JDSL
## 目标
通过设定json模板来提取实时数据中的关心的值尝试解释关系并组合成 TDengine 的插入语句

## 模板定义
json 中的不确定值(模板变量)使用美元符包裹 如 $gateway$
## 配置解析
```go
type Config struct {
	DBName          string            `json:"dbName"`
	STableName      string            `json:"sTableName"`
	TableName       string            `json:"tableName"`
	Keep            int               `json:"keep"`
	TagsMap         map[string]string `json:"tagsMap"`
	PayloadTemplate string            `json:"payloadTemplate"`
	Value           string            `json:"value"`
}
```
> * dbName:  TDengine 数据库名
> * sTableName:  TDengine 超级表名
> * tableName:  TDengine 表名(可为模板变量)
> * keep: 数据保存天数 
> * tagsMap: 最终数据的分组依据 键为超级表的tag定义,值可以为模板变量也可以为常量字符串
> * payloadTemplate: json 模板
> * value: 要存储的数据(必须为模板变量)
## 规则
> * 为避免表名为纯数字,默认添加前缀 "t_"
> * 暂不支持自定义写入时间
> * tag 的值只能为字符串,不可使用数组
> * value 必须为模板变量,支持 数字型 字符串 布尔值
> * tag 不能包含 value 的模板变量
> * 字符串限制最大长度64
## 样例
> 模板1
> ```json
> {
>   "dbName": "db",
>   "sTableName": "stb",
>   "tableName": "$pointName$",
>   "keep": 365,
>   "tagsMap": {
>     "device": "$device$"
>   },
>   "payloadTemplate": "{\"device\":\"$device$\",\"point\":{\"pointName\":\"$pointName$\",\"value\":[\"$value$\"]}}",
>   "value": "$value$"
> }
> ```
> 测试实时消息
> ```json
> {
>   "device": "d1",
>   "point": {
>     "pointName": "sunshine",
>     "value": [
>       92,
>       93,
>       94
>     ]
>   }
> }
> ```
> 生成语句
> ```sql
> CREATE DATABASE IF NOT EXISTS db KEEP 365
> CREATE TABLE IF NOT EXISTS db.stb_double (ts timestamp, value float) TAGS (device binary(64))
> CREATE TABLE IF NOT EXISTS db.stb_string (ts timestamp, value NCHAR(64)) TAGS (device binary(64))
> IMPORT INTO db.t_sunshine USING db.stb_double TAGS ("d1") VALUES (now,92)
> IMPORT INTO db.t_sunshine USING db.stb_double TAGS ("d1") VALUES (now,93)
> IMPORT INTO db.t_sunshine USING db.stb_double TAGS ("d1") VALUES (now,94)
> ```


> 模板2
> ```json
> {
>   "$gatewayID$": {
>     "$deviceID$": {
>       "$pointName$": "$value$"
>     }
>   }
> }
> ```
> 测试实时消息
> ```json
> {
>   "g1": {
>     "d1": {
>       "p1": "1"
>     },
>     "d2": {
>       "p2": 1
>     }
>   }
> }
> ```
> 生成语句
> ```sql
> CREATE DATABASE IF NOT EXISTS db KEEP 365
> CREATE TABLE IF NOT EXISTS db.stb_double (ts timestamp, value float) TAGS (device binary(64),gateway binary(64))
> CREATE TABLE IF NOT EXISTS db.stb_string (ts timestamp, value NCHAR(64)) TAGS (device binary(64),gateway binary(64))
> IMPORT INTO db.t_p1 USING db.stb_string TAGS ("d1","g1") VALUES (now,'1')
> IMPORT INTO db.t_p2 USING db.stb_double TAGS ("d2","g1") VALUES (now,1)
> IMPORT INTO db.t_p1 USING db.stb_string TAGS ("d1","g2") VALUES (now,'2')
> IMPORT INTO db.t_p2 USING db.stb_double TAGS ("d2","g2") VALUES (now,2)
> ```


> 模板3
> ```json
> {
>   "dbName": "db",
>   "sTableName": "stb",
>   "tableName": "zone",
>   "keep": 365,
>   "tagsMap": {
>     "type": "$ver$"
>   },
>   "payloadTemplate": "{\"$ver$\":\"$value$\"}",
>   "value": "$value$"
> }
> ```
> 测试实时消息
> ```json
> {
>   "temperature": "15",
>   "humidity": "17"
> }
> ```
> 生成语句
> ```sql
> CREATE DATABASE IF NOT EXISTS db KEEP 365
> CREATE TABLE IF NOT EXISTS db.stb_double (ts timestamp, value float) TAGS (type binary(64))
> CREATE TABLE IF NOT EXISTS db.stb_string (ts timestamp, value NCHAR(64)) TAGS (type binary(64))
> IMPORT INTO db.t_zone USING db.stb_string TAGS ("temperature") VALUES (now,'15')
> IMPORT INTO db.t_zone USING db.stb_string TAGS ("humidity") VALUES (now,'17')
> ```

更多有趣的格式自由发挥

## 使用
> * go: 
>   * logic.ParseJson 用于预处理 json 模板 
>   * logic.Read 解析json消息
> * c:
>  * 运行 build/build.py 生成动态库(需要 docker 和 python 2或3都可,需要环境变量 GOCACHE 和 GOCACHE)
> * java:
>  * jni 或 jna 调用动态库

## 测试

在 test/example 文件夹内创建测试文件夹,文件夹内创建模板文件 template.json 测试消息数据 msg.json, 在 test 文件夹内运行 go run main.go -p ./example/测试文件夹

会按顺序打印出 
> * 数据库创建语句
> * 超级表创建语句
> * 数据插入语句

同时写入文件 result.sql

## 编译要求
go 1.13+
python 2 or 3
docker 
环境变量 GOCACHE GOPATH