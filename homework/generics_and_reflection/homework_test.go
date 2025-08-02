package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go
const omitemptyTag = "omitempty"

type Person struct {
	Name    string `properties:"name"`
	Address string `properties:"address,omitempty"`
	Age     int    `properties:"age"`
	Married bool   `properties:"married"`
}

func Serialize(person Person) string {
	t := reflect.TypeOf(person)
	v := reflect.ValueOf(person)

	var lines []string

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		val := v.Field(i)

		if !val.CanInterface() {
			continue
		}

		tags, ok := field.Tag.Lookup("properties")
		if !ok {
			continue
		}

		var (
			tagName   string
			omitempty bool
		)

		for _, tag := range strings.Split(tags, ",") {
			if tag == omitemptyTag {
				omitempty = true
			} else {
				tagName = tag
			}
		}

		zero := reflect.Zero(field.Type)
		if omitempty && reflect.DeepEqual(val.Interface(), zero.Interface()) {
			continue
		}

		lines = append(lines, fmt.Sprintf("%s=%v", tagName, val.Interface()))
	}

	return strings.Join(lines, "\n")
}
func TestSerialization(t *testing.T) {
	tests := map[string]struct {
		person Person
		result string
	}{
		"test case with empty fields": {
			result: "name=\nage=0\nmarried=false",
		},
		"test case with fields": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
			},
			result: "name=John Doe\nage=30\nmarried=true",
		},
		"test case with omitempty field": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
				Address: "Paris",
			},
			result: "name=John Doe\naddress=Paris\nage=30\nmarried=true",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := Serialize(test.person)
			assert.Equal(t, test.result, result)
		})
	}
}
