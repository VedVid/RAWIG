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
	"math/rand"
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
	// For unmarshalling json data.
	Cells          []string
	Data           [][]int
	Layouts        [][][]string
	Char           map[string]string
	Name           map[string]string
	Color          map[string]string
	ColorDark      map[string]string
	Layer          map[string]int
	AlwaysVisible  map[string]bool
	Explored       map[string]bool
	Blocked        map[string]bool
	BlocksSight    map[string]bool
	MonstersCoords [][]int
	MonstersTypes  []string
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

func ReplaceTile(t *Tile, s string, m *MapJson) {
	/* ReplaceTile is function that takes tile, string (supposed to be
	   one-character-lenght - symbol of map tile, taken from json map) and
	   MapJson (ie unmarshalled json map).
	   It uses m's legend to overwrite old map values with data read from file. */
	t.Char = m.Char[s]
	t.Name = m.Name[s]
	t.Color = m.Color[s]
	t.ColorDark = m.ColorDark[s]
	t.Layer = m.Layer[s]
	t.AlwaysVisible = m.AlwaysVisible[s]
	t.Explored = m.Explored[s]
	t.Blocked = m.Blocked[s]
	t.BlocksSight = m.BlocksSight[s]
}

func LoadJsonMap(mapFile string) (Board, Creatures, error) {
	/* Function LoadJsonMap takes string (name of json map file) as argument,
	   and returns Board (ie map), Creatures (included in premade json maps)
	   and error.
	   It uses new type - struct MapJson - to store all values read from file.
	   Panics if unmarshalling encounters any error.
	   Other possible errors are about internal structure of json file:
	       - length of Data and Layouts has to be the same
	       - length of MonstersCoords and MonstersTypes has to be the same.
	   It is important because instead of using multi-type json lists
	   (it would be possible to store map monsters as [x: int, y: int, file: string])
	   there are independent structures. The reason is Go's limitations: bot lists
	   (slices) and dictionaries (maps) are strongly typed. (Un)Marshalling multi-type
	   lists would be cumbersome. On the other hand, it means that creating and editing
	   json maps require discipline.
	   After error checking, three major operations are queued.
	   At first, game reads json map (Cells) and modifies (previously initialized)
	   tiles regarding to json legend (Char, Name (...), BlocksSight).
	   Then it repeats this operation for every area marked as "randomly generated".
	   Some important points to make about these areas:
	       - they are not created *randomly*
	           = areas ("rooms") are specified in JsonMap.Data
	           = they are filled using prefabs (JsonMap.Layouts)
	   At the end, monsters are created and placed on map (their datas are stored
	   in json map as MonstersCoords (x, y) and MonstersTypes (their json files). */
	var jsonMap = &MapJson{}
	var err error
	err = MapFromJson(MapsPathJson+mapFile, jsonMap)
	if err != nil {
		fmt.Println(err)
		panic(-1)
	}
	cells := jsonMap.Cells
	data := jsonMap.Data
	layouts := jsonMap.Layouts
	// Number of items in data should match number of layouts.
	if len(data) != len(layouts) {
		txt := MapDataLayoutsError((len(data)), len(layouts), mapFile)
		err = errors.New("Length of data and layouts does not match. " + txt)
	}
	thisMap := InitializeEmptyMap()
	for x := 0; x < len(cells[0]); x++ {
		for y := 0; y < len(cells); y++ {
			// y,x because - due to 2darray nature - there is height first, width later...
			ReplaceTile(thisMap[x][y], string(cells[y][x]), jsonMap)
		}
	}
	for i, room := range data {
		layoutsToChoose := layouts[i]
		layout := layoutsToChoose[rand.Intn(len(layoutsToChoose))]
		for x := 0; x < len(layout[0]); x++ {
			for y := 0; y < len(layout); y++ {
				ReplaceTile(thisMap[room[0]+x][room[1]+y], string(layout[y][x]), jsonMap)
			}
		}
	}
	coords := jsonMap.MonstersCoords
	aiTypes := jsonMap.MonstersTypes
	if len(coords) != len(aiTypes) {
		txt := MapMonstersCoordsAiError(len(coords), len(aiTypes), mapFile)
		err = errors.New("Length of MonstersCoords and MonstersTypes does not match. " + txt)
	}
	var creatures = Creatures{}
	for j := 0; j < len(coords); j++ {
		monster, err := NewCreature(coords[j][0], coords[j][1], aiTypes[j]+".json")
		if err != nil {
			fmt.Println(err)
		}
		creatures = append(creatures, monster)
	}
	return thisMap, creatures, err
}
