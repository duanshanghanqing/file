package util

import "reflect"

func MapMapping(source map[string]interface{}, target interface{})  {
	// 原始Map
	_, s := reflect.TypeOf(source), reflect.ValueOf(source)
	// 必须Map类型
	if s.Kind() == reflect.Map {
		// 目标Map，指针类型
		_, Value := reflect.TypeOf(target), reflect.ValueOf(target)
		v := Value.Elem()

		// 映射赋值
		res := reflect.MakeMap(s.Type())
		keys := s.MapKeys()
		for _, k := range keys {
			key := k.Convert(res.Type().Key())
			value := s.MapIndex(key)
			v.SetMapIndex(key, value)
		}
	}
}

//EnumMapList := make(map[string]interface{})
//EnumMapList["FileState"] = 1
//
//EnumMapList2 := make(map[string]interface{})
//
//util.MapMapping(EnumMapList, &EnumMapList2)
//
//for key := range EnumMapList2 {
//	fmt.Println(key, EnumMapList2[key])
//}
