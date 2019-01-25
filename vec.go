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

import (
	blt "bearlibterminal"
	"errors"
)

const (
	vectorSymbol = "X" // Maybe vectorSymbol should be customized, as colors?
	VectorColorNeutral  = "white"
	VectorColorGood = "green"
	VectorColorBad = "red"
)

type Vector struct {
	/* Vector is struct that is supposed to help
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

func NewVector(sx, sy, tx, ty int) (*Vector, error) {
	/* Function NewVector creates new Vector with sx, sy as sources coords and
	   tx, ty as target coords. Vector has length also, and number of
	   "false" Values is equal to 1 + distance between source and target. */
	var err error
	if sx < 0 || sx >= MapSizeX || sy < 0 || sy >= MapSizeY ||
		tx < 0 || tx >= MapSizeX || ty < 0 || ty >= MapSizeY {
		txt := VectorCoordinatesOutOfMapBounds(sx, sy, tx, ty)
		err = errors.New("Vector coordinates are out of map bounds." + txt)
	}
	length := DistanceBetween(sx, sy, tx, ty)
	values := make([]bool, length+1)
	newVector := &Vector{sx, sy, tx, ty, values,
		[]int{}, []int{}}
	return newVector, err
}

func ComputeVector(vec *Vector) {
	/* Function ComputeVector takes *Vector as argument.
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
		stepY = (-1)
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
}

func ValidateVector(vec *Vector, b Board, c Creatures) bool {
	/* Function ValidateVector takes Vector and Board as arguments.
	   It is important function for ranged combat visualisation - function
	   checks if line is not blocked by map tiles or other creatures.
	   In future, it should check for objects as well.*/
	length := len(vec.TilesX)
Loop:
	for i := 0; i < length; i++ {
		x, y := vec.TilesX[i], vec.TilesY[i]
		if b[x][y].Blocked == true {
			break
		}
		for j := 0; j < len(c); j++ {
			if x == c[j].X && y == c[j].Y && c[j].Blocked == true {
				// Breaks on first enemy.
				vec.Values[i] = true
				break Loop
			}
		}
		vec.Values[i] = true
	}
	if vec.Values[len(vec.Values)-1] == true {
		// Vector is valid - path is passable.
		return true
	}
	// Vector is invalid - blocked tiles in path.
	return false
}

func PrintVector(vec *Vector, color1, color2 string, b Board, o Objects, c Creatures) {
	/* Function PrintVector has to take Vector, and (unfortunately,
	   due to flawed game architecture) Board, "global" Objects, and
	   Creatures.
	   At start, it clears whole screen and redraws it.
	   Then, it uses tile coords of Vector (ie TilesX and TilesY)
	   to set coordinates of printing line symbol. */
	blt.Clear()
	RenderAll(b, o, c)
	blt.Layer(LookLayer)
	ch1 := "[color=" + color1 + "]" + vectorSymbol
	ch2 := "[color=" + color2 + "]" + vectorSymbol
	length := len(vec.TilesX)
	for i := 0; i < length; i++ {
		if i == 0 && length > 1 {
			// Do not draw over player, unless he is targeting self.
			continue
		}
		x := vec.TilesX[i]
		y := vec.TilesY[i]
		if vec.Values[i] == true {
			blt.Print(x, y, ch1)
		} else {
			blt.Print(x, y, ch2)
		}
	}
	blt.Refresh()
}
