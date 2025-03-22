package golden

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var envUpdateFuncTestCases = []struct {
	name string
	env  map[string]string
	want bool
}{
	{
		name: "GOLDEN_UPDATE not set",
		want: false,
	},
	{
		name: "GOLDEN_UPDATE set to empty string",
		env:  map[string]string{"GOLDEN_UPDATE": ""},
		want: false,
	},
	{
		name: "GOLDEN_UPDATE set to 0",
		env:  map[string]string{"GOLDEN_UPDATE": "0"},
		want: false,
	},
	{
		name: "GOLDEN_UPDATE set to 1",
		env:  map[string]string{"GOLDEN_UPDATE": "1"},
		want: true,
	},
	{
		name: "GOLDEN_UPDATE set to 2",
		env:  map[string]string{"GOLDEN_UPDATE": "2"},
		want: false,
	},
	{
		name: "GOLDEN_UPDATE set to y",
		env:  map[string]string{"GOLDEN_UPDATE": "y"},
		want: true,
	},
	{
		name: "GOLDEN_UPDATE set to n",
		env:  map[string]string{"GOLDEN_UPDATE": "n"},
		want: false,
	},
	{
		name: "GOLDEN_UPDATE set to t",
		env:  map[string]string{"GOLDEN_UPDATE": "t"},
		want: true,
	},
	{
		name: "GOLDEN_UPDATE set to f",
		env:  map[string]string{"GOLDEN_UPDATE": "f"},
		want: false,
	},
	{
		name: "GOLDEN_UPDATE set to yes",
		env:  map[string]string{"GOLDEN_UPDATE": "yes"},
		want: true,
	},
	{
		name: "GOLDEN_UPDATE set to no",
		env:  map[string]string{"GOLDEN_UPDATE": "no"},
		want: false,
	},
	{
		name: "GOLDEN_UPDATE set to on",
		env:  map[string]string{"GOLDEN_UPDATE": "on"},
		want: true,
	},
	{
		name: "GOLDEN_UPDATE set to off",
		env:  map[string]string{"GOLDEN_UPDATE": "off"},
		want: false,
	},
	{
		name: "GOLDEN_UPDATE set to true",
		env:  map[string]string{"GOLDEN_UPDATE": "true"},
		want: true,
	},
	{
		name: "GOLDEN_UPDATE set to false",
		env:  map[string]string{"GOLDEN_UPDATE": "false"},
		want: false,
	},
	{
		name: "GOLDEN_UPDATE set to foobarnopebbq",
		env:  map[string]string{"GOLDEN_UPDATE": "foobarnopebbq"},
		want: false,
	},
}

func TestEnvUpdateFunc(t *testing.T) {
	for _, tt := range envUpdateFuncTestCases {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.env {
				t.Setenv(k, v)
			}

			got := EnvUpdateFunc()

			assert.Equal(t, tt.want, got)
		})
	}
}
