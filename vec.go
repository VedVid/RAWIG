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
	"errors"

	blt "bearlibterminal"
)

const (
	vectorSymbol = "X"
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
	   "false" Values is equal to distance between source and target. */
	var err error
	if sx < 0 || sx >= MapSizeX || sy < 0 || sy >= MapSizeY ||
		tx < 0 || tx >= MapSizeX || ty < 0 || ty >= MapSizeY {
			txt := VectorCoordinatesOutOfMapBounds(sx, sy, tx, ty)
			err = errors.New("Vector coordinates are out of map bounds." + txt)
	}
	length := DistanceBetween(sx, sy, tx, ty)
	values := make([]bool, length)
	newVector := &Vector{sx, sy, tx, ty, values,
	[]int{}, []int{}}
	return newVector, err
}

func ComputeVector(vec *Vector) {
	vec.TilesX = nil
	vec.TilesY = nil
	sx, sy := vec.StartX, vec.StartY
	tx, ty := vec.TargetX, vec.TargetY
	steep := AbsoluteValue(ty - sy) > AbsoluteValue(tx - sx)
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

func PrintVector(vec *Vector) {
	blt.Layer(LookLayer)
	ch := "[color=white]" + vectorSymbol
	for i, _ := range(vec.TilesX) {
		x := vec.TilesX[i]
		y := vec.TilesY[i]
		blt.Print(x, y, ch)
	}
	blt.Refresh()
}
