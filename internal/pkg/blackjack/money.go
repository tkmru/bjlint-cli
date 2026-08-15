package blackjack

import "fmt"

type Money int64

func (m Money) String() string {
	sign := ""
	if m < 0 {
		sign = "-"
		m = -m
	}
	return fmt.Sprintf("%s$%d.%02d", sign, m/100, m%100)
}
