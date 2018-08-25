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

type Tile struct {
	/*Tiles are map cells - floors, walls, doors*/
	Block    Basic
	Explored bool
	Blocked  bool
}

/*Board is map representation, that uses slice
  to hold data of its every cell*/
type Board []*Tile

func NewTile(layer, x, y int, colour, character string) *Tile {
	/*Function NewTile takes all values necessary by its struct,
	and creates then returns Tile.
	Every newly created tile is unexplored and not blocked by default.*/
	tileBlock := Basic{layer, x, y, colour, character}
	tileNew := &Tile{tileBlock, false, false}
	return tileNew
}

func FindTileByXY(b Board, x, y int) *Tile {
	/*Function FindTileByXY takes whole board as its argument, and
	desired x, y coords as well. It iterates through board, and
	returns tile that has same xy values as arguments;
	otherwise, it returns nil.
	It needs to be reworked to use idiomatic error boilerplate, thought.
	Also, for now, I'm not sure if range is worth trying - it makes copies;
	maybe basic for i :=0; i < len(b); i++ would make more sense?*/
	for _, v := range b {
		if x == v.Block.X && y == v.Block.Y {
			return v
		}
	}
	return nil
}
