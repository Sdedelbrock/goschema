# goschema
A library for validating and mutating JSON data into a struct by utilizing go struct tags.  Design to be a drop in replacement for json.Marshal/json.Unmarshal.

`go get github.com/Sdedelbrock/goschema`

##usage
```go
package main

import (
	"github.com/Sdedelbrock/goschema"
	"fmt"
)

type Person struct {
	// req means this field is required, if it is not present it will throw an error
	FirstName string `json:"first-name" schema:"req"`
	// truncate(n) will truncate the string to n characters
	LastName string `json:"last-name" schema:"req,truncate(4)"`
}

func main(){
	var i = Person{}
	err := schema.Unmarshal([]byte(`{}`), &i)
	fmt.Println(err)  // Schema: The Field FirstName is required
	err = schema.Unmarshal([]byte(`{"first-name":"Charlie", "last-name":"Chaplin"}`), &i)
	fmt.Println(err,i)  //<nil> {Charlie Chap}
}
```
