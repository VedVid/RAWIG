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

import "fmt"

const (
	//ai types
	PlayerAI = iota
	DumbAI
)

func (c *Creature) Attack(all Creatures) {
	/*Method Attack is, for now, only placeholder. It's called by
	CreaturesTakeTurn, for Creature that can attack player.*/
	if c != all[0] {
		fmt.Println("Attack!")
	}
}

func CreaturesTakeTurn(b Board, c Creatures) {
	/*Function CreaturesTakeTurn is supposed to handle all enemy creatures
	actions: movement, attacking, etc.
	It takes Board and Creatures as arguments.
	Iterates through all Creatures slice, and handles creature behaviour:
	if distance between creature and player is bigger than 1, creature
	moves towards player.
	It uses switch for matching AIType and behaviour.
	At first, I wanted to use map[int]METHOD, but it's not easy to implement.*/
	for _, v := range c {
		if v.DistanceTo(c[0].X, c[0].Y) > 1 {
			switch v.AIType {
			case DumbAI:
				v.MoveTowardsDumb(b, c[0].X, c[0].Y)
			}
		} else {
			v.Attack(c)
		}
	}
}
