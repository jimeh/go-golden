package golden

import (
	"bytes"
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type (
	NewEncoderFunc func(w io.Writer) MarshalEncoder
	NewDecoderFunc func(r io.Reader) MarshalDecoder
)

type MarshalEncoder interface {
	Encode(v interface{}) error
}

type MarshalDecoder interface {
	Decode(v interface{}) error
}

// MarshalAsserter allows building custom marshaling asserters, but providing
// functions which returns new encoder and decoders for the format to be
// asserted.
//
// All the Assert<format>Marshaling helper functions uses MarshalAsserter under
// the hood.
type MarshalAsserter struct {
	// Golden is the *Golden instance used to read/write golden files.
	Golden *Golden

	// Name of the format the MarshalAsserter handles.
	Format string

	// GoldName is the name of the golden file used when marshaling. This is by
	// default set based on Format using NewMarshalAsserter. For example if
	// Format is set to "JSON", GoldName will be set to "marshaled_json" by
	// default.
	GoldName string

	// NewEncoderFunc is the function used to create a new encoder for
	// marshaling objects.
	NewEncoderFunc NewEncoderFunc

	// NewDecoderFunc is the function used to create a new decoder for
	// unmarshaling objects.
	NewDecoderFunc NewDecoderFunc

	// NormalizeLineBreaks determines if Windows' CRLF (\r\n) and Mac Classic CR
	// (\r) line breaks are replaced with Unix's LF (\n) line breaks. This
	// ensure marshaling assertions works cross platform.
	NormalizeLineBreaks bool
}

func NewMarshalAsserter(
	golden *Golden,
	format string,
	newEncoderFunc NewEncoderFunc,
	newDecoderFunc NewDecoderFunc,
	normalizeLineBreaks bool,
) *MarshalAsserter {
	if golden == nil {
		golden = globalGolden
	}

	goldName := "marshaled_" + strings.ToLower(sanitizeFilename(format))

	return &MarshalAsserter{
		Golden:              golden,
		Format:              format,
		GoldName:            goldName,
		NewEncoderFunc:      newEncoderFunc,
		NewDecoderFunc:      newDecoderFunc,
		NormalizeLineBreaks: normalizeLineBreaks,
	}
}

// Marshaling asserts that the given "v" value marshals via the provided encoder
// to an expected value fetched from a golden file on disk, and then verifies
// that the marshaled result produces a value that is equal to "v" when
// unmarshaled.
//
// Used for objects that do NOT change when they are marshaled and unmarshaled.
func (s *MarshalAsserter) Marshaling(t *testing.T, v interface{}) {
	t.Helper()

	s.MarshalingP(t, v, v)
}

// MarshalingP asserts that the given "v" value marshals via the provided
// encoder to an expected value fetched from a golden file on disk, and then
// verifies that the marshaled result produces a value that is equal to "want"
// when unmarshaled.
//
// Used for objects that change when they are marshaled and unmarshaled.
func (s *MarshalAsserter) MarshalingP(
	t *testing.T,
	v interface{},
	want interface{},
) {
	t.Helper()

	if reflect.ValueOf(want).Kind() != reflect.Ptr {
		require.FailNowf(t,
			"only pointer types can be asserted",
			"%T is not a pointer type", want,
		)
	}

	var buf bytes.Buffer
	err := s.NewEncoderFunc(&buf).Encode(v)
	require.NoErrorf(t, err, "failed to %s marshal %T: %+v", s.Format, v, v)

	marshaled := buf.Bytes()
	if s.NormalizeLineBreaks {
		marshaled = normalizeLineBreaks(marshaled)
	}

	if s.Golden.Update() {
		s.Golden.SetP(t, s.GoldName, marshaled)
	}

	gold := s.Golden.GetP(t, s.GoldName)
	if s.NormalizeLineBreaks {
		gold = normalizeLineBreaks(gold)
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
	err = s.NewDecoderFunc(bytes.NewBuffer(gold)).Decode(got)
	require.NoErrorf(t, err,
		"failed to %s unmarshal %T from %s",
		s.Format, got, s.Golden.FileP(t, s.GoldName),
	)
	assert.Equal(t, want, got,
		"unmarshaling from golden file does not match expected object",
	)
}
