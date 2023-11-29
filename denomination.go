package main

import "strconv"

type denomination uint

func (d denomination) String() string {
	if d < 100 {
		return strconv.FormatFloat(float64(d)/100, 'f', 2, 32)
	}
	return strconv.Itoa(int(d / 100))
}

const (
	cent = 1

	twentyCents = 20 * cent
	fiftyCents  = 50 * cent
	one         = 100 * cent
	two         = 2 * one
	five        = 5 * one
	ten         = 10 * one
	twenty      = 20 * one
	fifty       = 50 * one
)

var orderedDenominations = []denomination{fifty, twenty, ten, five, two, one, fiftyCents, twentyCents, cent}
