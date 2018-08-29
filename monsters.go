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

type Creature struct {
	/*Creatures are living objects that
	  moves, attacks, dies, etc.*/
	BasicProperties
	VisibilityProperties
	CollisionProperties
}

/*Monsters holds every monster on map.*/
type Monsters []*Creature

func NewCreature(layer, x, y int, character, colour string,
	alwaysVisible, blocked, blocksSight bool) (*Creature, error) {
	/*Function NewCreture takes all values necessary by its struct,
	and creates then returns pointer to Creature*/
	var err error
	if layer < 0 {
		txt := LayerError(layer)
		err = errors.New("Creature layer is smaller than 0." + txt)
	}
	if x < 0 || x >= WindowSizeX || y < 0 || y >= WindowSizeY {
		txt := CoordsError(x, y)
		err = errors.New("Creature coords is out of window range." + txt)
	}
	if utf8.RuneCountInString(character) != 1 {
		txt := CharacterLengthError(character)
		err = errors.New("Creature character string length is not equal to 1." + txt)
	}
	creatureBasicProperties := BasicProperties{layer, x, y, character, colour}
	creatureVisibilityPropeties := VisibilityProperties{alwaysVisible}
	creatureCollisionProperties := CollisionProperties{blocked, blocksSight}
	creatureNew := &Creature{creatureBasicProperties,
		creatureVisibilityPropeties, creatureCollisionProperties}
	return creatureNew, err
}

func (c *Creature) Move(d Direction) {
	/*Move is method of Creature; it takes Direction type (tuple-like);
	check if next move won't put Creature off the screen, then updates
	Creature coords.*/
	if c.BasicProperties.X+d.X >= 0 &&
		c.BasicProperties.X+d.X <= WindowSizeX-1 &&
		c.BasicProperties.Y+d.Y >= 0 &&
		c.BasicProperties.Y+d.Y <= WindowSizeX-1 {
		c.BasicProperties.X += d.X
		c.BasicProperties.Y += d.Y
	}
}
