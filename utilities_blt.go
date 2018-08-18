package main

import (
	"strconv"

	blt "bearlibterminal"
)

func SetGlyph(path, number, filter string, size int) {
	/*Function SetTile allows to use special tiles (glyphs, bitmaps)
	as font elements;
	number variables has to be formatted in that way:
	U+<unicode-number>, like: U+E001
	Later, that U+E001 identifier may be used in printing functions, like
	wall := 0xE001 (note different format!); blt.Print(x, y, wall)*/
	blt.Set(number + ": " + path + ", resize=" + size + ", resize-filter=" + filter)
}
