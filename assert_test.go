package golden

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/jimeh/go-golden/marshal"
	"github.com/jimeh/go-golden/unmarshal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

//
// Helpers
//

type author struct {
	FirstName string `json:"first_name" yaml:"first_name" xml:"first_name"`
	LastName  string `json:"last_name" yaml:"last_name" xml:"last_name"`
}

type book struct {
	ID     string  `json:"id" yaml:"id" xml:"id"`
	Title  string  `json:"title" yaml:"title" xml:"title"`
	Author *author `json:"author,omitempty" yaml:"author,omitempty" xml:"author,omitempty"`
	Year   int     `json:"year,omitempty" yaml:"year,omitempty" xml:"year,omitempty"`
}

type article struct {
	ID     string     `json:"id" yaml:"id" xml:"id"`
	Title  string     `json:"title" yaml:"title" xml:"title"`
	Author *author    `json:"author" yaml:"author" xml:"author"`
	Date   *time.Time `json:"date,omitempty" yaml:"date,omitempty" xml:"date,omitempty"`

	Rank  int `json:"-" yaml:"-" xml:"-"`
	order int
}

// comic is used for testing custom marshal/unmarshal functions on a type.
type comic struct {
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

func (s *comic) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"%s":"%s=%s"}`, s.ID, s.Name, s.Issue)), nil
}

func (s *comic) UnmarshalJSON(data []byte) error {
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

func (s *comic) MarshalYAML() (interface{}, error) {
	return map[string]map[string]string{s.ID: {s.Name: s.Issue}}, nil
}

func (s *comic) UnmarshalYAML(value *yaml.Node) error {
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

func (s *comic) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(
		&xmlComic{ID: s.ID, Name: s.Name, Issue: s.Issue},
		start,
	)
}

func (s *comic) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	x := &xmlComic{}
	_ = d.DecodeElement(x, &start)

	v := comic{ID: x.ID, Name: x.Name, Issue: x.Issue}

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
		v:    &book{},
	},
	{
		name: "partial struct",
		v: &book{
			ID:    "cfda163c-d5c1-44a2-909b-5d2ce3a31979",
			Title: "The Traveler",
		},
	},
	{
		name: "full struct",
		v: &book{
			ID:    "cfda163c-d5c1-44a2-909b-5d2ce3a31979",
			Title: "The Traveler",
			Author: &author{
				FirstName: "John",
				LastName:  "Twelve Hawks",
			},
			Year: 2005,
		},
	},
	{
		name: "custom marshaling",
		v: &comic{
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
		v:    &article{},
		want: &article{},
	},
	{
		name: "partial struct",
		v: &book{
			ID:    "10eec54d-e30a-4428-be18-01095d889126",
			Title: "Time Travel",
		},
		want: &book{
			ID:    "10eec54d-e30a-4428-be18-01095d889126",
			Title: "Time Travel",
		},
	},
	{
		name: "full struct",
		v: &article{
			ID:    "10eec54d-e30a-4428-be18-01095d889126",
			Title: "Time Travel",
			Author: &author{
				FirstName: "Doc",
				LastName:  "Brown",
			},
			Date:  &articleDate,
			Rank:  8,
			order: 16,
		},
		want: &article{
			ID:    "10eec54d-e30a-4428-be18-01095d889126",
			Title: "Time Travel",
			Author: &author{
				FirstName: "Doc",
				LastName:  "Brown",
			},
			Date: &articleDate,
		},
	},
	{
		name: "custom marshaling",
		v: &comic{
			ID:      "2fd5af35-b85e-4f03-8eba-524be28d7a5b",
			Name:    "Hello World!",
			Issue:   "Forty Two",
			Ignored: "don't pay attention to this :)",
		},
		want: &comic{
			ID:    "2fd5af35-b85e-4f03-8eba-524be28d7a5b",
			Name:  "Hello World!",
			Issue: "Forty Two",
		},
	},
}

//
// Tests
//

func TestWithGolden(t *testing.T) {
	type args struct {
		golden Golden
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "nil",
			args: args{golden: nil},
		},
		{
			name: "non-nil",
			args: args{golden: New(WithSuffix(".my-custom-golden"))},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			original := New(WithSuffix(".original-golden"))
			o := &asserterOptions{golden: original}

			fn := WithGolden(tt.args.golden)
			fn.apply(o)

			if tt.args.golden == nil {
				assert.Equal(t, original, o.golden)
			} else {
				assert.Equal(t, tt.args.golden, o.golden)
			}
		})
	}
}

func TestNormalizedLineBreaks(t *testing.T) {
	type args struct {
		value bool
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "true",
			args: args{value: true},
		},
		{
			name: "false",
			args: args{value: false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &asserterOptions{normalizeLineBreaks: !tt.args.value}

			fn := WithNormalizedLineBreaks(tt.args.value)
			fn.apply(o)

			assert.Equal(t, tt.args.value, o.normalizeLineBreaks)
		})
	}
}

func TestNewAssert(t *testing.T) {
	otherGolden := New(WithSuffix(".other-golden"))

	type args struct {
		options []AsserterOption
	}

	tests := []struct {
		name string
		args args
		want *asserter
	}{
		{
			name: "no options",
			args: args{options: []AsserterOption{}},
			want: &asserter{
				json: NewMarshalingAsserter(
					defaultGolden, "JSON",
					marshal.JSON, unmarshal.JSON,
					true,
				),
				xml: NewMarshalingAsserter(
					defaultGolden, "XML",
					marshal.XML, unmarshal.XML,
					true,
				),
				yaml: NewMarshalingAsserter(
					defaultGolden, "YAML",
					marshal.YAML, unmarshal.YAML,
					true,
				),
			},
		},
		{
			name: "WithGlobal option",
			args: args{
				options: []AsserterOption{
					WithGolden(otherGolden),
				},
			},
			want: &asserter{
				json: NewMarshalingAsserter(
					otherGolden, "JSON",
					marshal.JSON, unmarshal.JSON,
					true,
				),
				xml: NewMarshalingAsserter(
					otherGolden, "XML",
					marshal.XML, unmarshal.XML,
					true,
				),
				yaml: NewMarshalingAsserter(
					otherGolden, "YAML",
					marshal.YAML, unmarshal.YAML,
					true,
				),
			},
		},
		{
			name: "WithNormalizedLineBreaks option",
			args: args{
				options: []AsserterOption{
					WithNormalizedLineBreaks(false),
				},
			},
			want: &asserter{
				json: NewMarshalingAsserter(
					defaultGolden, "JSON",
					marshal.JSON, unmarshal.JSON,
					false,
				),
				xml: NewMarshalingAsserter(
					defaultGolden, "XML",
					marshal.XML, unmarshal.XML,
					false,
				),
				yaml: NewMarshalingAsserter(
					defaultGolden, "YAML",
					marshal.YAML, unmarshal.YAML,
					false,
				),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewAsserter(tt.args.options...)
			assert.Implements(t, (*Asserter)(nil), a)

			got, ok := a.(*asserter)
			require.True(
				t, ok, "failed to type assert return value to a *asserter",
			)

			assert.Equal(t, tt.want.json.Golden, got.json.Golden)
			assert.Equal(t, tt.want.json.Format, got.json.Format)
			assert.Equal(t,
				tt.want.json.NormalizeLineBreaks, got.json.NormalizeLineBreaks,
			)

			assert.Equal(t, tt.want.xml.Golden, got.xml.Golden)
			assert.Equal(t, tt.want.xml.Format, got.xml.Format)
			assert.Equal(t,
				tt.want.xml.NormalizeLineBreaks, got.xml.NormalizeLineBreaks,
			)

			assert.Equal(t, tt.want.yaml.Golden, got.yaml.Golden)
			assert.Equal(t, tt.want.yaml.Format, got.yaml.Format)
			assert.Equal(t,
				tt.want.yaml.NormalizeLineBreaks, got.yaml.NormalizeLineBreaks,
			)
		})
	}
}

func Test_asserter(t *testing.T) {
	a := &asserter{}

	assert.Implements(t, (*Asserter)(nil), a)
}

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
			assert := NewAsserter()

			assert.JSONMarshaling(t, tt.v)
		})
	}
}

func TestAssert_JSONMarshalingP(t *testing.T) {
	for _, tt := range marshalingPTestCases {
		t.Run(tt.name, func(t *testing.T) {
			assert := NewAsserter()

			assert.JSONMarshalingP(t, tt.v, tt.want)
		})
	}
}

func TestAssert_XMLMarshaling(t *testing.T) {
	for _, tt := range marhalingTestCases {
		t.Run(tt.name, func(t *testing.T) {
			assert := NewAsserter()

			assert.XMLMarshaling(t, tt.v)
		})
	}
}

func TestAssert_XMLMarshalingP(t *testing.T) {
	for _, tt := range marshalingPTestCases {
		t.Run(tt.name, func(t *testing.T) {
			assert := NewAsserter()

			assert.XMLMarshalingP(t, tt.v, tt.want)
		})
	}
}

func TestAssert_YAMLMarshaling(t *testing.T) {
	for _, tt := range marhalingTestCases {
		t.Run(tt.name, func(t *testing.T) {
			assert := NewAsserter()

			assert.YAMLMarshaling(t, tt.v)
		})
	}
}

func TestAssert_YAMLMarshalingP(t *testing.T) {
	for _, tt := range marshalingPTestCases {
		t.Run(tt.name, func(t *testing.T) {
			assert := NewAsserter()

			assert.YAMLMarshalingP(t, tt.v, tt.want)
		})
	}
}
