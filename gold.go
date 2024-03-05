package golden

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/jimeh/go-golden/sanitize"
)

// gold is the underlying struct that implements the Golden interface.
type gold struct {
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
	fs FS

	// logOnWrite determines if a message is logged with t.Logf when a golden
	// file is written to with either of the set methods.
	logOnWrite bool
}

// Ensure golden satisfies Golden interface.
var _ Golden = &gold{}

func (g *gold) Do(t TestingT, data []byte) []byte {
	t.Helper()

	if g.Update() {
		g.Set(t, data)
	}

	return g.Get(t)
}

func (g *gold) DoP(t TestingT, name string, data []byte) []byte {
	t.Helper()

	if g.Update() {
		g.SetP(t, name, data)
	}

	return g.GetP(t, name)
}

func (g *gold) File(t TestingT) string {
	t.Helper()

	return g.file(t, "")
}

func (g *gold) FileP(t TestingT, name string) string {
	t.Helper()

	if name == "" {
		t.Fatalf("golden: name cannot be empty")
	}

	return g.file(t, name)
}

func (g *gold) file(t TestingT, name string) string {
	t.Helper()

	if t.Name() == "" {
		t.Fatalf("golden: could not determine filename for TestingT instance")
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

func (g *gold) Get(t TestingT) []byte {
	t.Helper()

	return g.get(t, "")
}

func (g *gold) GetP(t TestingT, name string) []byte {
	t.Helper()

	if name == "" {
		t.Fatalf("golden: name cannot be empty")
	}

	return g.get(t, name)
}

func (g *gold) get(t TestingT, name string) []byte {
	t.Helper()

	f := g.file(t, name)

	b, err := g.fs.ReadFile(f)
	if err != nil {
		t.Fatalf("golden: %s", err.Error())
	}

	return b
}

func (g *gold) Set(t TestingT, data []byte) {
	t.Helper()

	g.set(t, "", data)
}

func (g *gold) SetP(t TestingT, name string, data []byte) {
	t.Helper()

	if name == "" {
		t.Fatalf("golden: name cannot be empty")
	}

	g.set(t, name, data)
}

func (g *gold) set(t TestingT, name string, data []byte) {
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

	err = g.fs.WriteFile(f, data, g.fileMode)
	if err != nil {
		t.Fatalf("golden: filed to write file: %s", err.Error())
	}
}

func (g *gold) Update() bool {
	return g.updateFunc()
}
