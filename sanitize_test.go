package golden

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_sanitizeFilename(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     string
	}{
		{
			name:     "empty",
			filename: "",
			want:     "",
		},
		{
			name:     ".",
			filename: ".",
			want:     "_",
		},
		{
			name:     "..",
			filename: "..",
			want:     "__",
		},
		{
			name:     "...",
			filename: "...",
			want:     "___",
		},
		{
			name:     "clean",
			filename: "foo-bar-nope.golden",
			want:     "foo-bar-nope.golden",
		},
		{
			name:     "with spaces",
			filename: "foo  bar nope.golden",
			want:     "foo__bar_nope.golden",
		},
		{
			name:     "illegal chars",
			filename: `foo/?<>\:*|"bar.golden`,
			want:     "foo_________bar.golden",
		},
		{
			name: "control chars",
			filename: "foo\x00\x01\x02\x03\x04\x05\x06\x07\x08\x09\x0a\x0b" +
				"\x0c\x0d\x0e\x0f\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1a" +
				"\x1b\x1c\x1d\x1e\x1fbar.golden",
			want: "foo________________________________bar.golden",
		},
		{
			name:     "trailing whitespace",
			filename: "foobar.golden    ",
			want:     "foobar.golden",
		},
		{
			name:     "trailing dots",
			filename: "foobar.golden......",
			want:     "foobar.golden",
		},
		{
			name:     "trailing whitespace and dots",
			filename: "foobar.golden  ..  ..  ..  ",
			want:     "foobar.golden",
		},
		{name: "con", filename: "con", want: "___"},
		{name: "prn", filename: "prn", want: "___"},
		{name: "aux", filename: "aux", want: "___"},
		{name: "nul", filename: "nul", want: "___"},
		{name: "com1", filename: "com1", want: "____"},
		{name: "com2", filename: "com2", want: "____"},
		{name: "com3", filename: "com3", want: "____"},
		{name: "com4", filename: "com4", want: "____"},
		{name: "com5", filename: "com5", want: "____"},
		{name: "com6", filename: "com6", want: "____"},
		{name: "com7", filename: "com7", want: "____"},
		{name: "com8", filename: "com8", want: "____"},
		{name: "com9", filename: "com9", want: "____"},
		{name: "lpt1", filename: "lpt1", want: "____"},
		{name: "lpt2", filename: "lpt2", want: "____"},
		{name: "lpt3", filename: "lpt3", want: "____"},
		{name: "lpt4", filename: "lpt4", want: "____"},
		{name: "lpt5", filename: "lpt5", want: "____"},
		{name: "lpt6", filename: "lpt6", want: "____"},
		{name: "lpt7", filename: "lpt7", want: "____"},
		{name: "lpt8", filename: "lpt8", want: "____"},
		{name: "lpt9", filename: "lpt9", want: "____"},
		{name: "CON", filename: "CON", want: "___"},
		{name: "PRN", filename: "PRN", want: "___"},
		{name: "AUX", filename: "AUX", want: "___"},
		{name: "NUL", filename: "NUL", want: "___"},
		{name: "COM1", filename: "COM1", want: "____"},
		{name: "COM2", filename: "COM2", want: "____"},
		{name: "COM3", filename: "COM3", want: "____"},
		{name: "COM4", filename: "COM4", want: "____"},
		{name: "COM5", filename: "COM5", want: "____"},
		{name: "COM6", filename: "COM6", want: "____"},
		{name: "COM7", filename: "COM7", want: "____"},
		{name: "COM8", filename: "COM8", want: "____"},
		{name: "COM9", filename: "COM9", want: "____"},
		{name: "LPT1", filename: "LPT1", want: "____"},
		{name: "LPT2", filename: "LPT2", want: "____"},
		{name: "LPT3", filename: "LPT3", want: "____"},
		{name: "LPT4", filename: "LPT4", want: "____"},
		{name: "LPT5", filename: "LPT5", want: "____"},
		{name: "LPT6", filename: "LPT6", want: "____"},
		{name: "LPT7", filename: "LPT7", want: "____"},
		{name: "LPT8", filename: "LPT8", want: "____"},
		{name: "LPT9", filename: "LPT9", want: "____"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sanitizeFilename(tt.filename)

			assert.Equal(t, tt.want, got)
		})
	}
}
