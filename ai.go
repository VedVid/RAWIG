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
	RangedDumbAI
	RangedPatherAI
)

func CreaturesTakeTurn(b Board, c Creatures, o Objects) {
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
		HandleAI(b, c, o, v, v.AIType)
	}
}

func HandleAI(b Board, cs Creatures, o Objects, c *Creature, ai int) {
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
		/*case RangedPatherAI: //it will depend on ranged weapons and equipment implementation
			if c.DistanceTo(cs[0].X, cs[0].Y) > >>MONSTER_EQUIPPED_RANGED<< {
				c.MoveTowards(b, cs[0].X, cs[0].Y, ai)
			} else {
				vec, err := NewVector(c.X, c.Y, cs[0].X, cs[0].Y)
				if err != nil {
					fmt.Println(err)
				}
				_ := ComputeVector(vec)
				_, _, target, _ := ValidateVector(vec, b, cs, o)
				if target != cs[0] {
					c.MoveTowards(b, cs[0].X, cs[0].Y, ai)
				} else {
					c.AttackTarget(target)
				}
			} */
		/*case RangedDumbAI:
			if c.DistanceTo(cs[0].X, cs[0].Y) > >>MONSTER_EQUIPPED_RANGED<< {
				c.MoveTowards(b, cs[0].X, cs[0].Y, ai)
			} else {
				// DumbAI will not check if target is valid
				vec, err := NewVector(c.X, c.Y, cs[0].X, cs[0].Y)
				if err != nil {
					fmt.Println(err)
				}
				_ := ComputeVector(vec)
				_, _, target, _ := ValidateVector(vec, b, cs, o)
				if target != nil {
					c.AttackTarget(target)
				}
			} */
	}
}
