package golden

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/jimeh/envctl"
	"github.com/jimeh/go-mocktesting"
	"github.com/spf13/afero"
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
		t.Log("cleaning up golden files")
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
		t.Log("cleaning up golden files")
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

func TestNew(t *testing.T) {
	myUpdateFunc := func() bool { return false }

	type args struct {
		options []Option
	}
	tests := []struct {
		name string
		args args
		want *golden
	}{
		{
			name: "no options",
			args: args{options: nil},
			want: &golden{
				dirMode:    0o755,
				fileMode:   0o644,
				suffix:     ".golden",
				dirname:    "testdata",
				updateFunc: EnvUpdateFunc,
				fs:         afero.NewOsFs(),
				logOnWrite: true,
			},
		},
		{
			name: "all options",
			args: args{
				options: []Option{
					WithDirMode(0o777),
					WithFileMode(0o666),
					WithSuffix(".gold"),
					WithDirname("goldstuff"),
					WithUpdateFunc(myUpdateFunc),
					WithFs(afero.NewMemMapFs()),
					WithSilentWrites(),
				},
			},
			want: &golden{
				dirMode:    0o777,
				fileMode:   0o666,
				suffix:     ".gold",
				dirname:    "goldstuff",
				updateFunc: myUpdateFunc,
				fs:         afero.NewMemMapFs(),
				logOnWrite: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := New(tt.args.options...)
			got, ok := g.(*golden)
			require.True(t, ok, "New did not returns a *golden instance")

			gotUpdateFunc := runtime.FuncForPC(
				reflect.ValueOf(got.updateFunc).Pointer(),
			).Name()
			wantUpdateFunc := runtime.FuncForPC(
				reflect.ValueOf(tt.want.updateFunc).Pointer(),
			).Name()

			assert.Equal(t, tt.want.dirMode, got.dirMode)
			assert.Equal(t, tt.want.fileMode, got.fileMode)
			assert.Equal(t, tt.want.suffix, got.suffix)
			assert.Equal(t, tt.want.dirname, got.dirname)
			assert.Equal(t, tt.want.logOnWrite, got.logOnWrite)
			assert.Equal(t, wantUpdateFunc, gotUpdateFunc)
			assert.IsType(t, tt.want.fs, got.fs)
		})
	}
}

func Test_golden_File(t *testing.T) {
	type fields struct {
		suffix  *string
		dirname *string
	}
	tests := []struct {
		name           string
		testName       string
		fields         fields
		want           string
		wantAborted    bool
		wantFailCount  int
		wantTestOutput []string
	}{
		{
			name:     "top-level",
			testName: "TestFooBar",
			want:     filepath.Join("testdata", "TestFooBar.golden"),
		},
		{
			name:     "sub-test",
			testName: "TestFooBar/it_is_here",
			want: filepath.Join(
				"testdata", "TestFooBar", "it_is_here.golden",
			),
		},
		{
			name:          "blank test name",
			testName:      "",
			wantAborted:   true,
			wantFailCount: 1,
			wantTestOutput: []string{
				"golden: could not determine filename for given " +
					"*mocktesting.T instance\n",
			},
		},
		{
			name:     "custom dirname",
			testName: "TestFozBar",
			fields: fields{
				dirname: stringPtr("goldenfiles"),
			},
			want: filepath.Join("goldenfiles", "TestFozBar.golden"),
		},
		{
			name:     "custom suffix",
			testName: "TestFozBaz",
			fields: fields{
				suffix: stringPtr(".goldfile"),
			},
			want: filepath.Join("testdata", "TestFozBaz.goldfile"),
		},
		{
			name:     "custom dirname and suffix",
			testName: "TestFozBar",
			fields: fields{
				dirname: stringPtr("goldenfiles"),
				suffix:  stringPtr(".goldfile"),
			},
			want: filepath.Join("goldenfiles", "TestFozBar.goldfile"),
		},
		{
			name:     "invalid chars in test name",
			testName: `TestFooBar/foo?<>:*|"bar`,
			want: filepath.Join(
				"testdata", "TestFooBar", "foo_______bar.golden",
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.fields.suffix == nil {
				tt.fields.suffix = stringPtr(".golden")
			}
			if tt.fields.dirname == nil {
				tt.fields.dirname = stringPtr("testdata")
			}

			g := &golden{
				suffix:  *tt.fields.suffix,
				dirname: *tt.fields.dirname,
			}

			mt := mocktesting.NewT(tt.testName)

			var got string
			mocktesting.Go(func() {
				got = g.File(mt)
			})

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantAborted, mt.Aborted(), "aborted")
			assert.Equal(t,
				tt.wantFailCount, mt.FailedCount(), "failed count",
			)
			assert.Equal(t, tt.wantTestOutput, mt.Output(), "test output")
		})
	}
}

func Test_golden_Get(t *testing.T) {
	type fields struct {
		suffix  *string
		dirname *string
	}
	tests := []struct {
		name           string
		testName       string
		fields         fields
		files          map[string][]byte
		want           []byte
		wantAborted    bool
		wantFailCount  int
		wantTestOutput []string
	}{
		{
			name:     "file exists",
			testName: "TestFooBar",
			files: map[string][]byte{
				filepath.Join("testdata", "TestFooBar.golden"): []byte(
					"foo: bar\nhello: world",
				),
			},
			want: []byte("foo: bar\nhello: world"),
		},
		{
			name:          "file is missing",
			testName:      "TestFooBar",
			files:         map[string][]byte{},
			wantAborted:   true,
			wantFailCount: 1,
			wantTestOutput: []string{
				"golden: open " + filepath.Join(
					"testdata", "TestFooBar.golden",
				) + ": file does not exist\n",
			},
		},
		{
			name:     "sub-test file exists",
			testName: "TestFooBar/it_is_here",
			files: map[string][]byte{
				filepath.Join(
					"testdata", "TestFooBar", "it_is_here.golden",
				): []byte("this is really here ^_^\n"),
			},
			want: []byte("this is really here ^_^\n"),
		},
		{
			name:          "sub-test file is missing",
			testName:      "TestFooBar/not_really_here",
			files:         map[string][]byte{},
			wantAborted:   true,
			wantFailCount: 1,
			wantTestOutput: []string{
				"golden: open " + filepath.Join(
					"testdata", "TestFooBar", "not_really_here.golden",
				) + ": file does not exist\n",
			},
		},
		{
			name:          "blank test name",
			testName:      "",
			wantAborted:   true,
			wantFailCount: 1,
			wantTestOutput: []string{
				"golden: could not determine filename for given " +
					"*mocktesting.T instance\n",
			},
		},
		{
			name:     "custom dirname",
			testName: "TestFozBar",
			fields: fields{
				dirname: stringPtr("goldenfiles"),
			},
			files: map[string][]byte{
				filepath.Join("goldenfiles", "TestFozBar.golden"): []byte(
					"foo: bar\nhello: world",
				),
			},
			want: []byte("foo: bar\nhello: world"),
		},
		{
			name:     "custom suffix",
			testName: "TestFozBaz",
			fields: fields{
				suffix: stringPtr(".goldfile"),
			},
			files: map[string][]byte{
				filepath.Join("testdata", "TestFozBaz.goldfile"): []byte(
					"foo: bar\nhello: world",
				),
			},
			want: []byte("foo: bar\nhello: world"),
		},
		{
			name:     "custom dirname and suffix",
			testName: "TestFozBar",
			fields: fields{
				dirname: stringPtr("goldenfiles"),
				suffix:  stringPtr(".goldfile"),
			},
			files: map[string][]byte{
				filepath.Join("goldenfiles", "TestFozBar.goldfile"): []byte(
					"foo: bar\nhello: world",
				),
			},
			want: []byte("foo: bar\nhello: world"),
		},
		{
			name:     "invalid chars in test name",
			testName: `TestFooBar/foo?<>:*|"bar`,
			files: map[string][]byte{
				filepath.Join(
					"testdata", "TestFooBar", "foo_______bar.golden",
				): []byte("foo: bar\nhello: world"),
			},
			want: []byte("foo: bar\nhello: world"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			for f, b := range tt.files {
				_ = afero.WriteFile(fs, f, b, 0o644)
			}

			if tt.fields.suffix == nil {
				tt.fields.suffix = stringPtr(".golden")
			}
			if tt.fields.dirname == nil {
				tt.fields.dirname = stringPtr("testdata")
			}

			g := &golden{
				suffix:  *tt.fields.suffix,
				dirname: *tt.fields.dirname,
				fs:      fs,
			}

			mt := mocktesting.NewT(tt.testName)

			var got []byte
			mocktesting.Go(func() {
				got = g.Get(mt)
			})

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantAborted, mt.Aborted(), "aborted")
			assert.Equal(t,
				tt.wantFailCount, mt.FailedCount(), "failed count",
			)
			assert.Equal(t, tt.wantTestOutput, mt.Output(), "test output")
		})
	}
}

func Test_golden_FileP(t *testing.T) {
	type args struct {
		name string
	}
	type fields struct {
		suffix  *string
		dirname *string
	}
	tests := []struct {
		name           string
		testName       string
		args           args
		fields         fields
		want           string
		wantAborted    bool
		wantFailCount  int
		wantTestOutput []string
	}{
		{
			name:     "top-level",
			testName: "TestFooBar",
			args:     args{name: "yaml"},
			want:     filepath.Join("testdata", "TestFooBar", "yaml.golden"),
		},
		{
			name:     "sub-test",
			testName: "TestFooBar/it_is_here",
			args:     args{name: "json"},
			want: filepath.Join(
				"testdata", "TestFooBar", "it_is_here", "json.golden",
			),
		},
		{
			name:          "blank test name",
			testName:      "",
			args:          args{name: "json"},
			wantAborted:   true,
			wantFailCount: 1,
			wantTestOutput: []string{
				"golden: could not determine filename for given " +
					"*mocktesting.T instance\n",
			},
		},
		{
			name:     "custom dirname",
			testName: "TestFozBar",
			args:     args{name: "xml"},
			fields: fields{
				dirname: stringPtr("goldenfiles"),
			},
			want: filepath.Join("goldenfiles", "TestFozBar", "xml.golden"),
		},
		{
			name:     "custom suffix",
			testName: "TestFozBaz",
			args:     args{name: "toml"},
			fields: fields{
				suffix: stringPtr(".goldfile"),
			},
			want: filepath.Join("testdata", "TestFozBaz", "toml.goldfile"),
		},
		{
			name:     "custom dirname and suffix",
			testName: "TestFozBar",
			args:     args{name: "json"},
			fields: fields{
				dirname: stringPtr("goldenfiles"),
				suffix:  stringPtr(".goldfile"),
			},
			want: filepath.Join("goldenfiles", "TestFozBar", "json.goldfile"),
		},
		{
			name:     "invalid chars in test name",
			testName: `TestFooBar/foo?<>:*|"bar`,
			args:     args{name: "yml"},
			want: filepath.Join(
				"testdata", "TestFooBar", "foo_______bar", "yml.golden",
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.fields.suffix == nil {
				tt.fields.suffix = stringPtr(".golden")
			}
			if tt.fields.dirname == nil {
				tt.fields.dirname = stringPtr("testdata")
			}

			g := &golden{
				suffix:  *tt.fields.suffix,
				dirname: *tt.fields.dirname,
			}

			mt := mocktesting.NewT(tt.testName)

			var got string
			mocktesting.Go(func() {
				got = g.FileP(mt, tt.args.name)
			})

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantAborted, mt.Aborted(), "aborted")
			assert.Equal(t,
				tt.wantFailCount, mt.FailedCount(), "failed count",
			)
			assert.Equal(t, tt.wantTestOutput, mt.Output(), "test output")
		})
	}
}

func Test_golden_GetP(t *testing.T) {
	type args struct {
		name string
	}
	type fields struct {
		suffix  *string
		dirname *string
	}
	tests := []struct {
		name           string
		testName       string
		args           args
		fields         fields
		files          map[string][]byte
		want           []byte
		wantAborted    bool
		wantFailCount  int
		wantTestOutput []string
	}{
		{
			name:     "file exists",
			testName: "TestFooBar",
			args:     args{name: "yaml"},
			files: map[string][]byte{
				filepath.Join("testdata", "TestFooBar", "yaml.golden"): []byte(
					"foo: bar\nhello: world",
				),
			},
			want: []byte("foo: bar\nhello: world"),
		},
		{
			name:          "file is missing",
			testName:      "TestFooBar",
			args:          args{name: "yaml"},
			files:         map[string][]byte{},
			wantAborted:   true,
			wantFailCount: 1,
			wantTestOutput: []string{
				"golden: open " + filepath.Join(
					"testdata", "TestFooBar", "yaml.golden",
				) + ": file does not exist\n",
			},
		},
		{
			name:     "sub-test file exists",
			testName: "TestFooBar/it_is_here",
			args:     args{name: "plain"},
			files: map[string][]byte{
				filepath.Join(
					"testdata", "TestFooBar", "it_is_here", "plain.golden",
				): []byte("this is really here ^_^\n"),
			},
			want: []byte("this is really here ^_^\n"),
		},
		{
			name:          "sub-test file is missing",
			testName:      "TestFooBar/not_really_here",
			args:          args{name: "plain"},
			files:         map[string][]byte{},
			wantAborted:   true,
			wantFailCount: 1,
			wantTestOutput: []string{
				"golden: open " + filepath.Join(
					"testdata", "TestFooBar", "not_really_here", "plain.golden",
				) + ": file does not exist\n",
			},
		},
		{
			name:          "blank test name",
			testName:      "",
			args:          args{name: "plain"},
			wantAborted:   true,
			wantFailCount: 1,
			wantTestOutput: []string{
				"golden: could not determine filename for given " +
					"*mocktesting.T instance\n",
			},
		},
		{
			name:          "blank name",
			testName:      "TestFooBar",
			args:          args{name: ""},
			wantAborted:   true,
			wantFailCount: 1,
			wantTestOutput: []string{
				"golden: name cannot be empty\n",
			},
		},
		{
			name:     "custom dirname",
			testName: "TestFozBar",
			args:     args{name: "yaml"},
			fields: fields{
				dirname: stringPtr("goldenfiles"),
			},
			files: map[string][]byte{
				filepath.Join(
					"goldenfiles", "TestFozBar", "yaml.golden",
				): []byte("foo: bar\nhello: world"),
			},
			want: []byte("foo: bar\nhello: world"),
		},
		{
			name:     "custom suffix",
			testName: "TestFozBaz",
			args:     args{name: "yaml"},
			fields: fields{
				suffix: stringPtr(".goldfile"),
			},
			files: map[string][]byte{
				filepath.Join(
					"testdata", "TestFozBaz", "yaml.goldfile",
				): []byte("foo: bar\nhello: world"),
			},
			want: []byte("foo: bar\nhello: world"),
		},
		{
			name:     "custom dirname and suffix",
			testName: "TestFozBar",
			args:     args{name: "yaml"},
			fields: fields{
				dirname: stringPtr("goldenfiles"),
				suffix:  stringPtr(".goldfile"),
			},
			files: map[string][]byte{
				filepath.Join(
					"goldenfiles", "TestFozBar", "yaml.goldfile",
				): []byte("foo: bar\nhello: world"),
			},
			want: []byte("foo: bar\nhello: world"),
		},
		{
			name:     "invalid chars in test name",
			testName: `TestFooBar/foo?<>:*|"bar`,
			args:     args{name: "trash"},
			files: map[string][]byte{
				filepath.Join(
					"testdata", "TestFooBar", "foo_______bar", "trash.golden",
				): []byte("foo: bar\nhello: world"),
			},
			want: []byte("foo: bar\nhello: world"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			for f, b := range tt.files {
				_ = afero.WriteFile(fs, f, b, 0o644)
			}

			if tt.fields.suffix == nil {
				tt.fields.suffix = stringPtr(".golden")
			}
			if tt.fields.dirname == nil {
				tt.fields.dirname = stringPtr("testdata")
			}

			g := &golden{
				suffix:  *tt.fields.suffix,
				dirname: *tt.fields.dirname,
				fs:      fs,
			}

			mt := mocktesting.NewT(tt.testName)

			var got []byte
			mocktesting.Go(func() {
				got = g.GetP(mt, tt.args.name)
			})

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantAborted, mt.Aborted(), "aborted")
			assert.Equal(t,
				tt.wantFailCount, mt.FailedCount(), "failed count",
			)
			assert.Equal(t, tt.wantTestOutput, mt.Output(), "test output")
		})
	}
}
