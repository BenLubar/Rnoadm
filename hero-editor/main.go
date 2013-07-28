package main

import (
	"compress/gzip"
	"encoding/gob"
	"encoding/json"
	"flag"
	"os"
)

var Seed int64

func decode() {
	g, err := gzip.NewReader(os.Stdin)
	if err != nil {
		panic(err)
	}
	var player Player
	err = gob.NewDecoder(g).Decode(&player)
	if err != nil {
		g.Close()
		panic(err)
	}
	err = g.Close()
	if err != nil {
		panic(err)
	}

	json.NewEncoder(os.Stdout).Encode(player)
}

func encode() {
	var player Player
	err := json.NewDecoder(os.Stdin).Decode(&player)
	if err != nil {
		panic(err)
	}
	g, err := gzip.NewWriterLevel(os.Stdout, gzip.BestCompression)
	if err != nil {
		panic(err)
	}
	defer g.Close()
	err = gob.NewEncoder(g).Encode(&player)
	if err != nil {
		panic(err)
	}
}

func print() {
	g, err := gzip.NewReader(os.Stdin)
	if err != nil {
		panic(err)
	}
	var player Player
	err = gob.NewDecoder(g).Decode(&player)
	if err != nil {
		g.Close()
		panic(err)
	}
	err = g.Close()
	if err != nil {
		panic(err)
	}

	println(player.Name())
}

var (
	flagEncode = flag.Bool("e", false, "encode mode")
	flagDecode = flag.Bool("d", false, "decode mode")
	flagPrint  = flag.Bool("p", false, "print mode")
)

func main() {
	flag.Parse()
	count := 0
	if *flagEncode {
		count++
	}
	if *flagDecode {
		count++
	}
	if *flagPrint {
		count++
	}
	if count != 1 {
		flag.Usage()
		os.Exit(1)
	}
	if *flagEncode {
		encode()
	}
	if *flagDecode {
		decode()
	}
	if *flagPrint {
		print()
	}
}
