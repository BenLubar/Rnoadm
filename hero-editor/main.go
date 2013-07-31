package main

import (
	"compress/gzip"
	"encoding/gob"
	"flag"
	"fmt"
	"os"
)

var Seed int64

func read() *Player {
	f, err := os.Open(*file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	g, err := gzip.NewReader(f)
	if err != nil {
		panic(err)
	}
	defer g.Close()

	var player Player
	err = gob.NewDecoder(g).Decode(&player)
	if err != nil {
		panic(err)
	}

	return &player
}

func write(player *Player) {
	f, err := os.Create(*file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	g, err := gzip.NewWriterLevel(f, gzip.BestCompression)
	if err != nil {
		panic(err)
	}
	defer g.Close()

	err = gob.NewEncoder(g).Encode(&player)
	if err != nil {
		panic(err)
	}
}

func toggleadmin() {
	player := read()
	player.Admin = !player.Admin
	write(player)
}

func print() {
	player := read()

	fmt.Println("Name:", player.Name())
	fmt.Println("Admin:", player.Admin)
}

var (
	file      = flag.String("file", "", "file to edit [required]")
	flagAdmin = flag.Bool("a", false, "toggle player admin status")
	flagPrint = flag.Bool("p", false, "print info")
)

func main() {
	flag.Parse()
	if *file == "" {
		flag.Usage()
		os.Exit(1)
	}
	count := 0
	if *flagAdmin {
		count++
	}
	if *flagPrint {
		count++
	}
	if count != 1 {
		flag.Usage()
		os.Exit(1)
	}
	if *flagAdmin {
		toggleadmin()
	}
	if *flagPrint {
		print()
	}
}
