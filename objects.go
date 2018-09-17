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

const (
	//colors
	colorObject     = "blue"
	colorObjectDark = "dark blue"
)

type Object struct {
	/*Objects are every other things on map;
	  statues, tables, chairs; but also weapons,
	  armour parts, etc.*/
	BasicProperties
	VisibilityProperties
	CollisionProperties
}

/*Objects holds every object on map.*/
type Objects []*Object

func NewObject(layer, x, y int, character, colour, colourDark string,
	alwaysVisible, blocked, blocksSight bool) (*Object, error) {
	/*Function NewObject takes all values necessary by its struct,
	and creates then returns Object*/
	var err error
	if layer < 0 {
		txt := LayerError(layer)
		err = errors.New("Object layer is smaller than 0. " + txt)
	}
	if x < 0 || x >= WindowSizeX || y < 0 || y >= WindowSizeY {
		txt := CoordsError(x, y)
		err = errors.New("Object coords is out of window range." + txt)
	}
	if utf8.RuneCountInString(character) != 1 {
		txt := CharacterLengthError(character)
		err = errors.New("Object character string length is not equal to 1." + txt)
	}
	objectBasicProperties := BasicProperties{layer, x, y, character, colour,
		colourDark}
	objectVisibilityProperties := VisibilityProperties{alwaysVisible}
	objectCollisionProperties := CollisionProperties{blocked, blocksSight}
	objectNew := &Object{objectBasicProperties, objectVisibilityProperties,
		objectCollisionProperties}
	return objectNew, err
}
