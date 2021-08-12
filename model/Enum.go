package model

import (
	"reflect"
	"strconv"
)

// Enum 枚举结构
type Enum struct {
	Label    string `json:"label"`
	Value    string `json:"value"`
	Children []Enum `json:"children"`
}

// 接收任意类型
// 枚举map[int]string 转 [{ Label: "男"，value: 0 }, { Label: "女"，value: 1 }]
func toEnumList(s interface{}) []Enum {
	rValue := reflect.ValueOf(s) // value
	kd := rValue.Kind()          // 类别

	km := make([]Enum, 0) // 返回的数组

	if kd == reflect.Map { // 是 Map 类别
		iter := rValue.MapRange()
		for iter.Next() {
			k := iter.Key()
			v := iter.Value()
			// fmt.Println(int64(k.Int()) , v.String())
			km = append(km, Enum{Label: v.String(), Value: strconv.FormatInt(k.Int(), 10)})
		}
	}

	return km
}

// EnumMapListType 主业务
// 返回系统定义的全部枚举 { "sexType": [{ Label: "男"，value: 0 }, { Label: "女"，value: 1 }] }
type EnumMapListType map[string]interface{}

var EnumMapList EnumMapListType

func init() {
	EnumMapList = make(EnumMapListType)
	EnumMapList["FileState"] = toEnumList(FileState)
}

// EnumType 定义枚举类型
type EnumType int64

// 枚举值转字符串
func (i EnumType) String() string {
	return strconv.FormatInt(int64(i), 10)
}
