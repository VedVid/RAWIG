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
	//ai types
	NoAI = iota
	PlayerAI
	DumbAI
	PatherAI
)

func CreaturesTakeTurn(b Board, c Creatures) {
	/*Function CreaturesTakeTurn is supposed to handle all enemy creatures
	actions: movement, attacking, etc.
	It takes Board and Creatures as arguments.
	Iterates through all Creatures slice, and handles creature behavior:
	if distance between creature and player is bigger than 1, creature
	moves towards player. Else, it attacks.
	It uses switch for matching AIType and behavior. Skips Creatures with NoAI
	(ie corpses) and PlayerAI.
	At first, I wanted to use map[int]METHOD, but it's not easy to implement.*/
	for _, v := range c {
		switch v.AIType {
		case NoAI:
			continue
		case PlayerAI:
			continue
		case DumbAI:
			v.UseDumbAI(b, c, c[0].X, c[0].Y)
		case PatherAI:
			v.UsePatherAI(b, c, c[0].X, c[0].Y)
		default:
			continue
		}
	}
}

func (c *Creature) UseDumbAI(b Board, cs Creatures, tx, ty int) {
	if c.DistanceTo(cs[0].X, cs[0].Y) > 1 {
		c.MoveTowardsDumb(b, cs[0].X, cs[0].Y)
	} else {
		c.AttackTarget(cs[0])
	}
}

func (c *Creature) UsePatherAI(b Board, cs Creatures, tx, ty int) {
	if c.DistanceTo(cs[0].X, cs[0].Y) > 1 {
		c.MoveTowardsPath(b, cs[0].X, cs[0].Y)
	} else {
		c.AttackTarget(cs[0])
	}
}
