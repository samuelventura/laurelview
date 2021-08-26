package lvsdk

import (
	"os"
	"path/filepath"
	"strings"
)

//different filename same extension
func RelativeSibling(sibling string) string {
	exe := ExecutablePath()
	dir := filepath.Dir(exe)
	base := filepath.Base(exe)
	ext := filepath.Ext(base) //includes .
	file := sibling + ext
	return filepath.Join(dir, file)
}

//same file name different extension
func RelativeExtension(ext string) string {
	path := ExecutablePath()
	return ChangeExtension(path, ext)
}

func ChangeExtension(path string, next string) string {
	ext := filepath.Ext(path) //includes .
	npath := strings.TrimSuffix(path, ext)
	return npath + next
}

func ExecutablePath() string {
	exe, err := os.Executable()
	PanicIfError(err)
	return exe
}
