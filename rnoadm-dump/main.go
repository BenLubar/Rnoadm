package main

import (
	"compress/gzip"
	"encoding/gob"
	"flag"
	"fmt"
	"github.com/kr/pretty"
	"math/big"
	"os"
	"time"
)

func init() {
	gob.Register([]interface{}(nil))
	gob.Register(map[string]interface{}(nil))
	gob.Register(time.Time{})
	gob.Register(&big.Int{})
}

func main() {
	flag.Parse()

	for _, fn := range flag.Args() {
		dump(fn)
	}
}

func dump(fn string) {
	fmt.Println("===", fn, "===")
	defer fmt.Println()

	f, err := os.Open(fn)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	g, err := gzip.NewReader(f)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer g.Close()

	var data interface{}
	err = gob.NewDecoder(g).Decode(&data)
	if err != nil {
		fmt.Println(err)
		return
	}

	pretty.Println(data)
}
