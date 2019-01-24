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
	newVector := &Vector{sx, sy, tx, ty, values}
	return newVector, err
}