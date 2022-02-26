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
	blt "bearlibterminal"
	"errors"
)

const (
	BrensenhamColorNeutral = "white"
	BrensenhamColorGood    = "green"
	BrensenhamColorBad     = "red"
	BrensenhamWhyInspect   = "inspect"
	BrensenhamWhyTarget    = "target"
)

type Brensenham struct {
	/* Brensenham is struct that is supposed to help
	   with creating simple, straight lines between
	   two points. It has start point, target point,
	   and slice of bools.
	   These bools should be set to false by default.
	   For every passable tile from Start to Target,
	   one bool will be changed to true. */
	StartX  int
	StartY  int
	TargetX int
	TargetY int
	Values  []bool
	TilesX  []int
	TilesY  []int
}

func NewBrensenham(sx, sy, tx, ty int) (*Brensenham, error) {
	/* Function NewBrensenham creates new Brensenham with sx, sy as sources coords and
	   tx, ty as target coords. Brensenham has length also, and number of
	   "false" Values is equal to 1 + distance between source and target. */
	var err error
	if sx < 0 || sx >= MapSizeX || sy < 0 || sy >= MapSizeY ||
		tx < 0 || tx >= MapSizeX || ty < 0 || ty >= MapSizeY {
		txt := BrensenhamCoordinatesOutOfMapBounds(sx, sy, tx, ty)
		err = errors.New("Brensenham coordinates are out of map bounds." + txt)
	}
	length := DistanceBetween(sx, sy, tx, ty)
	values := make([]bool, length+1)
	newBrensenham := &Brensenham{sx, sy, tx, ty, values,
		[]int{}, []int{}}
	return newBrensenham, err
}

func ComputeBrensenham(vec *Brensenham) int {
	/* Function ComputeBrensenham takes *Brensenham as argument.
	   It uses Brensenham's Like algorithm to compute tile values
	   (stored in initially empty TilesX and TilesY) between
	   starting point (vec.StartX, vec.StartY) and goal
	   (vec.TargetX, vec.TargetY). */
	vec.TilesX = nil
	vec.TilesY = nil
	sx, sy := vec.StartX, vec.StartY
	tx, ty := vec.TargetX, vec.TargetY
	steep := AbsoluteValue(ty-sy) > AbsoluteValue(tx-sx)
	if steep == true {
		sx, sy = sy, sx
		tx, ty = ty, tx
	}
	rev := false
	if sx > tx {
		sx, tx = tx, sx
		sy, ty = ty, sy
		rev = true
	}
	dx := tx - sx
	dy := AbsoluteValue(ty - sy)
	errValue := dx / 2
	y := sy
	var stepY int
	if sy < ty {
		stepY = 1
	} else {
		stepY = -1
	}
	for x := sx; x <= tx; x++ {
		if steep == true {
			vec.TilesX = append(vec.TilesX, y)
			vec.TilesY = append(vec.TilesY, x)
		} else {
			vec.TilesX = append(vec.TilesX, x)
			vec.TilesY = append(vec.TilesY, y)
		}
		errValue -= dy
		if errValue < 0 {
			y += stepY
			errValue += dx
		}
	}
	if rev == true {
		vec.TilesX = ReverseIntSlice(vec.TilesX)
		vec.TilesY = ReverseIntSlice(vec.TilesY)
	}
	trueLength := len(vec.TilesX)
	return trueLength
}

func FindBrensenhamDirection(vec *Brensenham) ([]int, []int) {
	/* FindBrensenhamDirection is function that takes Brensenham as argument
	   and returns two slices of ints.
	   This function is supposed to find pattern of Brensenham's line,
	   and use it as direction indicator in ExtrapolateBrensenham function. */
	var dx = []int{}
	var dy = []int{}
	for x := 1; x < len(vec.TilesX); x++ {
		for y := 1; y < len(vec.TilesY); y++ {
			if vec.TilesX[x] > vec.TilesX[x-1] {
				dx = append(dx, 1)
			} else if vec.TilesX[x] < vec.TilesX[x-1] {
				dx = append(dx, -1)
			} else {
				dx = append(dx, 0)
			}
			if vec.TilesY[y] > vec.TilesY[y-1] {
				dy = append(dy, 1)
			} else if vec.TilesY[y] < vec.TilesY[y-1] {
				dy = append(dy, -1)
			} else {
				dy = append(dy, 0)
			}
		}
	}
	return dx, dy
}

func ExtrapolateBrensenham(vec *Brensenham, dx, dy []int) *Brensenham {
	/* Function ExtrapolateBrensenham takes Brensenham and two slices of ints
	   as arguments, and returns new Brensenham.
	   It uses slices as direction indicator, pattern - dx may look like
	   [0, 0, 1, 0, 0] - and while iterating ad infinitum, these values
	   will be added to existing ones. For example, if current vector
	   goes in horizontal x 10, 11, 12, it will indicate that every one
	   tile x grows. */
	startX, startY := vec.TargetX, vec.TargetY
	var newTilesX = vec.TilesX
	var newTilesY = vec.TilesY
	i := 0
	for {
		newX, newY := startX+dx[i], startY+dy[i]
		if newX < 0 || newX >= MapSizeX || newY < 0 || newY >= MapSizeY {
			break
		}
		newTilesX = append(newTilesX, newX)
		newTilesY = append(newTilesY, newY)
		startX, startY = newX, newY
		i++
		if i == len(dx) {
			i = 0
		}
	}
	values := make([]bool, len(newTilesX)+1)
	newBrensenham := &Brensenham{vec.StartY, vec.StartX,
		vec.TargetX, vec.TargetY, values, newTilesX, newTilesY}
	return newBrensenham
}

func ValidateBrensenham(vec *Brensenham, b Board, c Creatures,
	o Objects) (bool, *Tile, *Creature, *Object) {
	/* Function ValidateBrensenham takes Brensenham and Board as arguments.
	   It is important function for ranged combat visualisation - function
	   checks if line is not blocked by map tiles or other creatures,
	   or objects. Returns first blocked value.
	   Four values to return looks bad, but it may be better than
	   code duplication if there would be three different functions
	   for Tile, Creature and Object. */
	var tile *Tile
	var monster *Creature
	var object *Object
	valid := false
	length := len(vec.TilesX)
Loop:
	for i := 0; i < length; i++ {
		x, y := vec.TilesX[i], vec.TilesY[i]
		if x == vec.StartX && y == vec.StartY {
			continue
		}
		if b[x][y].Blocked == true {
			// Breaks on blocked tiles.
			tile = b[x][y]
			break
		}
		for j := 0; j < len(c); j++ {
			if x == c[j].X && y == c[j].Y && c[j].Blocked == true {
				// Breaks on first enemy.
				vec.Values[i] = true
				monster = c[j]
				break Loop
			}
		}
		for k := 0; k < len(o); k++ {
			if x == o[k].X && y == o[k].Y && o[k].Blocked == true {
				// Breaks on blocking objects.
				object = o[k]
				break Loop
			}
		}
		vec.Values[i] = true
	}
	if vec.Values[len(vec.Values)-1] == true {
		// Brensenham is valid - path is passable.
		valid = true
	}
	// Brensenham is invalid - blocked tiles in path.
	return valid, tile, monster, object
}

func PrintBrensenham(vec *Brensenham, why string, color1, color2 string, b Board, o Objects, c Creatures) {
	/* Function PrintBrensenham has to take Brensenham, and (unfortunately,
	   due to flawed game architecture) Board, "global" Objects, and
	   Creatures.
	   At start, it clears whole screen and redraws it.
	   Then, it uses tile coords of Brensenham (ie TilesX and TilesY)
	   to set coordinates of printing line symbol.*/
	blt.Clear()
	RenderAll(b, o, c)
	blt.Layer(LookLayer)
	length := len(vec.TilesX)
	for i := 0; i < length; i++ {
		if i == 0 && length > 1 {
			// Do not draw over player, unless he is targeting self.
			continue
		}
		x := vec.TilesX[i]
		y := vec.TilesY[i]
		if x >= 0 && x < MapSizeX && y >= 0 && y < MapSizeY {
			if why == BrensenhamWhyInspect {
				PrintRangedCharacter(x, y, BrensenhamColorNeutral, true)
				if i == 0 && length == 1 {
					break
				}
			} else if why == BrensenhamWhyTarget {
				if IsInFOV(b,
					vec.StartX, vec.StartY, vec.TargetX, vec.TargetY) == true {
					if vec.Values[i] == true {
						PrintRangedCharacter(x, y, color1, true)
					} else {
						PrintRangedCharacter(x, y, color2, false)
					}
				} else {
					if IsInFOV(b,
						vec.StartX, vec.StartY,
						vec.TilesX[i], vec.TilesY[i]) == true {
						PrintRangedCharacter(x, y, color1, true)
					} else {
						PrintRangedCharacter(x, y, color2, false)
					}
				}
			}
		}
	}
	blt.Refresh()
}

func PrintRangedCharacter(x, y int, color string, valid bool) {
	/* Function PrintRangedCharacter takes coords, color, and result
	   of validity test as arguments.
	   It draws green rectangle if cursor position is valid,
	   and red cross if it is not. */
	blt.Layer(LookLayer)
	if valid == true {
		var chars = []string{"▁", "▏", "▕", "▔"}
		for i, v := range chars {
			blt.Layer(LookLayer + i)
			ch := "[color=" + color + "]" + v + "[/color]"
			blt.Print(x, y, ch)
		}
	} else {
		ch := "[color=" + color + "]" + "X" + "[/color]"
		blt.Print(x, y, ch)
	}
}
