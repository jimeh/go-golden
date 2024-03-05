// Package golden is yet another package for working with *.golden test files,
// with a focus on simplicity through it's default behavior.
//
// Golden file names are based on the name of the test function and any subtest
// names by calling t.Name(). File names are sanitized to ensure they are
// compatible with Linux, macOS and Windows systems regardless of what
// characters might be in a subtest's name.
//
// # Usage
//
// Typical usage should look something like this:
//
//	func TestExampleMyStruct(t *testing.T) {
//		got, err := json.Marshal(&MyStruct{Foo: "Bar"})
//		require.NoError(t, err)
//
//		want := golden.Do(t, got)
//
//		assert.Equal(t, want, got)
//	}
//
// The above example will attempt to read/write to:
//
//	testdata/TestExampleMyStruct.golden
//
// The call to golden.Do() is equivalent to:
//
//	if golden.Update() {
//		golden.Set(t, got)
//	}
//	want := golden.Get(t)
//
// To update the golden file (have golden.Update() return true), simply set the
// GOLDEN_UPDATE environment variable to one of "1", "y", "t", "yes", "on", or
// "true" when running tests.
//
// # Sub-Tests
//
// As the golden filename is based on t.Name(), it works with sub-tests too,
// ensuring each sub-test gets it's own golden file. For example:
//
//	func TestExampleMyStructTabular(t *testing.T) {
//		tests := []struct {
//			name string
//			obj  *MyStruct
//		}{
//			{name: "empty struct", obj: &MyStruct{}},
//			{name: "full struct", obj: &MyStruct{Foo: "Bar"}},
//		}
//		for _, tt := range tests {
//			t.Run(tt.name, func(t *testing.T) {
//				got, err := json.Marshal(tt.obj)
//				require.NoError(t, err)
//
//				want := golden.Do(t, got)
//
//				assert.Equal(t, want, got)
//			})
//		}
//	}
//
// The above example will read/write to:
//
//	testdata/TestExampleMyStructTabular/empty_struct.golden
//	testdata/TestExampleMyStructTabular/full_struct.golden
//
// # Multiple Golden Files in a Single Test
//
// The "P" suffixed methods, GetP(), SetP(), DoP(), and FileP(), all take a name
// argument which allows using specific golden files within a given *testing.T
// instance.
//
//	func TestExampleMyStructP(t *testing.T) {
//		gotJSON, _ := json.Marshal(&MyStruct{Foo: "Bar"})
//		gotXML, _ := xml.Marshal(&MyStruct{Foo: "Bar"})
//
//		wantJSON := golden.DoP(t, "json", gotJSON)
//		wantXML := golden.DoP(t, "xml", gotXML)
//
//		assert.Equal(t, wantJSON, gotJSON)
//		assert.Equal(t, wantXML, gotXML)
//	}
//
// The above example will read/write to:
//
//	testdata/TestExampleMyStructP/json.golden
//	testdata/TestExampleMyStructP/xml.golden
//
// This works with tabular tests too of course:
//
//	 func TestExampleMyStructTabularP(t *testing.T) {
//		 tests := []struct {
//			 name string
//			 obj  *MyStruct
//		 }{
//			 {name: "empty struct", obj: &MyStruct{}},
//			 {name: "full struct", obj: &MyStruct{Foo: "Bar"}},
//		 }
//		 for _, tt := range tests {
//			 t.Run(tt.name, func(t *testing.T) {
//				 gotJSON, _ := json.Marshal(tt.obj)
//				 gotXML, _ := xml.Marshal(tt.obj)
//
//					wantJSON := golden.DoP(t, "json", gotJSON)
//					wantXML := golden.DoP(t, "xml", gotXML)
//
//				 assert.Equal(t, wantJSON, gotJSON)
//				 assert.Equal(t, wantXML, gotXML)
//			 })
//		 }
//	 }
//
// The above example will read/write to:
//
//	testdata/TestExampleMyStructTabularP/empty_struct/json.golden
//	testdata/TestExampleMyStructTabularP/empty_struct/xml.golden
//	testdata/TestExampleMyStructTabularP/full_struct/json.golden
//	testdata/TestExampleMyStructTabularP/full_struct/xml.golden
package golden

import "os"

var (
	// DefaultGolden is the default Golden instance used by all top-level
	// package functions.
	DefaultGolden = New()

	// DefaultDirMode is the default file system permissions used for any
	// created directories to hold golden files.
	DefaultDirMode = os.FileMode(0o755)

	// DefaultFileMode is the default file system permissions used for any
	// created or updated golden files written to disk.
	DefaultFileMode = os.FileMode(0o644)

	// DefaultSuffix is the default filename suffix used for all golden files.
	DefaultSuffix = ".golden"

	// DefaultDirname is the default name of the top-level directory used to
	// hold golden files.
	DefaultDirname = "testdata"

	// DefaultUpdateFunc is the default function used to determine if golden
	// files should be updated or not. It is called by Update().
	DefaultUpdateFunc = EnvUpdateFunc

	// DefaultLogOnWrite is the default value for logOnWrite on all Golden
	// instances.
	DefaultLogOnWrite = true
)

// Golden handles all interactions with golden files. The top-level package
// functions proxy through to a default global Golden instance.
type Golden interface {
	// File returns the filename of the golden file for the given testing.TB
	// instance as determined by t.Name().
	File(t TestingT) string

	// Get returns the content of the golden file for the given TestingT
	// instance as determined by t.Name(). If no golden file can be found/read,
	// it will fail the test by calling t.Fatal().
	Get(t TestingT) []byte

	// Set writes given data to the golden file for the given TestingT
	// instance as determined by t.Name(). If writing fails it will fail the
	// test by calling t.Fatal() with error details.
	Set(t TestingT, data []byte)

	// Do is a convenience function for calling Update(), Set(), and Get() in a
	// single call. If Update() returns true, data will be written to the golden
	// file using Set(), before reading it back with Get().
	Do(t TestingT, data []byte) []byte

	// FileP returns the filename of the specifically named golden file for the
	// given TestingT instance as determined by t.Name().
	FileP(t TestingT, name string) string

	// GetP returns the content of the specifically named golden file belonging
	// to the given TestingT instance as determined by t.Name(). If no golden
	// file can be found/read, it will fail the test with t.Fatal().
	//
	// This is very similar to Get(), but it allows multiple different golden
	// files to be used within the same one TestingT instance.
	GetP(t TestingT, name string) []byte

	// SetP writes given data of the specifically named golden file belonging to
	// the given TestingT instance as determined by t.Name(). If writing fails
	// it will fail the test with t.Fatal() detailing the error.
	//
	// This is very similar to Set(), but it allows multiple different golden
	// files to be used within the same one TestingT instance.
	SetP(t TestingT, name string, data []byte)

	// DoP is a convenience function for calling Update(), SetP(), and GetP() in
	// a single call. If Update() returns true, data will be written to the
	// golden file using SetP(), before reading it back with GetP().
	DoP(t TestingT, name string, data []byte) []byte

	// Update returns true when golden is set to update golden files. Should be
	// used to determine if golden.Set() or golden.SetP() should be called or
	// not.
	//
	// Default behavior uses EnvUpdateFunc() to check if the "GOLDEN_UPDATE"
	// environment variable is set to a truthy value. To customize set a new
	// UpdateFunc value on *Golden.
	Update() bool
}

// New returns a new Golden instance. Used to create custom Golden instances.
// See the the various Option functions for details of what can be customized.
func New(options ...Option) Golden {
	g := &gold{
		dirMode:    DefaultDirMode,
		fileMode:   DefaultFileMode,
		suffix:     DefaultSuffix,
		dirname:    DefaultDirname,
		updateFunc: DefaultUpdateFunc,
		fs:         DefaultFS,
		logOnWrite: DefaultLogOnWrite,
	}

	for _, opt := range options {
		opt.apply(g)
	}

	return g
}

// Do is a convenience function for calling Update(), Set(), and Get() in a
// single call. If Update() returns true, data will be written to the golden
// file using Set(), before reading it back with Get().
func Do(t TestingT, data []byte) []byte {
	t.Helper()

	return DefaultGolden.Do(t, data)
}

// DoP is a convenience function for calling Update(), SetP(), and GetP() in a
// single call. If Update() returns true, data will be written to the golden
// file using SetP(), before reading it back with GetP().
func DoP(t TestingT, name string, data []byte) []byte {
	t.Helper()

	return DefaultGolden.DoP(t, name, data)
}

// File returns the filename of the golden file for the given *testing.T
// instance as determined by t.Name().
func File(t TestingT) string {
	t.Helper()

	return DefaultGolden.File(t)
}

// FileP returns the filename of the specifically named golden file for the
// given *testing.T instance as determined by t.Name().
func FileP(t TestingT, name string) string {
	t.Helper()

	return DefaultGolden.FileP(t, name)
}

// Get returns the content of the golden file for the given *testing.T instance
// as determined by t.Name(). If no golden file can be found/read, it will fail
// the test by calling t.Fatal().
func Get(t TestingT) []byte {
	t.Helper()

	return DefaultGolden.Get(t)
}

// GetP returns the content of the specifically named golden file belonging
// to the given *testing.T instance as determined by t.Name(). If no golden file
// can be found/read, it will fail the test with t.Fatal().
//
// This is very similar to Get(), but it allows multiple different golden files
// to be used within the same one *testing.T instance.
func GetP(t TestingT, name string) []byte {
	t.Helper()

	return DefaultGolden.GetP(t, name)
}

// Set writes given data to the golden file for the given *testing.T instance as
// determined by t.Name(). If writing fails it will fail the test by calling
// t.Fatal() with error details.
func Set(t TestingT, data []byte) {
	t.Helper()

	DefaultGolden.Set(t, data)
}

// SetP writes given data of the specifically named golden file belonging to
// the given *testing.T instance as determined by t.Name(). If writing fails it
// will fail the test with t.Fatal() detailing the error.
//
// This is very similar to Set(), but it allows multiple different golden files
// to be used within the same one *testing.T instance.
func SetP(t TestingT, name string, data []byte) {
	t.Helper()

	DefaultGolden.SetP(t, name, data)
}

// Update returns true when golden is set to update golden files. Should be used
// to determine if golden.Set() or golden.SetP() should be called or not.
//
// Default behavior uses EnvUpdateFunc() to check if the "GOLDEN_UPDATE"
// environment variable is set to a truthy value. To customize create a custom
// *Golden instance with New() and set a new UpdateFunc value.
func Update() bool {
	return DefaultGolden.Update()
}
