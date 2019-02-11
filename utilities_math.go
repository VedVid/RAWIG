/*
Copyright (c) 2018, Tomasz "VedVid" Nowakowski
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this
   list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
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

func RandRange(min, max int) int {
	return RandInt(max-min) + min
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

func FindCreatureIndex(creature *Creature, arr Creatures) (int, error) {
	/* Function FindCreatureIndex works as FindObjectIndex,
	   but for monsters. */
	var err error
	index := WrongIndexValue
	for i := 0; i < len(arr); i++ {
		if arr[i] == creature {
			index = i
			break
		}
		if index == WrongIndexValue {
			err = errors.New("*Creature not found in []*Creature.")
		}
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

func (c *Creature) DistanceBetweenCreatures(c2 *Creature) int {
	distance := DistanceBetween(c.X, c.Y, c2.X, c2.Y)
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

func CreatureIsInSlice(c *Creature, arr Creatures) bool {
	/* Function CreatureIsInSlice takes Creature and slice of Creature
	   as arguments, and returns true, if c is present in arr.
	   Otherwise, returns false. */
	for _, v := range arr {
		if c == v {
			return true
		}
	}
	return false
}
