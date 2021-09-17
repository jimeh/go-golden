// Package golden is yet another package for working with *.golden test files,
// with a focus on simplicity through it's default behavior.
//
// Golden file names are based on the name of the test function and any subtest
// names by calling t.Name(). File names are sanitized to ensure they're
// compatible with Linux, macOS and Windows systems regardless of what crazy
// characters might be in a subtest's name.
//
// Usage
//
// Typical usage should look something like this:
//
//  func TestExampleMyStruct(t *testing.T) {
//      got, err := json.Marshal(&MyStruct{Foo: "Bar"})
//      require.NoError(t, err)
//
//      if golden.Update() {
//          golden.Set(t, got)
//      }
//      want := golden.Get(t)
//
//      assert.Equal(t, want, got)
//  }
//
// The above example will read/write to:
//
//  testdata/TestExampleMyStruct.golden
//
// To update the golden file (have golden.Update() return true), simply set the
// GOLDEN_UPDATE environment variable to one of "1", "y", "t", "yes", "on", or
// "true" when running tests.
//
// Sub-Tests
//
// As the golden filename is based on t.Name(), it works with sub-tests too,
// ensuring each sub-test gets it's own golden file. For example:
//
//  func TestExampleMyStructTabular(t *testing.T) {
//      tests := []struct {
//          name string
//          obj  *MyStruct
//      }{
//          {name: "empty struct", obj: &MyStruct{}},
//          {name: "full struct", obj: &MyStruct{Foo: "Bar"}},
//      }
//      for _, tt := range tests {
//          t.Run(tt.name, func(t *testing.T) {
//              got, err := json.Marshal(tt.obj)
//              require.NoError(t, err)
//
//              if golden.Update() {
//                  golden.Set(t, got)
//              }
//              want := golden.Get(t)
//
//              assert.Equal(t, want, got)
//          })
//      }
//  }
//
// The above example will read/write to:
//
//  testdata/TestExampleMyStructTabular/empty_struct.golden
//  testdata/TestExampleMyStructTabular/full_struct.golden
//
// Multiple Golden Files in a Single Test
//
// The "P" suffixed methods, GetP(), SetP(), and FileP(), all take a name
// argument which allows using specific golden files within a given *testing.T
// instance.
//
//  func TestExampleMyStructP(t *testing.T) {
//      gotJSON, _ := json.Marshal(&MyStruct{Foo: "Bar"})
//      gotXML, _ := xml.Marshal(&MyStruct{Foo: "Bar"})
//
//      if golden.Update() {
//          golden.SetP(t, "json", gotJSON)
//          golden.SetP(t, "xml", gotXML)
//      }
//
//      assert.Equal(t, golden.GetP(t, "json"), gotJSON)
//      assert.Equal(t, golden.GetP(t, "xml"), gotXML)
//  }
//
// The above example will read/write to:
//
//  testdata/TestExampleMyStructP/json.golden
//  testdata/TestExampleMyStructP/xml.golden
//
// This works with tabular tests too of course:
//
//  func TestExampleMyStructTabularP(t *testing.T) {
//      tests := []struct {
//          name string
//          obj  *MyStruct
//      }{
//          {name: "empty struct", obj: &MyStruct{}},
//          {name: "full struct", obj: &MyStruct{Foo: "Bar"}},
//      }
//      for _, tt := range tests {
//          t.Run(tt.name, func(t *testing.T) {
//              gotJSON, _ := json.Marshal(tt.obj)
//              gotXML, _ := xml.Marshal(tt.obj)
//
//              if golden.Update() {
//                  golden.SetP(t, "json", gotJSON)
//                  golden.SetP(t, "xml", gotXML)
//              }
//
//              assert.Equal(t, golden.GetP(t, "json"), gotJSON)
//              assert.Equal(t, golden.GetP(t, "xml"), gotXML)
//          })
//      }
//  }
//
// The above example will read/write to:
//
//  testdata/TestExampleMyStructTabularP/empty_struct/json.golden
//  testdata/TestExampleMyStructTabularP/empty_struct/xml.golden
//  testdata/TestExampleMyStructTabularP/full_struct/json.golden
//  testdata/TestExampleMyStructTabularP/full_struct/xml.golden
//
package golden

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const (
	DefaultDirMode  = 0o755
	DefaultFileMode = 0o644
	DefaultSuffix   = ".golden"
	DefaultDirname  = "testdata"
)

var DefaultUpdateFunc = EnvUpdateFunc

var global = New()

// File returns the filename of the golden file for the given *testing.T
// instance as determined by t.Name().
func File(t *testing.T) string {
	return global.File(t)
}

// Get returns the content of the golden file for the given *testing.T instance
// as determined by t.Name(). If no golden file can be found/read, it will fail
// the test by calling t.Fatal().
func Get(t *testing.T) []byte {
	return global.Get(t)
}

// Set writes given data to the golden file for the given *testing.T instance as
// determined by t.Name(). If writing fails it will fail the test by calling
// t.Fatal() with error details.
func Set(t *testing.T, data []byte) {
	global.Set(t, data)
}

// FileP returns the filename of the specifically named golden file for the
// given *testing.T instance as determined by t.Name().
func FileP(t *testing.T, name string) string {
	return global.FileP(t, name)
}

// GetP returns the content of the specifically named golden file belonging
// to the given *testing.T instance as determined by t.Name(). If no golden file
// can be found/read, it will fail the test with t.Fatal().
//
// This is very similar to Get(), but it allows multiple different golden files
// to be used within the same one *testing.T instance.
func GetP(t *testing.T, name string) []byte {
	return global.GetP(t, name)
}

// SetP writes given data of the specifically named golden file belonging to
// the given *testing.T instance as determined by t.Name(). If writing fails it
// will fail the test with t.Fatal() detailing the error.
//
// This is very similar to Set(), but it allows multiple different golden files
// to be used within the same one *testing.T instance.
func SetP(t *testing.T, name string, data []byte) {
	global.SetP(t, name, data)
}

// Update returns true when golden is set to update golden files. Used to
// determine if golden.Set() or golden.Write() should be called or not.
//
// Default behavior uses EnvUpdateFunc() to check if the "GOLDEN_UPDATE"
// environment variable is set to a truthy value. To customize create a custom
// *Golden instance with New() and set a new UpdateFunc value.
func Update() bool {
	return global.Update()
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
	// this should be ".golden", but can be changed here if needed.
	Suffix string

	// Dirname is the name of the top-level directory at the root of the package
	// which holds all golden files. Typically this should "testdata", but can
	// be changed here if needed.
	Dirname string

	// UpdateFunc is used to determine if golden files should be updated or
	// not. Its boolean return value is returned by Update().
	UpdateFunc UpdateFunc
}

// New returns a new *Golden instance with default values correctly
// populated. This is ideally how you should create a custom *Golden, and then
// modify the relevant fields as you see fit.
func New() *Golden {
	return &Golden{
		DirMode:    DefaultDirMode,
		FileMode:   DefaultFileMode,
		Suffix:     DefaultSuffix,
		Dirname:    DefaultDirname,
		UpdateFunc: DefaultUpdateFunc,
	}
}

// File returns the filename of the golden file for the given *testing.T
// instance as determined by t.Name().
func (s *Golden) File(t *testing.T) string {
	return s.file(t, "")
}

// Get returns the content of the golden file for the given *testing.T instance
// as determined by t.Name(). If no golden file can be found/read, it will fail
// the test by calling t.Fatal().
func (s *Golden) Get(t *testing.T) []byte {
	return s.get(t, "")
}

// Set writes given data to the golden file for the given *testing.T instance as
// determined by t.Name(). If writing fails it will fail the test by calling
// t.Fatal() with error details.
func (s *Golden) Set(t *testing.T, data []byte) {
	s.set(t, "", data)
}

// FileP returns the filename of the specifically named golden file for the
// given *testing.T instance as determined by t.Name().
func (s *Golden) FileP(t *testing.T, name string) string {
	if name == "" {
		if t != nil {
			t.Fatal("golden: name cannot be empty")
		}
		return ""
	}

	return s.file(t, name)
}

// GetP returns the content of the specifically named golden file belonging
// to the given *testing.T instance as determined by t.Name(). If no golden file
// can be found/read, it will fail the test with t.Fatal().
//
// This is very similar to Get(), but it allows multiple different golden files
// to be used within the same one *testing.T instance.
func (s *Golden) GetP(t *testing.T, name string) []byte {
	if name == "" {
		t.Fatal("golden: name cannot be empty")
		return nil
	}

	return s.get(t, name)
}

// SetP writes given data of the specifically named golden file belonging to
// the given *testing.T instance as determined by t.Name(). If writing fails it
// will fail the test with t.Fatal() detailing the error.
//
// This is very similar to Set(), but it allows multiple different golden files
// to be used within the same one *testing.T instance.
func (s *Golden) SetP(t *testing.T, name string, data []byte) {
	if name == "" {
		t.Fatal("golden: name cannot be empty")
	}

	s.set(t, name, data)
}

func (s *Golden) file(t *testing.T, name string) string {
	if t.Name() == "" {
		t.Fatalf("golden: could not determine filename for: %+v", t)
		return ""
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

func (s *Golden) get(t *testing.T, name string) []byte {
	f := s.file(t, name)

	b, err := ioutil.ReadFile(f)
	if err != nil {
		t.Fatalf("golden: failed reading %s: %s", f, err.Error())
	}

	return b
}

func (s *Golden) set(t *testing.T, name string, data []byte) {
	f := s.file(t, name)
	dir := filepath.Dir(f)

	t.Logf("golden: writing .golden file: %s", f)

	err := os.MkdirAll(dir, s.DirMode)
	if err != nil {
		t.Fatalf("golden: failed to create directory: %s", err.Error())
		return
	}

	err = ioutil.WriteFile(f, data, s.FileMode)
	if err != nil {
		t.Fatalf("golden: filed to write file: %s", err.Error())
	}
}

// Update returns true when golden is set to update golden files. Used to
// determine if golden.Set() or golden.Write() should be called or not.
//
// Default behavior uses EnvUpdateFunc() to check if the "GOLDEN_UPDATE"
// environment variable is set to a truthy value. To customize set a new
// UpdateFunc value on *Golden.
func (s *Golden) Update() bool {
	return s.UpdateFunc()
}
