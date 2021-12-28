package marshal_test

import (
	"testing"

	"github.com/jimeh/go-golden/marshal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type book struct {
	Title  string `json:"title" yaml:"title" xml:"title"`
	Author string `json:"author,omitempty" yaml:"author,omitempty" xml:"author,omitempty"`
	Price  int    `json:"price" yaml:"price" xml:"price"`
}

type shoe struct {
	Make  string `json:"make" yaml:"make" xml:"make"`
	Model string `json:"model,omitempty" yaml:"model,omitempty" xml:"model,omitempty"`
	Size  int    `json:"size" yaml:"size" xml:"size"`
}

func TestJSON(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name      string
		args      args
		want      []byte
		wantErr   string
		wantErrIs error
	}{
		{
			name: "nil",
			args: args{v: nil},
			want: []byte("null\n"),
		},
		{
			name: "empty struct (1)",
			args: args{v: &book{}},
			want: []byte(`{
  "title": "",
  "price": 0
}
`,
			),
		},
		{
			name: "empty struct (2)",
			args: args{v: &shoe{}},
			want: []byte(`{
  "make": "",
  "size": 0
}
`,
			),
		},
		{
			name: "full struct (1)",
			args: args{
				v: &book{
					Title:  "a",
					Author: "b",
					Price:  499,
				},
			},
			want: []byte(`{
  "title": "a",
  "author": "b",
  "price": 499
}
`,
			),
		},
		{
			name: "empty struct (2)",
			args: args{
				v: &shoe{
					Make:  "a",
					Model: "b",
					Size:  42,
				},
			},
			want: []byte(`{
  "make": "a",
  "model": "b",
  "size": 42
}
`,
			),
		},
		{
			name: "channel",
			args: args{
				v: make(chan int),
			},
			wantErr: "json: unsupported type: chan int",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := marshal.JSON(tt.args.v)

			if tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
			}

			if tt.wantErrIs != nil {
				assert.ErrorIs(t, err, tt.wantErrIs)
			}

			if tt.wantErr == "" && tt.wantErrIs == nil {
				require.NoError(t, err)
			}

			assert.Equal(t, string(tt.want), string(got))
		})
	}
}

func TestYAML(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name      string
		args      args
		want      []byte
		wantErr   string
		wantErrIs error
		wantPanic interface{}
	}{
		{
			name: "nil",
			args: args{v: nil},
			want: []byte("null\n"),
		},
		{
			name: "empty struct (1)",
			args: args{v: &book{}},
			want: []byte(`title: ""
price: 0
`,
			),
		},
		{
			name: "empty struct (2)",
			args: args{v: &shoe{}},
			want: []byte(`make: ""
size: 0
`,
			),
		},
		{
			name: "full struct (1)",
			args: args{
				v: &book{
					Title:  "a",
					Author: "b",
					Price:  499,
				},
			},
			want: []byte(`title: a
author: b
price: 499
`,
			),
		},
		{
			name: "empty struct (2)",
			args: args{
				v: &shoe{
					Make:  "a",
					Model: "b",
					Size:  42,
				},
			},
			want: []byte(`make: a
model: b
size: 42
`,
			),
		},
		{
			name: "channel",
			args: args{
				v: make(chan int),
			},
			wantPanic: "cannot marshal type: chan int",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := func() (got []byte, err error, p interface{}) {
				defer func() { p = recover() }()
				got, err = marshal.YAML(tt.args.v)

				return
			}

			got, err, p := f()

			if tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
			}

			if tt.wantErrIs != nil {
				assert.ErrorIs(t, err, tt.wantErrIs)
			}

			if tt.wantPanic != nil {
				assert.Equal(t, tt.wantPanic, p)
			}

			if tt.wantErr == "" && tt.wantErrIs == nil {
				require.NoError(t, err)
			}

			assert.Equal(t, string(tt.want), string(got))
		})
	}
}

func TestXML(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name      string
		args      args
		want      []byte
		wantErr   string
		wantErrIs error
	}{
		{
			name: "nil",
			args: args{v: nil},
			want: []byte(""),
		},
		{
			name: "empty struct (1)",
			args: args{v: &book{}},
			want: []byte(`<book>
  <title></title>
  <price>0</price>
</book>`,
			),
		},
		{
			name: "empty struct (2)",
			args: args{v: &shoe{}},
			want: []byte(`<shoe>
  <make></make>
  <size>0</size>
</shoe>`,
			),
		},
		{
			name: "full struct (1)",
			args: args{
				v: &book{
					Title:  "a",
					Author: "b",
					Price:  499,
				},
			},
			want: []byte(`<book>
  <title>a</title>
  <author>b</author>
  <price>499</price>
</book>`,
			),
		},
		{
			name: "empty struct (2)",
			args: args{
				v: &shoe{
					Make:  "a",
					Model: "b",
					Size:  42,
				},
			},
			want: []byte(`<shoe>
  <make>a</make>
  <model>b</model>
  <size>42</size>
</shoe>`,
			),
		},
		{
			name: "channel",
			args: args{
				v: make(chan int),
			},
			wantErr: "xml: unsupported type: chan int",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := marshal.XML(tt.args.v)

			if tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
			}

			if tt.wantErrIs != nil {
				assert.ErrorIs(t, err, tt.wantErrIs)
			}

			if tt.wantErr == "" && tt.wantErrIs == nil {
				require.NoError(t, err)
			}

			assert.Equal(t, string(tt.want), string(got))
		})
	}
}
