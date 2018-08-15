package main

import blt "bearlibterminal"

const (
	baseLayer = iota
	boardLayer
	objectsLayer
	monstersLayer
)

func PrintBoard(b Board) {
	for i := range b {
		blt.Layer(b[i].Block.Layer)
		glyph := "[color=" + b[i].Block.Color + "]" + b[i].Block.Char
		blt.Print(b[i].Block.X, b[i].Block.Y, glyph)
	}
}

func PrintObjects(o Objects) {
	for i := range o {
		blt.Layer(o[i].Block.Layer)
		glyph := "[color=" + o[i].Block.Color + "]" + o[i].Block.Char
		blt.Print(o[i].Block.X, o[i].Block.Y, glyph)
	}
}

func PrintMonsters(m Monsters) {
	for i := range m {
		blt.Layer(m[i].Block.Layer)
		glyph := "[color=" + m[i].Block.Color + "]" + m[i].Block.Char
		blt.Print(m[i].Block.X, m[i].Block.Y, glyph)
	}
}

func RenderAll(b Board, o Objects, m Monsters) {
	blt.Clear()
	PrintBoard(b)
	PrintObjects(o)
	PrintMonsters(m)
}
