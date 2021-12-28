package golden

import (
	"testing"

	"github.com/jimeh/go-golden/marshal"
	"github.com/jimeh/go-golden/unmarshal"
	"github.com/stretchr/testify/assert"
)

func TestNewMarshalingAsserter(t *testing.T) {
	type args struct {
		golden              Golden
		format              string
		marshalFunc         MarshalFunc
		unmarshalFunc       UnmarshalFunc
		normalizeLineBreaks bool
	}
	tests := []struct {
		name string
		args args
		want *MarshalingAsserter
	}{
		{
			name: "json",
			args: args{
				nil,
				"JSON",
				marshal.JSON,
				unmarshal.JSON,
				true,
			},
			want: &MarshalingAsserter{
				Golden:              nil,
				Format:              "JSON",
				GoldName:            "marshaled_json",
				MarshalFunc:         marshal.JSON,
				UnmarshalFunc:       unmarshal.JSON,
				NormalizeLineBreaks: true,
			},
		},
		{
			name: "xml",
			args: args{
				nil,
				"XML",
				marshal.XML,
				unmarshal.XML,
				true,
			},
			want: &MarshalingAsserter{
				Golden:              nil,
				Format:              "XML",
				GoldName:            "marshaled_xml",
				MarshalFunc:         marshal.XML,
				UnmarshalFunc:       unmarshal.XML,
				NormalizeLineBreaks: true,
			},
		},
		{
			name: "yaml",
			args: args{
				nil,
				"YAML",
				marshal.YAML,
				unmarshal.YAML,
				true,
			},
			want: &MarshalingAsserter{
				Golden:              nil,
				Format:              "YAML",
				GoldName:            "marshaled_yaml",
				MarshalFunc:         marshal.YAML,
				UnmarshalFunc:       unmarshal.YAML,
				NormalizeLineBreaks: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewMarshalingAsserter(
				tt.args.golden,
				tt.args.format,
				tt.args.marshalFunc,
				tt.args.unmarshalFunc,
				tt.args.normalizeLineBreaks,
			)

			assert.Equal(t, tt.want.Golden, got.Golden)
			assert.Equal(t, tt.want.Format, got.Format)
			assert.Equal(t, tt.want.GoldName, got.GoldName)
			assert.Equal(t,
				tt.want.NormalizeLineBreaks,
				got.NormalizeLineBreaks,
			)
		})
	}
}
