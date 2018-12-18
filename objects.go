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
	// Colors.
	colorObject     = "blue"
	colorObjectDark = "dark blue"
)

const (
	// Slots for inventory handling.
	SlotNA = iota

	SlotWeapon
)

type Object struct {
	/* Objects are every other things on map;
	   statues, tables, chairs; but also weapons,
	   armor parts, etc. */
	BasicProperties
	VisibilityProperties
	CollisionProperties
	ObjectProperties
}

// Objects holds every object on map.
type Objects []*Object

func NewObject(layer, x, y int, character, color, colorDark string,
	alwaysVisible, blocked, blocksSight bool, pickable, equippable bool, slot int) (*Object, error) {
	/* Function NewObject takes all values necessary by its struct,
	   and creates then returns Object. */
	var err error
	if layer < 0 {
		txt := LayerError(layer)
		err = errors.New("Object layer is smaller than 0. " + txt)
	}
	if x < 0 || x >= MapSizeX || y < 0 || y >= MapSizeY {
		txt := CoordsError(x, y)
		err = errors.New("Object coords is out of window range." + txt)
	}
	if utf8.RuneCountInString(character) != 1 {
		txt := CharacterLengthError(character)
		err = errors.New("Object character string length is not equal to 1." + txt)
	}
	if (equippable == false && slot != SlotNA) || (equippable == true && slot == SlotNA) {
		txt := EquippableSlotError(equippable, slot)
		err = errors.New("'equippable' and 'slot' values does not match." + txt)
	}
	objectBasicProperties := BasicProperties{layer, x, y, character, color,
		colorDark}
	objectVisibilityProperties := VisibilityProperties{alwaysVisible}
	objectCollisionProperties := CollisionProperties{blocked, blocksSight}
	objectProperties := ObjectProperties{pickable, equippable}
	objectNew := &Object{objectBasicProperties, objectVisibilityProperties,
		objectCollisionProperties, objectProperties}
	return objectNew, err
}
