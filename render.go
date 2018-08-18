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
	baseLayer = iota
	boardLayer
	objectsLayer
	monstersLayer
)

func PrintBoard(b Board) {
	for i, v := range b {
		blt.Layer(v.Block.Layer)
		glyph := "[color=" + v.Block.Color + "]" + v.Block.Char
		blt.Print(v.Block.X, v.Block.Y, glyph)
	}
}

func PrintObjects(o Objects) {
	for i, v := range o {
		blt.Layer(v.Block.Layer)
		glyph := "[color=" + v.Block.Color + "]" + v.Block.Char
		blt.Print(v.Block.X, v.Block.Y, glyph)
	}
}

func PrintMonsters(m Monsters) {
	for i, v := range m {
		blt.Layer(v.Block.Layer)
		glyph := "[color=" + v.Block.Color + "]" + v.Block.Char
		blt.Print(v.Block.X, v.Block.Y, glyph)
	}
}

func RenderAll(b Board, o Objects, m Monsters) {
	blt.Clear()
	PrintBoard(b)
	PrintObjects(o)
	PrintMonsters(m)
}
