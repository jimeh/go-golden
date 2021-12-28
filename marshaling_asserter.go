package golden

import (
	"reflect"
	"strings"

	"github.com/jimeh/go-golden/sanitize"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type (
	MarshalFunc   func(interface{}) ([]byte, error)
	UnmarshalFunc func([]byte, interface{}) error
)

// MarshalingAsserter allows building marshaling asserters by providing
// functions which marshal and unmarshal objects.
type MarshalingAsserter struct {
	// Golden is the *Golden instance used to read/write golden files.
	Golden Golden

	// Name of the format the MarshalAsserter handles.
	Format string

	// GoldName is the name of the golden file used when marshaling. This is by
	// default set based on Format using NewMarshalAsserter. For example if
	// Format is set to "JSON", GoldName will be set to "marshaled_json" by
	// default.
	GoldName string

	// MarshalFunc is the function used to marshal given objects.
	MarshalFunc MarshalFunc

	// UnmarshalFunc is the function used to unmarshal given objects.
	UnmarshalFunc UnmarshalFunc

	// NormalizeLineBreaks determines if Windows' CRLF (\r\n) and MacOS Classic
	// CR (\r) line breaks are replaced with Unix's LF (\n) line breaks. This
	// ensures marshaling assertions work across different platforms.
	NormalizeLineBreaks bool
}

// New returns a new MarshalingAsserter.
func NewMarshalingAsserter(
	golden Golden,
	format string,
	marshalFunc MarshalFunc,
	unmarshalFunc UnmarshalFunc,
	normalizeLineBreaks bool,
) *MarshalingAsserter {
	goldName := "marshaled_" + strings.ToLower(sanitize.Filename(format))

	return &MarshalingAsserter{
		Golden:              golden,
		Format:              format,
		GoldName:            goldName,
		MarshalFunc:         marshalFunc,
		UnmarshalFunc:       unmarshalFunc,
		NormalizeLineBreaks: normalizeLineBreaks,
	}
}

// Marshaling asserts that the given "v" value marshals via the provided encoder
// to an expected value fetched from a golden file on disk, and then verifies
// that the marshaled result produces a value that is equal to "v" when
// unmarshaled.
//
// Used for objects that do NOT change when they are marshaled and unmarshaled.
func (s *MarshalingAsserter) Marshaling(t TestingT, v interface{}) {
	t.Helper()

	s.MarshalingP(t, v, v)
}

// MarshalingP asserts that the given "v" value marshals via the provided
// encoder to an expected value fetched from a golden file on disk, and then
// verifies that the marshaled result produces a value that is equal to "want"
// when unmarshaled.
//
// Used for objects that change when they are marshaled and unmarshaled.
func (s *MarshalingAsserter) MarshalingP(
	t TestingT,
	v interface{},
	want interface{},
) {
	t.Helper()

	if reflect.ValueOf(want).Kind() != reflect.Ptr {
		require.FailNowf(t,
			"golden: only pointer types can be asserted",
			"%T is not a pointer type", want,
		)
	}

	marshaled, err := s.MarshalFunc(v)
	require.NoErrorf(t,
		err, "golden: failed to %s marshal %T: %+v", s.Format, v, v,
	)
	if s.NormalizeLineBreaks {
		marshaled = sanitize.LineBreaks(marshaled)
	}

	if s.Golden.Update() {
		s.Golden.SetP(t, s.GoldName, marshaled)
	}

	gold := s.Golden.GetP(t, s.GoldName)
	if s.NormalizeLineBreaks {
		gold = sanitize.LineBreaks(gold)
	}

	switch strings.ToLower(s.Format) {
	case "json":
		assert.JSONEq(t, string(gold), string(marshaled))
	case "yaml", "yml":
		assert.YAMLEq(t, string(gold), string(marshaled))
	default:
		assert.Equal(t, string(gold), string(marshaled))
	}

	got := reflect.New(reflect.TypeOf(want).Elem()).Interface()
	err = s.UnmarshalFunc(gold, got)

	f := s.Golden.FileP(t, s.GoldName)
	require.NoErrorf(t, err,
		"golden: failed to %s unmarshal %T from %s", s.Format, got, f,
	)
	assert.Equalf(t, want, got,
		"golden: unmarshaling from golden file does not match "+
			"expected object; golden file: %s", f,
	)
}
