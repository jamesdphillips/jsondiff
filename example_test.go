package jsondiff_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/jamesdphillips/jsondiff"
)

func ExampleCompare() {
	type Phone struct {
		Type   string `json:"type"`
		Number string `json:"number"`
	}
	type Person struct {
		Firstname string  `json:"firstName"`
		Lastname  string  `json:"lastName"`
		Gender    string  `json:"gender"`
		Age       int     `json:"age"`
		Phones    []Phone `json:"phoneNumbers"`
	}
	source, err := ioutil.ReadFile("testdata/examples/john.json")
	if err != nil {
		log.Fatal(err)
	}
	var john Person
	if err := json.Unmarshal(source, &john); err != nil {
		log.Fatal(err)
	}
	jack := john
	jack.Firstname = "Jack"
	jack.Phones = append(john.Phones, Phone{
		Type:   "mobile",
		Number: "209-212-0015",
	})
	patch, err := jsondiff.Compare(john, jack)
	if err != nil {
		log.Fatal(err)
	}
	for _, op := range patch {
		fmt.Printf("%s\n", op)
	}
	// Output:
	// {"op":"replace","path":"/firstName","value":"Jack"}
	// {"op":"add","path":"/phoneNumbers/-","value":{"number":"209-212-0015","type":"mobile"}}
}

func ExampleCompareJSON() {
	type Phone struct {
		Type   string `json:"type"`
		Number string `json:"number"`
	}
	type Person struct {
		Firstname string  `json:"firstName"`
		Lastname  string  `json:"lastName"`
		Gender    string  `json:"gender"`
		Age       int     `json:"age"`
		Phones    []Phone `json:"phoneNumbers"`
	}
	source, err := ioutil.ReadFile("testdata/examples/john.json")
	if err != nil {
		log.Fatal(err)
	}
	var john Person
	if err := json.Unmarshal(source, &john); err != nil {
		log.Fatal(err)
	}
	john.Age = 30
	john.Phones = append(john.Phones, Phone{
		Type:   "mobile",
		Number: "209-212-0015",
	})
	target, err := json.Marshal(john)
	if err != nil {
		log.Fatal(err)
	}
	patch, err := jsondiff.CompareJSON(source, target)
	if err != nil {
		log.Fatal(err)
	}
	for _, op := range patch {
		fmt.Printf("%s\n", op)
	}
	// Output:
	// {"op":"replace","path":"/age","value":30}
	// {"op":"add","path":"/phoneNumbers/-","value":{"number":"209-212-0015","type":"mobile"}}
}

func ExampleInvertible() {
	source := `{"a":"1","b":"2"}`
	target := `{"a":"3","c":"4"}`

	patch, err := jsondiff.CompareJSONOpts(
		[]byte(source),
		[]byte(target),
		jsondiff.Invertible(),
	)
	if err != nil {
		log.Fatal(err)
	}
	for _, op := range patch {
		fmt.Printf("%s\n", op)
	}
	// Output:
	// {"op":"test","path":"/a","value":"1"}
	// {"op":"replace","path":"/a","value":"3"}
	// {"op":"test","path":"/b","value":"2"}
	// {"op":"remove","path":"/b"}
	// {"op":"add","path":"/c","value":"4"}
}

func ExampleFactorize() {
	source := `{"a":[1,2,3],"b":{"foo":"bar"}}`
	target := `{"a":[1,2,3],"c":[1,2,3],"d":{"foo":"bar"}}`

	patch, err := jsondiff.CompareJSONOpts(
		[]byte(source),
		[]byte(target),
		jsondiff.Factorize(),
	)
	if err != nil {
		log.Fatal(err)
	}
	for _, op := range patch {
		fmt.Printf("%s\n", op)
	}
	// Output:
	// {"op":"copy","from":"/a","path":"/c"}
	// {"op":"move","from":"/b","path":"/d"}
}
