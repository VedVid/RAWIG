/*
Copyright (c) 2018 Tomasz "VedVid" Nowakowski

This software is provided 'as-is', without any express or implied
warranty. In no event will the authors be held liable for any damages
arising from the use of this software.

Permission is granted to anyone to use this software for any purpose,
including commercial applications, and to alter it and redistribute it
freely, subject to the following restrictions:

1. The origin of this software must not be misrepresented; you must not
   claim that you wrote the original software. If you use this software
   in a product, an acknowledgment in the product documentation would be
   appreciated but is not required.
2. Altered source versions must be plainly marked as such, and must not be
   misrepresented as being the original software.
3. This notice may not be removed or altered from any source distribution.
*/

package main

import (
	"errors"
	"math"
	"math/rand"

	blt "bearlibterminal"
)

const (
	// Values to catch errors.
	WrongIndexValue = -1
)

func RoundFloatToInt(x float64) int {
	/* Function RoundFloatToInt takes one float64 number,
	   rounds it to nearest 1.0, then returns it as a integer. */
	return int(math.Round(x))
}

func RandInt(max int) int {
	/* Function RandInt wraps rand.Intn function;
	   instead of returning 0..n-1 it returns 0..n. */
	return rand.Intn(max + 1)
}

func OrderToCharacter(i int) string {
	/* Function OrderToCharacter takes integer
	   and converts it to string. Typically,
	   it will be used with letters, but rune
	   is alias of int32 and support unicode
	   well.
	   Typically, one would like to return
	   string('a'-1+i)
	   to convert "1" to "a", but RAWIG will use
	   it to deal with bare slices that count
	   from 0.*/
	return string('a' + i)
}

func KeyToOrder(key int) int {
	/* Function KeyToOrder takes user input as integer
	   (in BearLibTerminal player input is passed as 0x...)
	   and return another int that is smaller by
	   first-key (ie "a" key).
	   It will need extensive error-checking
	   (or maybe just LBYL?) for wrong input. */
	return key - blt.TK_A
}

func FindObjectIndex(item *Object, arr Objects) (int, error) {
	/* Function FindObjectIndex takes object, and slice of objects
	   as arguments. It returns integer and error.
	   It is supposed to find index item in arr. If fails,
	   returns error. */
	var err error
	index := WrongIndexValue
	for i := 0; i < len(arr); i++ {
		if arr[i] == item {
			index = i
			break
		}
	}
	if index == WrongIndexValue {
		err = errors.New("*Object not found in []*Object.")
	}
	return index, err
}

func DistanceBetween(sourceX, sourceY, targetX, targetY int) int {
	/* Function DistanceBetween takes coords of source and target;
	   it computes distance between these two tiles.
	   As Go uses float64 for such a computations, it is necessary
	   to transform ints to float64 then round result to int. */
	dx := float64(targetX - sourceX)
	dy := float64(targetY - sourceY)
	distance := RoundFloatToInt(math.Sqrt(math.Pow(dx, 2) + math.Pow(dy, 2)))
	return distance
}

func AbsoluteValue(i int) int {
	/* Function AbsoluteValue returns absolute (ie "non negative") value. */
	if i < 0 {
		return -i
	}
	return i
}

func ReverseIntSlice(arr []int) []int {
	/* Function ReverseIntSlice takes slice of int and returns
	   it in reversed order. It is odd that "battery included"
	   language like Go does not have built-in functions for it. */
	var reversed = []int{}
	for i := len(arr) - 1; i >= 0; i-- {
		reversed = append(reversed, arr[i])
	}
	return reversed
}
