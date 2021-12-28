package unmarshal

import (
	"bytes"
	"encoding/json"
	"encoding/xml"

	"gopkg.in/yaml.v3"
)

// JSON parses the JSON-encoded data and stores the result in the value pointed
// to by v. Unknown fields in the JSON data is not allowed.
func JSON(data []byte, v interface{}) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()

	return dec.Decode(v)
}

// XML parses the XML-encoded data and stores the result in the value pointed
// to by v.
func XML(data []byte, v interface{}) error {
	dec := xml.NewDecoder(bytes.NewReader(data))

	return dec.Decode(v)
}

// YAML parses the YAML-encoded data and stores the result in the value pointed
// to by v. Unknown fields in the YAML data is not allowed.
func YAML(data []byte, v interface{}) error {
	dec := yaml.NewDecoder(bytes.NewReader(data))
	dec.KnownFields(true)

	return dec.Decode(v)
}
