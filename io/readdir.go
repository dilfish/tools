package io

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

// HasIndex find if target exists in list match
func HasIndex(target string, match []string) bool {
	for _, m := range match {
		if strings.Index(target, m) >= 0 {
			return true
		}
	}
	return false
}

// ListFiles list all files under directory dir
// extFilter drops the very extension, nameFilter drops the very names
func ListFiles(dir string, extFilter, nameFilter []string) []string {
	var names []string
	var ret []string
	if len(dir) == 0 {
		dir = "."
	}
	dirs, err := os.ReadDir(dir)
	if err != nil {
		log.Println("read dir error:", dir, err)
		return nil
	}
	for _, d := range dirs {
		if d.IsDir() {
			tmp := ListFiles(dir + "/" + d.Name(), extFilter, nameFilter)
			ret = append(ret, tmp...)
		} else {
			names = append(names, d.Name())
		}
	}
	for _, n := range names {
		ext := filepath.Ext(dir + "/" + n)
		nonExtName := n[:len(n) - len(ext)]
		if HasIndex(ext, extFilter) {
			continue
		}
		if HasIndex(nonExtName, nameFilter) {
			continue
		}
		ret = append(ret, dir + "/" + n)
	}
	return  ret
}