package schema

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func Unmarshal(data []byte, v interface{}) error {
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	return conform(reflect.ValueOf(v))
}

func Marshal(v interface{}) ([]byte, error) {
	j, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return j, nil
}

func conform(v reflect.Value) error {
	val := v.Elem()
	for i := 0; i < val.NumField(); i++ {
		valField := val.Field(i)
		switch valField.Kind() {
		case reflect.Struct:
			err := conform(valField.Addr())
			if err != nil {
				return err
			}
		case reflect.Slice:
			for j := 0; j < valField.Len(); j += 1 {
				if valField.Index(j).Kind() == reflect.Ptr {
					conform(valField.Index(j))
				} else {
					conform(valField.Index(j).Addr())
				}
			}
		default:
			if err := handleTags(val, i); err != nil {
				return err
			}
		}
	}
	return nil
}

func handleTags(val reflect.Value, i int) error {
	valField := val.Field(i)
	for _, tag := range strings.Split(val.Type().Field(i).Tag.Get("schema"), ",") {
		t := strings.TrimSpace(tag)
		switch {
		case t == "req":
			if isZero(valField) {
				return fmt.Errorf("Schema: Field %s is required", val.Type().Field(i).Name)
			}
		case strings.HasPrefix(t, "slen("):
			truncate(t, valField)
		}
	}
	return nil
}

func isZero(v reflect.Value) bool {
	zero := reflect.Zero(v.Type()).Interface()
	current := v.Interface()

	return (current == zero)
}

func truncate(t string, v reflect.Value) error {
	re := regexp.MustCompile(`^slen\((\d*)\)`)
	i, err := strconv.Atoi(re.FindStringSubmatch(t)[1])
	if err != nil {
		return err
	}
	switch v.Kind() {
	case reflect.String:
		val := v.String()
		if len(val) > i {
			v.SetString(val[:i])
		}
	}
	return nil
}
