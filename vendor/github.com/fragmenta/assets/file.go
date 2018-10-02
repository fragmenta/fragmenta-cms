package assets

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

const permissions = 0744

// File stores a filename and hash fingerprint for the asset file
type File struct {
	name  string
	hash  string
	path  string
	bytes []byte
}

// NewFile returns a new file object
func NewFile(p string) (*File, error) {

	// Load file from path to get bytes
	bytes, err := ioutil.ReadFile(p)
	if err != nil {
		return &File{}, err
	}

	// Calculate hash and save it
	file := &File{
		path:  p,
		name:  path.Base(p),
		hash:  bytesHash(bytes),
		bytes: bytes,
	}
	return file, nil
}

// Style returns true if this file is a CSS file
func (f *File) Style() bool {
	return strings.HasSuffix(f.name, ".css")
}

// Script returns true if this file is a js file
func (f *File) Script() bool {
	return strings.HasSuffix(f.name, ".js")
}

// MarshalJSON generates json for this file, of the form {group:{file:hash}}
func (f *File) MarshalJSON() ([]byte, error) {
	var b bytes.Buffer

	s := fmt.Sprintf("\"%s\":\"%s\"", f.path, f.hash)
	b.WriteString(s)

	return b.Bytes(), nil
}

// Newer returns true if file exists at path
func (f *File) Newer(dst string) bool {

	// Check mtimes
	stat, err := os.Stat(f.path)
	if err != nil {
		return false
	}
	srcM := stat.ModTime()
	stat, err = os.Stat(dst)

	// If the file doesn't exist, return true
	if os.IsNotExist(err) {
		return true
	}

	// Else check for other errors
	if err != nil {
		return false
	}

	dstM := stat.ModTime()

	return srcM.After(dstM)

}

// Copy our bytes to dstpath
func (f *File) Copy(dst string) error {
	err := ioutil.WriteFile(dst, f.bytes, permissions)
	if err != nil {
		return err
	}
	return nil
}

// LocalPath returns the relative path of this file
func (f *File) LocalPath() string {
	return f.path
}

// AssetPath returns the path of this file within the assets folder
func (f *File) AssetPath(dst string) string {
	folder := "styles"
	if f.Script() {
		folder = "scripts"
	}
	return path.Join(dst, "assets", folder, f.name)
}

// String returns a string representation of this object
func (f *File) String() string {
	return fmt.Sprintf("%s:%s", f.name, f.hash)
}
