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

import "math"

const (
	// Values for handling field of view algorithm execution.
	FOVRays   = 360 // Whole area around player; it may not work properly with other values.
	FOVLength = 5   // Sight range.
	FOVStep   = 1
)

var (
	// Slices to store FOV Rays values.
	// Should be immutable, but Go doesn't support immutable variables
	sinBase = []float64{}
	cosBase = []float64{}
)

func InitializeFOVTables() {
	// Function InitializeFOVTables creates data for raycasting.
	for i := 0; i < FOVRays; i++ {
		x := math.Sin(float64(i) / (float64(180) / math.Pi))
		y := math.Cos(float64(i) / (float64(180) / math.Pi))
		sinBase = append(sinBase, x)
		cosBase = append(cosBase, y)
	}
}

func CastRays(b Board, sx, sy int) {
	/* Function castRays is simple raycasting function for turning tiles to
	   explored.
	   It casts (fovRays / fovStep) rays (bigger fovStep means faster but
	   more error-prone raycasting) from player to coordinates in fovLength range.
	   Source of algorithm:
	   http://www.roguebasin.com/index.php?title=Raycasting_in_python [20170712] */
	for i := 0; i < FOVRays; i += FOVStep {
		rayX := sinBase[i]
		rayY := cosBase[i]
		x := float64(sx)
		y := float64(sy)
		bx, by := RoundFloatToInt(x), RoundFloatToInt(y)
		b[bx][by].Explored = true
		for j := 0; j < FOVLength; j++ {
			x -= rayX
			y -= rayY
			if x < 0 || y < 0 || x > MapSizeX-1 || y > MapSizeY-1 {
				break
			}
			bx2, by2 := RoundFloatToInt(x), RoundFloatToInt(y)
			b[bx2][by2].Explored = true
			if b[bx2][by2].BlocksSight == true {
				break
			}
		}
	}
}

func IsInFOV(b Board, sx, sy, tx, ty int) bool {
	/* Function isInFOV checks if target (tx, ty) is in fov of source (sx, sy).
	   Returns true if tx, ty == sx, sy; otherwise, it casts (FOVRays / fovStep)
	   rays (bigger fovStep means faster but more error-prone algorithm)
	   from source to tiles in fovLength range;
	   stops if cell has BlocksSight bool set to true.
	   Source of algorithm:
	   http://www.roguebasin.com/index.php?title=Raycasting_in_python [20170712]. */
	if sx == tx && sy == ty {
		return true
	}
	if sx < tx-FOVLength || sx > tx+FOVLength ||
		sy < ty-FOVLength || sy > ty+FOVLength {
		return false
	}
	for i := 0; i < FOVRays; i += FOVStep {
		rayX := sinBase[i]
		rayY := cosBase[i]
		x := float64(sx)
		y := float64(sy)
		for j := 0; j < FOVLength; j++ {
			x -= rayX
			y -= rayY
			if x < 0 || y < 0 || x > MapSizeX-1 || y > MapSizeY-1 {
				break
			}
			bx, by := RoundFloatToInt(x), RoundFloatToInt(y)
			if bx == tx && by == ty {
				return true
			}
			if b[bx][by].BlocksSight == true {
				break
			}
		}
	}
	return false
}

func (c *Creature) MonstersInFov(b Board, cs Creatures) Creatures {
	/* MonstersInFov is method of Creature. It takes global map, and
	   slice of creatures, as argument.
	   At first, new (empty) slice of creatures is made, to store
	   these monsters that are in c's field of view.
	   Then function iterates through Creatures passed as argument, and
	   adds every monster that is in c's fov, skipping source. */
	var inFov = Creatures{}
	for i := 0; i < len(cs); i++ {
		v := cs[i]
		if v == c {
			continue
		}
		if v.HPCurrent <= 0 {
			continue
		}
		if IsInFOV(b, c.X, c.Y, v.X, v.Y) == true {
			inFov = append(inFov, cs[i])
		}
	}
	return inFov
}

func (c *Creature) ObjectsInFov(b Board, o Objects) Objects {
	/* ObjectsInFov is method of Creature that works similar to
	   MonstersInFov. It returns slice of Objects that are present
	   in c's field of view. */
	var inFov = Objects{}
	for i := 0; i < len(o); i++ {
		v := o[i]
		if IsInFOV(b, c.X, c.Y, v.X, v.Y) == true {
			inFov = append(inFov, o[i])
		}
	}
	return inFov
}

func GetAllStringsFromTile(x, y int, b Board, c Creatures, o Objects) []string {
	/* GetAllStringsFromTile is function that takes coordinates, global map,
	   Creatures and Objects as arguments. It creates and then returns slice of
	   strings that contains names of all things on specific tile. It skips
	   tile names if there are objects present ("You see Monster and Objects here."),
	   otherwise it returns name of tile ("You see floor here."). */
	var s = []string{}
	for _, vc := range c {
		if vc.X == x && vc.Y == y {
			s = append(s, vc.Name)
		}
	}
	for _, vo := range o {
		if vo.X == x && vo.Y == y {
			s = append(s, vo.Name)
		}
	}
	if len(s) != 0 {
		return s
	}
	s = append(s, b[x][y].Name)
	return s
}

func GetAllStringsInFovTile(sx, sy, tx, ty int, b Board, c Creatures, o Objects) []string {
	/* GetAllStringInFovTile is function that uses IsInFOV and GetAllStringsFromTile
	   to create slice of strings of objects in field of view. */
	var s = []string{}
	if IsInFOV(b, sx, sy, tx, ty) == true {
		return GetAllStringsFromTile(tx, ty, b, c, o)
	}
	return s
}

func GetAliveCreatureFromTile(x, y int, c Creatures) *Creature {
	/* Function GetAliveCreatureFromTile takes coords and slice of Creature
	   as arguments, and returns Creature.
	   It iterates through all Creatures and find one that occupies specified tile.
	   This function could use []*Creature instead of *Creature, but monsters
	   should not overlap anyway. */
	var cs *Creature
	for i := 0; i < len(c); i++ {
		if c[i].X == x && c[i].Y == y && c[i].HPCurrent > 0 {
			cs = c[i]
		}
	}
	return cs
}

func GetAllThingsFromTile(x, y int, b Board, c Creatures, o Objects) (*Tile, Creatures, Objects) {
	/* GetAllThingsFromTile is function that takes coordinates, global map,
	   Creatures and Objects as arguments. It creates slice of Creature and
	   slice of Object that occupy coords, and returns them.
	   If these slices are empty, it returns board tile. */
	var cs = Creatures{}
	for i := 0; i < len(c); i++ {
		if c[i].X == x && c[i].Y == y {
			cs = append(cs, c[i])
		}
	}
	var os = Objects{}
	for j := 0; j < len(o); j++ {
		if o[j].X == x && o[j].Y == y {
			os = append(os, o[j])
		}
	}
	if len(cs) != 0 || len(os) != 0 {
		return nil, cs, os
	}
	return b[x][y], cs, os // cs and os are nil.
}
