package tools

import (
	"context"
	"fmt"
	"reflect"
)

func GetContextDetails(c context.Context) {
	rv := reflect.ValueOf(c)
	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		rv = rv.Elem()
	}

	if rv.Kind() == reflect.Struct {
		for i := 0; i < rv.NumField(); i++ {
			f := rv.Type().Field(i)

			if f.Name == "key" {
			}
			if f.Name == "Context" {
				// this is just a repetition of the above, so you can make a recursive
				// function from it, or for loop, that stops when there are no more
				// contexts to be inspected.
				rv := rv.Field(i)
				for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
					rv = rv.Elem()
				}

				if rv.Kind() == reflect.Struct {
					for i := 0; i < rv.NumField(); i++ {
						f := rv.Type().Field(i)

						if f.Name == "key" {
						} else {
							fmt.Printf("value: %+v\n", rv.Field(i))
						}
						// ...
					}
				}
			}
		}
	}
}
