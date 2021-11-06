package golden

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"regexp"
	"testing"
	"time"

	"gopkg.in/yaml.v3"
)

//
// Helpers
//

type Author struct {
	FirstName string `json:"first_name" yaml:"first_name" xml:"first_name"`
	LastName  string `json:"last_name" yaml:"last_name" xml:"last_name"`
}

type Book struct {
	ID     string  `json:"id" yaml:"id" xml:"id"`
	Title  string  `json:"title" yaml:"title" xml:"title"`
	Author *Author `json:"author,omitempty" yaml:"author,omitempty" xml:"author,omitempty"`
	Year   int     `json:"year,omitempty" yaml:"year,omitempty" xml:"year,omitempty"`
}

type Article struct {
	ID     string     `json:"id" yaml:"id" xml:"id"`
	Title  string     `json:"title" yaml:"title" xml:"title"`
	Author *Author    `json:"author" yaml:"author" xml:"author"`
	Date   *time.Time `json:"date,omitempty" yaml:"date,omitempty" xml:"date,omitempty"`

	Rank  int `json:"-" yaml:"-" xml:"-"`
	order int
}

type Comic struct {
	ID      string
	Name    string
	Issue   string
	Ignored string
}

type xmlComic struct {
	ID    string `xml:"id,attr"`
	Name  string `xml:",chardata"`
	Issue string `xml:"issue,attr"`
}

func (s *Comic) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"%s":"%s=%s"}`, s.ID, s.Name, s.Issue)), nil
}

func (s *Comic) UnmarshalJSON(data []byte) error {
	m := regexp.MustCompile(`^{\s*"(.*?)":\s*"(.*?)=(.*)"\s*}$`)
	matches := m.FindSubmatch(bytes.TrimSpace(data))
	if matches == nil {
		return nil
	}

	s.ID = string(matches[1])
	s.Name = string(matches[2])
	s.Issue = string(matches[3])

	return nil
}

func (s *Comic) MarshalYAML() (interface{}, error) {
	return map[string]map[string]string{s.ID: {s.Name: s.Issue}}, nil
}

func (s *Comic) UnmarshalYAML(value *yaml.Node) error {
	// Horribly hacky code, but it works and specifically only needs to extract
	// these specific three values.
	if len(value.Content) == 2 {
		s.ID = value.Content[0].Value
		if len(value.Content[1].Content) == 2 {
			s.Name = value.Content[1].Content[0].Value
			s.Issue = value.Content[1].Content[1].Value
		}
	}

	return nil
}

func (s *Comic) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(
		&xmlComic{ID: s.ID, Name: s.Name, Issue: s.Issue},
		start,
	)
}

func (s *Comic) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	x := &xmlComic{}
	_ = d.DecodeElement(x, &start)

	v := Comic{ID: x.ID, Name: x.Name, Issue: x.Issue}

	*s = v

	return nil
}

func boolPtr(b bool) *bool {
	return &b
}

func intPtr(i int) *int {
	return &i
}

func stringPtr(s string) *string {
	return &s
}

//
// Test cases
//

var marhalingTestCases = []struct {
	name string
	v    interface{}
}{
	{
		name: "true bool pointer",
		v:    boolPtr(true),
	},
	{
		name: "false bool pointer",
		v:    boolPtr(false),
	},
	{
		name: "int pointer",
		v:    intPtr(42),
	},
	{
		name: "string pointer",
		v:    stringPtr("hello world"),
	},
	{
		name: "empty struct",
		v:    &Book{},
	},
	{
		name: "partial struct",
		v: &Book{
			ID:    "cfda163c-d5c1-44a2-909b-5d2ce3a31979",
			Title: "The Traveler",
		},
	},
	{
		name: "full struct",
		v: &Book{
			ID:    "cfda163c-d5c1-44a2-909b-5d2ce3a31979",
			Title: "The Traveler",
			Author: &Author{
				FirstName: "John",
				LastName:  "Twelve Hawks",
			},
			Year: 2005,
		},
	},
	{
		name: "custom marshaling",
		v: &Comic{
			ID:    "2fd5af35-b85e-4f03-8eba-524be28d7a5b",
			Name:  "Hello World!",
			Issue: "Forty Two",
		},
	},
}

var articleDate = time.Date(
	2021, time.October, 27, 23, 30, 34, 0, time.FixedZone("", 1*60*60),
).UTC()

var marshalingPTestCases = []struct {
	name string
	v    interface{}
	want interface{}
}{
	{
		name: "true bool pointer",
		v:    boolPtr(true),
		want: boolPtr(true),
	},
	{
		name: "false bool pointer",
		v:    boolPtr(false),
		want: boolPtr(false),
	},
	{
		name: "int pointer",
		v:    intPtr(42),
		want: intPtr(42),
	},
	{
		name: "string pointer",
		v:    stringPtr("hello world"),
		want: stringPtr("hello world"),
	},
	{
		name: "empty struct",
		v:    &Article{},
		want: &Article{},
	},
	{
		name: "partial struct",
		v: &Book{
			ID:    "10eec54d-e30a-4428-be18-01095d889126",
			Title: "Time Travel",
		},
		want: &Book{
			ID:    "10eec54d-e30a-4428-be18-01095d889126",
			Title: "Time Travel",
		},
	},
	{
		name: "full struct",
		v: &Article{
			ID:    "10eec54d-e30a-4428-be18-01095d889126",
			Title: "Time Travel",
			Author: &Author{
				FirstName: "Doc",
				LastName:  "Brown",
			},
			Date:  &articleDate,
			Rank:  8,
			order: 16,
		},
		want: &Article{
			ID:    "10eec54d-e30a-4428-be18-01095d889126",
			Title: "Time Travel",
			Author: &Author{
				FirstName: "Doc",
				LastName:  "Brown",
			},
			Date: &articleDate,
		},
	},
	{
		name: "custom marshaling",
		v: &Comic{
			ID:      "2fd5af35-b85e-4f03-8eba-524be28d7a5b",
			Name:    "Hello World!",
			Issue:   "Forty Two",
			Ignored: "don't pay attention to this :)",
		},
		want: &Comic{
			ID:    "2fd5af35-b85e-4f03-8eba-524be28d7a5b",
			Name:  "Hello World!",
			Issue: "Forty Two",
		},
	},
}

//
// Tests
//

func TestAssertJSONMarshaling(t *testing.T) {
	for _, tt := range marhalingTestCases {
		t.Run(tt.name, func(t *testing.T) {
			AssertJSONMarshaling(t, tt.v)
		})
	}
}

func TestAssertJSONMarshalingP(t *testing.T) {
	for _, tt := range marshalingPTestCases {
		t.Run(tt.name, func(t *testing.T) {
			AssertJSONMarshalingP(t, tt.v, tt.want)
		})
	}
}

func TestAssertXMLMarshaling(t *testing.T) {
	for _, tt := range marhalingTestCases {
		t.Run(tt.name, func(t *testing.T) {
			AssertXMLMarshaling(t, tt.v)
		})
	}
}

func TestAssertXMLMarshalingP(t *testing.T) {
	for _, tt := range marshalingPTestCases {
		t.Run(tt.name, func(t *testing.T) {
			AssertXMLMarshalingP(t, tt.v, tt.want)
		})
	}
}

func TestAssertYAMLMarshaling(t *testing.T) {
	for _, tt := range marhalingTestCases {
		t.Run(tt.name, func(t *testing.T) {
			AssertYAMLMarshaling(t, tt.v)
		})
	}
}

func TestAssertYAMLMarshalingP(t *testing.T) {
	for _, tt := range marshalingPTestCases {
		t.Run(tt.name, func(t *testing.T) {
			AssertYAMLMarshalingP(t, tt.v, tt.want)
		})
	}
}

func TestAssert_JSONMarshaling(t *testing.T) {
	for _, tt := range marhalingTestCases {
		t.Run(tt.name, func(t *testing.T) {
			assert := NewAssert()

			assert.JSONMarshaling(t, tt.v)
		})
	}
}

func TestAssert_JSONMarshalingP(t *testing.T) {
	for _, tt := range marshalingPTestCases {
		t.Run(tt.name, func(t *testing.T) {
			assert := NewAssert()

			assert.JSONMarshalingP(t, tt.v, tt.want)
		})
	}
}

func TestAssert_XMLMarshaling(t *testing.T) {
	for _, tt := range marhalingTestCases {
		t.Run(tt.name, func(t *testing.T) {
			assert := NewAssert()

			assert.XMLMarshaling(t, tt.v)
		})
	}
}

func TestAssert_XMLMarshalingP(t *testing.T) {
	for _, tt := range marshalingPTestCases {
		t.Run(tt.name, func(t *testing.T) {
			assert := NewAssert()

			assert.XMLMarshalingP(t, tt.v, tt.want)
		})
	}
}

func TestAssert_YAMLMarshaling(t *testing.T) {
	for _, tt := range marhalingTestCases {
		t.Run(tt.name, func(t *testing.T) {
			assert := NewAssert()

			assert.YAMLMarshaling(t, tt.v)
		})
	}
}

func TestAssert_YAMLMarshalingP(t *testing.T) {
	for _, tt := range marshalingPTestCases {
		t.Run(tt.name, func(t *testing.T) {
			assert := NewAssert()

			assert.YAMLMarshalingP(t, tt.v, tt.want)
		})
	}
}
