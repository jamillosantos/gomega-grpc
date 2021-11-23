//go:build tools
// +build tools

package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

// main fixes the code generation from mockgen replacing the import `google.golang.org/grpc/internal/status` to `google.golang.org/grpc/status`
func main() {
	fileName := ""
	for _, fName := range os.Args[1:] {
		if fName == "--" {
			continue
		}
		fileName = fName
		break
	}
	if fileName == "" {
		panic("missing file: arg[1]")
	}
	f, err := os.OpenFile(fileName, os.O_RDWR, 0)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	var fileContent []byte
	fileContent, err = io.ReadAll(f)
	f.Seek(0, io.SeekStart)
	n, err := f.Write(bytes.ReplaceAll(fileContent, []byte("google.golang.org/grpc/internal/status"), []byte("google.golang.org/grpc/status")))
	f.Truncate(int64(n))
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed writing file:", err)
	}
}
