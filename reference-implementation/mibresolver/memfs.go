package mibresolver

import (
	"io"
	"io/fs"
	"sort"
	"strings"
	"time"
)

// memFS is an in-memory filesystem implementing fs.ReadDirFS.
// gosmi requires fs.ReadDirFS to discover and read MIB files.
// This allows MIB loading without OS filesystem access (e.g., in WASM).
type memFS struct {
	files map[string][]byte
}

func newMemFS() *memFS {
	return &memFS{files: make(map[string][]byte)}
}

func (m *memFS) Add(name string, data []byte) {
	m.files[name] = data
}

func (m *memFS) Open(name string) (fs.File, error) {
	if name == "." {
		return &memDir{fs: m}, nil
	}
	data, ok := m.files[name]
	if !ok {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
	}
	return &memFile{name: name, data: data, reader: strings.NewReader(string(data))}, nil
}

func (m *memFS) ReadDir(name string) ([]fs.DirEntry, error) {
	if name != "." {
		return nil, &fs.PathError{Op: "readdir", Path: name, Err: fs.ErrNotExist}
	}
	var entries []fs.DirEntry
	for fname, data := range m.files {
		entries = append(entries, &memDirEntry{name: fname, size: int64(len(data))})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})
	return entries, nil
}

// memFile implements fs.File for an in-memory file.
type memFile struct {
	name   string
	data   []byte
	reader *strings.Reader
}

func (f *memFile) Stat() (fs.FileInfo, error) {
	return &memFileInfo{name: f.name, size: int64(len(f.data))}, nil
}

func (f *memFile) Read(b []byte) (int, error) {
	return f.reader.Read(b)
}

func (f *memFile) Close() error {
	return nil
}

// memDir implements fs.File for the root directory, including fs.ReadDirFile.
type memDir struct {
	fs      *memFS
	entries []fs.DirEntry
	pos     int
}

func (d *memDir) Stat() (fs.FileInfo, error) {
	return &memFileInfo{name: ".", isDir: true}, nil
}

func (d *memDir) Read([]byte) (int, error) {
	return 0, io.EOF
}

func (d *memDir) Close() error {
	return nil
}

func (d *memDir) ReadDir(n int) ([]fs.DirEntry, error) {
	if d.entries == nil {
		var err error
		d.entries, err = d.fs.ReadDir(".")
		if err != nil {
			return nil, err
		}
	}
	if n <= 0 {
		remaining := d.entries[d.pos:]
		d.pos = len(d.entries)
		return remaining, nil
	}
	if d.pos >= len(d.entries) {
		return nil, io.EOF
	}
	end := d.pos + n
	if end > len(d.entries) {
		end = len(d.entries)
	}
	result := d.entries[d.pos:end]
	d.pos = end
	if d.pos >= len(d.entries) {
		return result, io.EOF
	}
	return result, nil
}

// memDirEntry implements fs.DirEntry.
type memDirEntry struct {
	name string
	size int64
}

func (e *memDirEntry) Name() string               { return e.name }
func (e *memDirEntry) IsDir() bool                 { return false }
func (e *memDirEntry) Type() fs.FileMode           { return 0 }
func (e *memDirEntry) Info() (fs.FileInfo, error)  { return &memFileInfo{name: e.name, size: e.size}, nil }

// memFileInfo implements fs.FileInfo.
type memFileInfo struct {
	name  string
	size  int64
	isDir bool
}

func (i *memFileInfo) Name() string      { return i.name }
func (i *memFileInfo) Size() int64       { return i.size }
func (i *memFileInfo) Mode() fs.FileMode { return 0444 }
func (i *memFileInfo) ModTime() time.Time { return time.Time{} }
func (i *memFileInfo) IsDir() bool       { return i.isDir }
func (i *memFileInfo) Sys() interface{}  { return nil }
