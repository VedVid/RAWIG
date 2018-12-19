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
	/* Function LayerError is helper function that returns string
	   to error; it takes layer integer as argument and returns string. */
	return "\n    <layer:  " + strconv.Itoa(layer) + ">"
}

func CoordsError(x, y int) string {
	/* Function CoordsError is helper function that returns string
	   to error; it takes coords x, y as arguments and returns string,
	   with use global MapSizeX and MapSizeY constants. */
	txt := "\n    <x: " + strconv.Itoa(x) + "; y: " + strconv.Itoa(y) +
		"; map width: " + strconv.Itoa(MapSizeX) + "; map height: " +
		strconv.Itoa(MapSizeY) + ">"
	return txt
}

func CharacterLengthError(character string) string {
	/* Function CharacterLengthError is helper function that returns string
	   to error; it takes character string as argument and returns string.
	   Character (as something's representation on map) is supposed to be
	   one-letter long. */
	txt := "\n    <length: " + strconv.Itoa(utf8.RuneCountInString(character)) +
		"; character: " + character + ">"
	return txt
}

func PlayerAIError(ai int) string {
	/* Function PlayerAIError is helper function that returns string to error;
	   it takes ai code (integer) as argument and returns string.
	   Player AI is supposed to be PlayerAI (defined in ai.go).
	   It's supposed to be warning, not error. */
	txt := "\n    <player ai code: " + strconv.Itoa(ai) + ">"
	return txt
}

func InitialHPError(hp int) string {
	/* Function InitialHPError is helper function that returns string to error;
	   it takes creature's HPMax as argument and returns string.
	   It will be warning instead of error sometimes - negative hp for newly created
	   creatures is unusual, but it is not bug per se. */
	txt := "\n    <fighter hp: " + strconv.Itoa(hp) + ">"
	return txt
}

func InitialAttackError(attack int) string {
	/* Function InitialAttackError is helper function that returns string
	   to error; it takes creature's attack value as argument and returns string.
	   Attack value should not be negative. */
	txt := "\n    <fighter attack: " + strconv.Itoa(attack) + ">"
	return txt
}

func InitialDefenseError(defense int) string {
	/* Function InitialDefenseError is helper function that returns string
	   to error; it takes creature's defense value as argument.
	   Defense value should not be negative. */
	txt := "\n    <fighter defense: " + strconv.Itoa(defense) + ">"
	return txt
}

func EquippableSlotError(equippable bool, slot int) string {
	/* Function EquippableSlotError is helper function that returns string
	   to error; it takes equippable bool and slot int as arguments.
	   Slot should not be 0 if equippable is set to true, and should be 0
	   if equippable is set to false. */
	txt := "\n    <equippable: " + strconv.FormatBool(equippable) + "; slot: " +
		strconv.Itoa(slot) + ">"
	return txt
}