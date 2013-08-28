// +build ignore

package main

import (
	"fmt"
	"hash/crc64"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

func main() {
	for _, fn := range os.Args[1:] {
		file(fn)
	}
}

func file(fn string) {
	ext := path.Ext(fn)
	if ext == ".go" {
		return
	}
	in, err := ioutil.ReadFile(fn)
	if err != nil {
		panic(err)
	}
	out, err := os.Create(fn + ".go")
	if err != nil {
		panic(err)
	}
	defer out.Close()

	fmt.Fprintf(out, `package resource

func init() {
	Resource[%q] = []byte{`, filepath.Base(fn))

	for i, b := range in {
		if i == 0 {
			fmt.Fprintf(out, `%d`, b)
		} else {
			fmt.Fprintf(out, `, %d`, b)
		}
	}
	hash := crc64.New(crc64.MakeTable(crc64.ISO))
	hash.Write(in)
	fmt.Fprintf(out, "}\n\tHash[%q] = \"W/\\\"%d\\\"\"\n}\n", filepath.Base(fn), hash.Sum64())
}
