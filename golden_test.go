package golden

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"reflect"
	"runtime"
	"sync"
	"testing"

	"github.com/jimeh/envctl"
	"github.com/jimeh/go-golden/test/testfs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//
// Test Helpers
//

func funcID(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

func prepareDefaultGoldenForTests(t *testing.T) *testfs.FS {
	realDefault := DefaultGolden
	t.Cleanup(func() { DefaultGolden = realDefault })

	fs := testfs.New()
	DefaultGolden = New(WithFS(fs))

	return fs
}

func testInGoroutine(t *testing.T, f func()) {
	t.Helper()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		f()
	}()
	wg.Wait()
}

//
// Tests
//

func TestDefault(t *testing.T) {
	require.IsType(t, &gold{}, DefaultGolden)

	dg := DefaultGolden.(*gold)

	assert.Equal(t, fs.FileMode(0o755), dg.dirMode)
	assert.Equal(t, fs.FileMode(0o644), dg.fileMode)
	assert.Equal(t, ".golden", dg.suffix)
	assert.Equal(t, "testdata", dg.dirname)
	assert.Equal(t, funcID(EnvUpdateFunc), funcID(dg.updateFunc))
	assert.Equal(t, NewFS(), dg.fs)
	assert.Equal(t, true, dg.logOnWrite)
}

func TestDo(t *testing.T) {
	tests := []struct {
		name               string
		testName           string
		content            []byte
		existing           []byte
		wantFilepath       string
		wantNoUpdateLogs   []string
		wantNoUpdateFatals []string
		wantUpdateLogs     []string
		wantUpdateFatals   []string
	}{
		{
			name:     "empty test name",
			testName: "",
			wantUpdateFatals: []string{
				"golden: could not determine filename for TestingT instance",
			},
			wantNoUpdateFatals: []string{
				"golden: could not determine filename for TestingT instance",
			},
		},
		{
			name:         "without slashes",
			testName:     "TestFoo",
			content:      []byte("new content"),
			existing:     []byte("old content"),
			wantFilepath: filepath.Join("testdata", "TestFoo.golden"),
			wantUpdateLogs: []string{
				fmt.Sprintf(
					"golden: writing golden file: %s",
					filepath.Join("testdata", "TestFoo.golden"),
				),
			},
		},
		{
			name:         "with slashes",
			testName:     "TestFoo/bar",
			content:      []byte("new stuff with slashes"),
			existing:     []byte("old stuff"),
			wantFilepath: filepath.Join("testdata", "TestFoo", "bar.golden"),
			wantUpdateLogs: []string{
				fmt.Sprintf(
					"golden: writing golden file: %s",
					filepath.Join("testdata", "TestFoo", "bar.golden"),
				),
			},
		},
		{
			name:     "with spaces and special characters",
			testName: `TestFoo/John's "lost" flip-flop?<>:*|"`,
			content:  []byte("Did John lose his flip-flop again?"),
			existing: []byte("Where is the flip-flop?"),
			wantFilepath: filepath.Join(
				"testdata", "TestFoo", "John's__lost__flip-flop_______.golden",
			),
			wantUpdateLogs: []string{
				fmt.Sprintf(
					"golden: writing golden file: %s",
					filepath.Join(
						"testdata", "TestFoo",
						"John's__lost__flip-flop_______.golden",
					),
				),
			},
		},
		{
			name:         "does not exist",
			testName:     "TestFoo/nope",
			content:      []byte("new stuff with slashes"),
			existing:     nil,
			wantFilepath: filepath.Join("testdata", "TestFoo", "nope.golden"),
			wantUpdateLogs: []string{
				fmt.Sprintf(
					"golden: writing golden file: %s",
					filepath.Join("testdata", "TestFoo", "nope.golden"),
				),
			},
			wantNoUpdateFatals: []string{
				fmt.Sprintf(
					"golden: open %s: no such file or directory",
					filepath.Join(
						"/root", "testdata", "TestFoo", "nope.golden",
					),
				),
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name+"/no update", func(t *testing.T) {
			require.False(t, Update())

			fs := prepareDefaultGoldenForTests(t)
			ft := &fakeTestingT{name: tt.testName}

			if tt.existing != nil {
				err := fs.MkdirAll(filepath.Dir(tt.wantFilepath), 0o755)
				require.NoError(t, err)

				err = fs.WriteFile(tt.wantFilepath, tt.existing, 0o600)
				require.NoError(t, err)
			}

			var got []byte
			testInGoroutine(t, func() {
				got = Do(ft, tt.content)
			})

			if tt.existing == nil {
				assert.Equal(t, tt.existing, got)
			} else {
				assert.GreaterOrEqual(t, 1, len(ft.fatals))
			}

			assert.Equal(t, tt.wantNoUpdateFatals, ft.fatals)
			assert.Equal(t, tt.wantNoUpdateLogs, ft.logs)
		})
		t.Run(tt.name+"/update", func(t *testing.T) {
			envctl.WithClean(map[string]string{"GOLDEN_UPDATE": "1"}, func() {
				require.True(t, Update())

				fs := prepareDefaultGoldenForTests(t)
				ft := &fakeTestingT{name: tt.testName}

				if tt.existing != nil {
					err := fs.MkdirAll(filepath.Dir(tt.wantFilepath), 0o755)
					require.NoError(t, err)

					err = fs.WriteFile(tt.wantFilepath, tt.existing, 0o600)
					require.NoError(t, err)
				}

				var got []byte
				testInGoroutine(t, func() {
					got = Do(ft, tt.content)
				})

				assert.Equal(t, tt.content, got)
				assert.Equal(t, tt.wantUpdateFatals, ft.fatals)
				assert.Equal(t, tt.wantUpdateLogs, ft.logs)
			})
		})
	}
}

func TestDoP(t *testing.T) {
	tests := []struct {
		name               string
		testName           string
		goldenName         string
		content            []byte
		existing           []byte
		wantFilepath       string
		wantNoUpdateLogs   []string
		wantNoUpdateFatals []string
		wantUpdateLogs     []string
		wantUpdateFatals   []string
	}{
		{
			name:       "empty test name",
			testName:   "",
			goldenName: "junk",
			wantUpdateFatals: []string{
				"golden: could not determine filename for TestingT instance",
			},
			wantNoUpdateFatals: []string{
				"golden: could not determine filename for TestingT instance",
			},
		},
		{
			name:       "empty golden name",
			testName:   "TestBar",
			goldenName: "",
			wantUpdateFatals: []string{
				"golden: name cannot be empty",
			},
			wantNoUpdateFatals: []string{
				"golden: name cannot be empty",
			},
		},
		{
			name:         "without slashes",
			testName:     "TestBar",
			goldenName:   "foo",
			content:      []byte("new content"),
			existing:     []byte("old content"),
			wantFilepath: filepath.Join("testdata", "TestBar", "foo.golden"),
			wantUpdateLogs: []string{
				fmt.Sprintf(
					"golden: writing golden file: %s",
					filepath.Join("testdata", "TestBar", "foo.golden"),
				),
			},
		},
		{
			name:       "with slashes in test name",
			testName:   "TestBar/foo",
			goldenName: "junk",
			content:    []byte("new stuff with slashes"),
			existing:   []byte("old stuff"),
			wantFilepath: filepath.Join(
				"testdata", "TestBar", "foo", "junk.golden",
			),
			wantUpdateLogs: []string{
				fmt.Sprintf(
					"golden: writing golden file: %s",
					filepath.Join("testdata", "TestBar", "foo", "junk.golden"),
				),
			},
		},
		{
			name:       "with slashes in golden name",
			testName:   "TestBar",
			goldenName: "foo/junk",
			content:    []byte("new stuff with slashes"),
			existing:   []byte("old stuff"),
			wantFilepath: filepath.Join(
				"testdata", "TestBar", "foo", "junk.golden",
			),
			wantUpdateLogs: []string{
				fmt.Sprintf(
					"golden: writing golden file: %s",
					filepath.Join("testdata", "TestBar", "foo", "junk.golden"),
				),
			},
		},
		{
			name:       "with spaces and special characters",
			testName:   `TestBar/John's "lost" flip-flop?<>:*|"`,
			goldenName: "junk/*plastic*",
			content:    []byte("Did John lose his flip-flop again?"),
			existing:   []byte("Where is the flip-flop?"),
			wantFilepath: filepath.Join(
				"testdata", "TestBar", "John's__lost__flip-flop_______",
				"junk", "_plastic_.golden",
			),
			wantUpdateLogs: []string{
				fmt.Sprintf(
					"golden: writing golden file: %s",
					filepath.Join(
						"testdata", "TestBar", "John's__lost__flip-flop_______",
						"junk", "_plastic_.golden",
					),
				),
			},
		},
		{
			name:         "does not exist",
			testName:     "TestBar",
			goldenName:   "junk",
			content:      []byte("new stuff with slashes"),
			existing:     nil,
			wantFilepath: filepath.Join("testdata", "TestBar", "junk.golden"),
			wantUpdateLogs: []string{
				fmt.Sprintf(
					"golden: writing golden file: %s",
					filepath.Join("testdata", "TestBar", "junk.golden"),
				),
			},
			wantNoUpdateFatals: []string{
				fmt.Sprintf(
					"golden: open %s: no such file or directory",
					filepath.Join(
						"/root", "testdata", "TestBar", "junk.golden",
					),
				),
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name+"/no update", func(t *testing.T) {
			require.False(t, Update())

			fs := prepareDefaultGoldenForTests(t)
			ft := &fakeTestingT{name: tt.testName}

			if tt.existing != nil {
				err := fs.MkdirAll(filepath.Dir(tt.wantFilepath), 0o755)
				require.NoError(t, err)

				err = fs.WriteFile(tt.wantFilepath, tt.existing, 0o600)
				require.NoError(t, err)
			}

			var got []byte
			testInGoroutine(t, func() {
				got = DoP(ft, tt.goldenName, tt.content)
			})

			if tt.existing == nil {
				assert.Equal(t, tt.existing, got)
			} else {
				assert.GreaterOrEqual(t, 1, len(ft.fatals))
			}

			assert.Equal(t, tt.wantNoUpdateFatals, ft.fatals)
			assert.Equal(t, tt.wantNoUpdateLogs, ft.logs)
		})
		t.Run(tt.name+"/update", func(t *testing.T) {
			envctl.WithClean(map[string]string{"GOLDEN_UPDATE": "1"}, func() {
				require.True(t, Update())

				fs := prepareDefaultGoldenForTests(t)
				ft := &fakeTestingT{name: tt.testName}

				if tt.existing != nil {
					err := fs.MkdirAll(filepath.Dir(tt.wantFilepath), 0o755)
					require.NoError(t, err)

					err = fs.WriteFile(tt.wantFilepath, tt.existing, 0o600)
					require.NoError(t, err)
				}

				var got []byte
				testInGoroutine(t, func() {
					got = DoP(ft, tt.goldenName, tt.content)
				})

				assert.Equal(t, tt.content, got)
				assert.Equal(t, tt.wantUpdateFatals, ft.fatals)
				assert.Equal(t, tt.wantUpdateLogs, ft.logs)
			})
		})
	}
}

func TestFile(t *testing.T) {
	tests := []struct {
		name       string
		testName   string
		want       string
		wantFatals []string
	}{
		{
			name:     "empty test name",
			testName: "",
			wantFatals: []string{
				"golden: could not determine filename for TestingT instance",
			},
		},
		{
			name:     "without slashes",
			testName: "TestFoo",
			want:     filepath.Join("testdata", "TestFoo.golden"),
		},
		{
			name:     "with slashes",
			testName: "TestFoo/bar",
			want:     filepath.Join("testdata", "TestFoo", "bar.golden"),
		},
		{
			name:     "with slashes and special characters",
			testName: `TestFoo/John's "lost" flip-flop?<>:*|"`,
			want: filepath.Join(
				"testdata", "TestFoo", "John's__lost__flip-flop_______.golden",
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ft := &fakeTestingT{name: tt.testName}

			var got string
			testInGoroutine(t, func() {
				got = File(ft)
			})

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantFatals, ft.fatals)
		})
	}
}

func TestFileP(t *testing.T) {
	tests := []struct {
		name       string
		testName   string
		goldenName string
		want       string
		wantFatals []string
	}{
		{
			name:       "empty test name",
			testName:   "",
			goldenName: "junk",
			wantFatals: []string{
				"golden: could not determine filename for TestingT instance",
			},
		},
		{
			name:       "empty golden name",
			testName:   "TestFoo",
			goldenName: "",
			wantFatals: []string{
				"golden: name cannot be empty",
			},
		},
		{
			name:       "without slashes",
			testName:   "TestFoo",
			goldenName: "bar",
			want:       filepath.Join("testdata", "TestFoo", "bar.golden"),
		},
		{
			name:       "slashes in test name",
			testName:   "TestFoo/bar",
			goldenName: "junk",
			want: filepath.Join(
				"testdata", "TestFoo", "bar", "junk.golden",
			),
		},
		{
			name:       "slashes in golden name",
			testName:   "TestFoo",
			goldenName: "bar/junk",
			want: filepath.Join(
				"testdata", "TestFoo", "bar", "junk.golden",
			),
		},
		{
			name:       "slashes in test and golden name",
			testName:   "TestFoo/bar",
			goldenName: "junk/plastic",
			want: filepath.Join(
				"testdata", "TestFoo", "bar", "junk", "plastic.golden",
			),
		},
		{
			name:       "slashes and special characters",
			testName:   `TestFoo/John's "lost" flip-flop?<>:*|"`,
			goldenName: "junk/*plastic*",
			want: filepath.Join(
				"testdata", "TestFoo", "John's__lost__flip-flop_______",
				"junk", "_plastic_.golden",
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ft := &fakeTestingT{name: tt.testName}

			var got string
			testInGoroutine(t, func() {
				got = FileP(ft, tt.goldenName)
			})

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantFatals, ft.fatals)
		})
	}
}

func TestGet(t *testing.T) {
	tests := []struct {
		name       string
		testName   string
		files      map[string][]byte
		want       []byte
		wantFatals []string
	}{
		{
			name:     "empty test name",
			testName: "",
			wantFatals: []string{
				"golden: could not determine filename for TestingT instance",
			},
		},
		{
			name:     "without slashes",
			testName: "TestFoo",
			files: map[string][]byte{
				filepath.Join("testdata", "TestFoo.golden"): []byte("bar\n"),
			},
			want: []byte("bar\n"),
		},
		{
			name:     "with slashes",
			testName: "TestFoo/bar",
			files: map[string][]byte{
				filepath.Join("testdata", "TestFoo", "bar.golden"): []byte(
					"bar\n",
				),
			},
			want: []byte("bar\n"),
		},
		{
			name:     "with slashes and special characters",
			testName: `TestFoo/John's "lost" flip-flop?<>:*|"`,
			files: map[string][]byte{
				filepath.Join(
					"testdata", "TestFoo",
					"John's__lost__flip-flop_______.golden",
				): []byte("bar nope\n"),
			},
			want: []byte("bar nope\n"),
		},
		{
			name:     "file does not exist",
			testName: "TestFoo",
			files:    map[string][]byte{},
			wantFatals: []string{
				fmt.Sprintf(
					"golden: open %s: no such file or directory",
					filepath.Join("/root", "testdata", "TestFoo.golden"),
				),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := prepareDefaultGoldenForTests(t)
			ft := &fakeTestingT{name: tt.testName}

			for file, content := range tt.files {
				err := fs.MkdirAll(filepath.Dir(file), 0o755)
				require.NoError(t, err)

				err = fs.WriteFile(file, content, 0o600)
				require.NoError(t, err)
			}

			var got []byte
			testInGoroutine(t, func() {
				got = Get(ft)
			})

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantFatals, ft.fatals)
		})
	}
}

func TestGetP(t *testing.T) {
	tests := []struct {
		name       string
		testName   string
		goldenName string
		files      map[string][]byte
		want       []byte
		wantFatals []string
	}{
		{
			name:       "empty test name",
			testName:   "",
			goldenName: "junk",
			wantFatals: []string{
				"golden: could not determine filename for TestingT instance",
			},
		},
		{
			name:       "empty golden name",
			testName:   "TestBar",
			goldenName: "",
			wantFatals: []string{
				"golden: name cannot be empty",
			},
		},
		{
			name:       "without slashes",
			testName:   "TestBar",
			goldenName: "junk",
			files: map[string][]byte{
				filepath.Join("testdata", "TestBar", "junk.golden"): []byte(
					"foo junk\n",
				),
			},
			want: []byte("foo junk\n"),
		},
		{
			name:       "with slashes in test name",
			testName:   "TestBar/foo",
			goldenName: "junk",
			files: map[string][]byte{
				filepath.Join(
					"testdata", "TestBar", "foo", "junk.golden",
				): []byte("foo\n"),
			},
			want: []byte("foo\n"),
		},
		{
			name:       "with slashes in golden name",
			testName:   "TestBar",
			goldenName: "foo/junk",
			files: map[string][]byte{
				filepath.Join(
					"testdata", "TestBar", "foo", "junk.golden",
				): []byte("foo\n"),
			},
			want: []byte("foo\n"),
		},
		{
			name:       "slashes and special characters",
			testName:   `TestFoo/John's "lost" flip-flop?<>:*|"`,
			goldenName: "junk/*plastic*",
			files: map[string][]byte{
				filepath.Join(
					"testdata", "TestFoo", "John's__lost__flip-flop_______",
					"junk", "_plastic_.golden",
				): []byte("junk here\n"),
			},
			want: []byte("junk here\n"),
		},
		{
			name:       "file does not exist",
			testName:   "TestBar",
			goldenName: "junk",
			files:      map[string][]byte{},
			wantFatals: []string{
				fmt.Sprintf(
					"golden: open %s: no such file or directory",
					filepath.Join(
						"/root", "testdata", "TestBar", "junk.golden",
					),
				),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := prepareDefaultGoldenForTests(t)
			ft := &fakeTestingT{name: tt.testName}

			for file, content := range tt.files {
				err := fs.MkdirAll(filepath.Dir(file), 0o755)
				require.NoError(t, err)

				err = fs.WriteFile(file, content, 0o600)
				require.NoError(t, err)
			}

			var got []byte
			testInGoroutine(t, func() {
				got = GetP(ft, tt.goldenName)
			})

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantFatals, ft.fatals)
		})
	}
}

func TestSet(t *testing.T) {
	tests := []struct {
		name         string
		testName     string
		wantFilepath string
		content      []byte
		wantLogs     []string
		wantFatals   []string
	}{
		{
			name:     "empty test name",
			testName: "",
			wantFatals: []string{
				"golden: could not determine filename for TestingT instance",
			},
		},
		{
			name:         "without slashes",
			testName:     "TestFoo",
			content:      []byte("foobar here"),
			wantFilepath: filepath.Join("testdata", "TestFoo.golden"),
			wantLogs: []string{
				fmt.Sprintf(
					"golden: writing golden file: %s",
					filepath.Join("testdata", "TestFoo.golden"),
				),
			},
		},
		{
			name:         "with slashes",
			testName:     "TestFoo/bar",
			content:      []byte("foo/bar style sub-sub-folders works too"),
			wantFilepath: filepath.Join("testdata", "TestFoo", "bar.golden"),
			wantLogs: []string{
				fmt.Sprintf(
					"golden: writing golden file: %s",
					filepath.Join("testdata", "TestFoo", "bar.golden"),
				),
			},
		},
		{
			name:     "with spaces and special characters",
			testName: `TestFoo/John's "lost" flip-flop?<>:*|"`,
			content:  []byte("Did John lose his flip-flop again?"),
			wantFilepath: filepath.Join(
				"testdata", "TestFoo", "John's__lost__flip-flop_______.golden",
			),
			wantLogs: []string{
				fmt.Sprintf(
					"golden: writing golden file: %s",
					filepath.Join(
						"testdata", "TestFoo",
						"John's__lost__flip-flop_______.golden",
					),
				),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := prepareDefaultGoldenForTests(t)
			ft := &fakeTestingT{name: tt.testName}

			testInGoroutine(t, func() {
				Set(ft, tt.content)
			})

			assert.Equal(t, tt.wantLogs, ft.logs)
			if len(tt.wantFatals) == 0 {
				got, err := fs.ReadFile(tt.wantFilepath)
				require.NoError(t, err)

				assert.Equal(t, tt.content, got)

				filePerms, err := fs.FileMode(tt.wantFilepath)
				require.NoError(t, err)

				dirPerms, err := fs.FileMode(filepath.Dir(tt.wantFilepath))
				require.NoError(t, err)

				assert.Equal(t, filePerms, DefaultFileMode)
				assert.Equal(t, dirPerms, DefaultDirMode)
			} else {
				assert.Equal(t, tt.wantFatals, ft.fatals)
				assert.False(t,
					fs.Exists(tt.wantFilepath),
					"file should not exist",
				)
			}
		})
	}
}

func TestSetP(t *testing.T) {
	tests := []struct {
		name         string
		testName     string
		goldenName   string
		wantFilepath string
		content      []byte
		wantLogs     []string
		wantFatals   []string
	}{
		{
			name:       "empty test name",
			testName:   "",
			goldenName: "junk",
			wantFatals: []string{
				"golden: could not determine filename for TestingT instance",
			},
		},
		{
			name:       "empty golden name",
			testName:   "TestBar",
			goldenName: "",
			wantFatals: []string{
				"golden: name cannot be empty",
			},
		},
		{
			name:         "without slashes",
			testName:     "TestBar",
			goldenName:   "junk",
			content:      []byte("junk here"),
			wantFilepath: filepath.Join("testdata", "TestBar", "junk.golden"),
			wantLogs: []string{
				fmt.Sprintf(
					"golden: writing golden file: %s",
					filepath.Join("testdata", "TestBar", "junk.golden"),
				),
			},
		},
		{
			name:       "with slashes in test name",
			testName:   "TestBar/foo",
			goldenName: "junk",
			content:    []byte("foo/bar style sub-sub-folders works too"),
			wantFilepath: filepath.Join(
				"testdata", "TestBar", "foo", "junk.golden",
			),
			wantLogs: []string{
				fmt.Sprintf(
					"golden: writing golden file: %s",
					filepath.Join("testdata", "TestBar", "foo", "junk.golden"),
				),
			},
		},
		{
			name:       "with slashes in golden name",
			testName:   "TestBar",
			goldenName: "foo/junk",
			content:    []byte("foo/bar style sub-sub-folders works too"),
			wantFilepath: filepath.Join(
				"testdata", "TestBar", "foo", "junk.golden",
			),
			wantLogs: []string{
				fmt.Sprintf(
					"golden: writing golden file: %s",
					filepath.Join("testdata", "TestBar", "foo", "junk.golden"),
				),
			},
		},
		{
			name:       "slashes and special characters",
			testName:   `TestFoo/John's "lost" flip-flop?<>:*|"`,
			goldenName: "junk/*plastic*",
			content:    []byte("Did John lose his flip-flop again?"),
			wantFilepath: filepath.Join(
				"testdata", "TestFoo", "John's__lost__flip-flop_______",
				"junk", "_plastic_.golden",
			),
			wantLogs: []string{
				fmt.Sprintf(
					"golden: writing golden file: %s",
					filepath.Join(
						"testdata", "TestFoo", "John's__lost__flip-flop_______",
						"junk", "_plastic_.golden",
					),
				),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := prepareDefaultGoldenForTests(t)
			ft := &fakeTestingT{name: tt.testName}

			testInGoroutine(t, func() {
				SetP(ft, tt.goldenName, tt.content)
			})

			assert.Equal(t, tt.wantLogs, ft.logs)
			if len(tt.wantFatals) == 0 {
				got, err := fs.ReadFile(tt.wantFilepath)
				require.NoError(t, err)

				assert.Equal(t, tt.content, got)

				filePerms, err := fs.FileMode(tt.wantFilepath)
				require.NoError(t, err)

				dirPerms, err := fs.FileMode(filepath.Dir(tt.wantFilepath))
				require.NoError(t, err)

				assert.Equal(t, filePerms, DefaultFileMode)
				assert.Equal(t, dirPerms, DefaultDirMode)
			} else {
				assert.Equal(t, tt.wantFatals, ft.fatals)
				assert.False(t, fs.Exists(tt.wantFilepath))
			}
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
		want *gold
	}{
		{
			name: "no options",
			args: args{options: nil},
			want: &gold{
				dirMode:    0o755,
				fileMode:   0o644,
				suffix:     ".golden",
				dirname:    "testdata",
				updateFunc: EnvUpdateFunc,
				fs:         NewFS(),
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
					WithFS(testfs.New()),
					WithSilentWrites(),
				},
			},
			want: &gold{
				dirMode:    0o777,
				fileMode:   0o666,
				suffix:     ".gold",
				dirname:    "goldstuff",
				updateFunc: myUpdateFunc,
				fs:         testfs.New(),
				logOnWrite: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := New(tt.args.options...)
			got, ok := g.(*gold)
			require.True(t, ok, "New did not returns a *gold type")

			assert.Equal(t, tt.want.dirMode, got.dirMode)
			assert.Equal(t, tt.want.fileMode, got.fileMode)
			assert.Equal(t, tt.want.suffix, got.suffix)
			assert.Equal(t, tt.want.dirname, got.dirname)
			assert.Equal(t, funcID(tt.want.updateFunc), funcID(got.updateFunc))
			assert.IsType(t, tt.want.fs, got.fs)
			assert.Equal(t, tt.want.logOnWrite, got.logOnWrite)
		})
	}
}
