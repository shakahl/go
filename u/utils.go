package u

import (
	"reflect"
	"strings"
)

func IsMissing(x interface{}) bool {
	if x == nil {
		return true
	}

	v := reflect.ValueOf(x)
	k := v.Kind()

	for {
		if k == reflect.Ptr || k == reflect.Interface {
			v = v.Elem()
			k = v.Kind()
		} else {
			break
		}
	}

	if k == reflect.Invalid { // nil
		return true
	} else if k == reflect.String {
		return strings.TrimSpace(x.(string)) == ""
	} else if k == reflect.Slice || k == reflect.Array || k == reflect.Chan || k == reflect.Map {
		return v.Len() == 0
	}

	return false
}
