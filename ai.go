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

const (
	// Types of AI.
	NoAI = iota
	PlayerAI
	DumbAI
	PatherAI
)

func CreaturesTakeTurn(b Board, c Creatures) {
	/* Function CreaturesTakeTurn is supposed to handle all enemy creatures
	   actions: movement, attacking, etc.
	   It takes Board and Creatures as arguments.
	   Iterates through all Creatures slice, and handles creature behavior:
	   if distance between creature and player is bigger than 1, creature
	   moves towards player. Else, it attacks.
	   It passed Creature's ai type as argument of MoveTowards to force
	   different behavior. */
	var ai int
	for _, v := range c {
		ai = v.AIType
		if ai == NoAI || ai == PlayerAI {
			continue
		} else {
			if v.DistanceTo(c[0].X, c[0].Y) > 1 {
				v.MoveTowards(b, c[0].X, c[0].Y, ai)
			} else {
				v.AttackTarget(c[0])
			}
		}
	}
}
