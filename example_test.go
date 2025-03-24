package golden_test

import (
	"encoding/json"
	"encoding/xml"
	"testing"

	"github.com/jimeh/go-golden"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The tests in this file are examples from the README and the package-level Go
// documentation.

type MyStruct struct {
	Foo string `json:"foo,omitempty"`
}

// TestExampleMyStruct reads/writes the following golden file:
//
//	testdata/TestExampleMyStruct.golden
func TestExampleMyStruct(t *testing.T) {
	got, err := json.Marshal(&MyStruct{Foo: "Bar"})
	require.NoError(t, err)

	want := golden.Do(t, got)

	assert.Equal(t, want, got)
}

// TestExampleMyStructTabular reads/writes the following golden files:
//
//	testdata/TestExampleMyStructTabular/empty_struct.golden
//	testdata/TestExampleMyStructTabular/full_struct.golden
func TestExampleMyStructTabular(t *testing.T) {
	tests := []struct {
		name string
		obj  *MyStruct
	}{
		{name: "empty struct", obj: &MyStruct{}},
		{name: "full struct", obj: &MyStruct{Foo: "Bar"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.obj)
			require.NoError(t, err)

			want := golden.Do(t, got)

			assert.Equal(t, want, got)
		})
	}
}

// TestExampleMyStructP reads/writes the following golden file:
//
//	testdata/TestExampleMyStructP/json.golden
//	testdata/TestExampleMyStructP/xml.golden
func TestExampleMyStructP(t *testing.T) {
	gotJSON, _ := json.Marshal(&MyStruct{Foo: "Bar"})
	gotXML, _ := xml.Marshal(&MyStruct{Foo: "Bar"})

	wantJSON := golden.DoP(t, "json", gotJSON)
	wantXML := golden.DoP(t, "xml", gotXML)

	assert.Equal(t, wantJSON, gotJSON)
	assert.Equal(t, wantXML, gotXML)
}

// TestExampleMyStructTabularP reads/writes the following golden file:
//
//	testdata/TestExampleMyStructTabularP/empty_struct/json.golden
//	testdata/TestExampleMyStructTabularP/empty_struct/xml.golden
//	testdata/TestExampleMyStructTabularP/full_struct/json.golden
//	testdata/TestExampleMyStructTabularP/full_struct/xml.golden
func TestExampleMyStructTabularP(t *testing.T) {
	tests := []struct {
		name string
		obj  *MyStruct
	}{
		{name: "empty struct", obj: &MyStruct{}},
		{name: "full struct", obj: &MyStruct{Foo: "Bar"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotJSON, _ := json.Marshal(tt.obj)
			gotXML, _ := xml.Marshal(tt.obj)

			wantJSON := golden.DoP(t, "json", gotJSON)
			wantXML := golden.DoP(t, "xml", gotXML)

			assert.Equal(t, wantJSON, gotJSON)
			assert.Equal(t, wantXML, gotXML)
		})
	}
}
