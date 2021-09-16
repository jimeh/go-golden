package golden

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/jimeh/envctl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdating(t *testing.T) {
	tests := []struct {
		name string
		env  map[string]string
		want bool
	}{
		{
			name: "GOLDEN_UPDATE not set",
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			envctl.WithClean(tt.env, func() {
				got := Updating()

				assert.Equal(t, tt.want, got)
			})
		})
	}
}

func TestFile(t *testing.T) {
	got := File(t)

	assert.Equal(t, filepath.Join("testdata", "TestFile.golden"), got)

	tests := []struct {
		name string
		want string
	}{
		{
			name: "",
			want: filepath.Join("testdata", "TestFile", "#00.golden"),
		},
		{
			name: "foobar",
			want: filepath.Join("testdata", "TestFile", "foobar.golden"),
		},
		{
			name: "foo/bar",
			want: filepath.Join("testdata", "TestFile", "foo/bar.golden"),
		},
		{
			name: `"foobar"`,
			want: filepath.Join("testdata", "TestFile", "\"foobar\".golden"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := File(t)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGet(t *testing.T) {
	t.Cleanup(func() {
		err := os.RemoveAll(filepath.Join("testdata", t.Name()))
		require.NoError(t, err)
		err = os.Remove(filepath.Join("testdata", t.Name()+".golden"))
		require.NoError(t, err)
	})

	err := os.MkdirAll("testdata", 0o755)
	require.NoError(t, err)

	content := []byte("foobar\nhello world :)")
	err = ioutil.WriteFile( //nolint:gosec
		filepath.Join("testdata", "TestGet.golden"), content, 0o644,
	)
	require.NoError(t, err)

	got := Get(t)
	assert.Equal(t, content, got)

	tests := []struct {
		name string
		file string
		want []byte
	}{
		{
			name: "",
			file: filepath.Join("testdata", "TestGet", "#00.golden"),
			want: []byte("number double-zero here"),
		},
		{
			name: "foobar",
			file: filepath.Join("testdata", "TestGet", "foobar.golden"),
			want: []byte("foobar here"),
		},
		{
			name: "foo/bar",
			file: filepath.Join("testdata", "TestGet", "foo", "bar.golden"),
			want: []byte("foo/bar style sub-sub-folders works too"),
		},
		{
			name: "john's lost flip-flop",
			file: filepath.Join(
				"testdata", "TestGet", "john's_lost_flip-flop.golden",
			),
			want: []byte("Did John lose his flip-flop again?"),
		},
		{
			name: "thing: it's a thing!",
			file: filepath.Join(
				"testdata", "TestGet", "thing:_it's_a_thing!.golden",
			),
			want: []byte("A thing? Really? Are we getting lazy? :P"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := File(t)
			dir := filepath.Dir(f)

			err := os.MkdirAll(dir, 0o755)
			require.NoError(t, err)

			err = ioutil.WriteFile(f, tt.want, 0o644) //nolint:gosec
			require.NoError(t, err)

			got := Get(t)

			assert.Equal(t, tt.file, f)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSet(t *testing.T) {
	t.Cleanup(func() {
		err := os.RemoveAll(filepath.Join("testdata", t.Name()))
		require.NoError(t, err)
		err = os.Remove(filepath.Join("testdata", t.Name()+".golden"))
		require.NoError(t, err)
	})

	content := []byte("This is the default golden file for TestSet ^_^")
	Set(t, content)

	b, err := ioutil.ReadFile(filepath.Join("testdata", "TestSet.golden"))
	require.NoError(t, err)

	assert.Equal(t, content, b)

	tests := []struct {
		name    string
		file    string
		content []byte
	}{
		{
			name:    "",
			file:    filepath.Join("testdata", "TestSet", "#00.golden"),
			content: []byte("number double-zero strikes again"),
		},
		{
			name:    "foobar",
			file:    filepath.Join("testdata", "TestSet", "foobar.golden"),
			content: []byte("foobar here"),
		},
		{
			name:    "foo/bar",
			file:    filepath.Join("testdata", "TestSet", "foo", "bar.golden"),
			content: []byte("foo/bar style sub-sub-folders works too"),
		},
		{
			name: "john's lost flip-flop",
			file: filepath.Join(
				"testdata", "TestSet", "john's_lost_flip-flop.golden",
			),
			content: []byte("Did John lose his flip-flop again?"),
		},
		{
			name: "thing: it's a thing!",
			file: filepath.Join(
				"testdata", "TestSet", "thing:_it's_a_thing!.golden",
			),
			content: []byte("A thing? Really? Are we getting lazy? :P"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := File(t)

			Set(t, tt.content)

			got, err := ioutil.ReadFile(f)
			require.NoError(t, err)

			assert.Equal(t, tt.file, f)
			assert.Equal(t, tt.content, got)
		})
	}
}

func TestGetNamed(t *testing.T) {
	t.Cleanup(func() {
		err := os.RemoveAll(filepath.Join("testdata", t.Name()))
		require.NoError(t, err)
		err = os.Remove(filepath.Join("testdata", t.Name()+".golden"))
		require.NoError(t, err)
	})

	err := os.MkdirAll(filepath.Join("testdata", "TestGetNamed"), 0o755)
	require.NoError(t, err)

	content := []byte("this is the default golden file for TestGetNamed")
	err = ioutil.WriteFile( //nolint:gosec
		filepath.Join("testdata", "TestGetNamed.golden"), content, 0o644,
	)
	require.NoError(t, err)

	got := GetNamed(t, "")
	assert.Equal(t, content, got)

	content = []byte("this is the named golden file for TestGetNamed")
	err = ioutil.WriteFile( //nolint:gosec
		filepath.Join("testdata", "TestGetNamed", "sub-name.golden"),
		content, 0o644,
	)
	require.NoError(t, err)

	got = GetNamed(t, "sub-name")
	assert.Equal(t, content, got)

	tests := []struct {
		name  string
		named string
		file  string
		want  []byte
	}{
		{
			name: "",
			file: filepath.Join("testdata", "TestGetNamed", "#00.golden"),
			want: []byte("number double-zero here"),
		},
		{
			name:  "",
			named: "sub-zero-one",
			file: filepath.Join(
				"testdata", "TestGetNamed", "#01/sub-zero-one.golden",
			),
			want: []byte("number zero-one here"),
		},
		{
			name:  "foobar",
			named: "email",
			file: filepath.Join(
				"testdata", "TestGetNamed", "foobar/email.golden",
			),
			want: []byte("foobar email here"),
		},
		{
			name:  "foobar",
			named: "json",
			file: filepath.Join(
				"testdata", "TestGetNamed", "foobar#01/json.golden",
			),
			want: []byte("foobar json here"),
		},
		{
			name:  "foo/bar",
			named: "hello/world",
			file: filepath.Join(
				"testdata", "TestGetNamed",
				"foo", "bar",
				"hello", "world.golden",
			),
			want: []byte("foo/bar style sub-sub-folders works too"),
		},
		{
			name:  "john's lost flip-flop",
			named: "left",
			file: filepath.Join(
				"testdata", "TestGetNamed", "john's_lost_flip-flop",
				"left.golden",
			),
			want: []byte("Did John lose his left flip-flop again?"),
		},
		{
			name:  "john's lost flip-flop",
			named: "right",
			file: filepath.Join(
				"testdata", "TestGetNamed", "john's_lost_flip-flop#01",
				"right.golden",
			),
			want: []byte("Did John lose his right flip-flop again?"),
		},
		{
			name:  "thing: it's",
			named: "a thing!",
			file: filepath.Join(
				"testdata", "TestGetNamed", "thing:_it's", "a thing!.golden",
			),
			want: []byte("A thing? Really? Are we getting lazy? :P"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NamedFile(t, tt.named)
			dir := filepath.Dir(f)

			err := os.MkdirAll(dir, 0o755)
			require.NoError(t, err)

			err = ioutil.WriteFile(f, tt.want, 0o644) //nolint:gosec
			require.NoError(t, err)

			got := GetNamed(t, tt.named)

			assert.Equal(t, filepath.FromSlash(tt.file), f)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSetNamed(t *testing.T) {
	t.Cleanup(func() {
		err := os.RemoveAll(filepath.Join("testdata", t.Name()))
		require.NoError(t, err)
		err = os.Remove(filepath.Join("testdata", t.Name()+".golden"))
		require.NoError(t, err)
	})

	content := []byte("This is the default golden file for TestSetNamed ^_^")
	SetNamed(t, "", content)

	b, err := ioutil.ReadFile(filepath.Join("testdata", "TestSetNamed.golden"))
	require.NoError(t, err)

	assert.Equal(t, content, b)

	content = []byte("This is the named golden file for TestSetNamed ^_^")
	SetNamed(t, "sub-name", content)

	b, err = ioutil.ReadFile(
		filepath.Join("testdata", "TestSetNamed", "sub-name.golden"),
	)
	require.NoError(t, err)

	assert.Equal(t, content, b)

	tests := []struct {
		name    string
		named   string
		file    string
		content []byte
	}{
		{
			name:    "",
			file:    filepath.Join("testdata", "TestSetNamed", "#00.golden"),
			content: []byte("number double-zero strikes again"),
		},
		{
			name:  "",
			named: "sub-zero-one",
			file: filepath.Join(
				"testdata", "TestSetNamed", "#01", "sub-zero-one.golden",
			),
			content: []byte("number zero-one sub-zero-one strikes again"),
		},
		{
			name:  "foobar",
			named: "email",
			file: filepath.Join(
				"testdata", "TestSetNamed", "foobar", "email.golden",
			),
			content: []byte("foobar here"),
		},
		{
			name:  "foobar",
			named: "json",
			file: filepath.Join(
				"testdata", "TestSetNamed", "foobar#01", "json.golden",
			),
			content: []byte("foobar here"),
		},
		{
			name: "foo/bar",
			file: filepath.Join(
				"testdata", "TestSetNamed", "foo", "bar.golden",
			),
			content: []byte("foo/bar style sub-sub-folders works too"),
		},
		{
			name:  "john's lost flip-flop",
			named: "left",
			file: filepath.Join(
				"testdata", "TestSetNamed", "john's_lost_flip-flop",
				"left.golden",
			),
			content: []byte("Did John lose his left flip-flop again?"),
		},
		{
			name:  "john's lost flip-flop",
			named: "right",
			file: filepath.Join(
				"testdata", "TestSetNamed", "john's_lost_flip-flop#01",
				"right.golden",
			),
			content: []byte("Did John lose his right flip-flop again?"),
		},
		{
			name:  "thing: it's",
			named: "a thing!",
			file: filepath.Join(
				"testdata", "TestSetNamed", "thing:_it's", "a thing!.golden",
			),
			content: []byte("A thing? Really? Are we getting lazy? :P"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NamedFile(t, tt.named)

			SetNamed(t, tt.named, tt.content)

			got, err := ioutil.ReadFile(f)
			require.NoError(t, err)

			assert.Equal(t, tt.file, f)
			assert.Equal(t, tt.content, got)
		})
	}
}

func TestNamedFile(t *testing.T) {
	got := NamedFile(t, "")
	assert.Equal(t, "testdata/TestNamedFile.golden", got)

	got = NamedFile(t, "sub-name")
	assert.Equal(t, "testdata/TestNamedFile/sub-name.golden", got)

	tests := []struct {
		name  string
		named string
		want  string
	}{
		{
			name:  "",
			named: "",
			want:  "testdata/TestNamedFile/#00.golden",
		},
		{
			name:  "",
			named: "sub-thing",
			want:  "testdata/TestNamedFile/#01/sub-thing.golden",
		},
		{
			name: "foobar",
			want: "testdata/TestNamedFile/foobar.golden",
		},
		{
			name:  "fozbaz",
			named: "email",
			want:  "testdata/TestNamedFile/fozbaz/email.golden",
		},
		{
			name:  "fozbaz",
			named: "json",
			want:  "testdata/TestNamedFile/fozbaz#01/json.golden",
		},
		{
			name:  "foo/bar",
			named: "hello/world",
			want:  "testdata/TestNamedFile/foo/bar/hello/world.golden",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NamedFile(t, tt.named)

			assert.Equal(t, tt.want, got)
		})
	}
}
