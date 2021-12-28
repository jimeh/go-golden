package sanitize

import "bytes"

var (
	lf   = []byte{10}
	cr   = []byte{13}
	crlf = []byte{13, 10}
)

// LineBreaks replaces Windows CRLF (\r\n) and MacOS Classic CR (\r)
// line-breaks with Unix LF (\n) line breaks.
func LineBreaks(data []byte) []byte {
	// Replace Windows CRLF (\r\n) with Unix LF (\n)
	result := bytes.ReplaceAll(data, crlf, lf)

	// Replace Classic MacOS CR (\r) with Unix LF (\n)
	result = bytes.ReplaceAll(result, cr, lf)

	return result
}
