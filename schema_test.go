package schema

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
)

const jsBadReqString = `{"age":25}`
const jsBadReqInt = `{"name":"John"}`
const jsBadReqNested = `{"name":"John","age":25, "hair":{}}`

const jsTruncateString = `{"name":"Jonathon","age":25, "hair":{"color":"brown"}}`

const jsTruncateStringPtr = `{"name":"Jonathon","age":25, "hair":[{"color":"brown"},{"color":"red"}]}`

type Parent struct {
	TestPtrs1   []TestPtr    `json:"test_ptrs1" `
	TestSlices1 []TestSlice  `json:"test_slices1" `
	Tests1      []Test       `json:"tests1" `
	TestPtrs2   []*TestPtr   `json:"test_ptrs2" `
	TestSlices2 []*TestSlice `json:"test_slices2" `
	Tests2      []*Test      `json:"tests2" `
	TestPtrs3   *[]TestPtr   `json:"test_ptrs3" `
	TestSlices3 *[]TestSlice `json:"test_slices3" `
	Tests3      *[]Test      `json:"tests3" `
}

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

type TestSlice struct {
	Name string `json:"name" schema:"req,truncate(4)"`
	Age  int    `json:"age" schema:"req"`
	Hair []Hair `json:"hair" schema:req`
}

type Hair struct {
	Color string `json:"color" schema:"req,truncate(2)"`
}

func TestTruncateReallyComplicatedStruct(t *testing.T) {
	var p = &Parent{}

	parent := `{"test_ptrs1":[{"name":"John","age":25,"hair":[{"color":"brown"}]},{"name":"John","age":25,"hair":[{"color":"brown"}]}],"test_ptrs2":[{"name":"John","age":25,"hair":[{"color":"brown"}]},{"name":"John","age":25,"hair":[{"color":"brown"}]}],"test_ptrs3":[{"name":"John","age":25,"hair":[{"color":"brown"}]},{"name":"John","age":25,"hair":[{"color":"brown"}]}],"test_slices1":[{"name":"John","age":25,"hair":[{"color":"brown"}]},{"name":"John","age":25,"hair":[{"color":"brown"}]}],"test_slices2":[{"name":"John","age":25,"hair":[{"color":"brown"}]},{"name":"John","age":25,"hair":[{"color":"brown"}]}],"test_slices3":[{"name":"John","age":25,"hair":[{"color":"brown"}]},{"name":"John","age":25,"hair":[{"color":"brown"}]}],"tests1":[{"name":"John","age":25,"hair":{"color":"brown"}},{"name":"John","age":25,"hair":{"color":"brown"}}],"tests2":[{"name":"John","age":25,"hair":{"color":"brown"}},{"name":"John","age":25,"hair":{"color":"brown"}}],"tests3":[{"name":"John","age":25,"hair":{"color":"brown"}},{"name":"John","age":25,"hair":{"color":"brown"}}]}`
	expected := `{"test_ptrs1":[{"name":"John","age":25,"hair":[{"color":"br"}]},{"name":"John","age":25,"hair":[{"color":"br"}]}],"test_slices1":[{"name":"John","age":25,"hair":[{"color":"br"}]},{"name":"John","age":25,"hair":[{"color":"br"}]}],"tests1":[{"name":"John","age":25,"hair":{"color":"br"}},{"name":"John","age":25,"hair":{"color":"br"}}],"test_ptrs2":[{"name":"John","age":25,"hair":[{"color":"br"}]},{"name":"John","age":25,"hair":[{"color":"br"}]}],"test_slices2":[{"name":"John","age":25,"hair":[{"color":"br"}]},{"name":"John","age":25,"hair":[{"color":"br"}]}],"tests2":[{"name":"John","age":25,"hair":{"color":"br"}},{"name":"John","age":25,"hair":{"color":"br"}}],"test_ptrs3":[{"name":"John","age":25,"hair":[{"color":"br"}]},{"name":"John","age":25,"hair":[{"color":"br"}]}],"test_slices3":[{"name":"John","age":25,"hair":[{"color":"br"}]},{"name":"John","age":25,"hair":[{"color":"br"}]}],"tests3":[{"name":"John","age":25,"hair":{"color":"br"}},{"name":"John","age":25,"hair":{"color":"br"}}]}`

	err := Unmarshal(json.Unmarshal, []byte(parent), p)
	if err != nil {
		t.Error(err)
	}

	output, err := Marshal(json.Marshal, p)

	if !bytes.Equal([]byte(expected), output) {
		var dst bytes.Buffer
		json.Indent(&dst, output, "=", "\t")
		t.Log(dst.String())
		t.Error("Failed")
	}
}

func TestTruncateComplicatedStruct(t *testing.T) {
	var s = &TestSlice{}
	testPtr := `{"name": "John", "age" : 25, "hair": [{"color": "brown"}, {"color": "blonde"}]}`

	err := Unmarshal(json.Unmarshal, []byte(testPtr), s)
	if err != nil {
		t.Error(err)
	}

	if s.Hair[0].Color != "br" || s.Hair[1].Color != "bl" {
		t.Error("Fail")
	}

}

func TestTruncateComplicatedStruct2(t *testing.T) {
	var s = &TestPtr{}
	testPtr := `{"name": "John", "age" : 25, "hair": [{"color": "brown"}]}`

	err := Unmarshal(json.Unmarshal, []byte(testPtr), s)
	if err != nil {
		t.Error(err)
	}

	if s.Hair[0].Color != "br" {
		t.Error("Fail")
	}

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
		err := Unmarshal(json.Unmarshal, []byte(f.json), f.s)
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
		err := Unmarshal(json.Unmarshal, []byte(f.json), f.s)
		if fmt.Sprint(err) != fmt.Sprint(f.ExpectedErr) {
			t.Error(f.msg, err)
		}
		if f.Expected(f.s.(*Test)) != true {
			t.Error(f.msg, f.s)
		}
	}
}
