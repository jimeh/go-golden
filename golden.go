// Package golden is yet another package for working with *.golden test files.
package golden

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

const (
	DefaultDirMode  = 0o755
	DefaultFileMode = 0o644
	DefaultSuffix   = ".golden"
	DefaultDirname  = "testdata"
)

type Golden struct {
	DirMode      os.FileMode
	FileMode     os.FileMode
	Suffix       string
	Dirname      string
	UpdatingFunc UpdatingFunc
}

func New() *Golden {
	return &Golden{
		DirMode:      DefaultDirMode,
		FileMode:     DefaultFileMode,
		Suffix:       DefaultSuffix,
		Dirname:      DefaultDirname,
		UpdatingFunc: EnvVarUpdatingFunc,
	}
}

// Updating returns true when the function assigned to UpdatingFunc returns
// true.
func (s *Golden) Updating() bool {
	return s.UpdatingFunc()
}

// Get returns the content of the default golden file for the given *testing.T
// instance as determined by t.Name(). If no golden file can be found/read, it
// will fail the test with t.Fatal().
func (s *Golden) Get(t *testing.T) []byte {
	return s.GetNamed(t, "")
}

// Set writes given data of the default golden file for the given *testing.T
// instance as determined by t.Name(). If writing fails it will fail the test
// with t.Fatal() detailing the error.
func (s *Golden) Set(t *testing.T, data []byte) {
	s.SetNamed(t, "", data)
}

func (s *Golden) File(t *testing.T) string {
	return s.NamedFile(t, "")
}

func (s *Golden) GetNamed(t *testing.T, name string) []byte {
	if t == nil {
		return nil
	}

	f := s.NamedFile(t, name)

	b, err := ioutil.ReadFile(f)
	if err != nil {
		t.Fatalf("golden: failed reading %s: %s", f, err.Error())
	}

	return b
}

func (s *Golden) SetNamed(t *testing.T, name string, data []byte) {
	if t == nil {
		return
	}

	f := s.NamedFile(t, name)
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

func (s *Golden) NamedFile(t *testing.T, name string) string {
	if t == nil || t.Name() == "" {
		t.Fatalf("golden: could not determine filename for: %+v", t)
		return ""
	}

	base := []string{s.Dirname, filepath.FromSlash(t.Name())}
	if name != "" {
		base = append(base, name)
	}

	return filepath.Clean(filepath.Join(base...) + s.Suffix)
}
