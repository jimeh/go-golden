package golden

import "os"

type FS interface {
	// MkdirAll creates a directory named path, along with any necessary
	// parents, and returns nil, or else returns an error. The permission bits
	// perm (before umask) are used for all directories that MkdirAll creates.
	MkdirAll(path string, perm os.FileMode) error

	// ReadFile reads the named file and returns the contents. A successful call
	// returns err == nil, not err == EOF. Because ReadFile reads the whole
	// file, it does not treat an EOF from Read as an error to be reported.
	ReadFile(filename string) ([]byte, error)

	// WriteFile writes data to a file named by filename. If the file does not
	// exist, WriteFile creates it with permissions perm; otherwise WriteFile
	// truncates it before writing, without changing permissions.
	WriteFile(name string, data []byte, perm os.FileMode) error
}

type fsImpl struct{}

var _ FS = fsImpl{}

// NewFS returns a new FS instance which operates against the host file system
// via calls to functions in the os package.
func NewFS() FS {
	return fsImpl{}
}

// DefaultFS is the default FS instance used by all top-level package functions,
// including the Default Golden instance, and also the New function.
var DefaultFS = NewFS()

func (fsImpl) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (fsImpl) ReadFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

func (fsImpl) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return os.WriteFile(filename, data, perm)
}
