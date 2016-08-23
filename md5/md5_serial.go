package main

import (
	"crypto/md5"
	"io/ioutil"
	"os"
	"path/filepath"
)

// MD5All reads all the files in the file tree rooted at root and returns a map
// from file path to the MD5 sum of the file's contents.  If the directory walk
// fails or any read operation fails, MD5All returns an error.
func MD5All(root string) (map[string][md5.Size]byte, error) {
	// make an empty map between string and [md5.Size]byte
	// which is basically a map between file path and its md5 sum
	m := make(map[string][md5.Size]byte)

	// this is the function to be applied onto every file in the directory
	walkFunc := func(path string, info os.FileInfo, err error) error {
		// if the walk returns an error, return the error
		if err != nil {
			return err
		}
		// if the file is not a regular file, skip it by returning nil
		if !info.Mode().IsRegular() {
			return nil
		}
		// read the file's data
		data, err := ioutil.ReadFile(path)
		// if the read is not successful, return the read error
		if err != nil {
			return err
		}
		// if the read is successful, compute the md5 checksum of the data
		// and use the result to populate the map
		m[path] = md5.Sum(data)
		// upon successfully computing the checksum of the file, return nil
		return nil
	}

	err := filepath.Walk(root, walkFunc)
	// if there is any error walking through the files, return nil and report
	// the error
	if err != nil {
		return nil, err
	}
	// if there is no error, return the map and nil
	return m, nil
}
