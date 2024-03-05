package golden

import (
	"fmt"
	"runtime"
)

type fakeTestingT struct {
	helper bool
	name   string
	logs   []string
	fatals []string
}

func (m *fakeTestingT) Helper() {
	m.helper = true
}

func (m *fakeTestingT) Fatalf(format string, args ...interface{}) {
	m.fatals = append(m.fatals, fmt.Sprintf(format, args...))
	runtime.Goexit()
}

func (m *fakeTestingT) Logf(format string, args ...interface{}) {
	m.logs = append(m.logs, fmt.Sprintf(format, args...))
}

func (m *fakeTestingT) Name() string {
	return m.name
}
