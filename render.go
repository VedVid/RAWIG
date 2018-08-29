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

import blt "bearlibterminal"

const (
	BaseLayer = iota
	BoardLayer
	ObjectsLayer
	MonstersLayer
	PlayerLayer
)

func PrintBoard(b Board) {
	/*Function PrintBoard is used in RenderAll;
	it takes level map as arguments and iterates through that slice;
	prints every tile on its coords*/
	for _, v := range b {
		blt.Layer(v.BasicProperties.Layer)
		glyph := "[color=" + v.BasicProperties.Color + "]" +
			v.BasicProperties.Char
		blt.Print(v.BasicProperties.X, v.BasicProperties.Y, glyph)
	}
}

func PrintObjects(o Objects) {
	/*Function PrintObjects is used in RenderAll;
	it takes slice of objects as argument and iterates through it;
	prints every object on its coords*/
	for _, v := range o {
		blt.Layer(v.BasicProperties.Layer)
		glyph := "[color=" + v.BasicProperties.Color + "]" +
			v.BasicProperties.Char
		blt.Print(v.BasicProperties.X, v.BasicProperties.Y, glyph)
	}
}

func PrintMonsters(b Board, m Monsters) {
	/*Function PrintMonsters is used in RenderAll;
	it takes slice of monsters as argument and iterates through it;
	checks for every monster if is in player's (assuming that first monster
	is player) FOV by calling IsInFOV;
	prints monster if that function returns true*/
	for _, v := range m {
		if IsInFOV(b, m[0].BasicProperties.X, m[0].BasicProperties.Y,
			v.BasicProperties.X, v.BasicProperties.Y) == true {
			blt.Layer(v.BasicProperties.Layer)
			glyph := "[color=" + v.BasicProperties.Color + "]" +
				v.BasicProperties.Char
			blt.Print(v.BasicProperties.X, v.BasicProperties.Y, glyph)
		}
	}
}

func RenderAll(b Board, o Objects, m Monsters) {
	/*Function RenderAll prints every tile and character on game screen;
	takes board slice (ie level map), slice of objects, and slice of monsters
	as arguments;
	at first, it clears whole terminal window, then uses arguments:
	CastRays (for raycasting FOV) of first object (assuming that it is player),
	then:
	call functions for printing map, objects and monsters;
	at the end, RenderAll calls blt.Refresh() that makes
	changes to the game window visible*/
	blt.Clear()
	CastRays(b, o[0].BasicProperties.X, o[0].BasicProperties.Y)
	PrintBoard(b)
	PrintObjects(o)
	PrintMonsters(b, m)
	blt.Refresh()
}
