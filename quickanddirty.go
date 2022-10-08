// Package qad provides a few quick and dirty funcs
// that panic on error.

package qad

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"time"
)

func FileMode(filename string) fs.FileMode {
	st, err := os.Stat(filename)
	if err != nil {
		log.Panicf("filemode file:%q err:%v", filename, err)
	}
	return st.Mode()
}

func FileExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func FileModTime(filename string) time.Time {
	st, err := os.Stat(filename)
	Assert(err, "stat", filename)
	return st.ModTime()
}

func NewFile(filename string, mode fs.FileMode) *os.File {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_EXCL, mode)
	if err != nil {
		log.Panicf("newfile file:%q err:%v", filename, err)
	}
	return f
}

func MoveFile(src, dst string) {
	err := os.Rename(src, dst)
	if err == nil {
		return
	}

	// copy file
	Assert(!FileExist(dst), "dst file exist", dst)

	srcf, err := os.Open(src)
	Assert(err, "open src", src)
	defer srcf.Close()

	tmpfile := dst + ".tmp"
	dstf := NewFile(tmpfile, FileMode(src))
	defer dstf.Close() // ensure close

	_, err = io.Copy(dstf, srcf)
	Assert(err, "copy", src, "=>", dst)
	dstf.Close()

	mtime := FileModTime(src)
	err = os.Chtimes(tmpfile, mtime, mtime)
	Assert(err, "chtimes", dst, mtime)

	err = os.Rename(tmpfile, dst)
	Assert(err, "rename tmpfile", tmpfile, "=>", dst)
}

func Assert(errOk any, args ...any) {
	switch v := errOk.(type) {
	case error:
		if v != nil {
			log.Fatal(append([]any{fmt.Sprintf("ASSERT err:%v:", v)}, args...)...)
		}
	case bool:
		if !v {
			log.Fatal(append([]any{"ASSERT false:"}, args...)...)
		}
	default:
		log.Fatal(append([]any{fmt.Sprintf("ASSERT unknown:%v", v)}, args...))
	}
}
