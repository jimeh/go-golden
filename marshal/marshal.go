package marshal

import (
	"bytes"
	"encoding/json"
	"encoding/xml"

	"gopkg.in/yaml.v3"
)

// JSON returns the JSON encoding of v. Returned JSON is intended by two spaces
// (pretty formatted), and is not HTML escaped.
func JSON(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)

	err := enc.Encode(v)

	return buf.Bytes(), err
}

// XML returns the XML encoding of v. Returned XML is intended by two spaces
// (pretty formatted).
func XML(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := xml.NewEncoder(&buf)
	enc.Indent("", "  ")

	err := enc.Encode(v)

	return buf.Bytes(), err
}

// YAML returns the YAML encoding of v. Returned YAML is intended by two spaces.
func YAML(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)

	err := enc.Encode(v)

	return buf.Bytes(), err
}
