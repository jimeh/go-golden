package golden

import (
	"github.com/jimeh/go-golden/marshal"
	"github.com/jimeh/go-golden/unmarshal"
)

// Asserter exposes a series of JSON, YAML, and XML marshaling assertion
// helpers.
type Asserter interface {
	// JSONMarshaling asserts that the given "v" value JSON marshals to an
	// expected value fetched from a golden file on disk, and then verifies that
	// the marshaled result produces a value that is equal to "v" when
	// unmarshaled.
	//
	// Used for objects that do NOT change when they are marshaled and
	// unmarshaled.
	JSONMarshaling(t TestingT, v interface{})

	// JSONMarshalingP asserts that the given "v" value JSON marshals to an
	// expected value fetched from a golden file on disk, and then verifies that
	// the marshaled result produces a value that is equal to "want" when
	// unmarshaled.
	//
	// Used for objects that change when they are marshaled and unmarshaled.
	JSONMarshalingP(t TestingT, v interface{}, want interface{})

	// XMLMarshaling asserts that the given "v" value XML marshals to an
	// expected value fetched from a golden file on disk, and then verifies that
	// the marshaled result produces a value that is equal to "v" when
	// unmarshaled.
	//
	// Used for objects that do NOT change when they are marshaled and
	// unmarshaled.
	XMLMarshaling(t TestingT, v interface{})

	// XMLMarshalingP asserts that the given "v" value XML marshals to an
	// expected value fetched from a golden file on disk, and then verifies that
	// the marshaled result produces a value that is equal to "want" when
	// unmarshaled.
	//
	// Used for objects that change when they are marshaled and unmarshaled.
	XMLMarshalingP(t TestingT, v interface{}, want interface{})

	// YAMLMarshaling asserts that the given "v" value YAML marshals to an
	// expected value fetched from a golden file on disk, and then verifies that
	// the marshaled result produces a value that is equal to "v" when
	// unmarshaled.
	//
	// Used for objects that do NOT change when they are marshaled and
	// unmarshaled.
	YAMLMarshaling(t TestingT, v interface{})

	// YAMLMarshalingP asserts that the given "v" value YAML marshals to an
	// expected value fetched from a golden file on disk, and then verifies that
	// the marshaled result produces a value that is equal to "want" when
	// unmarshaled.
	//
	// Used for objects that change when they are marshaled and unmarshaled.
	YAMLMarshalingP(t TestingT, v interface{}, want interface{})
}

// NewAsserter returns a new Asserter which exposes a number of marshaling
// assertion helpers for JSON, YAML and XML.
//
// The default encoders all specify indentation of two spaces, essentially
// enforcing pretty formatting for JSON and XML.
//
// The default decoders for JSON and YAML prohibit unknown fields which are not
// present on the provided struct.
func NewAsserter(options ...AsserterOption) Asserter {
	o := &asserterOptions{
		golden:              defaultGolden,
		normalizeLineBreaks: true,
	}

	for _, opt := range options {
		opt.apply(o)
	}

	return &asserter{
		json: NewMarshalingAsserter(
			o.golden, "JSON",
			marshal.JSON, unmarshal.JSON,
			o.normalizeLineBreaks,
		),
		xml: NewMarshalingAsserter(
			o.golden, "XML",
			marshal.XML, unmarshal.XML,
			o.normalizeLineBreaks,
		),
		yaml: NewMarshalingAsserter(
			o.golden, "YAML",
			marshal.YAML, unmarshal.YAML,
			o.normalizeLineBreaks,
		),
	}
}

type asserterOptions struct {
	golden              Golden
	normalizeLineBreaks bool
}

type AsserterOption interface {
	apply(*asserterOptions)
}

type asserterOptionFunc func(*asserterOptions)

func (fn asserterOptionFunc) apply(c *asserterOptions) {
	fn(c)
}

// WithGolden allows setting a custom *Golden instance when calling NewAssert().
func WithGolden(golden Golden) AsserterOption {
	return asserterOptionFunc(func(a *asserterOptions) {
		if golden != nil {
			a.golden = golden
		}
	})
}

// WithNormalizedLineBreaks allows turning off line-break normalization which
// replaces Windows' CRLF (\r\n) and Mac Classic CR (\r) line breaks with Unix's
// LF (\n) line breaks.
func WithNormalizedLineBreaks(value bool) AsserterOption {
	return asserterOptionFunc(func(a *asserterOptions) {
		a.normalizeLineBreaks = value
	})
}

// asserter implements the Assert interface.
type asserter struct {
	json *MarshalingAsserter
	xml  *MarshalingAsserter
	yaml *MarshalingAsserter
}

func (s *asserter) JSONMarshaling(t TestingT, v interface{}) {
	t.Helper()

	s.json.Marshaling(t, v)
}

func (s *asserter) JSONMarshalingP(
	t TestingT,
	v interface{},
	want interface{},
) {
	t.Helper()

	s.json.MarshalingP(t, v, want)
}

func (s *asserter) XMLMarshaling(t TestingT, v interface{}) {
	t.Helper()

	s.xml.Marshaling(t, v)
}

func (s *asserter) XMLMarshalingP(t TestingT, v, want interface{}) {
	t.Helper()

	s.xml.MarshalingP(t, v, want)
}

func (s *asserter) YAMLMarshaling(t TestingT, v interface{}) {
	t.Helper()

	s.yaml.Marshaling(t, v)
}

func (s *asserter) YAMLMarshalingP(t TestingT, v, want interface{}) {
	t.Helper()

	s.yaml.MarshalingP(t, v, want)
}
