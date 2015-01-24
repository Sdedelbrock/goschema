package schema

import "testing"

const jsGood = `{"name":"John","age":25, "hair":{"color":"brown"}}`
const jsBadReqString = `{"age":25}`
const jsBadReqInt = `{"name":"John"}`
const jsBadReqNested = `{"name":"John","age":25, "hair":{}}`

const jsTruncateString = `{"name":"Jonathon","age":25, "hair":{"color":"brown"}}`

const jsTruncateStringPtr = `{"name":"Jonathon","age":25, "hair":[{"color":"brown"},{"color":"red"}]}`

type Test struct {
	Name string `json:"name" schema:"req,slen(4)"`
	Age  int    `json:"age" schema:"req"`
	Hair Hair   `json:"hair" schema:req`
}
type TestPtr struct {
	Name string  `json:"name" schema:"req,slen(4)"`
	Age  int     `json:"age" schema:"req"`
	Hair []*Hair `json:"hair" schema:req`
}

type Hair struct {
	Color string `json:"color" schema:"req,slen(2)"`
}

func TestUnmarshalGood(t *testing.T) {
	v := Test{}
	err := Unmarshal([]byte(jsGood), &v)
	if err != nil {
		t.Error("Could not marshal proper JSON: ", err)
	}
}

func TestUnmarshalBadReq(t *testing.T) {
	err := Unmarshal([]byte(jsBadReqString), &Test{})
	t.Log(err)
	if err == nil {
		t.Error("Did not throw error on required field Name (string)")
	}
	err = Unmarshal([]byte(jsBadReqInt), &Test{})
	t.Log(err)
	if err == nil {
		t.Error("Did not throw error on required field Age (int)")
	}
	err = Unmarshal([]byte(jsBadReqNested), &Test{})
	t.Log(err)
	if err == nil {
		t.Error("Did not throw error on required nested field Hair (struct)")
	}
}
func TestUnmarshalTruncateString(t *testing.T) {
	v := Test{}
	err := Unmarshal([]byte(jsTruncateString), &v)
	if err != nil {
		t.Error("Could not marshal proper JSON:", err)
	}
	if v.Name != "Jona" {
		t.Error("slen tag found and string not truncated: expected Jona got: ", v.Name)
	}
}
func TestUnmarshalTruncateStringPtr(t *testing.T) {
	v := TestPtr{}
	err := Unmarshal([]byte(jsTruncateStringPtr), &v)
	if err != nil {
		t.Error("Could not marshal proper JSON:", err)
	}
	if v.Name != "Jona" {
		t.Error("slen tag found and string not truncated: expected Jona got: ", v.Name)
	}
	if v.Hair[0].Color != "br" {
		t.Error("slen tag found and *[]struct not truncated: expected br got: ", v.Hair[0].Color)
	}

}
