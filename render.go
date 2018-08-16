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
