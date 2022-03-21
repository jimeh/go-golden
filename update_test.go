package golden

import (
	"testing"

	"github.com/jimeh/envctl"
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
		name: "GOLDEN_UPDATE set to Y",
		env:  map[string]string{"GOLDEN_UPDATE": "Y"},
		want: true,
	},
	{
		name: "GOLDEN_UPDATE set to n",
		env:  map[string]string{"GOLDEN_UPDATE": "n"},
		want: false,
	},
	{
		name: "GOLDEN_UPDATE set to N",
		env:  map[string]string{"GOLDEN_UPDATE": "N"},
		want: false,
	},
	{
		name: "GOLDEN_UPDATE set to t",
		env:  map[string]string{"GOLDEN_UPDATE": "t"},
		want: true,
	},
	{
		name: "GOLDEN_UPDATE set to T",
		env:  map[string]string{"GOLDEN_UPDATE": "T"},
		want: true,
	},
	{
		name: "GOLDEN_UPDATE set to f",
		env:  map[string]string{"GOLDEN_UPDATE": "f"},
		want: false,
	},
	{
		name: "GOLDEN_UPDATE set to F",
		env:  map[string]string{"GOLDEN_UPDATE": "F"},
		want: false,
	},
	{
		name: "GOLDEN_UPDATE set to yes",
		env:  map[string]string{"GOLDEN_UPDATE": "yes"},
		want: true,
	},
	{
		name: "GOLDEN_UPDATE set to Yes",
		env:  map[string]string{"GOLDEN_UPDATE": "Yes"},
		want: true,
	},
	{
		name: "GOLDEN_UPDATE set to YeS",
		env:  map[string]string{"GOLDEN_UPDATE": "YeS"},
		want: true,
	},
	{
		name: "GOLDEN_UPDATE set to YES",
		env:  map[string]string{"GOLDEN_UPDATE": "YES"},
		want: true,
	},
	{
		name: "GOLDEN_UPDATE set to no",
		env:  map[string]string{"GOLDEN_UPDATE": "no"},
		want: false,
	},
	{
		name: "GOLDEN_UPDATE set to No",
		env:  map[string]string{"GOLDEN_UPDATE": "No"},
		want: false,
	},
	{
		name: "GOLDEN_UPDATE set to nO",
		env:  map[string]string{"GOLDEN_UPDATE": "nO"},
		want: false,
	},
	{
		name: "GOLDEN_UPDATE set to NO",
		env:  map[string]string{"GOLDEN_UPDATE": "NO"},
		want: false,
	},
	{
		name: "GOLDEN_UPDATE set to on",
		env:  map[string]string{"GOLDEN_UPDATE": "on"},
		want: true,
	},
	{
		name: "GOLDEN_UPDATE set to oN",
		env:  map[string]string{"GOLDEN_UPDATE": "oN"},
		want: true,
	},
	{
		name: "GOLDEN_UPDATE set to On",
		env:  map[string]string{"GOLDEN_UPDATE": "On"},
		want: true,
	},
	{
		name: "GOLDEN_UPDATE set to ON",
		env:  map[string]string{"GOLDEN_UPDATE": "ON"},
		want: true,
	},
	{
		name: "GOLDEN_UPDATE set to off",
		env:  map[string]string{"GOLDEN_UPDATE": "off"},
		want: false,
	},
	{
		name: "GOLDEN_UPDATE set to Off",
		env:  map[string]string{"GOLDEN_UPDATE": "Off"},
		want: false,
	},
	{
		name: "GOLDEN_UPDATE set to oFF",
		env:  map[string]string{"GOLDEN_UPDATE": "oFF"},
		want: false,
	},
	{
		name: "GOLDEN_UPDATE set to OFF",
		env:  map[string]string{"GOLDEN_UPDATE": "OFF"},
		want: false,
	},
	{
		name: "GOLDEN_UPDATE set to true",
		env:  map[string]string{"GOLDEN_UPDATE": "true"},
		want: true,
	},
	{
		name: "GOLDEN_UPDATE set to True",
		env:  map[string]string{"GOLDEN_UPDATE": "True"},
		want: true,
	},
	{
		name: "GOLDEN_UPDATE set to TruE",
		env:  map[string]string{"GOLDEN_UPDATE": "TruE"},
		want: true,
	},
	{
		name: "GOLDEN_UPDATE set to TRUE",
		env:  map[string]string{"GOLDEN_UPDATE": "TRUE"},
		want: true,
	},
	{
		name: "GOLDEN_UPDATE set to false",
		env:  map[string]string{"GOLDEN_UPDATE": "false"},
		want: false,
	},
	{
		name: "GOLDEN_UPDATE set to False",
		env:  map[string]string{"GOLDEN_UPDATE": "False"},
		want: false,
	},
	{
		name: "GOLDEN_UPDATE set to FaLsE",
		env:  map[string]string{"GOLDEN_UPDATE": "FaLsE"},
		want: false,
	},
	{
		name: "GOLDEN_UPDATE set to FALSE",
		env:  map[string]string{"GOLDEN_UPDATE": "FALSE"},
		want: false,
	},
	{
		name: "GOLDEN_UPDATE set to foobarnopebbq",
		env:  map[string]string{"GOLDEN_UPDATE": "foobarnopebbq"},
		want: false,
	},
	{
		name: "GOLDEN_UPDATE set to FOOBARNOPEBBQ",
		env:  map[string]string{"GOLDEN_UPDATE": "FOOBARNOPEBBQ"},
		want: false,
	},
}

func TestEnvUpdateFunc(t *testing.T) {
	for _, tt := range envUpdateFuncTestCases {
		t.Run(tt.name, func(t *testing.T) {
			envctl.WithClean(tt.env, func() {
				got := EnvUpdateFunc()

				assert.Equal(t, tt.want, got)
			})
		})
	}
}
