package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Person struct {
	Name    string `properties:"name"`
	Address string `properties:"address,omitempty"`
	Age     int    `properties:"age"`
	Married bool   `properties:"married"`
}

func Serialize(person Person) string {
	v := reflect.ValueOf(person)
	t := v.Type()

	var builder strings.Builder
	first := true

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		tag, ok := field.Tag.Lookup("properties")
		if !ok {
			continue
		}

		parts := strings.Split(tag, ",")
		name := parts[0]
		omitempty := len(parts) > 1 && parts[1] == "omitempty"

		if omitempty && value.IsZero() {
			continue
		}

		if !first {
			builder.WriteByte('\n')
		}
		first = false

		switch value.Kind() {
		case reflect.Bool:
			fmt.Fprintf(&builder, "%s=%v", name, value.Bool())
		case reflect.String:
			fmt.Fprintf(&builder, "%s=%s", name, value.String())
		case reflect.Int, reflect.Int8, reflect.Int32:
			fmt.Fprintf(&builder, "%s=%d", name, value.Int())
		default:
			panic("unhandled default case")
		}
	}

	return builder.String()
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
