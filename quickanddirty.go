// Package qad provides a few quick and dirty funcs
// that panic on error.

package qad

import (
	"fmt"
	"io/fs"
	"log"
	"os"
)

func FileMode(filename string) fs.FileMode {
	st, err := os.Stat(filename)
	if err != nil {
		log.Panicf("filemode file:%q err:%v", filename, err)
	}
	return st.Mode()
}

func NewFile(filename string, mode fs.FileMode) *os.File {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_EXCL, mode)
	if err != nil {
		log.Panicf("newfile file:%q err:%v", filename, err)
	}
	return f
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
