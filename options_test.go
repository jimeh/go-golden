package golden

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithDirMode(t *testing.T) {
	customMode := os.FileMode(0o700)
	g := &Golden{}

	opt := WithDirMode(customMode)
	opt(g)

	assert.Equal(t, customMode, g.DirMode)
}

func TestWithFileMode(t *testing.T) {
	customMode := os.FileMode(0o600)
	g := &Golden{}

	opt := WithFileMode(customMode)
	opt(g)

	assert.Equal(t, customMode, g.FileMode)
}

func TestWithSuffix(t *testing.T) {
	customSuffix := ".custom"
	g := &Golden{}

	opt := WithSuffix(customSuffix)
	opt(g)

	assert.Equal(t, customSuffix, g.Suffix)
}

func TestWithDirname(t *testing.T) {
	customDirname := "custom-testdata"
	g := &Golden{}

	opt := WithDirname(customDirname)
	opt(g)

	assert.Equal(t, customDirname, g.Dirname)
}

func TestWithUpdateFunc(t *testing.T) {
	customUpdateFunc := func() bool { return true }
	g := &Golden{}

	opt := WithUpdateFunc(customUpdateFunc)
	opt(g)

	assertSameFunc(t, customUpdateFunc, g.UpdateFunc)
}
