package main

import (
	"fmt"
	"math/rand"
)

func main() {
	fmt.Println(NewSpell(rand.New(rand.NewSource(0))))
}
