package golden

var defaultAsserter = NewAsserter()

// AssertJSONMarshaling asserts that the given "v" value JSON marshals to an
// expected value fetched from a golden file on disk, and then verifies that the
// marshaled result produces a value that is equal to "v" when unmarshaled.
//
// Used for objects that do NOT change when they are marshaled and unmarshaled.
func AssertJSONMarshaling(t TestingT, v interface{}) {
	t.Helper()

	defaultAsserter.JSONMarshaling(t, v)
}

// AssertJSONMarshalingP asserts that the given "v" value JSON marshals to an
// expected value fetched from a golden file on disk, and then verifies that the
// marshaled result produces a value that is equal to "want" when unmarshaled.
//
// Used for objects that change when they are marshaled and unmarshaled.
func AssertJSONMarshalingP(t TestingT, v, want interface{}) {
	t.Helper()

	defaultAsserter.JSONMarshalingP(t, v, want)
}

// AssertXMLMarshaling asserts that the given "v" value XML marshals to an
// expected value fetched from a golden file on disk, and then verifies that the
// marshaled result produces a value that is equal to "v" when unmarshaled.
//
// Used for objects that do NOT change when they are marshaled and unmarshaled.
func AssertXMLMarshaling(t TestingT, v interface{}) {
	t.Helper()

	defaultAsserter.XMLMarshaling(t, v)
}

// AssertXMLMarshalingP asserts that the given "v" value XML marshals to an
// expected value fetched from a golden file on disk, and then verifies that the
// marshaled result produces a value that is equal to "want" when unmarshaled.
//
// Used for objects that change when they are marshaled and unmarshaled.
func AssertXMLMarshalingP(t TestingT, v, want interface{}) {
	t.Helper()

	defaultAsserter.XMLMarshalingP(t, v, want)
}

// AssertYAMLMarshaling asserts that the given "v" value YAML marshals to an
// expected value fetched from a golden file on disk, and then verifies that the
// marshaled result produces a value that is equal to "v" when unmarshaled.
//
// Used for objects that do NOT change when they are marshaled and unmarshaled.
func AssertYAMLMarshaling(t TestingT, v interface{}) {
	t.Helper()

	defaultAsserter.YAMLMarshaling(t, v)
}

// AssertYAMLMarshalingP asserts that the given "v" value YAML marshals to an
// expected value fetched from a golden file on disk, and then verifies that the
// marshaled result produces a value that is equal to "want" when unmarshaled.
//
// Used for objects that change when they are marshaled and unmarshaled.
func AssertYAMLMarshalingP(t TestingT, v, want interface{}) {
	t.Helper()

	defaultAsserter.YAMLMarshalingP(t, v, want)
}
