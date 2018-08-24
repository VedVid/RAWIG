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

import blt "bearlibterminal"

const (
	//does player action took a turn?
	takeTurn = iota
	didntTakeTurn
)

var (
	//directions; should be immutable!
	N = Direction{0, 1}
	E = Direction{1, 0}
	S = Direction{0, -1}
	W = Direction{-1, 0}
)

type Direction struct {
	//Direction is type introduced due to lack of tuples or pair in Go.
	X, Y int
}

func Controls(k int, p Creature) int {
	/*Function Controls is input handler; it takes integer k
	(keycodes are basically numbers, but creating new "type key int"
	is not convenient) and Monster p (which is player);
	Controls handle input, then returns integer value that depends
	if player spent turn by action or not.*/
	switch k {
	case blt.TK_UP:
		p.Move(N)
	case blt.TK_RIGHT:
		p.Move(E)
	case blt.TK_DOWN:
		p.Move(S)
	case blt.TK_LEFT:
		p.Move(W)
	}
	return takeTurn
}
