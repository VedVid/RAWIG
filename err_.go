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
	"strconv"
	"unicode/utf8"
)

func LayerError(layer int) string {
	return "\n    <layer:  " + strconv.Itoa(layer) + ">"
}

func CoordsError(x, y int) string {
	txt := "\n    <x: " + strconv.Itoa(x) + "; y: " + strconv.Itoa(y) +
		"; map width: " + strconv.Itoa(WindowSizeX) + "; map height: " +
		strconv.Itoa(WindowSizeY) + ">"
	return txt
}

func CharacterLengthError(character string) string {
	txt := "\n    <length: " + strconv.Itoa(utf8.RuneCountInString(character)) +
		"; character: " + character + ">"
	return txt
}
