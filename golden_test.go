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
			want: filepath.Join("testdata", "TestFile", "foo", "bar.golden"),
		},
		{
			name: `"foobar"`,
			want: filepath.Join("testdata", "TestFile", "_foobar_.golden"),
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
		err := os.RemoveAll(filepath.Join("testdata", "TestGet"))
		require.NoError(t, err)
		err = os.Remove(filepath.Join("testdata", "TestGet.golden"))
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
				"testdata", "TestGet", "thing__it's_a_thing!.golden",
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
		err := os.RemoveAll(filepath.Join("testdata", "TestSet"))
		require.NoError(t, err)
		err = os.Remove(filepath.Join("testdata", "TestSet.golden"))
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
				"testdata", "TestSet", "thing__it's_a_thing!.golden",
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

func TestFileP(t *testing.T) {
	got := FileP(t, "sub-name")
	assert.Equal(t,
		filepath.Join("testdata", "TestFileP", "sub-name.golden"), got,
	)

	tests := []struct {
		name  string
		named string
		want  string
	}{
		{
			name:  "",
			named: "sub-thing",
			want: filepath.Join(
				"testdata", "TestFileP", "#00", "sub-thing.golden",
			),
		},
		{
			name:  "fozbaz",
			named: "email",
			want: filepath.Join(
				"testdata", "TestFileP", "fozbaz", "email.golden",
			),
		},
		{
			name:  "fozbaz",
			named: "json",
			want: filepath.Join(
				"testdata", "TestFileP", "fozbaz#01", "json.golden",
			),
		},
		{
			name:  "foo/bar",
			named: "hello/world",
			want: filepath.Join(
				"testdata", "TestFileP",
				"foo", "bar",
				"hello", "world.golden",
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FileP(t, tt.named)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetP(t *testing.T) {
	t.Cleanup(func() {
		err := os.RemoveAll(filepath.Join("testdata", "TestGetP"))
		require.NoError(t, err)
	})

	err := os.MkdirAll(filepath.Join("testdata", "TestGetP"), 0o755)
	require.NoError(t, err)

	content := []byte("this is the named golden file for TestGetP")
	err = ioutil.WriteFile( //nolint:gosec
		filepath.Join("testdata", "TestGetP", "sub-name.golden"),
		content, 0o644,
	)
	require.NoError(t, err)

	got := GetP(t, "sub-name")
	assert.Equal(t, content, got)

	tests := []struct {
		name  string
		named string
		file  string
		want  []byte
	}{
		{
			name:  "",
			named: "sub-zero-one",
			file: filepath.Join(
				"testdata", "TestGetP", "#00", "sub-zero-one.golden",
			),
			want: []byte("number zero-one here"),
		},
		{
			name:  "foobar",
			named: "email",
			file: filepath.Join(
				"testdata", "TestGetP", "foobar", "email.golden",
			),
			want: []byte("foobar email here"),
		},
		{
			name:  "foobar",
			named: "json",
			file: filepath.Join(
				"testdata", "TestGetP", "foobar#01", "json.golden",
			),
			want: []byte("foobar json here"),
		},
		{
			name:  "foo/bar",
			named: "hello/world",
			file: filepath.Join(
				"testdata", "TestGetP",
				"foo", "bar",
				"hello", "world.golden",
			),
			want: []byte("foo/bar style sub-sub-folders works too"),
		},
		{
			name:  "john's lost flip-flop",
			named: "left",
			file: filepath.Join(
				"testdata", "TestGetP", "john's_lost_flip-flop",
				"left.golden",
			),
			want: []byte("Did John lose his left flip-flop again?"),
		},
		{
			name:  "john's lost flip-flop",
			named: "right",
			file: filepath.Join(
				"testdata", "TestGetP", "john's_lost_flip-flop#01",
				"right.golden",
			),
			want: []byte("Did John lose his right flip-flop again?"),
		},
		{
			name:  "thing: it's",
			named: "a thing!",
			file: filepath.Join(
				"testdata", "TestGetP", "thing__it's", "a_thing!.golden",
			),
			want: []byte("A thing? Really? Are we getting lazy? :P"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := FileP(t, tt.named)
			dir := filepath.Dir(f)

			err := os.MkdirAll(dir, 0o755)
			require.NoError(t, err)

			err = ioutil.WriteFile(f, tt.want, 0o644) //nolint:gosec
			require.NoError(t, err)

			got := GetP(t, tt.named)

			assert.Equal(t, filepath.FromSlash(tt.file), f)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSetP(t *testing.T) {
	t.Cleanup(func() {
		err := os.RemoveAll(filepath.Join("testdata", "TestSetP"))
		require.NoError(t, err)
	})

	content := []byte("This is the named golden file for TestSetP ^_^")
	SetP(t, "sub-name", content)

	b, err := ioutil.ReadFile(
		filepath.Join("testdata", "TestSetP", "sub-name.golden"),
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
			name:  "",
			named: "sub-zero-one",
			file: filepath.Join(
				"testdata", "TestSetP", "#00", "sub-zero-one.golden",
			),
			content: []byte("number zero-one sub-zero-one strikes again"),
		},
		{
			name:  "foobar",
			named: "email",
			file: filepath.Join(
				"testdata", "TestSetP", "foobar", "email.golden",
			),
			content: []byte("foobar here"),
		},
		{
			name:  "foobar",
			named: "json",
			file: filepath.Join(
				"testdata", "TestSetP", "foobar#01", "json.golden",
			),
			content: []byte("foobar here"),
		},
		{
			name:  "john's lost flip-flop",
			named: "left",
			file: filepath.Join(
				"testdata", "TestSetP", "john's_lost_flip-flop",
				"left.golden",
			),
			content: []byte("Did John lose his left flip-flop again?"),
		},
		{
			name:  "john's lost flip-flop",
			named: "right",
			file: filepath.Join(
				"testdata", "TestSetP", "john's_lost_flip-flop#01",
				"right.golden",
			),
			content: []byte("Did John lose his right flip-flop again?"),
		},
		{
			name:  "thing: it's",
			named: "a thing!",
			file: filepath.Join(
				"testdata", "TestSetP", "thing__it's", "a_thing!.golden",
			),
			content: []byte("A thing? Really? Are we getting lazy? :P"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := FileP(t, tt.named)

			SetP(t, tt.named, tt.content)

			got, err := ioutil.ReadFile(f)
			require.NoError(t, err)

			assert.Equal(t, tt.file, f)
			assert.Equal(t, tt.content, got)
		})
	}
}

func TestUpdate(t *testing.T) {
	for _, tt := range envUpdateFuncTestCases {
		t.Run(tt.name, func(t *testing.T) {
			envctl.WithClean(tt.env, func() {
				got := Update()

				assert.Equal(t, tt.want, got)
			})
		})
	}
}
