package golden

import (
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaults(t *testing.T) {
	t.Run("Default", func(t *testing.T) {
		assert.IsType(t, &Golden{}, Default)

		assert.Equal(t, DefaultDirMode, Default.DirMode)
		assert.Equal(t, DefaultFileMode, Default.FileMode)
		assert.Equal(t, DefaultSuffix, Default.Suffix)
		assert.Equal(t, DefaultDirname, Default.Dirname)

		// Use runtime.FuncForPC() to verify the UpdateFunc value is set to
		// the EnvUpdateFunc function by default.
		gotFP := reflect.ValueOf(Default.UpdateFunc).Pointer()
		gotFuncName := runtime.FuncForPC(gotFP).Name()
		wantFP := reflect.ValueOf(EnvUpdateFunc).Pointer()
		wantFuncName := runtime.FuncForPC(wantFP).Name()

		assert.Equal(t, wantFuncName, gotFuncName)
	})

	t.Run("DefaultDirMode", func(t *testing.T) {
		assert.Equal(t, os.FileMode(0o755), DefaultDirMode)
	})

	t.Run("DefaultFileMode", func(t *testing.T) {
		assert.Equal(t, os.FileMode(0o644), DefaultFileMode)
	})

	t.Run("DefaultSuffix", func(t *testing.T) {
		assert.Equal(t, ".golden", DefaultSuffix)
	})

	t.Run("DefaultDirname", func(t *testing.T) {
		assert.Equal(t, "testdata", DefaultDirname)
	})

	t.Run("DefaultUpdateFunc", func(t *testing.T) {
		gotFP := reflect.ValueOf(DefaultUpdateFunc).Pointer()
		gotFuncName := runtime.FuncForPC(gotFP).Name()
		wantFP := reflect.ValueOf(EnvUpdateFunc).Pointer()
		wantFuncName := runtime.FuncForPC(wantFP).Name()
		assert.Equal(t, wantFuncName, gotFuncName)
	})
}

// TestNew is a horribly hack to test that the New() function uses the
// package-level Default* variables.
func TestNew(t *testing.T) {
	// Capture the default values before we change them.
	defaultDirMode := DefaultDirMode
	defaultFileMode := DefaultFileMode
	defaultSuffix := DefaultSuffix
	defaultDirname := DefaultDirname
	defaultUpdateFunc := DefaultUpdateFunc

	// Restore the default values after the test.
	t.Cleanup(func() {
		DefaultDirMode = defaultDirMode
		DefaultFileMode = defaultFileMode
		DefaultSuffix = defaultSuffix
		DefaultDirname = defaultDirname
		DefaultUpdateFunc = defaultUpdateFunc
	})

	// Set all the default values to new values.
	DefaultDirMode = os.FileMode(0o700)
	DefaultFileMode = os.FileMode(0o600)
	DefaultSuffix = ".gold"
	DefaultDirname = "goldenfiles"

	updateFunc := func() bool { return true }
	DefaultUpdateFunc = updateFunc

	// Create a new Golden instance with the new values.
	got := New()

	assert.Equal(t, DefaultDirMode, got.DirMode)
	assert.Equal(t, DefaultFileMode, got.FileMode)
	assert.Equal(t, DefaultSuffix, got.Suffix)
	assert.Equal(t, DefaultDirname, got.Dirname)

	// Verify the UpdateFunc value is set to the new value.
	gotFP := reflect.ValueOf(got.UpdateFunc).Pointer()
	gotFuncName := runtime.FuncForPC(gotFP).Name()
	wantFP := reflect.ValueOf(updateFunc).Pointer()
	wantFuncName := runtime.FuncForPC(wantFP).Name()

	assert.Equal(t, wantFuncName, gotFuncName)
}

func TestDo(t *testing.T) {
	t.Cleanup(func() {
		err := os.RemoveAll(filepath.Join("testdata", "TestDo"))
		require.NoError(t, err)
		err = os.Remove(filepath.Join("testdata", "TestDo.golden"))
		require.NoError(t, err)
	})

	//
	// Test when Update is false
	//
	content := []byte("This is the golden file for TestDo")
	err := os.MkdirAll("testdata", 0o755)
	require.NoError(t, err)

	err = os.WriteFile(
		filepath.Join("testdata", "TestDo.golden"),
		content, 0o600,
	)
	require.NoError(t, err)

	newContent := []byte("This should not be written")
	t.Setenv("GOLDEN_UPDATE", "false")
	got := Do(t, newContent)
	assert.Equal(t, content, got)

	// Verify file wasn't changed
	fileContent, err := os.ReadFile(
		filepath.Join("testdata", "TestDo.golden"),
	)
	require.NoError(t, err)
	assert.Equal(t, content, fileContent)

	//
	// Test when Update is true
	//
	updatedContent := []byte("This is the updated content for TestDo")
	t.Setenv("GOLDEN_UPDATE", "true")
	got = Do(t, updatedContent)
	assert.Equal(t, updatedContent, got)

	// Verify file was updated
	fileContent, err = os.ReadFile(
		filepath.Join("testdata", "TestDo.golden"),
	)
	require.NoError(t, err)
	assert.Equal(t, updatedContent, fileContent)

	//
	// Test with sub-tests
	//
	tests := []struct {
		name    string
		content []byte
	}{
		{
			name:    "simple",
			content: []byte("Simple content for sub-test"),
		},
		{
			name:    "complex/path",
			content: []byte("Complex path content for sub-test"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test with Update true
			t.Setenv("GOLDEN_UPDATE", "true")
			got := Do(t, tt.content)
			assert.Equal(t, tt.content, got)

			// Verify file was written with correct content
			f := File(t)
			fileContent, err := os.ReadFile(f)
			require.NoError(t, err)
			assert.Equal(t, tt.content, fileContent)

			// Test with Update false
			t.Setenv("GOLDEN_UPDATE", "false")

			newContent := []byte(
				"This should not be written in sub-test",
			)
			got = Do(t, newContent)
			assert.Equal(t, tt.content, got)

			// Verify file wasn't changed
			f = File(t)
			fileContent, err = os.ReadFile(f)
			require.NoError(t, err)
			assert.Equal(t, tt.content, fileContent)
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
	err = os.WriteFile(
		filepath.Join("testdata", "TestGet.golden"), content, 0o600,
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

			err = os.WriteFile(f, tt.want, 0o600)
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

	b, err := os.ReadFile(filepath.Join("testdata", "TestSet.golden"))
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

			got, err := os.ReadFile(f)
			require.NoError(t, err)

			assert.Equal(t, tt.file, f)
			assert.Equal(t, tt.content, got)
		})
	}
}

func TestDoP(t *testing.T) {
	t.Cleanup(func() {
		err := os.RemoveAll(filepath.Join("testdata", "TestDoP"))
		require.NoError(t, err)
	})

	//
	// Test when Update is false
	//
	name := "test-format"
	content := []byte("This is the golden file for TestDoP")
	err := os.MkdirAll(filepath.Join("testdata", "TestDoP"), 0o755)
	require.NoError(t, err)

	goldenFile := filepath.Join("testdata", "TestDoP", name+".golden")
	err = os.WriteFile(goldenFile, content, 0o600)
	require.NoError(t, err)

	newContent := []byte("This should not be written")
	t.Setenv("GOLDEN_UPDATE", "false")
	got := DoP(t, name, newContent)
	assert.Equal(t, content, got)

	// Verify file wasn't changed
	fileContent, err := os.ReadFile(goldenFile)
	require.NoError(t, err)
	assert.Equal(t, content, fileContent)

	//
	// Test when Update is true
	//
	updatedContent := []byte("This is the updated content for TestDoP")
	t.Setenv("GOLDEN_UPDATE", "true")
	got = DoP(t, name, updatedContent)
	assert.Equal(t, updatedContent, got)

	// Verify file was updated
	fileContent, err = os.ReadFile(goldenFile)
	require.NoError(t, err)
	assert.Equal(t, updatedContent, fileContent)

	//
	// Test with sub-tests
	//
	tests := []struct {
		testName string
		name     string
		content  []byte
	}{
		{
			testName: "json format",
			name:     "json",
			content:  []byte(`{"key": "value"}`),
		},
		{
			testName: "xml format",
			name:     "xml",
			content:  []byte(`<root><key>value</key></root>`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			// Test with Update true
			t.Setenv("GOLDEN_UPDATE", "true")
			got := DoP(t, tt.name, tt.content)
			assert.Equal(t, tt.content, got)

			// Verify file was written with correct content
			f := FileP(t, tt.name)
			fileContent, err := os.ReadFile(f)
			require.NoError(t, err)
			assert.Equal(t, tt.content, fileContent)

			// Test with Update false
			t.Setenv("GOLDEN_UPDATE", "false")
			newContent := []byte(
				"This should not be written in sub-test",
			)
			got = DoP(t, tt.name, newContent)
			assert.Equal(t, tt.content, got)

			// Verify file wasn't changed
			f = FileP(t, tt.name)
			fileContent, err = os.ReadFile(f)
			require.NoError(t, err)
			assert.Equal(t, tt.content, fileContent)
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
	err = os.WriteFile(
		filepath.Join("testdata", "TestGetP", "sub-name.golden"),
		content, 0o600,
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

			err = os.WriteFile(f, tt.want, 0o600)
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

	b, err := os.ReadFile(
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

			got, err := os.ReadFile(f)
			require.NoError(t, err)

			assert.Equal(t, tt.file, f)
			assert.Equal(t, tt.content, got)
		})
	}
}

func TestUpdate(t *testing.T) {
	for _, tt := range envUpdateFuncTestCases {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.env {
				t.Setenv(k, v)
			}

			got := Update()

			assert.Equal(t, tt.want, got)
		})
	}
}
