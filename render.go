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

func PrintBoard(b Board, m Monsters) {
	/*Function PrintBoard is used in RenderAll;
	it takes level map as arguments and iterates through that slice;
	prints every tile on its coords if certain conditions are met:
	is Explored already, and:
	- is in player's field of view (prints "normal" color) or
	- is AlwaysVisible (prints dark color).*/
	for _, v := range b {
		if v.Explored == true {
			blt.Layer(v.Layer)
			if IsInFOV(b, m[0].X, m[0].Y, v.X, v.Y) {
				glyph := "[color=" + v.Color + "]" + v.Char
				blt.Print(v.X, v.Y, glyph)
			} else {
				if v.AlwaysVisible == true {
					glyph := "[color=" + v.ColorDark + "]" + v.Char
					blt.Print(v.X, v.Y, glyph)
				}
			}
		}
	}
}

func PrintObjects(b Board, o Objects, m Monsters) {
	/*Function PrintObjects is used in RenderAll;
	it takes slice of objects as argument and iterates through it;
	prints every object on its coords if certain conditions are met:
	AlwaysVisible bool is set to true, or is in player fov.*/
	for _, v := range o {
		if (IsInFOV(b, m[0].X, m[0].Y, v.X, v.Y) == true) ||
			(v.AlwaysVisible == true) {
			blt.Layer(v.Layer)
			glyph := "[color=" + v.Color + "]" + v.Char
			blt.Print(v.X, v.Y, glyph)
		}
	}
}

func PrintMonsters(b Board, m Monsters) {
	/*Function PrintMonsters is used in RenderAll;
	it takes slice of monsters as argument and iterates through it;
	checks for every monster on its coords if certain conditions are met:
	AlwaysVisible bool is set to true, or is in player fov.*/
	for _, v := range m {
		if (IsInFOV(b, m[0].X, m[0].Y, v.X, v.Y) == true) ||
			(v.AlwaysVisible == true) {
			blt.Layer(v.Layer)
			glyph := "[color=" + v.Color + "]" + v.Char
			blt.Print(v.X, v.Y, glyph)
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
	CastRays(b, m[0].X, m[0].Y)
	PrintBoard(b, m)
	PrintObjects(b, o, m)
	PrintMonsters(b, m)
	blt.Refresh()
}
