package sanitize_test

import (
	"testing"

	"github.com/jimeh/go-golden/sanitize"
	"github.com/stretchr/testify/assert"
)

func TestLineBreaks(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "nil",
			args: args{data: nil},
			want: nil,
		},
		{
			name: "empty",
			args: args{data: []byte{}},
			want: nil,
		},
		{
			name: "no line breaks",
			args: args{data: []byte("hello world")},
			want: []byte("hello world"),
		},
		{
			name: "UNIX line breaks",
			args: args{data: []byte("hello\nworld\nhow are you?")},
			want: []byte("hello\nworld\nhow are you?"),
		},
		{
			name: "Windows line breaks",
			args: args{data: []byte("hello\r\nworld\r\nhow are you?")},
			want: []byte("hello\nworld\nhow are you?"),
		},
		{
			name: "MacOS Classic line breaks",
			args: args{data: []byte("hello\rworld\rhow are you?")},
			want: []byte("hello\nworld\nhow are you?"),
		},
		{
			name: "Windows and MacOS Classic line breaks",
			args: args{data: []byte("hello\r\nworld\rhow are you?")},
			want: []byte("hello\nworld\nhow are you?"),
		},
		{
			name: "Windows, MacOS Classic, and UNIX line breaks",
			args: args{data: []byte("hello\r\nworld\rhow are you?\nGood!")},
			want: []byte("hello\nworld\nhow are you?\nGood!"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sanitize.LineBreaks(tt.args.data)

			assert.Equal(t, tt.want, got)
		})
	}
}
