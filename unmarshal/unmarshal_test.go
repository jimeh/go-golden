package unmarshal_test

import (
	"io"
	"testing"

	"github.com/jimeh/go-golden/unmarshal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type book struct {
	Title  string `json:"title" yaml:"title" xml:"title"`
	Author string `json:"author" yaml:"author" xml:"author"`
	Price  int    `json:"price" yaml:"price" xml:"price"`
}

type shoe struct {
	Make  string `json:"make" yaml:"make" xml:"make"`
	Model string `json:"model" yaml:"model" xml:"model"`
	Size  int    `json:"size" yaml:"size" xml:"size"`
}

func TestJSON(t *testing.T) {
	type args struct {
		data []byte
		v    interface{}
	}
	tests := []struct {
		name      string
		args      args
		want      interface{}
		wantErr   string
		wantErrIs error
	}{
		{
			name:      "nil",
			args:      args{data: nil, v: nil},
			wantErrIs: io.EOF,
		},
		{
			name: "empty string (1)",
			args: args{
				data: []byte(""),
				v:    &book{},
			},
			wantErrIs: io.EOF,
		},
		{
			name: "empty string (2)",
			args: args{
				data: []byte(""),
				v:    &shoe{},
			},
			wantErrIs: io.EOF,
		},
		{
			name: "no fields (1)",
			args: args{
				data: []byte("{}"),
				v:    &book{},
			},
			want: &book{},
		},
		{
			name: "no fields (2)",
			args: args{
				data: []byte("{}"),
				v:    &shoe{},
			},
			want: &shoe{},
		},
		{
			name: "empty fields (1)",
			args: args{
				data: []byte(`{"title":"","author":"","price":0}`),
				v:    &book{},
			},
			want: &book{},
		},
		{
			name: "empty fields (2)",
			args: args{
				data: []byte(`{"make":"","model":"","size":0}`),
				v:    &shoe{},
			},
			want: &shoe{},
		},
		{
			name: "populated fields (1)",
			args: args{
				data: []byte(`{"title":"a","author":"b","price":499}`),
				v:    &book{},
			},
			want: &book{Title: "a", Author: "b", Price: 499},
		},
		{
			name: "populated fields (2)",
			args: args{
				data: []byte(`{"Make":"a","model":"b","size":42}`),
				v:    &shoe{},
			},
			want: &shoe{Make: "a", Model: "b", Size: 42},
		},
		{
			name: "unknown field (1)",
			args: args{
				data: []byte(`{"title":"a","summary":"b","price":499}`),
				v:    &book{},
			},
			wantErr: `json: unknown field "summary"`,
		},
		{
			name: "unknown field (2)",
			args: args{
				data: []byte(`{"make":"a","inventory":"b","size":42}`),
				v:    &shoe{},
			},
			wantErr: `json: unknown field "inventory"`,
		},
		{
			name: "to channel",
			args: args{
				data: []byte(`{"make":"a","model":"b","size":42}`),
				v:    make(chan int),
			},
			wantErr: `json: Unmarshal(non-pointer chan int)`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := unmarshal.JSON(tt.args.data, tt.args.v)

			if tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
			}

			if tt.wantErrIs != nil {
				assert.ErrorIs(t, err, tt.wantErrIs)
			}

			if tt.wantErr == "" && tt.wantErrIs == nil {
				require.NoError(t, err)
				assert.Equal(t, tt.want, tt.args.v)
			}
		})
	}
}

func TestYAML(t *testing.T) {
	type args struct {
		data []byte
		v    interface{}
	}
	tests := []struct {
		name      string
		args      args
		want      interface{}
		wantErr   string
		wantErrIs error
	}{
		{
			name:      "nil",
			args:      args{data: nil, v: nil},
			wantErrIs: io.EOF,
		},
		{
			name: "empty string (1)",
			args: args{
				data: []byte(""),
				v:    &book{},
			},
			wantErrIs: io.EOF,
		},
		{
			name: "empty string (2)",
			args: args{
				data: []byte(""),
				v:    &shoe{},
			},
			wantErrIs: io.EOF,
		},
		{
			name: "no fields (1)",
			args: args{
				data: []byte("{}"),
				v:    &book{},
			},
			want: &book{},
		},
		{
			name: "no fields (2)",
			args: args{
				data: []byte("{}"),
				v:    &shoe{},
			},
			want: &shoe{},
		},
		{
			name: "empty fields (1)",
			args: args{
				data: []byte("title:\nauthor:\nprice: 0"),
				v:    &book{},
			},
			want: &book{},
		},
		{
			name: "empty fields (2)",
			args: args{
				data: []byte("make:\nmodel:\nsize: 0"),
				v:    &shoe{},
			},
			want: &shoe{},
		},
		{
			name: "populated fields (1)",
			args: args{
				data: []byte("title: a\nauthor: b\nprice: 499"),
				v:    &book{},
			},
			want: &book{Title: "a", Author: "b", Price: 499},
		},
		{
			name: "populated fields (2)",
			args: args{
				data: []byte("make: a\nmodel: b\nsize: 42"),
				v:    &shoe{},
			},
			want: &shoe{Make: "a", Model: "b", Size: 42},
		},
		{
			name: "unknown field (1)",
			args: args{
				data: []byte("title: a\nsummary: b\nprice: 499"),
				v:    &book{},
			},
			wantErr: "yaml: unmarshal errors:\n  " +
				"line 2: field summary not found in type unmarshal_test.book",
		},
		{
			name: "unknown field (2)",
			args: args{
				data: []byte("make: a\ninventory: b\nsize: 42"),
				v:    &shoe{},
			},
			wantErr: "yaml: unmarshal errors:\n  " +
				"line 2: field inventory not found in type unmarshal_test.shoe",
		},
		{
			name: "to channel",
			args: args{
				data: []byte("make: a\nmodel: b\nsize: 42"),
				v:    make(chan int),
			},
			wantErr: "yaml: unmarshal errors:\n  " +
				"line 1: cannot unmarshal !!map into chan int",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := unmarshal.YAML(tt.args.data, tt.args.v)

			if tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
			}

			if tt.wantErrIs != nil {
				assert.ErrorIs(t, err, tt.wantErrIs)
			}

			if tt.wantErr == "" && tt.wantErrIs == nil {
				require.NoError(t, err)
				assert.Equal(t, tt.want, tt.args.v)
			}
		})
	}
}

func TestXML(t *testing.T) {
	type args struct {
		data []byte
		v    interface{}
	}
	tests := []struct {
		name      string
		args      args
		want      interface{}
		wantErr   string
		wantErrIs error
	}{
		{
			name:    "nil",
			args:    args{data: nil, v: nil},
			wantErr: "non-pointer passed to Unmarshal",
		},
		{
			name: "empty string (1)",
			args: args{
				data: []byte(""),
				v:    &book{},
			},
			wantErrIs: io.EOF,
		},
		{
			name: "empty string (2)",
			args: args{
				data: []byte(""),
				v:    &shoe{},
			},
			wantErrIs: io.EOF,
		},
		{
			name: "no fields (1)",
			args: args{
				data: []byte("<book></book>"),
				v:    &book{},
			},
			want: &book{},
		},
		{
			name: "no fields (2)",
			args: args{
				data: []byte("<shoe></shoe>"),
				v:    &shoe{},
			},
			want: &shoe{},
		},
		{
			name: "empty fields (1)",
			args: args{
				data: []byte("<book>" +
					"<title></title>" +
					"<author></author>" +
					"<price></price>" +
					"</book>"),
				v: &book{},
			},
			want: &book{},
		},
		{
			name: "empty fields (2)",
			args: args{
				data: []byte("<shoe>" +
					"<make></make>" +
					"<model></model>" +
					"<size></size>" +
					"</shoe>"),
				v: &shoe{},
			},
			want: &shoe{},
		},
		{
			name: "populated fields (1)",
			args: args{
				data: []byte("<book>" +
					"<title>a</title>" +
					"<author>b</author>" +
					"<price>499</price>" +
					"</book>"),
				v: &book{},
			},
			want: &book{Title: "a", Author: "b", Price: 499},
		},
		{
			name: "populated fields (2)",
			args: args{
				data: []byte("<shoe>" +
					"<make>a</make>" +
					"<model>b</model>" +
					"<size>42</size>" +
					"</shoe>"),
				v: &shoe{},
			},
			want: &shoe{Make: "a", Model: "b", Size: 42},
		},
		{
			name: "unknown field (1)",
			args: args{
				data: []byte("<book>" +
					"<title>a</title>" +
					"<summary>b</summary>" +
					"<price>499</price>" +
					"</book>"),
				v: &book{},
			},
			want: &book{Title: "a", Author: "", Price: 499},
		},
		{
			name: "unknown field (2)",
			args: args{
				data: []byte("<shoe>" +
					"<make>a</make>" +
					"<inventory>b</inventory>" +
					"<size>42</size>" +
					"</shoe>"),
				v: &shoe{},
			},
			want: &shoe{Make: "a", Model: "", Size: 42},
		},
		{
			name: "to channel",
			args: args{
				data: []byte("<shoe>" +
					"<make>a</make>" +
					"<model>b</model>" +
					"<size>42</size>" +
					"</shoe>"),
				v: make(chan int),
			},
			wantErr: "non-pointer passed to Unmarshal",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := unmarshal.XML(tt.args.data, tt.args.v)

			if tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
			}

			if tt.wantErrIs != nil {
				assert.ErrorIs(t, err, tt.wantErrIs)
			}

			if tt.wantErr == "" && tt.wantErrIs == nil {
				require.NoError(t, err)
				assert.Equal(t, tt.want, tt.args.v)
			}
		})
	}
}
