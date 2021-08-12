package util

import "reflect"

// EntityMapping p_接收的是实体，f_ 接收的是指针
func EntityMapping(p_ interface{}, f_ interface{}) {
	pType, pValue := reflect.TypeOf(p_), reflect.ValueOf(p_)
	_, fValue := reflect.TypeOf(f_), reflect.ValueOf(f_)
	v := fValue.Elem()

	for i := 0; i < pType.NumField(); i++ {
		// 左侧的key, value
		field := pType.Field(i) // key字段, 使用 mode 对象的key
		value := pValue.FieldByName(field.Name)

		// 左侧的key在右侧中也有
		// 使用左侧的key在右侧中检查，判断右侧中有没有
		if v.FieldByName(field.Name).String() != "<invalid Value>" {
			v.FieldByName(field.Name).Set(value)
		}
	}
}
