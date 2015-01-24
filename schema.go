package schema

import (
	"encoding/json"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type SchemaError struct {
	error
	Field   string
	ErrType string
}

func (s *SchemaError) Error() string {
	switch s.ErrType {
	case "req":
		return "Schema: The Field " + s.Field + " is required"
	}
	return "Schema: Unknown error on Field " + s.Field
}

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
			err := handleTags(val, i)
			if err != nil {
				return err
			}
			return conform(valField.Addr())
		case reflect.Slice:
			err := handleTags(val, i)
			if err != nil {
				return err
			}
			for j := 0; j < valField.Len(); j += 1 {
				if valField.Index(j).Kind() == reflect.Ptr {
					return conform(valField.Index(j))
				} else {
					return conform(valField.Index(j).Addr())
				}
			}
		//TODO: Add map
		default:
			err := handleTags(val, i)
			if err != nil {
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
				return &SchemaError{Field: val.Type().Field(i).Name, ErrType: "req"}
			}
		case strings.HasPrefix(t, "truncate("):
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
	re := regexp.MustCompile(`^truncate\((\d*)\)`)
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
