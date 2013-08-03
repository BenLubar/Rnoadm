// +build ignore

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
)

func main() {
	for _, fn := range os.Args[1:] {
		file(fn)
	}
}

func file(fn string) {
	ext := path.Ext(fn)
	if ext == ".go" || ext == ".jar" {
		return
	}
	var in []byte
	var err error
	if ext == ".js" {
		cmd := exec.Command("java", "-jar", "closure-compiler.jar",
			"--accept_const_keyword",
			"--compilation_level", "ADVANCED_OPTIMIZATIONS",
			"--formatting", "SINGLE_QUOTES",
			"--js", fn)
		cmd.Stderr = os.Stderr
		in, err = cmd.Output()
	} else {
		in, err = ioutil.ReadFile(fn)
	}
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
	Resource[%q] = []byte{`, path.Base(fn))

	for i, b := range in {
		if i == 0 {
			fmt.Fprintf(out, `%d`, b)
		} else {
			fmt.Fprintf(out, `, %d`, b)
		}
	}
	fmt.Fprintf(out, "}\n}\n")
}
