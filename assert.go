package golden

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"testing"

	"gopkg.in/yaml.v3"
)

var globalAssert = NewAssert()

// AssertJSONMarshaling asserts that the given "v" value JSON marshals to an
// expected value fetched from a golden file on disk, and then verifies that the
// marshaled result produces a value that is equal to "v" when unmarshaled.
//
// Used for objects that do NOT change when they are marshaled and unmarshaled.
func AssertJSONMarshaling(t *testing.T, v interface{}) {
	t.Helper()

	globalAssert.JSONMarshaling(t, v)
}

// AssertJSONMarshalingP asserts that the given "v" value JSON marshals to an
// expected value fetched from a golden file on disk, and then verifies that the
// marshaled result produces a value that is equal to "want" when unmarshaled.
//
// Used for objects that change when they are marshaled and unmarshaled.
func AssertJSONMarshalingP(t *testing.T, v, want interface{}) {
	t.Helper()

	globalAssert.JSONMarshalingP(t, v, want)
}

// AssertXMLMarshaling asserts that the given "v" value XML marshals to an
// expected value fetched from a golden file on disk, and then verifies that the
// marshaled result produces a value that is equal to "v" when unmarshaled.
//
// Used for objects that do NOT change when they are marshaled and unmarshaled.
func AssertXMLMarshaling(t *testing.T, v interface{}) {
	t.Helper()

	globalAssert.XMLMarshaling(t, v)
}

// AssertXMLMarshalingP asserts that the given "v" value XML marshals to an
// expected value fetched from a golden file on disk, and then verifies that the
// marshaled result produces a value that is equal to "want" when unmarshaled.
//
// Used for objects that change when they are marshaled and unmarshaled.
func AssertXMLMarshalingP(t *testing.T, v, want interface{}) {
	t.Helper()

	globalAssert.XMLMarshalingP(t, v, want)
}

// AssertYAMLMarshaling asserts that the given "v" value YAML marshals to an
// expected value fetched from a golden file on disk, and then verifies that the
// marshaled result produces a value that is equal to "v" when unmarshaled.
//
// Used for objects that do NOT change when they are marshaled and unmarshaled.
func AssertYAMLMarshaling(t *testing.T, v interface{}) {
	t.Helper()

	globalAssert.YAMLMarshaling(t, v)
}

// AssertYAMLMarshalingP asserts that the given "v" value YAML marshals to an
// expected value fetched from a golden file on disk, and then verifies that the
// marshaled result produces a value that is equal to "want" when unmarshaled.
//
// Used for objects that change when they are marshaled and unmarshaled.
func AssertYAMLMarshalingP(t *testing.T, v, want interface{}) {
	t.Helper()

	globalAssert.YAMLMarshalingP(t, v, want)
}

// Assert exposes a series of JSON, YAML, and XML marshaling assertion helpers.
type Assert interface {
	// JSONMarshaling asserts that the given "v" value JSON marshals to an
	// expected value fetched from a golden file on disk, and then verifies that
	// the marshaled result produces a value that is equal to "v" when
	// unmarshaled.
	//
	// Used for objects that do NOT change when they are marshaled and
	// unmarshaled.
	JSONMarshaling(t *testing.T, v interface{})

	// JSONMarshalingP asserts that the given "v" value JSON marshals to an
	// expected value fetched from a golden file on disk, and then verifies that
	// the marshaled result produces a value that is equal to "want" when
	// unmarshaled.
	//
	// Used for objects that change when they are marshaled and unmarshaled.
	JSONMarshalingP(t *testing.T, v interface{}, want interface{})

	// XMLMarshaling asserts that the given "v" value XML marshals to an
	// expected value fetched from a golden file on disk, and then verifies that
	// the marshaled result produces a value that is equal to "v" when
	// unmarshaled.
	//
	// Used for objects that do NOT change when they are marshaled and
	// unmarshaled.
	XMLMarshaling(t *testing.T, v interface{})

	// XMLMarshalingP asserts that the given "v" value XML marshals to an
	// expected value fetched from a golden file on disk, and then verifies that
	// the marshaled result produces a value that is equal to "want" when
	// unmarshaled.
	//
	// Used for objects that change when they are marshaled and unmarshaled.
	XMLMarshalingP(t *testing.T, v interface{}, want interface{})

	// YAMLMarshaling asserts that the given "v" value YAML marshals to an
	// expected value fetched from a golden file on disk, and then verifies that
	// the marshaled result produces a value that is equal to "v" when
	// unmarshaled.
	//
	// Used for objects that do NOT change when they are marshaled and
	// unmarshaled.
	YAMLMarshaling(t *testing.T, v interface{})

	// YAMLMarshalingP asserts that the given "v" value YAML marshals to an
	// expected value fetched from a golden file on disk, and then verifies that
	// the marshaled result produces a value that is equal to "want" when
	// unmarshaled.
	//
	// Used for objects that change when they are marshaled and unmarshaled.
	YAMLMarshalingP(t *testing.T, v interface{}, want interface{})
}

type AssertOption interface {
	apply(*asserter)
}

type assertOptionFunc func(*asserter)

func (fn assertOptionFunc) apply(c *asserter) {
	fn(c)
}

// WithGolden allows setting a custom *Golden instance when calling NewAssert().
func WithGolden(golden *Golden) AssertOption {
	return assertOptionFunc(func(a *asserter) {
		a.golden = golden
	})
}

// WithNormalizedLineBreaks allows turning off line-break normalization which
// replaces Windows' CRLF (\r\n) and Mac Classic CR (\r) line breaks with Unix's
// LF (\n) line breaks.
func WithNormalizedLineBreaks(value bool) AssertOption {
	return assertOptionFunc(func(a *asserter) {
		a.normalizeLineBreaks = value
	})
}

// NewAssert returns a new Assert which exposes a number of marshaling assertion
// helpers for JSON, YAML and XML.
//
// The default encoders all specify indentation of two spaces, essentially
// enforcing pretty formatting for JSON and XML.
//
// The default decoders for JSON and YAML prohibit unknown fields which are not
// present on the provided struct.
func NewAssert(options ...AssertOption) Assert {
	a := &asserter{
		golden:              globalGolden,
		normalizeLineBreaks: true,
	}

	for _, opt := range options {
		opt.apply(a)
	}

	a.JSONAsserter = NewMarshalAsserter(
		a.golden, "JSON",
		newJSONEncoder, newJSONDecoder,
		a.normalizeLineBreaks,
	)
	a.XMLAsserter = NewMarshalAsserter(
		a.golden, "XML",
		newXMLEncoder, newXMLDecoder,
		a.normalizeLineBreaks,
	)
	a.YAMLAsserter = NewMarshalAsserter(
		a.golden, "YAML",
		newYAMLEncoder, newYAMLDecoder,
		a.normalizeLineBreaks,
	)

	return a
}

// asserter implements the Assert interface.
type asserter struct {
	golden              *Golden
	normalizeLineBreaks bool

	JSONAsserter *MarshalAsserter
	XMLAsserter  *MarshalAsserter
	YAMLAsserter *MarshalAsserter
}

func (s *asserter) JSONMarshaling(t *testing.T, v interface{}) {
	t.Helper()

	s.JSONAsserter.Marshaling(t, v)
}

func (s *asserter) JSONMarshalingP(
	t *testing.T,
	v interface{},
	want interface{},
) {
	t.Helper()

	s.JSONAsserter.MarshalingP(t, v, want)
}

func (s *asserter) XMLMarshaling(t *testing.T, v interface{}) {
	t.Helper()

	s.XMLAsserter.Marshaling(t, v)
}

func (s *asserter) XMLMarshalingP(t *testing.T, v, want interface{}) {
	t.Helper()

	s.XMLAsserter.MarshalingP(t, v, want)
}

func (s *asserter) YAMLMarshaling(t *testing.T, v interface{}) {
	t.Helper()

	s.YAMLAsserter.Marshaling(t, v)
}

func (s *asserter) YAMLMarshalingP(t *testing.T, v, want interface{}) {
	t.Helper()

	s.YAMLAsserter.MarshalingP(t, v, want)
}

// newJSONEncoder is the default JSONEncoderFunc used by Assert. It returns a
// *json.Encoder which is set to indent with two spaces.
func newJSONEncoder(w io.Writer) MarshalEncoder {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")

	return enc
}

// newJSONDecoder is the default JSONDecoderFunc used by Assert. It returns a
// *json.Decoder which disallows unknown fields.
func newJSONDecoder(r io.Reader) MarshalDecoder {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	return dec
}

// newXMLEncoder is the default XMLEncoderFunc used by Assert. It returns a
// *xml.Encoder which is set to indent with two spaces.
func newXMLEncoder(w io.Writer) MarshalEncoder {
	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")

	return enc
}

// newXMLDecoder is the default XMLDecoderFunc used by Assert.
func newXMLDecoder(r io.Reader) MarshalDecoder {
	return xml.NewDecoder(r)
}

// newYAMLEncoder is the default YAMLEncoderFunc used by Assert. It returns a
// *yaml.Encoder which is set to indent with two spaces.
func newYAMLEncoder(w io.Writer) MarshalEncoder {
	enc := yaml.NewEncoder(w)
	enc.SetIndent(2)

	return enc
}

// newYAMLDecoder is the default YAMLDecoderFunc used by Assert. It returns a
// *yaml.Decoder which disallows unknown fields.
func newYAMLDecoder(r io.Reader) MarshalDecoder {
	dec := yaml.NewDecoder(r)
	dec.KnownFields(true)

	return dec
}

// normalizeLineBreaks replaces Windows CRLF (\r\n) and Classic MacOS CR (\r)
// line-breaks with Unix LF (\n) line breaks.
func normalizeLineBreaks(data []byte) []byte {
	// Replace Windows CRLF (\r\n) with Unix LF (\n)
	result := bytes.ReplaceAll(data, []byte{13, 10}, []byte{10})
	// Replace Classic MacOS CR (\r) with Unix LF (\n)
	result = bytes.ReplaceAll(result, []byte{13}, []byte{10})

	return result
}
