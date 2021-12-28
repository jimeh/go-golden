package golden_test

import (
	"testing"

	"github.com/jimeh/go-golden"
)

type MyBook struct {
	FooBar string `json:"foo_bar,omitempty" yaml:"fooBar,omitempty" xml:"Foo_Bar,omitempty"`
	Bar    string `json:"-" yaml:"-" xml:"-"`
	baz    string
}

// TestExampleMyBookMarshaling reads/writes the following golden files:
//
//  testdata/TestExampleMyBookMarshaling/marshaled_json.golden
//  testdata/TestExampleMyBookMarshaling/marshaled_xml.golden
//  testdata/TestExampleMyBookMarshaling/marshaled_yaml.golden
//
func TestExampleMyBookMarshaling(t *testing.T) {
	obj := &MyBook{FooBar: "Hello World!"}

	golden.AssertJSONMarshaling(t, obj)
	golden.AssertYAMLMarshaling(t, obj)
	golden.AssertXMLMarshaling(t, obj)
}

// TestExampleMyBookMarshalingP reads/writes the following golden files:
//
//  testdata/TestExampleMyBookMarshalingP/marshaled_json.golden
//  testdata/TestExampleMyBookMarshalingP/marshaled_xml.golden
//  testdata/TestExampleMyBookMarshalingP/marshaled_yaml.golden
//
func TestExampleMyBookMarshalingP(t *testing.T) {
	obj := &MyBook{FooBar: "Hello World!", Bar: "Oops", baz: "nope!"}
	want := &MyBook{FooBar: "Hello World!"}

	golden.AssertJSONMarshalingP(t, obj, want)
	golden.AssertYAMLMarshalingP(t, obj, want)
	golden.AssertXMLMarshalingP(t, obj, want)
}

// TestExampleMyBookMarshalingTabular reads/writes the following golden files:
//
//  testdata/TestExampleMyBookMarshalingTabular/empty/marshaled_json.golden
//  testdata/TestExampleMyBookMarshalingTabular/empty/marshaled_xml.golden
//  testdata/TestExampleMyBookMarshalingTabular/empty/marshaled_yaml.golden
//  testdata/TestExampleMyBookMarshalingTabular/full/marshaled_json.golden
//  testdata/TestExampleMyBookMarshalingTabular/full/marshaled_xml.golden
//  testdata/TestExampleMyBookMarshalingTabular/full/marshaled_yaml.golden
//
func TestExampleMyBookMarshalingTabular(t *testing.T) {
	tests := []struct {
		name string
		obj  *MyBook
	}{
		{name: "empty", obj: &MyBook{}},
		{name: "full", obj: &MyBook{FooBar: "Hello World!"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			golden.AssertJSONMarshaling(t, tt.obj)
			golden.AssertYAMLMarshaling(t, tt.obj)
			golden.AssertXMLMarshaling(t, tt.obj)
		})
	}
}
