package dotty

import (
	"errors"
	"strconv"
	"strings"
)

func Exists(src interface{}, key string) (exists bool, err error) {
	bits := strings.Split(key, ".")

	exists = true

	value := src
	for _, bit := range bits {
		switch value.(type) {
		case map[string]interface{}:
			_, ok := value.(map[string]interface{})[bit]
			if ok != true {
				exists = false
				return
			}

			value = value.(map[string]interface{})[bit]
		case []interface{}:
			i, ierr := strconv.ParseInt(bit, 10, 32)
			if ierr != nil {
				err = ierr
				return
			}

			if len(value.([]interface{})) <= int(i) {
				exists = false
				return
			}

			value = value.([]interface{})[i]
		default:
			err = errors.New("encountered invalid type in object")
			return
		}
	}

	return
}

func Get(src interface{}, key string) (value interface{}, err error) {
	bits := strings.Split(key, ".")

	value = src
	for _, bit := range bits {
		switch value.(type) {
		case map[string]interface{}:
			value = value.(map[string]interface{})[bit]
		case []interface{}:
			i, ierr := strconv.ParseInt(bit, 10, 32)
			if ierr != nil {
				err = ierr
				return
			}
			value = value.([]interface{})[i]
		default:
			err = errors.New("encountered invalid type in object")
			return
		}
	}

	return
}

func Put(dst interface{}, key string, value interface{}) (err error) {
	bits := strings.Split(key, ".")

	target := dst
	for i, bit := range bits {
		if i == len(bits)-1 {
			switch target.(type) {
			case map[string]interface{}:
				target.(map[string]interface{})[bit] = value
				err = nil
				return
			default:
				err = errors.New("final object target was not a map")
			}
		}

		switch target.(type) {
		case map[string]interface{}:
			_, ok := target.(map[string]interface{})[bit]
			if ok != true {
				target.(map[string]interface{})[bit] = map[string]interface{}{}
			}

			target = target.(map[string]interface{})[bit]
		case []interface{}:
			n, nerr := strconv.ParseInt(bit, 10, 32)
			if nerr != nil {
				err = nerr
				return
			}

			if len(target.([]interface{})) <= int(n) {
				err = errors.New("encountered an array that was shorter than the key required")
				return
			}

			target = target.([]interface{})[n]
		default:
			err = errors.New("encountered invalid type in object")
			return
		}
	}

	return
}
