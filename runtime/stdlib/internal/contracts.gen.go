// Code generated by go-bindata. DO NOT EDIT.
// sources:
// contracts/crypto.cdc (5.269kB)

package internal

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data, name string) ([]byte, error) {
	gz, err := gzip.NewReader(strings.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("read %q: %w", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)

	if err != nil {
		return nil, fmt.Errorf("read %q: %w", name, err)
	}

	clErr := gz.Close()
	if clErr != nil {
		return nil, clErr
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes  []byte
	info   os.FileInfo
	digest [sha256.Size]byte
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _contractsCryptoCdc = "\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xb4\x58\xdd\x6e\xdb\x46\x13\xbd\xd7\x53\xcc\xe7\x2b\x1b\x9f\xc2\xb8\x40\x51\x14\x02\x18\x23\x6d\x93\xd6\x70\x8a\x16\x76\x6c\x5f\x18\x46\xb2\x16\x47\xe4\xc2\xcc\x52\xd8\x1d\xca\x66\x05\xbd\x7b\xb1\x5c\xfe\xec\x1f\x6d\xb9\x75\x75\x23\x4a\x7b\x66\xf6\xcc\xd9\x99\x21\x87\xb3\x75\x7d\x07\x8a\x64\xbd\x24\xe0\x82\x50\xae\xd8\x12\xe1\x82\xe7\x82\x51\x2d\xf1\x0a\x25\x5f\x71\x94\x00\xdb\xd9\x0c\x00\x40\xc3\x57\xb5\x80\x8d\x5e\x68\x0e\xdb\xff\xf4\x47\xf5\x16\x0b\xb8\xb9\x3c\x15\xf4\xe3\xed\x7c\x58\x23\x96\x2f\xe0\x82\x24\x17\xf9\xdc\x31\xc0\xec\x17\x46\x2c\x62\xb1\xae\xef\x4a\xbe\x3c\xc3\x26\xb2\x36\xec\xf4\xbe\xcc\x2b\xc9\xa9\xf8\xb6\x18\xf9\x0e\xff\x8d\xf8\x82\xa9\xc2\x82\xfe\x66\xff\x6c\x41\x47\x0b\xf8\xa9\xaa\xca\xd9\x6e\x16\x17\x43\x5b\x44\x14\xd0\x7e\xc7\xf8\xb3\x78\x20\xd1\xd0\xd9\x33\x6c\x3a\x2f\x3d\xa1\x65\x25\x48\xb2\x25\xc1\xcf\xb2\x59\x53\x15\xe5\xf1\xc5\x23\x30\xb9\xc7\xe8\x1e\xb6\x03\x21\x89\x54\x4b\x01\x0a\xcb\x55\xa2\xdd\x5d\x73\x2a\x3e\xb3\xfc\x50\xfb\x9c\x9b\x18\x0e\x0e\x1c\x9f\xc3\xe5\x51\xeb\x64\x17\x52\xea\x7d\x04\xcc\x6c\x49\xfe\x1d\x4d\x94\xed\xd7\xa1\xd9\xc0\x22\x4b\x2c\xdf\x93\x6d\x77\xd8\x67\xd8\x7c\xe2\x8a\x3e\x08\x92\x8d\xb5\xa1\x46\x94\x48\x70\x8f\xcd\xa9\xc8\xf0\x71\x01\xa7\x82\x82\x55\x2b\x59\xff\xec\x2f\x03\xd0\xb3\x59\x68\x83\x1f\x90\xe7\x05\x2d\xe0\xf2\x23\x7f\xfc\xe1\xfb\x60\x99\xab\x73\xdc\x54\xf7\x98\x75\x89\x3b\x00\xb8\xe0\x34\xa6\xa4\xfe\x38\xcc\xe7\xce\x52\x8c\xb6\x8b\x78\x8a\xb3\x8b\x74\x09\xbb\x6b\x3e\xdb\xfe\xff\x23\x4b\xe8\xb6\xae\xf5\xb1\xf6\x7c\x21\x1d\xa8\x87\xa0\x81\x39\xa4\x63\x14\x21\xcc\xa1\x0f\xa9\x1b\x4e\x08\x37\x31\x40\xda\x05\x13\x02\x86\x40\x20\x1d\x83\x1a\x60\xbb\xa7\x53\xab\xaf\xda\x76\x55\xf2\x4d\x7b\x90\x28\x48\x72\x54\x0b\xb8\xb1\x13\xf0\xd6\x3b\xd0\xa8\x50\x9d\x29\xa4\x70\x73\x6b\x71\x18\x2e\xdf\xbe\x7d\x0b\xef\xb3\x4c\x01\x03\x81\x0f\x5a\x4c\x78\xe0\x54\x00\x15\x08\x39\xdf\xa0\xf0\xc3\xec\x6b\x97\x65\x99\x9b\x42\x5f\xfe\xc3\x4c\x19\x93\x61\xe1\xd7\xa0\x63\x65\x17\x21\xa4\x8e\x02\x49\x89\x22\xa7\x22\x80\x63\xeb\x27\x75\xdc\xba\x81\x81\x53\x1f\xfd\xd5\x3c\xc0\x58\xe1\xaf\xe3\xe1\x47\x24\x28\xa6\x25\xb0\x65\x30\xdf\xe1\xba\x55\x34\x2b\x56\x2a\x74\x00\x47\x93\xe9\x90\xb0\xf5\x1a\x45\x76\xd8\x06\xef\xc2\xba\xe6\xd9\xae\x4c\x25\xcc\x79\x8b\x51\x6d\x92\xe8\x8c\x61\x64\xe5\x0b\x6f\xd5\x01\xbe\x02\x4e\x80\x8f\x5c\x91\x4a\x3c\x6b\x53\x1d\xf7\xd8\x28\x60\x12\x81\x95\x0f\xac\x51\xdd\xce\x98\xcd\xe1\xae\x6e\x1d\x36\x50\xb0\x0d\xc2\xd7\x21\xc8\xaf\xb0\xe2\x58\x66\xa0\x90\x80\x2a\x20\x59\x63\x90\x97\x39\xd2\xa1\xd3\xcd\xbc\x94\x39\xf1\xaa\x84\xaf\xc6\x8c\x79\x17\x4d\x19\xcf\xc0\x12\x49\xf0\xd2\x59\xda\xcd\x62\x52\xda\x2e\x6f\xfa\xbd\x26\x8b\xf1\x77\x26\xef\x9f\x52\x16\xa4\xd1\xc2\xa8\x94\x55\xa8\x40\x54\x04\x19\x96\x48\x08\x3c\x2c\x54\x83\xf7\x34\x79\x3d\x11\x3c\x01\xec\x5f\xba\xbc\x96\xb5\x94\x28\xba\x6a\x4d\x9f\xd3\x02\xbc\x3c\x1d\x21\x2f\xa8\x50\x7b\xcb\x64\xaf\x72\x75\x2c\xf6\xaf\x5d\xc7\x6c\xcf\x42\x76\x6c\xf6\xa8\x6a\x27\xc7\xc1\x29\xea\xa9\xa2\x94\x35\xea\x03\x1d\xf3\x66\x78\x12\x36\xe5\xb6\x61\x25\xcf\x60\x55\x49\x0f\x82\x59\xfb\x6c\x14\x24\x10\x57\x57\xda\xc2\x95\x7c\xf0\x79\x81\x34\xde\x97\x86\xe7\xeb\xdb\x79\x80\xf6\x1e\xe2\xed\x96\xae\xef\xf7\x7e\x2b\xdf\x30\x69\x98\x5e\xb7\x22\xa9\xfe\x56\x00\x29\x1c\x27\xc7\x61\xdb\x57\x88\xe2\xac\x3d\x6c\xbe\xd4\xb7\xca\xed\xa9\x20\xe3\x79\x07\x29\x6c\xbd\xd2\xd4\xd1\x0f\x21\x00\x17\x4e\x3c\x3e\x15\xa3\x30\x7c\x10\x4a\x83\xfb\xd2\x34\xd5\xc8\x95\x61\x19\x5a\xf0\xd5\xe8\x34\x79\x79\x75\xc1\xd8\x40\xc2\xce\x0e\x61\xb3\xf1\x49\x72\x65\xb1\x2c\x98\xe9\x12\xac\x94\xc8\xb2\x06\xee\x50\x9f\x39\xa2\x88\xd3\x76\xa4\xbc\x09\xa3\xb8\x85\x93\x13\xc3\xea\xf5\x88\x9f\xe3\xb2\x92\x99\xa7\xee\x03\x53\x13\x34\xf7\xe0\x98\x9a\xe2\x89\x6d\xf6\x2b\x9a\xbe\xca\x96\x54\xb3\x52\x6f\x18\xc2\xba\x67\x09\xbf\x69\x45\x76\xda\x2b\x5b\xcc\x01\x74\xcd\x3b\x2a\xfb\x3d\x36\xd6\xb3\xe3\xeb\x67\x04\xda\x19\xff\x54\xda\xfe\xcf\x0c\x90\x89\xf2\xa7\xfb\xc4\x9f\xe6\x9d\x13\x19\x27\xfb\x51\xa3\xe1\x2a\x6c\x72\xd0\x8f\xbd\xdd\x6e\x59\xf5\x8d\x71\x71\x81\x6b\x26\x19\xf1\x4a\x7c\x66\xf9\xa5\x42\x19\x37\xb4\x1b\xca\x78\x1d\xc7\x5a\x8d\x5e\x4b\x3c\xfc\x7c\xaa\xd3\xc3\xc4\x0b\x04\xd7\x41\x88\x88\x7b\x72\xef\x1b\xda\xc7\xf4\x8c\x01\x91\x91\xa7\xff\xbc\xf4\xf4\xed\xfe\x09\xa9\xfb\xf3\xff\x6d\x2c\x91\x19\x26\xfe\x10\xe3\xd8\xbe\x4b\xe1\xbb\xe4\x78\xdf\x91\x66\xb8\x29\xbc\x78\x62\x0e\x5f\x16\xcd\x1c\x4c\x3b\xf6\xb8\xa3\x6b\xc4\xe6\x9f\x0f\x90\x63\xbd\xa4\xa3\xdf\x89\xb0\xfb\x59\x6d\x22\x8d\xfb\x57\x19\x1e\x3a\x28\xb1\x05\x6c\x83\x97\x6a\x3b\xd7\xc6\xbc\xd2\x58\xc0\xd6\xbc\x70\xea\x08\xb4\x5a\xec\xe7\x6e\x1e\xb8\x38\xb2\x6f\x79\x6e\xe8\xc3\x9b\xbd\x34\x24\xeb\x9a\x18\xa7\xdd\x00\x8d\xd2\x7e\x3a\x81\x53\xc1\x89\xb3\x92\xff\x85\xb0\xac\x84\x22\x26\x48\x79\x3b\x4e\x08\x07\x29\x1c\x7c\xfc\xf4\xc7\xf5\x9b\xab\xe3\xe4\xf8\x4d\xad\x50\x1e\x74\xba\xef\x66\x7f\x07\x00\x00\xff\xff\x80\x1a\xbe\x1b\x95\x14\x00\x00"

func contractsCryptoCdcBytes() ([]byte, error) {
	return bindataRead(
		_contractsCryptoCdc,
		"contracts/crypto.cdc",
	)
}

func contractsCryptoCdc() (*asset, error) {
	bytes, err := contractsCryptoCdcBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "contracts/crypto.cdc", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x3, 0x5c, 0x57, 0xc, 0x6d, 0x79, 0x84, 0x16, 0xa6, 0x85, 0x10, 0xad, 0xce, 0x20, 0x97, 0xcc, 0xc, 0x7e, 0x7f, 0x4f, 0x81, 0xad, 0x5b, 0xfe, 0x6a, 0x4d, 0x56, 0x9d, 0xb1, 0xf, 0x87, 0x33}}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// AssetString returns the asset contents as a string (instead of a []byte).
func AssetString(name string) (string, error) {
	data, err := Asset(name)
	return string(data), err
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// MustAssetString is like AssetString but panics when Asset would return an
// error. It simplifies safe initialization of global variables.
func MustAssetString(name string) string {
	return string(MustAsset(name))
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetDigest returns the digest of the file with the given name. It returns an
// error if the asset could not be found or the digest could not be loaded.
func AssetDigest(name string) ([sha256.Size]byte, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return [sha256.Size]byte{}, fmt.Errorf("AssetDigest %s can't read by error: %v", name, err)
		}
		return a.digest, nil
	}
	return [sha256.Size]byte{}, fmt.Errorf("AssetDigest %s not found", name)
}

// Digests returns a map of all known files and their checksums.
func Digests() (map[string][sha256.Size]byte, error) {
	mp := make(map[string][sha256.Size]byte, len(_bindata))
	for name := range _bindata {
		a, err := _bindata[name]()
		if err != nil {
			return nil, err
		}
		mp[name] = a.digest
	}
	return mp, nil
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"contracts/crypto.cdc": contractsCryptoCdc,
}

// AssetDebug is true if the assets were built with the debug flag enabled.
const AssetDebug = false

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"},
// AssetDir("data/img") would return []string{"a.png", "b.png"},
// AssetDir("foo.txt") and AssetDir("notexist") would return an error, and
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		canonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(canonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"contracts": {nil, map[string]*bintree{
		"crypto.cdc": {contractsCryptoCdc, map[string]*bintree{}},
	}},
}}

// RestoreAsset restores an asset under the given directory.
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	return os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
}

// RestoreAssets restores an asset under the given directory recursively.
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(canonicalName, "/")...)...)
}
