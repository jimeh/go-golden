package testfs

import (
	"errors"
	"os"
	"path"
	"strings"
)

type Node struct {
	data  []byte
	perm  os.FileMode
	isDir bool
}

type FS struct {
	Pwd   string
	Nodes map[string]*Node
}

func New() *FS {
	return &FS{
		Pwd: "/root",
		Nodes: map[string]*Node{
			"/":     {perm: 0o755, isDir: true},
			"/root": {perm: 0o700, isDir: true},
		},
	}
}

func (fs *FS) MkdirAll(name string, perm os.FileMode) error {
	if !path.IsAbs(name) && name != "" {
		name = path.Join(fs.Pwd, name)
	}

	dirs := []string{name}
	for d := path.Dir(name); d != "/"; d = path.Dir(d) {
		dirs = append(dirs, d)
	}
	dirs = append(dirs, "/")

	for i := len(dirs) - 1; i >= 0; i-- {
		dir := dirs[i]
		parent := path.Dir(dir)

		if info, ok := fs.Nodes[dir]; ok {
			if !info.isDir {
				return &os.PathError{
					Op:   "mkdir",
					Path: dir,
					Err:  errors.New("not a directory"),
				}
			}

			continue
		}

		parentInfo, ok := fs.Nodes[parent]
		if !ok {
			return &os.PathError{
				Op:   "mkdir",
				Path: parent,
				Err:  errors.New("no such file or directory"),
			}
		}
		if !parentInfo.isDir {
			return &os.PathError{
				Op:   "mkdir",
				Path: parent,
				Err:  errors.New("not a directory"),
			}
		}
		// Ensure all parent directories have execute permissions, and direct
		// parent also has write permission.
		if parentInfo.perm&0o100 == 0 || i == 1 && parentInfo.perm&0o200 == 0 {
			return &os.PathError{
				Op:   "mkdir",
				Path: dir,
				Err:  errors.New("permission denied"),
			}
		}

		fs.Nodes[dir] = &Node{perm: perm, isDir: true}
	}

	return nil
}

func (fs *FS) ReadFile(name string) ([]byte, error) {
	if !path.IsAbs(name) && name != "" {
		name = path.Join(fs.Pwd, name)
	}

	_, err := fs.checkParents(name, false)
	if err != nil {
		return nil, err
	}

	info, ok := fs.Nodes[name]
	if !ok {
		return nil, &os.PathError{
			Op:   "open",
			Path: name,
			Err:  errors.New("no such file or directory"),
		}
	}
	if info.isDir {
		return nil, &os.PathError{
			Op:   "open",
			Path: name,
			Err:  errors.New("is a directory"),
		}
	}
	if info.perm&0o400 == 0 {
		return nil, &os.PathError{
			Op:   "open",
			Path: name,
			Err:  errors.New("permission denied"),
		}
	}

	return info.data, nil
}

func (fs *FS) WriteFile(name string, data []byte, perm os.FileMode) error {
	if !path.IsAbs(name) && name != "" {
		name = path.Join(fs.Pwd, name)
	}

	parent, err := fs.checkParents(name, true)
	if err != nil {
		return err
	}

	info, ok := fs.Nodes[name]
	if ok {
		if info.isDir {
			return &os.PathError{
				Op:   "open",
				Path: name,
				Err:  errors.New("is a directory"),
			}
		}
	}
	// Return error if file exists and has no write permission, or if the file
	// does not exist and the direct parent has no write permission.
	if ok && info.perm&0o200 == 0 || !ok && parent.perm&0o200 == 0 {
		return &os.PathError{
			Op:   "open",
			Path: name,
			Err:  errors.New("permission denied"),
		}
	}

	fs.Nodes[name] = &Node{data: data, perm: perm}

	return nil
}

func (fs *FS) Remove(name string) error {
	if !path.IsAbs(name) && name != "" {
		name = path.Join(fs.Pwd, name)
	}

	parent, err := fs.checkParents(name, false)
	if err != nil {
		return err
	}

	if parent != nil && parent.perm&0o200 == 0 {
		return &os.PathError{
			Op:   "remove",
			Path: name,
			Err:  errors.New("permission denied"),
		}
	}

	info, ok := fs.Nodes[name]
	if !ok {
		return &os.PathError{
			Op:   "remove",
			Path: name,
			Err:  errors.New("no such file or directory"),
		}
	}
	if info.perm&0o200 == 0 {
		return &os.PathError{
			Op:   "remove",
			Path: name,
			Err:  errors.New("permission denied"),
		}
	}
	if info.isDir {
		for p := range fs.Nodes {
			if strings.HasPrefix(p, name) && p != name {
				return &os.PathError{
					Op:   "remove",
					Path: name,
					Err:  errors.New("directory not empty"),
				}
			}
		}
	}

	delete(fs.Nodes, name)

	return nil
}

func (fs *FS) Exists(name string) bool {
	if !path.IsAbs(name) && name != "" {
		name = path.Join(fs.Pwd, name)
	}

	_, ok := fs.Nodes[name]

	return ok
}

func (fs *FS) FileMode(name string) (os.FileMode, error) {
	if !path.IsAbs(name) && name != "" {
		name = path.Join(fs.Pwd, name)
	}

	if info, ok := fs.Nodes[name]; ok {
		return info.perm, nil
	}

	return 0, &os.PathError{
		Op:   "open",
		Path: name,
		Err:  os.ErrNotExist,
	}
}

func (fs *FS) checkParents(absPath string, noExistError bool) (*Node, error) {
	var parents []string
	for d := path.Dir(absPath); d != "/"; d = path.Dir(d) {
		parents = append(parents, d)
	}
	parents = append(parents, "/")
	var directParent *Node

	for i := 0; i < len(parents); i++ {
		dir := parents[i]
		info, ok := fs.Nodes[dir]
		if !ok && noExistError {
			return nil, &os.PathError{
				Op:   "open",
				Path: dir,
				Err:  errors.New("no such file or directory"),
			}
		}
		if info != nil && !info.isDir {
			return nil, &os.PathError{
				Op:   "open",
				Path: dir,
				Err:  errors.New("not a directory"),
			}
		}
		// Ensure all parent directories have execute permissions.
		if info != nil && info.perm&0o100 == 0 {
			return nil, &os.PathError{
				Op:   "open",
				Path: dir,
				Err:  errors.New("permission denied"),
			}
		}
		if i == 0 {
			directParent = info
		}
	}

	return directParent, nil
}
