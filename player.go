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
	"unicode/utf8"
)

func NewPlayer(layer, x, y int, colour, character string) (*Creature, error) {
	/*Function NewPlayer takes all values necessary by its struct,
	and creates then returns pointer to Creature;
	so, it's basically NewMonster function.*/
	var err error
	if layer < 0 {
		err = errors.New("Player layer is smaller than 0.")
	}
	if x < 0 || x >= WindowSizeX || y < 0 || y >= WindowSizeY {
		err = errors.New("Player coords is out of window range.")
	}
	if utf8.RuneCountInString(character) != 1 {
		err = errors.New("Player character string length is not equal to 1.")
	}
	playerBlock := Basic{layer, x, y, colour, character}
	playerNew := &Creature{playerBlock}
	return playerNew, err
}
