package dotty

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func Exists(src interface{}, key string) (exists bool, err error) {
	bits := strings.Split(key, ".")

	exists = true

	value := src
	for _, bit := range bits {
		_type := reflect.TypeOf(value)

		if _type.Kind() == reflect.Map && _type.Key().Kind() == reflect.String {
			if currentValue := reflect.ValueOf(value).MapIndex(reflect.ValueOf(bit)); currentValue == reflect.Zero(_type) {
				exists = false
				return
			} else {
				value = currentValue.Interface()
			}
		} else if _type.Kind() == reflect.Array {
			if i, _err := strconv.ParseInt(bit, 10, 32); _err != nil {
				err = _err
				return
			} else if reflect.ValueOf(value).Len() <= int(i) {
				exists = false
				return
			} else {
				value = reflect.ValueOf(value).Index(int(i)).Interface()
			}
		} else {
			err = errors.New(fmt.Sprintf("Exists() encountered invalid type (%s) in object (at %s's parent)", _type, bit))
		}
	}

	return
}

func Get(src interface{}, key string) (value interface{}, err error) {
	bits := strings.Split(key, ".")

	value = src
	for _, bit := range bits {
		_type := reflect.TypeOf(value)

		if _type.Kind() == reflect.Map && _type.Key().Kind() == reflect.String {
			value = reflect.ValueOf(value).MapIndex(reflect.ValueOf(bit)).Interface()
		} else if _type.Kind() == reflect.Array {
			if i, _err := strconv.ParseInt(bit, 10, 32); _err != nil {
				err = _err
				return
			} else {
				value = reflect.ValueOf(value).Index(int(i)).Interface()
			}
		} else {
			err = errors.New(fmt.Sprintf("Get() encountered invalid type (%s) in object (at %s's parent)", _type, bit))
			return
		}
	}

	return
}

func Put(dst interface{}, key string, value interface{}) (err error) {
	bits := strings.Split(key, ".")

	target := dst
	for _, bit := range bits[0 : len(bits)-1] {
		_type := reflect.TypeOf(target)
		_value := reflect.ValueOf(target)

		if _type == nil {
			err = errors.New(fmt.Sprintf("Put() encountered nil value in object (at %s's parent)", bit))
			return
		} else if _type.Kind() == reflect.Map && _type.Key().Kind() == reflect.String {
			if _target := reflect.ValueOf(value).MapIndex(reflect.ValueOf(bit)); _target == reflect.Zero(_type) {
				target = map[string]interface{}{}
				_value.SetMapIndex(reflect.ValueOf(bit), reflect.ValueOf(target))
			} else {
				target = _target.Interface()
			}
		} else if _type.Kind() == reflect.Array || _type.Kind() == reflect.Slice {
			if n, _err := strconv.ParseInt(bit, 10, 32); _err != nil {
				err = _err
				return
			} else if _value.Len() <= int(n) {
				err = errors.New(fmt.Sprintf("Put() encountered an array that was shorter (%d) than the key required (%d)", _value.Len(), int(n)))
				return
			} else {
				target = _value.Index(int(n))
			}
		} else {
			err = errors.New(fmt.Sprintf("Put() encountered invalid type (%s) in object (at %s's parent)", _type, bit))
			return
		}
	}

	_type := reflect.TypeOf(target)
	_value := reflect.ValueOf(target)

	if _type == nil || _type.Kind() != reflect.Map || _type.Key().Kind() != reflect.String {
		err = errors.New(fmt.Sprintf("Put() target key (%s)'s parent was not a map (%s)", key, _type))
		return
	}

	bit := bits[len(bits)-1]

	_value.SetMapIndex(reflect.ValueOf(bit), reflect.ValueOf(value))

	return
}
