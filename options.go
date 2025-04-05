package golden

import (
	"os"
)

// Option is a function that modifies a Golden instance.
type Option func(*Golden)

// WithDirMode sets the directory mode for a Golden instance.
func WithDirMode(mode os.FileMode) Option {
	return func(g *Golden) {
		g.DirMode = mode
	}
}

// WithFileMode sets the file mode for a Golden instance.
func WithFileMode(mode os.FileMode) Option {
	return func(g *Golden) {
		g.FileMode = mode
	}
}

// WithSuffix sets the file suffix for a Golden instance.
func WithSuffix(suffix string) Option {
	return func(g *Golden) {
		g.Suffix = suffix
	}
}

// WithDirname sets the directory name for a Golden instance.
func WithDirname(dirname string) Option {
	return func(g *Golden) {
		g.Dirname = dirname
	}
}

// WithUpdateFunc sets the update function for a Golden instance.
func WithUpdateFunc(updateFunc UpdateFunc) Option {
	return func(g *Golden) {
		g.UpdateFunc = updateFunc
	}
}
