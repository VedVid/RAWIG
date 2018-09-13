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
	"math"
)

const (
	//values for handling field of view algorithm execution
	FOVRays   = 360 //whole area around player; it may not work properly with other values
	FOVLength = 5   //sight range
	FOVStep   = 1
)

var (
	sinBase = []float64{}
	cosBase = []float64{}
)

func InitializeFOVTables() {
	for i := 0; i < FOVRays; i++ {
		x := math.Sin(float64(i) / (float64(180) / math.Pi))
		y := math.Cos(float64(i) / (float64(180) / math.Pi))
		sinBase = append(sinBase, x)
		cosBase = append(cosBase, y)
	}
}

func CastRays(b Board, sx, sy int) {
	/* Function castRays is simple raycasting function for turning tiles to explored.
	   It casts (fovRays / fovStep) rays (bigger fovStep means faster but
	   more error-prone raycasting) from player to coordinates in fovLength range.
	   Source of algorithm:
	   http://www.roguebasin.com/index.php?title=Raycasting_in_python [20170712]*/
	for i := 0; i < FOVRays; i += FOVStep {
		rayX := sinBase[i]
		rayY := cosBase[i]
		x := float64(sx)
		y := float64(sy)
		t1, _ := FindTileByXY(b, RoundFloatToInt(x), RoundFloatToInt(y))
		t1.Explored = true
		for j := 0; j < FOVLength; j++ {
			x -= rayX
			y -= rayY
			if x < 0 || y < 0 || x > WindowSizeX-1 || y > WindowSizeY-1 {
				break
			}
			t2, _ := FindTileByXY(b, RoundFloatToInt(x), RoundFloatToInt(y))
			t2.Explored = true
			if t2.Blocked {
				break
			}
		}
	}
}

func IsInFOV(b Board, sx, sy, tx, ty int) bool {
	/* Function isInFOV checks if target (tx, ty) is in fov of source (sx, sy).
	   Returns true if tx, ty == sx, sy; otherwise, it casts (FOVRays / fovStep)
	   rays (bigger fovStep means faster but more error-prone algorithm)
	   from source to tiles in fovLength range; stops if cell is blocked.
	   Source of algorithm:
	   http://www.roguebasin.com/index.php?title=Raycasting_in_python [20170712].*/
	if sx == tx && sy == ty {
		return true
	}
	for i := 0; i < FOVRays; i += FOVStep {
		rayX := sinBase[i]
		rayY := cosBase[i]
		x := float64(sx)
		y := float64(sy)
		for j := 0; j < FOVLength; j++ {
			x -= rayX
			y -= rayY
			if x < 0 || y < 0 || x > WindowSizeX-1 || y > WindowSizeY-1 {
				break
			}
			bx, by := RoundFloatToInt(x), RoundFloatToInt(y)
			if bx == tx && by == ty {
				return true
			}
			t, _ := FindTileByXY(b, bx, by)
			if t.Blocked {
				break
			}
		}
	}
	return false
}
