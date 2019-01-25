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
	"fmt"
)

func (c *Creature) Look(b Board, o Objects, cs Creatures) {
	/* Look is method of Creature (that is supposed to be player).
	   It has to take Board, "global" Objects and Creatures as arguments,
	   because function PrintVector need to call RenderAll function.
	   At first, Look creates new para-vector, with player coords as
	   starting point, and dynamic end position.
	   Then ComputeVector checks what tiles are present
	   between Start and End, and adds their coords to vector values.
	   Line from Vector is drawn, then game waits for player input,
	   that will change position of "looking" cursors.
	   Loop breaks with Escape key as input. */
	startX, startY := c.X, c.Y
	targetX, targetY := startX, startY
	for {
		vec, err := NewVector(startX, startY, targetX, targetY)
		if err != nil {
			fmt.Println(err)
		}
		ComputeVector(vec)
		_ = ValidateVector(vec, b)
		PrintVector(vec, VectorColorNeutral, VectorColorNeutral, b, o, cs)
		key := blt.Read()
		if key == blt.TK_ESCAPE {
			break
		}
		switch key {
		case blt.TK_UP:
			targetY--
		case blt.TK_RIGHT:
			targetX++
		case blt.TK_DOWN:
			targetY++
		case blt.TK_LEFT:
			targetX--
		}
	}
}

func (c *Creature) Target(b Board, o Objects, cs Creatures) {
	length := FOVLength //hardcoded for now; will be passed as argument later
	startX, startY := c.X, c.Y
	targets := c.MonstersInFov(b, cs)
	targetable := MonstersInRange(targets, length) //use ValidateVector
}
