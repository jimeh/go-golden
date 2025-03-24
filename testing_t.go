package golden

type TestingT interface {
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Helper()
	Logf(format string, args ...interface{})
	Name() string
}
