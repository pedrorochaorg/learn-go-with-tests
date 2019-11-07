package reflection

import "reflect"

func walk(x interface{}, fn func(input string)) {
	val := reflect.ValueOf(x)


	// If the field is a pointer we need to extract the underlying value before we can use it as a regular field.
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)

		switch field.Kind() {
		case reflect.String:
			fn(field.String())
		case reflect.Struct:
			walk(field.Interface(), fn)
		}
	}
}
