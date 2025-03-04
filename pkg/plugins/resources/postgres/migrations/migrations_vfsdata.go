// Code generated by vfsgen; DO NOT EDIT.

// +build !dev

package migrations

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	pathpkg "path"
	"time"
)

// Migrations statically implements the virtual filesystem provided to vfsgen.
var Migrations = func() http.FileSystem {
	fs := vfsgen۰FS{
		"/": &vfsgen۰DirInfo{
			name:    "/",
			modTime: time.Date(2020, 6, 15, 14, 49, 38, 390867092, time.UTC),
		},
		"/1579518998_create_resources.up.sql": &vfsgen۰CompressedFileInfo{
			name:             "1579518998_create_resources.up.sql",
			modTime:          time.Date(2020, 4, 2, 12, 24, 0, 540178165, time.UTC),
			uncompressedSize: 299,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x8e\x4f\x0b\x82\x40\x10\x47\xef\x7e\x8a\xdf\x51\x61\x0f\x76\xee\x64\xb1\x81\x64\x16\xba\x41\x1e\x97\x65\x48\x0f\xea\xb2\xb3\x49\x7d\xfb\x68\xfb\x07\x1d\xc2\x39\x0e\x6f\xde\xbc\x75\x25\x33\x25\xa1\xb2\x55\x21\x91\x6f\x50\xee\x15\xe4\x29\xaf\x55\x0d\x47\x3c\x5e\x9c\x21\x46\x1c\x01\xc0\xa0\x7b\xc2\x6b\x26\xed\x4c\xab\x5d\xbc\x48\xd3\x24\xdc\x94\xc7\xa2\x10\x1f\x8c\xad\x36\xf4\x1f\xeb\x89\xdb\x19\x36\x7f\xb3\x73\x9e\x4e\xe4\xb8\x1b\x87\x80\x75\x83\xa7\x33\xb9\x1f\x82\x2d\x99\xb7\xc8\xd3\xd5\x3f\xb7\x87\x2a\xdf\x65\x55\x83\xad\x6c\x10\x3f\xca\xc5\xb7\x5f\x84\x46\x11\x12\x92\x28\x59\xde\x03\x00\x00\xff\xff\x88\x1c\x8d\x52\x2b\x01\x00\x00"),
		},
		"/1580128050_add_creation_modification_time.up.sql": &vfsgen۰CompressedFileInfo{
			name:             "1580128050_add_creation_modification_time.up.sql",
			modTime:          time.Date(2020, 4, 2, 12, 24, 0, 540282811, time.UTC),
			uncompressedSize: 165,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x72\xf4\x09\x71\x0d\x52\x08\x71\x74\xf2\x71\x55\x28\x4a\x2d\xce\x2f\x2d\x4a\x4e\x2d\x56\x70\x74\x71\x51\x70\xf6\xf7\x09\xf5\xf5\x53\x48\x2e\x4a\x4d\x2c\xc9\xcc\xcf\x8b\x2f\xc9\xcc\x4d\x55\x08\xf1\xf4\x75\x0d\x0e\x71\xf4\x0d\x50\xf0\xf3\x0f\x51\xf0\x0b\xf5\xf1\x51\x70\x71\x75\x73\x0c\xf5\x09\x51\xc8\xcb\x2f\xd7\xd0\xb4\xe6\x22\x68\x60\x6e\x7e\x4a\x66\x5a\x66\x32\x29\x86\x02\x02\x00\x00\xff\xff\x56\x69\x01\xb8\xa5\x00\x00\x00"),
		},
		"/1589041445_add_unique_id_and_owner.up.sql": &vfsgen۰CompressedFileInfo{
			name:             "1589041445_add_unique_id_and_owner.up.sql",
			modTime:          time.Date(2020, 5, 12, 16, 39, 56, 292245031, time.UTC),
			uncompressedSize: 973,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x9c\x91\xc1\x6e\xdb\x30\x10\x44\xef\xfc\x8a\xe9\xc9\x31\x20\x17\x71\xae\x46\x0e\xac\x49\xb7\x42\x25\xca\xa0\x64\x14\x39\x05\x84\xb2\xb1\x8c\x28\x92\x40\x2a\x4d\xf3\xf7\x85\x28\xa5\x66\x52\xd7\x45\x7b\xdd\xe5\xec\xbc\x19\xf2\xa4\x90\x1a\x05\xff\x94\x48\x58\x72\xed\x93\x2d\xc9\x31\x00\x10\x3a\xdb\x62\x9d\xa9\xbc\xd0\x3c\x56\xc5\x71\x7b\xdb\x3d\xd0\x4b\xe4\xdf\x70\x21\xfe\xfc\x04\x5b\x1d\xa7\x5c\xdf\xe0\xab\xbc\xc1\x45\x63\x1e\x29\xc2\x23\xb9\x2a\x42\xff\xd2\xd1\x7c\xc5\xfe\xea\x9d\xec\x52\x85\x41\xe8\x3a\x53\xd2\x39\xc1\x08\xe2\xdf\xb7\xcf\x0d\xd9\xdb\x41\x85\xef\xc6\x96\x95\xb1\x17\xcb\xcb\xcb\xb3\x76\xbf\xa9\x07\xcc\xff\x57\x0f\xf1\xfe\x51\xfd\xab\xc2\xf1\xc2\xfd\x03\x36\x99\x96\xf1\x67\x35\x96\x77\xcc\x14\x05\x84\x51\xe0\x37\x87\x96\x1b\xa9\xa5\x5a\xcb\xfc\xe8\x70\xa2\x76\x64\x0a\x42\x26\xb2\x90\x58\xf3\x7c\xcd\x85\x5c\x31\x36\x0e\xd8\x46\x67\x69\x80\xf7\xed\x8b\xd4\xd2\xab\x70\x8d\x99\x30\xbd\xe9\x6a\xd3\x50\xdc\xb8\xc3\xbe\xea\x67\x2b\xc6\x16\x0b\x38\xea\x47\x0c\xdc\xb7\x16\xa6\xae\x03\x77\xfa\x51\x52\xd7\x23\x25\x57\xb1\xdd\x56\xf0\x22\x08\x0f\xbb\x64\xb9\x2c\xc2\xef\xba\x86\xbd\xfa\xe8\x81\x7d\x33\xc1\x57\xf8\x8d\x0f\x11\x6c\x5e\xc1\x86\xf3\xb3\x77\xec\xb0\x57\x13\xfe\x74\x72\x38\xb1\xf4\x27\x18\xc0\x95\x18\xe6\x6f\x0f\x4c\xe3\xe5\x38\xfe\xf0\x3a\xf7\x29\x17\xb8\xa3\x9a\x7a\xc2\x9d\x69\xf6\xf5\xa1\xd9\x87\x15\xb7\x0d\x39\x3c\x1f\xfa\xaa\x7d\x9a\xaa\x98\x9f\x2d\x34\x48\x1c\xe7\x50\xbb\x24\x99\xbc\x83\xc0\xa7\x16\x9e\xeb\xed\xe2\x1d\xea\xcf\x00\x00\x00\xff\xff\x17\x5c\x38\xd4\xcd\x03\x00\x00"),
		},
		"/1592232449_add_leader_table.up.sql": &vfsgen۰CompressedFileInfo{
			name:             "1592232449_add_leader_table.up.sql",
			modTime:          time.Date(2020, 6, 15, 14, 49, 38, 390783584, time.UTC),
			uncompressedSize: 217,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x74\x8e\x31\xab\xc2\x30\x1c\xc4\xf7\x7c\x8a\x1b\x5b\x78\xbc\xe1\x41\xa7\x4e\x49\xde\x9f\x1a\xd4\xa8\x31\x2a\x99\x4a\x6d\x33\x88\x36\x81\x54\xeb\xd7\x17\x6c\x47\x1d\xef\xf8\x71\xf7\x93\x86\xb8\x25\x58\x2e\x56\x84\x5b\x6c\xaf\x03\x32\x06\x00\xa1\xe9\x3d\xe4\x82\x1b\x2e\x2d\x19\x1c\xb9\x71\x4a\x57\xd9\x5f\x51\xe4\xd8\x1a\xb5\xe6\xc6\x61\x49\xee\xe7\x0d\x27\xdf\xc6\xd4\xd5\xa3\x4f\xc3\x25\x86\x3a\x3c\xfa\xb3\x4f\x10\xaa\x52\xda\x4e\x44\xd7\xdc\x1b\x08\x67\x89\x4f\x39\x3e\x83\x4f\x5f\xf6\x59\x5e\x32\x36\x8b\xed\x69\x77\x20\x2d\x67\xb7\x3a\x8d\x01\x9b\x93\xa6\x7f\x08\x37\x55\xbf\x1f\xbf\x4b\xf6\x0a\x00\x00\xff\xff\x9b\x4e\x1b\xbc\xd9\x00\x00\x00"),
		},
	}
	fs["/"].(*vfsgen۰DirInfo).entries = []os.FileInfo{
		fs["/1579518998_create_resources.up.sql"].(os.FileInfo),
		fs["/1580128050_add_creation_modification_time.up.sql"].(os.FileInfo),
		fs["/1589041445_add_unique_id_and_owner.up.sql"].(os.FileInfo),
		fs["/1592232449_add_leader_table.up.sql"].(os.FileInfo),
	}

	return fs
}()

type vfsgen۰FS map[string]interface{}

func (fs vfsgen۰FS) Open(path string) (http.File, error) {
	path = pathpkg.Clean("/" + path)
	f, ok := fs[path]
	if !ok {
		return nil, &os.PathError{Op: "open", Path: path, Err: os.ErrNotExist}
	}

	switch f := f.(type) {
	case *vfsgen۰CompressedFileInfo:
		gr, err := gzip.NewReader(bytes.NewReader(f.compressedContent))
		if err != nil {
			// This should never happen because we generate the gzip bytes such that they are always valid.
			panic("unexpected error reading own gzip compressed bytes: " + err.Error())
		}
		return &vfsgen۰CompressedFile{
			vfsgen۰CompressedFileInfo: f,
			gr:                        gr,
		}, nil
	case *vfsgen۰DirInfo:
		return &vfsgen۰Dir{
			vfsgen۰DirInfo: f,
		}, nil
	default:
		// This should never happen because we generate only the above types.
		panic(fmt.Sprintf("unexpected type %T", f))
	}
}

// vfsgen۰CompressedFileInfo is a static definition of a gzip compressed file.
type vfsgen۰CompressedFileInfo struct {
	name              string
	modTime           time.Time
	compressedContent []byte
	uncompressedSize  int64
}

func (f *vfsgen۰CompressedFileInfo) Readdir(count int) ([]os.FileInfo, error) {
	return nil, fmt.Errorf("cannot Readdir from file %s", f.name)
}
func (f *vfsgen۰CompressedFileInfo) Stat() (os.FileInfo, error) { return f, nil }

func (f *vfsgen۰CompressedFileInfo) GzipBytes() []byte {
	return f.compressedContent
}

func (f *vfsgen۰CompressedFileInfo) Name() string       { return f.name }
func (f *vfsgen۰CompressedFileInfo) Size() int64        { return f.uncompressedSize }
func (f *vfsgen۰CompressedFileInfo) Mode() os.FileMode  { return 0444 }
func (f *vfsgen۰CompressedFileInfo) ModTime() time.Time { return f.modTime }
func (f *vfsgen۰CompressedFileInfo) IsDir() bool        { return false }
func (f *vfsgen۰CompressedFileInfo) Sys() interface{}   { return nil }

// vfsgen۰CompressedFile is an opened compressedFile instance.
type vfsgen۰CompressedFile struct {
	*vfsgen۰CompressedFileInfo
	gr      *gzip.Reader
	grPos   int64 // Actual gr uncompressed position.
	seekPos int64 // Seek uncompressed position.
}

func (f *vfsgen۰CompressedFile) Read(p []byte) (n int, err error) {
	if f.grPos > f.seekPos {
		// Rewind to beginning.
		err = f.gr.Reset(bytes.NewReader(f.compressedContent))
		if err != nil {
			return 0, err
		}
		f.grPos = 0
	}
	if f.grPos < f.seekPos {
		// Fast-forward.
		_, err = io.CopyN(ioutil.Discard, f.gr, f.seekPos-f.grPos)
		if err != nil {
			return 0, err
		}
		f.grPos = f.seekPos
	}
	n, err = f.gr.Read(p)
	f.grPos += int64(n)
	f.seekPos = f.grPos
	return n, err
}
func (f *vfsgen۰CompressedFile) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		f.seekPos = 0 + offset
	case io.SeekCurrent:
		f.seekPos += offset
	case io.SeekEnd:
		f.seekPos = f.uncompressedSize + offset
	default:
		panic(fmt.Errorf("invalid whence value: %v", whence))
	}
	return f.seekPos, nil
}
func (f *vfsgen۰CompressedFile) Close() error {
	return f.gr.Close()
}

// vfsgen۰DirInfo is a static definition of a directory.
type vfsgen۰DirInfo struct {
	name    string
	modTime time.Time
	entries []os.FileInfo
}

func (d *vfsgen۰DirInfo) Read([]byte) (int, error) {
	return 0, fmt.Errorf("cannot Read from directory %s", d.name)
}
func (d *vfsgen۰DirInfo) Close() error               { return nil }
func (d *vfsgen۰DirInfo) Stat() (os.FileInfo, error) { return d, nil }

func (d *vfsgen۰DirInfo) Name() string       { return d.name }
func (d *vfsgen۰DirInfo) Size() int64        { return 0 }
func (d *vfsgen۰DirInfo) Mode() os.FileMode  { return 0755 | os.ModeDir }
func (d *vfsgen۰DirInfo) ModTime() time.Time { return d.modTime }
func (d *vfsgen۰DirInfo) IsDir() bool        { return true }
func (d *vfsgen۰DirInfo) Sys() interface{}   { return nil }

// vfsgen۰Dir is an opened dir instance.
type vfsgen۰Dir struct {
	*vfsgen۰DirInfo
	pos int // Position within entries for Seek and Readdir.
}

func (d *vfsgen۰Dir) Seek(offset int64, whence int) (int64, error) {
	if offset == 0 && whence == io.SeekStart {
		d.pos = 0
		return 0, nil
	}
	return 0, fmt.Errorf("unsupported Seek in directory %s", d.name)
}

func (d *vfsgen۰Dir) Readdir(count int) ([]os.FileInfo, error) {
	if d.pos >= len(d.entries) && count > 0 {
		return nil, io.EOF
	}
	if count <= 0 || count > len(d.entries)-d.pos {
		count = len(d.entries) - d.pos
	}
	e := d.entries[d.pos : d.pos+count]
	d.pos += count
	return e, nil
}
