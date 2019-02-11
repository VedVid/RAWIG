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
	MeleeDumbAI
	MeleePatherAI
)

func CreaturesTakeTurn(b Board, c Creatures) {
	/* Function CreaturesTakeTurn is supposed to handle all enemy creatures
	   actions: movement, attacking, etc.
	   It takes Board and Creatures as arguments.
	   Iterates through all Creatures slice, and calls HandleAI function with
	   specific parameters.
	   It skips NoAI and PlayerAI. */
	var ai int
	for _, v := range c {
		ai = v.AIType
		if ai == NoAI || ai == PlayerAI {
			continue
		}
		HandleAI(b, c, v, v.AIType)
	}
}

func HandleAI(b Board, cs Creatures, c *Creature, ai int) {
		switch ai {
		case MeleeDumbAI:
			if c.DistanceTo(cs[0].X, cs[0].Y) > 1 {
				c.MoveTowards(b, cs[0].X, cs[0].Y, ai)
			}
		case MeleePatherAI:
			// The same set of functions as for DumbAI.
			// Just for clarity.
			if c.DistanceTo(cs[0].X, cs[0].Y) > 1 {
				c.MoveTowards(b, cs[0].X, cs[0].Y, ai)
			}
	}
}
