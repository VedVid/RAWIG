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

import "math"

const (
	//values for handling field of view algorithm execution
	FOVRays   = 360 //whole area around player; it may not work properly with other values
	FOVLength = 5   //sight range
	FOVStep   = 1
)

var (
	sintable = []float64{}
	costable = []float64{}
)

func InitializeFOVTables() {
	/*Function InitializeFOVTables fills sintable and costable with proper
	values; it's provided because it's less ugly than hardcoding all
	necessary floats. It has to be called in init function of main file*/
	for i := 0; i < FOVRays; i++ {
		x := math.Sin(float64(i) / (float64(180) / math.Pi))
		y := math.Cos(float64(i) / (float64(180) / math.Pi))
		sintable = append(sintable, x)
		costable = append(costable, y)
	}
}

func CastRays(b Board, sx, sy int) {
	/*Function CastRays is simple raycasting function for field of view.
	It casts (rays div step) rays from source, in specified range.
	Bigger step means faster execution, but also more errors and artifacts.
	Since Go is relatively fast language, it's safe to use small steps.
	It's translation of released to public domain python algorithm, published
	by [init. initd5@gmail.com] via roguebasin (access: 20180825):
	http://www.roguebasin.com/index.php?title=Raycasting_in_python
	This implementation uses floating point numbers to make it easy to
	adapt for even more precise raycasting.
	I may change raycasting to shadowcasting (or other algorithm) later.*/
	for i := 0; i < FOVRays; i += FOVStep {
		rayX := sintable[i]
		rayY := costable[i]
		x := float64(sx)
		y := float64(sy)
		t1 := FindTileByXY(b, sx, sy)
		t1.Explored = true
		for j := 0; j < FOVLength; j++ {
			x -= rayX
			y -= rayY
			if x < 0 || y < 0 || x >= WindowSizeX || y >= WindowSizeY {
				break
			}
			bx, by := int(math.Round(x)), int(math.Round(y))
			t2 := FindTileByXY(b, bx, by)
			t2.Explored = true
			if t2.Blocked == true {
				break
			}
		}
	}
}

func IsInFOV(b Board, sx, sy, tx, ty int) bool {
	/*Function IsInFOV checks if target coords (tx, ty) are within
	source (sx, sy) field of view. Returns true if target is on the same tile
	as source; otherwise, it casts (FOVRays / FOVStep) rays from source to
	tiles in FOVLength range and stops if cells is blocked.
	It's nasty code duplication from CastRays function and it may be
	addressed later, but I think it's readable enough.
	Part of translation of init's python raycasting algorithm.*/
	if sx == tx && sy == ty {
		return true
	}
	for i := 0; i < FOVRays; i += FOVStep {
		rayX := sintable[i]
		rayY := sintable[i]
		x := float64(sx)
		y := float64(sy)
		for j := 0; j < FOVLength; j++ {
			x -= rayX
			y -= rayY
			if x < 0 || y < 0 || x >= WindowSizeX || y >= WindowSizeY {
				break
			}
			bx, by := int(math.Round(x)), int(math.Round(y))
			if bx == tx && by == ty {
				return true
			}
			t := FindTileByXY(b, bx, by)
			if t.Blocked == true {
				break
			}
		}
	}
	return false
}
