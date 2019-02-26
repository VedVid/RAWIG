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

type MapJson struct {
	Cells []string
	Data [][]int
	Layouts [][][]string
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
	tileBasicProperties := BasicProperties{x, y, character, name, color,
		colorDark}
	tileVisibilityProperties := VisibilityProperties{layer, alwaysVisible}
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
			b[x][y], err = NewTile(BoardLayer, x, y, ".", "floor", "light gray",
				"dark gray", true, false, false, false)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	return b
}

func LoadJsonMap(mapFile string) (Board, error) {
	var jsonMap = &MapJson{}
	err := MapFromJson(MapsPathJson+mapFile, jsonMap)
	if err != nil {
		fmt.Println(err)
		panic(-1)
	}
	cells := jsonMap.Cells
	data := jsonMap.Data
	layouts := jsonMap.Layouts
	//number of items in data should match 2nd items in layouts
	var err2 error
	if len(data) != len(layouts) {
		txt := MapDataLayoutsError((len(data)), len(layouts), mapFile)
		err2 = errors.New("Length of data and layouts does not match. " + txt)
	}
	thisMap := InitializeEmptyMap()
	for x := 0; x < len(cells[0]); x++ {
		for y := 0; y < len(cells); y++ {
			continue
			//thisMap[x][y].Char = string(cells[y][x]) //y,x because - due to 2darray nature - there is height first, width later...
		}
	}
	fmt.Println(len(cells))
	fmt.Println(len(cells[0]))
	return thisMap, err2
}
