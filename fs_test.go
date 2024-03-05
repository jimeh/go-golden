package golden

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMkdirAll(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name    string
		path    string
		perm    os.FileMode
		wantErr bool
	}{
		{"create new dir", "newdir", 0o755, false},
		{"create nested dirs", "nested/dir/structure", 0o755, false},
		{"invalid path", string([]byte{0, 0}), 0o755, true},
	}

	fs := NewFS()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(tempDir, tt.path)
			err := fs.MkdirAll(path, tt.perm)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				_, err := os.Stat(path)
				assert.NoError(t, err)
			}
		})
	}
}

func TestReadFile(t *testing.T) {
	tempDir := t.TempDir()

	sampleFilePath := filepath.Join(tempDir, "sample.txt")
	sampleContent := []byte("Hello, world!")
	err := os.WriteFile(sampleFilePath, sampleContent, 0o600)
	require.NoError(t, err)

	tests := []struct {
		name     string
		filename string
		want     []byte
		wantErr  bool
	}{
		{"read existing file", sampleFilePath, sampleContent, false},
		{"file does not exist", "nonexistent.txt", nil, true},
	}

	fs := NewFS()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fs.ReadFile(tt.filename)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, string(tt.want), string(got))
			}
		})
	}
}

func TestWriteFile(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name     string
		filename string
		data     []byte
		perm     os.FileMode
		wantErr  bool
	}{
		{
			"write to new file",
			"newfile.txt",
			[]byte("new content"),
			0o644,
			false,
		},
		{
			"overwrite existing file",
			"existing.txt",
			[]byte("overwritten content"),
			0o644,
			false,
		},
		{
			"invalid filename",
			string([]byte{0, 0}),
			[]byte("invalid filename"),
			0o644,
			true,
		},
		{
			"non-existent directory",
			"nonexistentdir/newfile.txt",
			[]byte("this will fail"),
			0o644,
			true,
		},
	}

	fs := NewFS()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := filepath.Join(tempDir, tt.filename)
			err := fs.WriteFile(filePath, tt.data, tt.perm)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				content, err := os.ReadFile(filePath)
				assert.NoError(t, err)
				assert.Equal(t, tt.data, content)
			}
		})
	}
}
