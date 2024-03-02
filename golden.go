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

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jimeh/go-golden/sanitize"
	"github.com/spf13/afero"
)

//go:generate mockgen -source=golden.go -destination=golden_mock_test.go -package golden -self_package github.com/jimeh/go-golden

// TestingT is a interface describing a sub-set of methods of *testing.T which
// golden uses.
type TestingT interface {
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Helper()
	Log(args ...interface{})
	Logf(format string, args ...interface{})
	Name() string
}

// Default is the default Golden instance used by all top-level package
// functions.
var Default = New()

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
func Set(t *testing.T, data []byte) {
	t.Helper()

	Default.Set(t, data)
}

// Do is a convenience function for calling Update(), Set(), and Get() in a
// single call. If Update() returns true, data will be written to the golden
// file using Set(), before reading it back with Get().
func Do(t TestingT, data []byte) []byte {
	t.Helper()

	return Default.Do(t, data)
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
func GetP(t *testing.T, name string) []byte {
	t.Helper()

	return Default.GetP(t, name)
}

// SetP writes given data of the specifically named golden file belonging to
// the given *testing.T instance as determined by t.Name(). If writing fails it
// will fail the test with t.Fatal() detailing the error.
//
// This is very similar to Set(), but it allows multiple different golden files
// to be used within the same one *testing.T instance.
func SetP(t *testing.T, name string, data []byte) {
	t.Helper()

	Default.SetP(t, name, data)
}

// DoP is a convenience function for calling Update(), SetP(), and GetP() in a
// single call. If Update() returns true, data will be written to the golden
// file using SetP(), before reading it back with GetP().
func DoP(t TestingT, name string, data []byte) []byte {
	t.Helper()

	return Default.DoP(t, name, data)
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
	g := &golden{
		dirMode:    0o755,
		fileMode:   0o644,
		suffix:     ".golden",
		dirname:    "testdata",
		updateFunc: EnvUpdateFunc,
		fs:         afero.NewOsFs(),
		logOnWrite: true,
	}

	for _, opt := range options {
		opt.apply(g)
	}

	return g
}

type Option interface {
	apply(*golden)
}

type optionFunc func(*golden)

func (fn optionFunc) apply(g *golden) {
	fn(g)
}

// WithDirMode sets the file system permissions used for any folders created to
// hold golden files.
//
// When this option is not provided, the default value is 0o755.
func WithDirMode(mode os.FileMode) Option {
	return optionFunc(func(g *golden) {
		g.dirMode = mode
	})
}

// WithFileMode sets the file system permissions used for any created or updated
// golden files written to.
//
// When this option is not provided, the default value is 0o644.
func WithFileMode(mode os.FileMode) Option {
	return optionFunc(func(g *golden) {
		g.fileMode = mode
	})
}

// WithSuffix sets the filename suffix used for all golden files.
//
// When this option is not provided, the default value is ".golden".
func WithSuffix(suffix string) Option {
	return optionFunc(func(g *golden) {
		g.suffix = suffix
	})
}

// WithDirname sets the name of the top-level directory used to hold golden
// files.
//
// When this option is not provided, the default value is "testdata".
func WithDirname(name string) Option {
	return optionFunc(func(g *golden) {
		g.dirname = name
	})
}

// WithUpdateFunc sets the function used to determine if golden files should be
// updated or not. Essentially the provided UpdateFunc is called by Update().
//
// When this option is not provided, the default value is EnvUpdateFunc.
func WithUpdateFunc(fn UpdateFunc) Option {
	return optionFunc(func(g *golden) {
		g.updateFunc = fn
	})
}

// WithFs sets the afero.Fs instance which is used for all file system
// operations to read/write golden files.
//
// When this option is not provided, the default value is afero.NewOsFs().
func WithFs(fs afero.Fs) Option {
	return optionFunc(func(g *golden) {
		g.fs = fs
	})
}

// WithSilentWrites silences the "golden: writing [...]" log messages whenever
// set functions write a golden file to disk.
func WithSilentWrites() Option {
	return optionFunc(func(g *golden) {
		g.logOnWrite = false
	})
}

// golden is the underlying struct that implements the Golden interface.
type golden struct {
	// dirMode determines the file system permissions of any folders created to
	// hold golden files.
	dirMode os.FileMode

	// fileMode determines the file system permissions of any created or updated
	// golden files written to disk.
	fileMode os.FileMode

	// suffix determines the filename suffix for all golden files. Typically
	// this should be ".golden", but can be changed here if needed.
	suffix string

	// dirname is the name of the top-level directory at the root of the package
	// which holds all golden files. Typically this should be "testdata", but
	// can be changed here if needed.
	dirname string

	// updateFunc is used to determine if golden files should be updated or
	// not. Its boolean return value is returned by Update().
	updateFunc UpdateFunc

	// fs is used for all file system operations. This enables providing custom
	// afero.fs instances which can be useful for testing purposes.
	fs afero.Fs

	// logOnWrite determines if a message is logged with t.Logf when a golden
	// file is written to with either of the set methods.
	logOnWrite bool
}

// Ensure golden satisfies Golden interface.
var _ Golden = &golden{}

func (g *golden) File(t TestingT) string {
	t.Helper()

	return g.file(t, "")
}

func (g *golden) Get(t TestingT) []byte {
	t.Helper()

	return g.get(t, "")
}

func (g *golden) Set(t TestingT, data []byte) {
	t.Helper()

	g.set(t, "", data)
}

func (g *golden) Do(t TestingT, data []byte) []byte {
	t.Helper()

	if g.Update() {
		g.Set(t, data)
	}

	return g.Get(t)
}

func (g *golden) FileP(t TestingT, name string) string {
	t.Helper()

	if name == "" {
		t.Fatalf("golden: test name cannot be empty")
	}

	return g.file(t, name)
}

func (g *golden) GetP(t TestingT, name string) []byte {
	t.Helper()

	if name == "" {
		t.Fatal("golden: name cannot be empty")
	}

	return g.get(t, name)
}

func (g *golden) SetP(t TestingT, name string, data []byte) {
	t.Helper()

	if name == "" {
		t.Fatal("golden: name cannot be empty")
	}

	g.set(t, name, data)
}

func (g *golden) DoP(t TestingT, name string, data []byte) []byte {
	t.Helper()

	if g.Update() {
		g.SetP(t, name, data)
	}

	return g.GetP(t, name)
}

func (g *golden) file(t TestingT, name string) string {
	t.Helper()

	if t.Name() == "" {
		t.Fatalf(
			"golden: could not determine filename for given %T instance", t,
		)
	}

	base := []string{g.dirname, filepath.FromSlash(t.Name())}
	if name != "" {
		base = append(base, name)
	}

	f := filepath.Clean(filepath.Join(base...) + g.suffix)

	dirty := strings.Split(f, string(os.PathSeparator))
	clean := make([]string, 0, len(dirty))
	for _, s := range dirty {
		clean = append(clean, sanitize.Filename(s))
	}

	return strings.Join(clean, string(os.PathSeparator))
}

func (g *golden) get(t TestingT, name string) []byte {
	t.Helper()

	f := g.file(t, name)

	b, err := afero.ReadFile(g.fs, f)
	if err != nil {
		t.Fatalf("golden: %s", err.Error())
	}

	return b
}

func (g *golden) set(t TestingT, name string, data []byte) {
	t.Helper()

	f := g.file(t, name)
	dir := filepath.Dir(f)

	if g.logOnWrite {
		t.Logf("golden: writing golden file: %s", f)
	}

	err := g.fs.MkdirAll(dir, g.dirMode)
	if err != nil {
		t.Fatalf("golden: failed to create directory: %s", err.Error())
	}

	err = afero.WriteFile(g.fs, f, data, g.fileMode)
	if err != nil {
		t.Fatalf("golden: filed to write file: %s", err.Error())
	}
}

func (g *golden) Update() bool {
	return g.updateFunc()
}
