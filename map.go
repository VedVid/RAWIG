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
	"fmt"
	"unicode/utf8"
)

type Tile struct {
	// Tiles are map cells - floors, walls, doors.
	BasicProperties
	VisibilityProperties
	Explored bool
	CollisionProperties
}

/* Board is map representation, that uses 2d slice
   to hold data of its every cell. */
type Board [][]*Tile

func NewTile(layer, x, y int, character, name, color, colorDark string,
	alwaysVisible, explored, blocked, blocksSight bool) (*Tile, error) {
	/* Function NewTile takes all values necessary by its struct,
	   and creates then returns Tile. */
	var err error
	if layer < 0 {
		txt := LayerError(layer)
		err = errors.New("Tile layer is smaller than 0." + txt)
	}
	if x < 0 || x >= MapSizeX || y < 0 || y >= MapSizeY {
		txt := CoordsError(x, y)
		err = errors.New("Tile coords is out of window range." + txt)
	}
	if utf8.RuneCountInString(character) != 1 {
		txt := CharacterLengthError(character)
		err = errors.New("Tile character string length is not equal to 1." + txt)
	}
	tileBasicProperties := BasicProperties{layer, x, y, character, name, color,
		colorDark}
	tileVisibilityProperties := VisibilityProperties{alwaysVisible}
	tileCollisionProperties := CollisionProperties{blocked, blocksSight}
	tileNew := &Tile{tileBasicProperties, tileVisibilityProperties,
		explored, tileCollisionProperties}
	return tileNew, err
}

func InitializeEmptyMap() Board {
	/* Function InitializeEmptyMap returns new Board, filled with
	   generic (ie "empty") tiles.
	   It starts by declaring 2d slice of *Tile - unfortunately, Go seems to
	   lack simple way to do it, therefore it's necessary to use
	   the first for loop.
	   The second, nested loop initializes specific Tiles within Board bounds. */
	b := make([][]*Tile, MapSizeX)
	for i := range b {
		b[i] = make([]*Tile, MapSizeY)
	}
	for x := 0; x < MapSizeX; x++ {
		for y := 0; y < MapSizeY; y++ {
			//workaround to missing _, err := ... idiom that won't work here
			var err error
			b[x][y], err = NewTile(BoardLayer, x, y, ".", "", "light gray",
				"dark gray", true, false, false, false)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	return b
}
