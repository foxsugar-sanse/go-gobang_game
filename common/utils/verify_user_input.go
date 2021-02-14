package utils

import (
	"reflect"
)

type Info struct {
	FieldNumber int
	Empty 		int
	Exists 		int
}



func UserInput(data interface{}) *Info {
	switch value := reflect.ValueOf(data); value.Kind() {
	case reflect.Struct:
		num := 0
		empty := func() int {
			for i := 0; i < value.NumField(); i++ {
				switch val := reflect.ValueOf(value.Field(i).Interface()) ; val.Kind() {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					if value.Field(i).Int() == 0 {
						num++
					}
				case reflect.String:
					if value.Field(i).String() == "" {
						num++
					}
				case reflect.Uint, reflect.Uint8, reflect.Uint16,reflect.Uint32, reflect.Uint64:
					if value.Field(i).Uint() == 0 {
						num++
					}
				case reflect.Float32, reflect.Float64:
					if value.Field(i).Float() == 0 {
						num++
					}
				case reflect.Map:
					// 匹配map

				case reflect.Slice:
					// 匹配Slice
					if val := value.Field(i).Interface().([]interface{}); len(val) == 0 {
						num++
					}
				}
			}
			return num
		} ()
		exits := value.NumField() - empty
		return &Info{
			FieldNumber: value.NumField(),
			Empty: empty,
			Exists: exits,
		}
	default:
		return nil
	}
}
