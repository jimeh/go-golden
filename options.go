package golden

import "os"

type Option interface {
	apply(*gold)
}

type optionFunc func(*gold)

func (fn optionFunc) apply(g *gold) {
	fn(g)
}

// WithDirMode sets the file system permissions used for any folders created to
// hold golden files.
//
// When this option is not provided, the default value is 0o755.
func WithDirMode(mode os.FileMode) Option {
	return optionFunc(func(g *gold) {
		g.dirMode = mode
	})
}

// WithFileMode sets the file system permissions used for any created or updated
// golden files written to.
//
// When this option is not provided, the default value is 0o644.
func WithFileMode(mode os.FileMode) Option {
	return optionFunc(func(g *gold) {
		g.fileMode = mode
	})
}

// WithSuffix sets the filename suffix used for all golden files.
//
// When this option is not provided, the default value is ".golden".
func WithSuffix(suffix string) Option {
	return optionFunc(func(g *gold) {
		g.suffix = suffix
	})
}

// WithDirname sets the name of the top-level directory used to hold golden
// files.
//
// When this option is not provided, the default value is "testdata".
func WithDirname(name string) Option {
	return optionFunc(func(g *gold) {
		g.dirname = name
	})
}

// WithUpdateFunc sets the function used to determine if golden files should be
// updated or not. Essentially the provided UpdateFunc is called by Update().
//
// When this option is not provided, the default value is EnvUpdateFunc.
func WithUpdateFunc(fn UpdateFunc) Option {
	return optionFunc(func(g *gold) {
		g.updateFunc = fn
	})
}

// WithFS sets the afero.Fs instance which is used for all file system
// operations to read/write golden files.
//
// When this option is not provided, the default value is afero.NewOsFs().
func WithFS(fs FS) Option {
	return optionFunc(func(g *gold) {
		g.fs = fs
	})
}

// WithSilentWrites silences the "golden: writing [...]" log messages whenever
// set functions write a golden file to disk.
func WithSilentWrites() Option {
	return optionFunc(func(g *gold) {
		g.logOnWrite = false
	})
}
