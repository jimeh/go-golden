package golden

// TestingT is a interface describing a sub-set of methods of *testing.T which
// golden uses.
type TestingT interface {
	Fatalf(format string, args ...interface{})
	Helper()
	Logf(format string, args ...interface{})
	Name() string
}
