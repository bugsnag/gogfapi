package gfapi

// This file includes higher level operations on files, such as those provided by the 'os' package

// #cgo pkg-config: glusterfs-api
// #include "glusterfs/api/glfs.h"
// #include <stdlib.h>
// #include <sys/stat.h>
import "C"
import (
	"errors"
	"io"
	"os"
	"syscall"
)

// File is the gluster file object.
type File struct {
	name string
	Fd
	isDir bool
}

// Close closes an open File.
// Close is similar to os.Close in its functioning.
//
// Returns an Error on failure.
func (f *File) Close() error {
	var err error
	var ret C.int

	if f.isDir {
		ret, err = C.glfs_closedir(f.Fd.fd)
	} else {
		ret, err = C.glfs_close(f.Fd.fd)
	}
	if ret < 0 {
		return err
	}

	return nil
}

// Chdir has not been implemented yet
func (f *File) Chdir() error {
	return errors.New("Chdir has not been implemented yet")
}

// Chmod changes the mode of the file to the given mode
//
// Returns an error on failure
func (f *File) Chmod(mode os.FileMode) error {
	return f.Fd.Fchmod(posixMode(mode))
}

// Chown has not been implemented yet
func (f *File) Chown(uid, gid int) error {
	return errors.New("Chown has not been implemented yet")
}

// Name returns the name of the opened file
func (f *File) Name() string {
	return f.name
}

// Read reads atmost len(b) bytes into b
//
// Returns number of bytes read and an error if any
func (f *File) Read(b []byte) (n int, err error) {
	if f == nil {
		return 0, os.ErrInvalid
	}
	n, e := f.Fd.Read(b)
	if n == 0 && len(b) > 0 && e == nil {
		return 0, io.EOF
	}
	if e != nil {
		err = &os.PathError{"read", f.name, e}
	}
	return n, err
}

// Readdir returns the information of files in a directory.
//
// n is the maximum number of items to return. If there are more items than
// the maximum they can be obtained in successive calls. If maximum is 0
// then all the items will be returned.
func (f *File) Readdir(n int) ([]os.FileInfo, error) {
	if !f.isDir {
		return nil, &os.PathError{"readdir", f.name, syscall.ENOTDIR}
	}
	return f.Fd.Readdir(n)
}

// Readdirnames returns the names of files in a directory.
//
// n is the maximum number of items to return and works the same way as Readdir.
func (f *File) Readdirnames(n int) ([]string, error) {
	if !f.isDir {
		return nil, &os.PathError{"readdirnames", f.name, syscall.ENOTDIR}
	}
	return f.Fd.Readdirnames(n)
}

// Seek sets the offset for the next read or write on the file based on whence,
// 0 - relative to beginning of file, 1 - relative to current offset, 2 - relative to end
//
// Returns new offset and an error if any
func (f *File) Seek(offset int64, whence int) (int64, error) {
	return f.Fd.lseek(offset, whence)
}

// Stat returns an os.FileInfo object describing the file
//
// Returns an error on failure
func (f *File) Stat() (os.FileInfo, error) {
	var stat syscall.Stat_t
	err := f.Fd.Fstat(&stat)

	if err != nil {
		return nil, err
	}
	return fileInfoFromStat(&stat, f.name), nil
}

// Sync commits the file to the storage
//
// Returns error on failure
func (f *File) Sync() error {
	return f.Fd.Fsync()
}

// DataSync commits data changes only to storage. Does not commit metadata changes.
// Returns error on failure
func (f *File) DataSync() error {
	return f.Fd.FDatasync()
}

// Write writes len(b) bytes to the file
//
// Returns number of bytes written and an error if any
func (f *File) Write(b []byte) (n int, err error) {
	if f == nil {
		return 0, os.ErrInvalid
	}
	n, e := f.Fd.Write(b)

	if n != len(b) {
		err = io.ErrShortWrite
	}
	if e != nil {
		err = &os.PathError{"write", f.name, e}
	}
	return n, err
}

// WriteString writes the contents of string s to the file
//
// Returns number of bytes written and an error if any
func (f *File) WriteString(s string) (int, error) {
	return f.Write([]byte(s))
}

// Fallocate manipulates the allocated disk space for the file
//
// Returns error on failure
func (f *File) Fallocate(mode int, offset int64, len int64) error {
	return f.Fd.Fallocate(mode, offset, len)
}

// Getxattr gets value of the extended attribute 'attr' and place it in 'dest'
//
// Returns number of bytes placed in 'dest' and error if any
func (f *File) Getxattr(attr string, dest []byte) (int64, error) {
	return f.Fd.Fgetxattr(attr, dest)
}

// Setxattr sets extended attribute with key 'attr' and value 'data'
//
// Returns error on failure
func (f *File) Setxattr(attr string, data []byte, flags int) error {
	return f.Fd.Fsetxattr(attr, data, flags)
}

// Removexattr removes extended attribute named 'attr'
//
// Returns error on failure
func (f *File) Removexattr(attr string) error {
	return f.Fd.Fremovexattr(attr)
}
