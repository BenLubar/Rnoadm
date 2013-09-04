package material

import (
	"math/big"
	"strings"
)

// Ported to math/big.Int from github.com/dustin/go-humanize
func Comma(v *big.Int) string {
	{
		var copy big.Int
		copy.Set(v)
		v = &copy
	}
	sign := ""
	if v.Sign() < 0 {
		sign = "-"
		v.Abs(v)
	}

	tmp := &big.Int{}
	herman := big.NewInt(999)
	thousand := big.NewInt(1000)
	var parts []string

	for v.Cmp(herman) > 0 {
		part := tmp.Mod(v, thousand).String()
		switch len(part) {
		case 2:
			part = "0" + part
		case 1:
			part = "00" + part
		}
		v.Div(v, thousand)
		parts = append(parts, part)
	}
	parts = append(parts, v.String())
	for i, j := 0, len(parts)-1; i < j; i, j = i+1, j-1 {
		parts[i], parts[j] = parts[j], parts[i]
	}
	return sign + strings.Join(parts, ",")
}
