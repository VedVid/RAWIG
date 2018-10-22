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
	"math"
	"math/rand"
)

func RoundFloatToInt(x float64) int {
	/*Function RoundFloatToInt takes one float64 number,
	rounds it to nearest 1.0, then returns it as a integer.*/
	return int(math.Round(x))
}

func RandInt(max int) int {
	/*Function RandInt wraps rand.Intn function;
	instead of returning 0..n-1 it returns 0..n.*/
	return rand.Intn(max + 1)
}
