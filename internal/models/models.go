package models

import (
	"fmt"
	"reflect"
	"strings"
)

type StatusReturn struct {
	Status string `json:"status"`
}

type CrawlArgs struct {
	Url             string `json:"url"`
	MaxDepth        int    `json:"maxDepth,omitempty"`
	MaxPathSegments int    `json:"maxPathSegments,omitempty"`
}

type CrawlReturn struct {
	Links   []string          `json:"links"`
	Results map[string]string `json:"results"`
}

type GenerateTestsArgs struct {
	Url string `json:"url"`
}

type GenerateTestsReturn struct {
	TestFiles    []TestFile `json:"testFiles" jsonschema_description:"Array of test files to be generated"`
	Dependencies []string   `json:"dependencies" jsonschema_description:"NPM packages required for these tests"`
}

type TestFile struct {
	Filename string `json:"filename" jsonschema_description:"Name of the test file (e.g., 'login.spec.ts')"`
	Content  string `json:"content" jsonschema_description:"Complete content of the test file"`
	FilePath string `json:"-" jsonschema_description:"Absolute path to the temporary file where the test content is stored"`
}

func (r *GenerateTestsReturn) Validate() error {
	var missingFields []string
	val := reflect.ValueOf(*r)
	typ := val.Type()

	for i := range val.NumField() {
		field := val.Field(i)
		fieldTyp := typ.Field(i)

		isValid := true

		switch field.Kind() {
		case reflect.String:
			isValid = field.String() != ""
		case reflect.Slice, reflect.Array:
			isValid = field.Len() > 0
			for j := range field.Len() {
				el := field.Index(j)
				switch el.Kind() {
				case reflect.String:
					isValid = el.String() != ""
				case reflect.Struct:
					if e, ok := el.Interface().(TestFile); ok {
						isValid = e.Filename != "" && e.Content != ""
					}
				}
			}
		case reflect.Map:
			isValid = field.Len() > 0
		case reflect.Ptr, reflect.Interface:
			isValid = !field.IsNil()
		}

		if !isValid {
			jsonTag := fieldTyp.Tag.Get("json")
			fieldName := strings.Split(jsonTag, ",")[0]
			if fieldName == "" {
				fieldName = strings.ToLower(fieldTyp.Name)
			}
			missingFields = append(missingFields, fieldName)
		}

	}

	if len(missingFields) > 0 {
		return fmt.Errorf("required fields: %s", strings.Join(missingFields, ", "))
	}

	return nil
}

type ErrorReturn struct {
	Error string `json:"error"`
}
