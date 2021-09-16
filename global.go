package golden

import "testing"

var global = New()

// Updating returns true when golden is set to update golden files. Used to
// determine if golden.Set() should be called or not.
func Updating() bool {
	return global.Updating()
}

// Get returns the content of the default golden file for the given *testing.T
// instance as determined by t.Name(). If no golden file can be found/read, it
// will fail the test with t.Fatal().
func Get(t *testing.T) []byte {
	return global.Get(t)
}

// Set writes given data of the default golden file for the given *testing.T
// instance as determined by t.Name(). If writing fails it will fail the test
// with t.Fatal() detailing the error.
func Set(t *testing.T, data []byte) {
	global.Set(t, data)
}

// File returns the filename for the default golden file for the given
// *testing.T instance as determined by t.Name().
func File(t *testing.T) string {
	return global.File(t)
}

// GetNamed return the content of the specifically named golden file belonging
// to the given *testing.T instance as determined by t.Name(). If no golden file
// can be found/read, it will fail the test with t.Fatal().
func GetNamed(t *testing.T, name string) []byte {
	return global.GetNamed(t, name)
}

// SetNamed writes given data of the specifically named golden file belonging to
// the given *testing.T instance as determined by t.Name(). If writing fails it
// will fail the test with t.Fatal() detailing the error.
func SetNamed(t *testing.T, name string, data []byte) {
	global.SetNamed(t, name, data)
}

func NamedFile(t *testing.T, name string) string {
	return global.NamedFile(t, name)
}
