package main

import (
	"math/rand"
)

type Name struct {
	Nickname string
}

func (n *Name) String() string {
	s := n.Nickname
	return s
}

func GenerateHumanName(r *rand.Rand, gender Gender) *Name {
	return &Name{"TODO"}
}
