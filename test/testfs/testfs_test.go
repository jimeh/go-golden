package testfs

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFSMkdirAll(t *testing.T) {
	type args struct {
		path string
		perm os.FileMode
	}

	tests := []struct {
		name    string
		args    args
		nodes   map[string]*Node
		want    map[string]*Node
		wantErr bool
	}{
		{
			name: "create relative new dir",
			args: args{path: "newdir", perm: 0o755},
			want: map[string]*Node{
				"/root/newdir": {perm: 0o755, isDir: true},
			},
		},
		{
			name: "create absolute new dir",
			args: args{path: "/opt/newdir", perm: 0o755},
			want: map[string]*Node{
				"/opt":        {perm: 0o755, isDir: true},
				"/opt/newdir": {perm: 0o755, isDir: true},
			},
		},

		{
			name: "create relative nested dirs",
			args: args{path: "nested/dir/structure", perm: 0o755},
			want: map[string]*Node{
				"/root/nested":               {perm: 0o755, isDir: true},
				"/root/nested/dir":           {perm: 0o755, isDir: true},
				"/root/nested/dir/structure": {perm: 0o755, isDir: true},
			},
		},
		{
			name: "create absolute nested dirs",
			args: args{path: "/opt/nested/dir/structure", perm: 0o755},
			want: map[string]*Node{
				"/opt":                      {perm: 0o755, isDir: true},
				"/opt/nested":               {perm: 0o755, isDir: true},
				"/opt/nested/dir":           {perm: 0o755, isDir: true},
				"/opt/nested/dir/structure": {perm: 0o755, isDir: true},
			},
		},
		{
			name: "create relative nested dirs with other perms",
			args: args{path: "nested/dir/structure", perm: 0o750},
			want: map[string]*Node{
				"/root/nested":               {perm: 0o750, isDir: true},
				"/root/nested/dir":           {perm: 0o750, isDir: true},
				"/root/nested/dir/structure": {perm: 0o750, isDir: true},
			},
		},
		{
			name: "create absolute nested dirs with other perms",
			args: args{path: "/opt/nested/dir/structure", perm: 0o750},
			want: map[string]*Node{
				"/opt":                      {perm: 0o750, isDir: true},
				"/opt/nested":               {perm: 0o750, isDir: true},
				"/opt/nested/dir":           {perm: 0o750, isDir: true},
				"/opt/nested/dir/structure": {perm: 0o750, isDir: true},
			},
		},
		{
			name: "create relative nested dirs with existing dirs",
			args: args{path: "nested/dir/structure", perm: 0o755},
			want: map[string]*Node{
				"/root/nested":               {perm: 0o755, isDir: true},
				"/root/nested/dir":           {perm: 0o755, isDir: true},
				"/root/nested/dir/structure": {perm: 0o755, isDir: true},
			},
		},
		{
			name: "create absolute nested dirs with existing dirs",
			args: args{path: "/root/nested/dir/structure", perm: 0o755},
			want: map[string]*Node{
				"/root/nested":               {perm: 0o755, isDir: true},
				"/root/nested/dir":           {perm: 0o755, isDir: true},
				"/root/nested/dir/structure": {perm: 0o755, isDir: true},
			},
		},
		{
			name: "create relative under file",
			args: args{path: "file/newdir", perm: 0o755},
			nodes: map[string]*Node{
				"/root/file": {perm: 0o644},
			},
			wantErr: true,
		},
		{
			name: "create absolute under file",
			args: args{path: "/root/file/newdir", perm: 0o755},
			nodes: map[string]*Node{
				"/root/file": {perm: 0o644},
			},
			wantErr: true,
		},
		{
			name: "create relative directory without execute permission",
			args: args{path: "dir/newdir", perm: 0o755},
			nodes: map[string]*Node{
				"/root": {perm: 0o644},
			},
			wantErr: true,
		},
		{
			name: "create absolute directory without execute permission",
			args: args{path: "/root/dir/newdir", perm: 0o755},
			nodes: map[string]*Node{
				"/root": {perm: 0o644},
			},
			wantErr: true,
		},
		{
			name: "create relative directory without write permission",
			args: args{path: "dir/newdir", perm: 0o755},
			nodes: map[string]*Node{
				"/root": {perm: 0o444},
			},
			wantErr: true,
		},
		{
			name: "create absolute directory without write permission",
			args: args{path: "/root/dir/newdir", perm: 0o755},
			nodes: map[string]*Node{
				"/root": {perm: 0o444},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := &FS{
				Pwd: "/root",
				Nodes: map[string]*Node{
					"/":     {perm: 0o755, isDir: true},
					"/root": {perm: 0o700, isDir: true},
				},
			}

			for fp, info := range tt.nodes {
				fs.Nodes[fp] = info
			}

			err := fs.MkdirAll(tt.args.path, tt.args.perm)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				for fp, info := range tt.want {
					got := fs.Nodes[fp]
					assert.Equal(t, info, got, "path: %s", fp)
				}
			}
		})
	}
}

func TestFSReadFile(t *testing.T) {
	type args struct {
		name string
	}

	tests := []struct {
		name    string
		args    args
		nodes   map[string]*Node
		want    []byte
		wantErr bool
	}{
		{
			name: "relative read existing file",
			args: args{name: "file.txt"},
			nodes: map[string]*Node{
				"/":              {perm: 0o755, isDir: true},
				"/root":          {perm: 0o755, isDir: true},
				"/root/file.txt": {data: []byte("file content"), perm: 0o644},
			},
			want: []byte("file content"),
		},
		{
			name: "absolute read existing file",
			args: args{name: "/opt/file.txt"},
			nodes: map[string]*Node{
				"/":             {perm: 0o755, isDir: true},
				"/opt":          {perm: 0o755, isDir: true},
				"/opt/file.txt": {data: []byte("file content"), perm: 0o644},
			},
			want: []byte("file content"),
		},
		{
			name: "relative file does not exist",
			args: args{name: "nonexistent.txt"},
			nodes: map[string]*Node{
				"/":     {perm: 0o755, isDir: true},
				"/root": {perm: 0o755, isDir: true},
			},
			wantErr: true,
		},
		{
			name: "absolute file does not exist",
			args: args{name: "/opt/nonexistent.txt"},
			nodes: map[string]*Node{
				"/":    {perm: 0o755, isDir: true},
				"/opt": {perm: 0o755, isDir: true},
			},
			wantErr: true,
		},
		{
			name: "relative file is a directory",
			args: args{name: "dir"},
			nodes: map[string]*Node{
				"/":         {perm: 0o755, isDir: true},
				"/root":     {perm: 0o755, isDir: true},
				"/root/dir": {perm: 0o755, isDir: true},
			},
			wantErr: true,
		},
		{
			name: "absolute file is a directory",
			args: args{name: "/opt/dir"},
			nodes: map[string]*Node{
				"/":        {perm: 0o755, isDir: true},
				"/opt":     {perm: 0o755, isDir: true},
				"/opt/dir": {perm: 0o755, isDir: true},
			},
			wantErr: true,
		},
		{
			name: "relative file permission denied",
			args: args{name: "file.txt"},
			nodes: map[string]*Node{
				"/":              {perm: 0o755, isDir: true},
				"/root":          {perm: 0o755, isDir: true},
				"/root/file.txt": {data: []byte("file content"), perm: 0o200},
			},
			wantErr: true,
		},
		{
			name: "relative no directory read permission",
			args: args{name: "file.txt"},
			nodes: map[string]*Node{
				"/":              {perm: 0o755, isDir: true},
				"/root":          {perm: 0o355, isDir: true},
				"/root/file.txt": {data: []byte("file content"), perm: 0o644},
			},
			want: []byte("file content"),
		},
		{
			name: "relative no directory execute permission",
			args: args{name: "file.txt"},
			nodes: map[string]*Node{
				"/":              {perm: 0o755, isDir: true},
				"/root":          {perm: 0o655, isDir: true},
				"/root/file.txt": {data: []byte("file content"), perm: 0o200},
			},
			wantErr: true,
		},
		{
			name: "relative no grandparent directory execute permission",
			args: args{name: "foo/file.txt"},
			nodes: map[string]*Node{
				"/":                  {perm: 0o755, isDir: true},
				"/root":              {perm: 0o655, isDir: true},
				"/root/foo":          {perm: 0o755, isDir: true},
				"/root/foo/file.txt": {data: []byte("hello"), perm: 0o200},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := &FS{
				Pwd:   "/root",
				Nodes: tt.nodes,
			}

			got, err := fs.ReadFile(tt.args.name)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestFSWriteFile(t *testing.T) {
	type args struct {
		name string
		data []byte
		perm os.FileMode
	}

	tests := []struct {
		name     string
		args     args
		nodes    map[string]*Node
		wantPath string
		wantErr  bool
	}{
		{
			name: "relative write to new file",
			args: args{
				name: "newfile.txt",
				data: []byte("new content"),
				perm: 0o644,
			},
			wantPath: "/tmp/newfile.txt",
			nodes: map[string]*Node{
				"/":    {perm: 0o755, isDir: true},
				"/tmp": {perm: 0o755, isDir: true},
			},
		},
		{
			name: "absolute write to new file",
			args: args{
				name: "/opt/newfile.txt",
				data: []byte("new content"),
				perm: 0o644,
			},
			wantPath: "/opt/newfile.txt",
			nodes: map[string]*Node{
				"/":    {perm: 0o755, isDir: true},
				"/opt": {perm: 0o755, isDir: true},
			},
		},
		{
			name: "relative overwrite existing file",
			args: args{
				name: "existing.txt",
				data: []byte("overwritten"),
				perm: 0o644,
			},
			wantPath: "/tmp/existing.txt",
			nodes: map[string]*Node{
				"/":             {perm: 0o755, isDir: true},
				"/tmp":          {perm: 0o755, isDir: true},
				"/tmp/existing": {data: []byte("existing"), perm: 0o644},
			},
		},
		{
			name: "absolute overwrite existing file",
			args: args{
				name: "/opt/existing.txt",
				data: []byte("overwritten"),
				perm: 0o644,
			},
			wantPath: "/opt/existing.txt",
			nodes: map[string]*Node{
				"/":             {perm: 0o755, isDir: true},
				"/opt":          {perm: 0o755, isDir: true},
				"/opt/existing": {data: []byte("existing"), perm: 0o644},
			},
		},
		{
			name: "relative overwrite file permissions denied",
			args: args{
				name: "existing.txt",
				data: []byte("overwritten"),
				perm: 0o644,
			},
			wantPath: "/tmp/existing.txt",
			nodes: map[string]*Node{
				"/":                 {perm: 0o755, isDir: true},
				"/tmp":              {perm: 0o755, isDir: true},
				"/tmp/existing.txt": {data: []byte("existing"), perm: 0o400},
			},
			wantErr: true,
		},
		{
			name: "absolute overwrite file permissions denied",
			args: args{
				name: "/opt/existing.txt",
				data: []byte("overwritten"),
				perm: 0o644,
			},
			wantPath: "/opt/existing.txt",
			nodes: map[string]*Node{
				"/":                 {perm: 0o755, isDir: true},
				"/opt":              {perm: 0o755, isDir: true},
				"/opt/existing.txt": {data: []byte("existing"), perm: 0o400},
			},
			wantErr: true,
		},
		{
			name: "relative overwrite directory",
			args: args{
				name: "dir",
				data: []byte("overwritten"),
				perm: 0o644,
			},
			wantPath: "/tmp/dir",
			nodes: map[string]*Node{
				"/":        {perm: 0o755, isDir: true},
				"/tmp":     {perm: 0o755, isDir: true},
				"/tmp/dir": {perm: 0o644, isDir: true},
			},
			wantErr: true,
		},
		{
			name: "absolute overwrite directory",
			args: args{
				name: "/opt/dir",
				data: []byte("overwritten"),
				perm: 0o644,
			},
			wantPath: "/opt/dir",
			nodes: map[string]*Node{
				"/":        {perm: 0o755, isDir: true},
				"/opt":     {perm: 0o755, isDir: true},
				"/opt/dir": {perm: 0o644, isDir: true},
			},
			wantErr: true,
		},
		{
			name: "relative write to non-existent directory",
			args: args{
				name: "nonexistentdir/newfile.txt",
				data: []byte("this will fail"),
				perm: 0o644,
			},
			wantPath: "/tmp/nonexistentdir/newfile.txt",
			nodes: map[string]*Node{
				"/":    {perm: 0o755, isDir: true},
				"/tmp": {perm: 0o755, isDir: true},
			},
			wantErr: true,
		},
		{
			name: "absolute write to non-existent directory",
			args: args{
				name: "/opt/nonexistentdir/newfile.txt",
				data: []byte("this will fail"),
				perm: 0o644,
			},
			wantPath: "/opt/nonexistentdir/newfile.txt",
			nodes: map[string]*Node{
				"/":    {perm: 0o755, isDir: true},
				"/opt": {perm: 0o755, isDir: true},
			},
			wantErr: true,
		},
		{
			name: "relative write parent directory is a file",
			args: args{
				name: "file/newfile.txt",
				data: []byte("this will fail"),
				perm: 0o644,
			},
			wantPath: "/tmp/file/newfile.txt",
			nodes: map[string]*Node{
				"/":         {perm: 0o755, isDir: true},
				"/tmp":      {perm: 0o755, isDir: true},
				"/tmp/file": {data: []byte("file content"), perm: 0o644},
			},
			wantErr: true,
		},
		{
			name: "relative no parent directory write permission denied",
			args: args{
				name: "dir/newfile.txt",
				data: []byte("this will fail"),
				perm: 0o644,
			},
			wantPath: "/tmp/dir/newfile.txt",
			nodes: map[string]*Node{
				"/":        {perm: 0o755, isDir: true},
				"/tmp":     {perm: 0o755, isDir: true},
				"/tmp/dir": {perm: 0o500, isDir: true},
			},
			wantErr: true,
		},
		{
			name: "relative no parent directory execute permission denied",
			args: args{
				name: "dir/newfile.txt",
				data: []byte("this will fail"),
				perm: 0o644,
			},
			wantPath: "/tmp/dir/newfile.txt",
			nodes: map[string]*Node{
				"/":        {perm: 0o755, isDir: true},
				"/tmp":     {perm: 0o755, isDir: true},
				"/tmp/dir": {perm: 0o600, isDir: true},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := &FS{
				Pwd:   "/tmp",
				Nodes: tt.nodes,
			}

			err := fs.WriteFile(tt.args.name, tt.args.data, tt.args.perm)

			if tt.wantErr {
				assert.Error(t, err)
				if _, ok := tt.nodes[tt.wantPath]; ok {
					assert.Equal(t,
						tt.nodes[tt.wantPath],
						fs.Nodes[tt.wantPath],
					)
				} else {
					assert.NotContains(t, fs.Nodes, tt.wantPath)
				}
			} else {
				assert.NoError(t, err)

				got := fs.Nodes[tt.wantPath]
				assert.Equal(t, tt.args.data, got.data)
				assert.Equal(t, tt.args.perm, got.perm)
				assert.Equal(t, false, got.isDir)
			}
		})
	}
}

func TestFSRemove(t *testing.T) {
	type args struct {
		name string
	}

	tests := []struct {
		name    string
		args    args
		nodes   map[string]*Node
		wantErr bool
	}{
		{
			name: "relative remove existing file",
			args: args{name: "file.txt"},
			nodes: map[string]*Node{
				"/":              {perm: 0o755, isDir: true},
				"/root":          {perm: 0o755, isDir: true},
				"/root/file.txt": {data: []byte("file content"), perm: 0o644},
			},
		},
		{
			name: "absolute remove existing file",
			args: args{name: "/opt/file.txt"},
			nodes: map[string]*Node{
				"/":             {perm: 0o755, isDir: true},
				"/opt":          {perm: 0o755, isDir: true},
				"/opt/file.txt": {data: []byte("file content"), perm: 0o644},
			},
		},
		{
			name: "relative file does not exist",
			args: args{name: "nonexistent.txt"},
			nodes: map[string]*Node{
				"/":     {perm: 0o755, isDir: true},
				"/root": {perm: 0o755, isDir: true},
			},
			wantErr: true,
		},
		{
			name: "absolute file does not exist",
			args: args{name: "/opt/nonexistent.txt"},
			nodes: map[string]*Node{
				"/":    {perm: 0o755, isDir: true},
				"/opt": {perm: 0o755, isDir: true},
			},
			wantErr: true,
		},
		{
			name: "relative file is a directory",
			args: args{name: "dir"},
			nodes: map[string]*Node{
				"/":         {perm: 0o755, isDir: true},
				"/root":     {perm: 0o755, isDir: true},
				"/root/dir": {perm: 0o755, isDir: true},
			},
		},
		{
			name: "absolute file is a directory",
			args: args{name: "/opt/dir"},
			nodes: map[string]*Node{
				"/":        {perm: 0o755, isDir: true},
				"/opt":     {perm: 0o755, isDir: true},
				"/opt/dir": {perm: 0o755, isDir: true},
			},
		},
		{
			name: "relative file permission denied",
			args: args{name: "file.txt"},
			nodes: map[string]*Node{
				"/root":          {perm: 0o755, isDir: true},
				"/root/file.txt": {data: []byte("file content"), perm: 0o400},
			},
			wantErr: true,
		},
		{
			name: "absolute file permission denied",
			args: args{name: "/opt/file.txt"},
			nodes: map[string]*Node{
				"/":             {perm: 0o755, isDir: true},
				"/opt":          {perm: 0o755, isDir: true},
				"/opt/file.txt": {data: []byte("file content"), perm: 0o400},
			},
			wantErr: true,
		},
		{
			name: "relative no directory write permission",
			args: args{name: "file.txt"},
			nodes: map[string]*Node{
				"/root":          {perm: 0o555, isDir: true},
				"/root/file.txt": {data: []byte("file content"), perm: 0o644},
			},
			wantErr: true,
		},
		{
			name: "relative no directory execute permission",
			args: args{name: "file.txt"},
			nodes: map[string]*Node{
				"/root":          {perm: 0o655, isDir: true},
				"/root/file.txt": {data: []byte("file content"), perm: 0o644},
			},
			wantErr: true,
		},
		{
			name: "relative no grandparent directory execute permission",
			args: args{name: "file.txt"},
			nodes: map[string]*Node{
				"/root":     {perm: 0o655, isDir: true},
				"/root/dir": {perm: 0o755, isDir: true},
				"/root/dir/file.txt": {
					data: []byte("file content"), perm: 0o644,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := &FS{
				Pwd:   "/root",
				Nodes: tt.nodes,
			}

			err := fs.Remove(tt.args.name)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
