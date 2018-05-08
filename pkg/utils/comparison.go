package utils

import "reflect"

func InArray(array interface{}, value interface{}) (bool, int) {

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for index := 0; index < s.Len(); index++ {
			if reflect.DeepEqual(value, s.Index(index).Interface()) == true {
				return true, index
			}
		}
	}

	return false, -1
}
