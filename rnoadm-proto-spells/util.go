package main

import (
	"math/big"
	"strings"
)

var zero = big.NewInt(0)
var one_thousand = big.NewInt(1000)

func Comma(n *big.Int) string {
	n = (&big.Int{}).Set(n) // copy

	var negative string
	if n.Cmp(zero) < 0 {
		negative = "-"
	}
	n.Abs(n)

	tmp := &big.Int{}
	var s []string
	for n.Cmp(one_thousand) >= 0 {
		tmp.Mod(n, one_thousand)
		tmp.Add(tmp, one_thousand)
		s = append(s, tmp.String()[1:])
		n.Div(n, one_thousand)
	}
	s = append(s, n.String())
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return negative + strings.Join(s, ",")
}

func CommaPlus(n *big.Int) string {
	if n.Cmp(zero) > 0 {
		return "+" + Comma(n)
	}
	return Comma(n)

}
