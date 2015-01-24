package schema

import (
	"fmt"
	"testing"
)

const jsBadReqString = `{"age":25}`
const jsBadReqInt = `{"name":"John"}`
const jsBadReqNested = `{"name":"John","age":25, "hair":{}}`

const jsTruncateString = `{"name":"Jonathon","age":25, "hair":{"color":"brown"}}`

const jsTruncateStringPtr = `{"name":"Jonathon","age":25, "hair":[{"color":"brown"},{"color":"red"}]}`

type Test struct {
	Name string `json:"name" schema:"req,truncate(4)"`
	Age  int    `json:"age" schema:"req"`
	Hair Hair   `json:"hair" schema:req`
}
type TestPtr struct {
	Name string  `json:"name" schema:"req,truncate(4)"`
	Age  int     `json:"age" schema:"req"`
	Hair []*Hair `json:"hair" schema:req`
}

type Hair struct {
	Color string `json:"color" schema:"req,truncate(2)"`
}

func TestRequired(t *testing.T) {
	fixtures := []struct {
		Expected error
		s        interface{}
		json     string
		name     string
		msg      string
	}{
		{nil, &Test{}, `{"name":"John","age":25, "hair":{"color":"brown"}}`, "Unmarshal proper json", "expect nil error, got: "},
		{nil, &TestPtr{}, `{"name":"John","age":25, "hair":[{"color":"brown"}]}`, "Unmarshal proper json with slice pointer", "expect nil error, got: "},
		{&SchemaError{Field: "Name", ErrType: "req"}, &Test{}, `{"age":25, "hair":{"color":"brown"}}`, "Omit required field 'name' (string)", "expect require error, got: "},
		{&SchemaError{Field: "Color", ErrType: "req"}, &Test{}, `{"name":"John","age":25, "hair":{}}`, "Omit required nested field 'hair:color' (string)", "expect require error, got: "},
		{&SchemaError{Field: "Color", ErrType: "req"}, &TestPtr{}, `{"name":"John","age":25, "hair":[{}]}`, "Omit required nested field pointer 'hair:color' (string)", "expect require error, got: "},
	}

	for _, f := range fixtures {
		t.Log(f.name)
		err := Unmarshal([]byte(f.json), f.s)
		if fmt.Sprint(err) != fmt.Sprint(f.Expected) {
			t.Error(f.name, f.msg, err)
		}
	}
}

func TestTruncateString(t *testing.T) {
	fixtures := []struct {
		ExpectedErr error
		Expected    func(*Test) bool
		s           interface{}
		json        string
		name        string
		msg         string
	}{
		{nil, func(t *Test) bool { return t.Name == "Jona" }, &Test{}, `{"name":"Jonathon","age":25, "hair":{"color":"brown"}}`, "Truncate String", "expect nil error & Name=Jona got:"},
		{nil, func(t *Test) bool { return t.Hair.Color == "br" }, &Test{}, `{"name":"Jonathon","age":25, "hair":{"color":"brown"}}`, "Truncate String", "expect nil error & Color=br got:"},
	}

	for _, f := range fixtures {
		t.Log(f.name)
		err := Unmarshal([]byte(f.json), f.s)
		if fmt.Sprint(err) != fmt.Sprint(f.ExpectedErr) {
			t.Error(f.msg, err)
		}
		if f.Expected(f.s.(*Test)) != true {
			t.Error(f.msg, f.s)
		}
	}
}
