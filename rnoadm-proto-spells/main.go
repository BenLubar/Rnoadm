package main

import (
	"fmt"
	"math/rand"
)

func main() {
	r := rand.New(rand.NewSource(0))

	for i := 0; i < 10; i++ {
		fmt.Println("\n\n\nSpell", i)
		fmt.Println(NewSpell(r))
	}
}
