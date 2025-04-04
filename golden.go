// Package golden is yet another package for working with *.golden test files,
// with a focus on simplicity through it's default behavior.
//
// Golden file names are based on the name of the test function and any sub-test
// names by calling t.Name(). File names are sanitized to ensure they're
// compatible with Linux, macOS and Windows systems regardless of what crazy
// characters might be in a sub-test's name.
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
// The above example will read/write to:
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
// ensuring each sub-test gets its own golden file. For example:
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
// The "P" suffixed methods, GetP(), SetP(), and FileP(), all take a name
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
//	func TestExampleMyStructTabularP(t *testing.T) {
//		tests := []struct {
//			name string
//			obj  *MyStruct
//		}{
//			{name: "empty struct", obj: &MyStruct{}},
//			{name: "full struct", obj: &MyStruct{Foo: "Bar"}},
//		}
//		for _, tt := range tests {
//			t.Run(tt.name, func(t *testing.T) {
//				gotJSON, _ := json.Marshal(tt.obj)
//				gotXML, _ := xml.Marshal(tt.obj)
//
//				wantJSON := golden.DoP(t, "json", gotJSON)
//				wantXML := golden.DoP(t, "xml", gotXML)
//
//				assert.Equal(t, wantJSON, gotJSON)
//				assert.Equal(t, wantXML, gotXML)
//			})
//		}
//	}
//
// The above example will read/write to:
//
//	testdata/TestExampleMyStructTabularP/empty_struct/json.golden
//	testdata/TestExampleMyStructTabularP/empty_struct/xml.golden
//	testdata/TestExampleMyStructTabularP/full_struct/json.golden
//	testdata/TestExampleMyStructTabularP/full_struct/xml.golden
package golden

import (
	"os"
	"path/filepath"
	"strings"
)

var (
	// Default is the default *Golden instance. All package-level functions use
	// the Default instance.
	Default = New()

	// DefaultDirMode is the default DirMode value used by New().
	DefaultDirMode = os.FileMode(0o755)

	// DefaultFileMode is the default FileMode value used by New().
	DefaultFileMode = os.FileMode(0o644)

	// DefaultSuffix is the default Suffix value used by New().
	DefaultSuffix = ".golden"

	// DefaultDirname is the default Dirname value used by New().
	DefaultDirname = "testdata"

	// DefaultUpdateFunc is the default UpdateFunc value used by New().
	DefaultUpdateFunc = EnvUpdateFunc
)

// Do is a convenience function for calling Update(), Set(), and Get() in a
// single call. If Update() returns true, data will be written to the golden
// file using Set(), before reading it back with Get().
func Do(t TestingT, data []byte) []byte {
	t.Helper()

	return Default.Do(t, data)
}

// File returns the filename of the golden file for the given *testing.T
// instance as determined by t.Name().
func File(t TestingT) string {
	t.Helper()

	return Default.File(t)
}

// Get returns the content of the golden file for the given *testing.T instance
// as determined by t.Name(). If no golden file can be found/read, it will fail
// the test by calling t.Fatal().
func Get(t TestingT) []byte {
	t.Helper()

	return Default.Get(t)
}

// Set writes given data to the golden file for the given *testing.T instance as
// determined by t.Name(). If writing fails it will fail the test by calling
// t.Fatal() with error details.
func Set(t TestingT, data []byte) {
	t.Helper()

	Default.Set(t, data)
}

// DoP is a convenience function for calling Update(), SetP(), and GetP() in a
// single call. If Update() returns true, data will be written to the golden
// file using SetP(), before reading it back with GetP().
func DoP(t TestingT, name string, data []byte) []byte {
	t.Helper()

	return Default.DoP(t, name, data)
}

// FileP returns the filename of the specifically named golden file for the
// given *testing.T instance as determined by t.Name().
func FileP(t TestingT, name string) string {
	t.Helper()

	return Default.FileP(t, name)
}

// GetP returns the content of the specifically named golden file belonging
// to the given *testing.T instance as determined by t.Name(). If no golden file
// can be found/read, it will fail the test with t.Fatal().
//
// This is very similar to Get(), but it allows multiple different golden files
// to be used within the same one *testing.T instance.
func GetP(t TestingT, name string) []byte {
	t.Helper()

	return Default.GetP(t, name)
}

// SetP writes given data of the specifically named golden file belonging to
// the given *testing.T instance as determined by t.Name(). If writing fails it
// will fail the test with t.Fatal() detailing the error.
//
// This is very similar to Set(), but it allows multiple different golden files
// to be used within the same one *testing.T instance.
func SetP(t TestingT, name string, data []byte) {
	t.Helper()

	Default.SetP(t, name, data)
}

// Update returns true when golden is set to update golden files. Should be used
// to determine if golden.Set() or golden.SetP() should be called or not.
//
// Default behavior uses EnvUpdateFunc() to check if the "GOLDEN_UPDATE"
// environment variable is set to a truthy value. To customize create a custom
// *Golden instance with New() and set a new UpdateFunc value.
func Update() bool {
	return Default.Update()
}

// Golden handles all interactions with golden files. The top-level package
// functions all just proxy through to a default global *Golden instance.
type Golden struct {
	// DirMode determines the file system permissions of any folders created to
	// hold golden files.
	DirMode os.FileMode

	// FileMode determines the file system permissions of any created or updated
	// golden files written to disk.
	FileMode os.FileMode

	// Suffix determines the filename suffix for all golden files. Typically
	// this would be ".golden".
	Suffix string

	// Dirname is the name of the top-level directory at the root of the package
	// which holds all golden files. Typically this should "testdata", but can
	// be changed here if needed.
	Dirname string

	// UpdateFunc is used to determine if golden files should be updated or
	// not. Its boolean return value is returned by Update().
	UpdateFunc UpdateFunc
}

// New returns a new *Golden instance with default values correctly populated.
// It accepts zero or more Option functions that can modify the default values.
func New(opts ...Option) *Golden {
	g := &Golden{
		DirMode:    DefaultDirMode,
		FileMode:   DefaultFileMode,
		Suffix:     DefaultSuffix,
		Dirname:    DefaultDirname,
		UpdateFunc: DefaultUpdateFunc,
	}

	for _, opt := range opts {
		opt(g)
	}

	return g
}

// Do is a convenience function for calling Update(), Set(), and Get() in a
// single call. If Update() returns true, data will be written to the golden
// file using Set(), before reading it back with Get().
func (s *Golden) Do(t TestingT, data []byte) []byte {
	t.Helper()

	if s.Update() {
		s.Set(t, data)
	}

	return s.Get(t)
}

// File returns the filename of the golden file for the given *testing.T
// instance as determined by t.Name().
func (s *Golden) File(t TestingT) string {
	t.Helper()

	return s.file(t, "")
}

// Get returns the content of the golden file for the given *testing.T instance
// as determined by t.Name(). If no golden file can be found/read, it will fail
// the test by calling t.Fatal().
func (s *Golden) Get(t TestingT) []byte {
	t.Helper()

	return s.get(t, "")
}

// Set writes given data to the golden file for the given *testing.T instance as
// determined by t.Name(). If writing fails it will fail the test by calling
// t.Fatal() with error details.
func (s *Golden) Set(t TestingT, data []byte) {
	t.Helper()

	s.set(t, "", data)
}

// DoP is a convenience function for calling Update(), SetP(), and GetP() in a
// single call. If Update() returns true, data will be written to the golden
// file using SetP(), before reading it back with GetP().
func (s *Golden) DoP(t TestingT, name string, data []byte) []byte {
	t.Helper()

	if name == "" {
		t.Fatalf("golden: name cannot be empty")
	}

	if s.Update() {
		s.SetP(t, name, data)
	}

	return s.GetP(t, name)
}

// FileP returns the filename of the specifically named golden file for the
// given *testing.T instance as determined by t.Name().
func (s *Golden) FileP(t TestingT, name string) string {
	t.Helper()

	if name == "" {
		t.Fatalf("golden: name cannot be empty")
	}

	return s.file(t, name)
}

// GetP returns the content of the specifically named golden file belonging
// to the given *testing.T instance as determined by t.Name(). If no golden file
// can be found/read, it will fail the test with t.Fatal().
//
// This is very similar to Get(), but it allows multiple different golden files
// to be used within the same one *testing.T instance.
func (s *Golden) GetP(t TestingT, name string) []byte {
	t.Helper()

	if name == "" {
		t.Fatalf("golden: name cannot be empty")
	}

	return s.get(t, name)
}

// SetP writes given data of the specifically named golden file belonging to
// the given *testing.T instance as determined by t.Name(). If writing fails it
// will fail the test with t.Fatal() detailing the error.
//
// This is very similar to Set(), but it allows multiple different golden files
// to be used within the same one *testing.T instance.
func (s *Golden) SetP(t TestingT, name string, data []byte) {
	t.Helper()

	if name == "" {
		t.Fatalf("golden: name cannot be empty")
	}

	s.set(t, name, data)
}

// Update returns true when golden is set to update golden files. Should be used
// to determine if golden.Set() or golden.SetP() should be called or not.
//
// Default behavior uses EnvUpdateFunc() to check if the "GOLDEN_UPDATE"
// environment variable is set to a truthy value. To customize set a new
// UpdateFunc value on *Golden.
func (s *Golden) Update() bool {
	return s.UpdateFunc()
}

func (s *Golden) file(t TestingT, name string) string {
	if t.Name() == "" {
		t.Fatalf("golden: could not determine filename")
	}

	base := []string{s.Dirname, filepath.FromSlash(t.Name())}
	if name != "" {
		base = append(base, name)
	}

	f := filepath.Clean(filepath.Join(base...) + s.Suffix)

	dirty := strings.Split(f, string(os.PathSeparator))
	clean := make([]string, 0, len(dirty))
	for _, s := range dirty {
		clean = append(clean, sanitizeFilename(s))
	}

	return strings.Join(clean, string(os.PathSeparator))
}

func (s *Golden) get(t TestingT, name string) []byte {
	f := s.file(t, name)

	b, err := os.ReadFile(f)
	if err != nil {
		t.Fatalf("golden: failed reading %s: %s", f, err.Error())
	}

	return b
}

func (s *Golden) set(t TestingT, name string, data []byte) {
	f := s.file(t, name)
	dir := filepath.Dir(f)

	t.Logf("golden: writing .golden file: %s", f)

	err := os.MkdirAll(dir, s.DirMode)
	if err != nil {
		t.Fatalf("golden: failed to create directory: %s", err.Error())
	}

	err = os.WriteFile(f, data, s.FileMode)
	if err != nil {
		t.Fatalf("golden: filed to write file: %s", err.Error())
	}
}
