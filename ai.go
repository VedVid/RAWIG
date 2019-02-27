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

import "fmt"

const (
	// Types of AI.
	NoAI = iota
	PlayerAI
	MeleeDumbAI
	MeleePatherAI
	RangedDumbAI
	RangedPatherAI
)

const (
	// Probability of triggering AI
	AITrigger = 92
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
		HandleAI(b, c, o, v)
		TriggerAI(b, c[0], v)
	}
}

func TriggerAI(b Board, p, c *Creature) {
	/* TriggerAI is function that takes Board and two Creatures as arguments.
	   First Creature is supposed to be player, second one - enemy.
	   Enemy with AITriggered set to false will ignore the player existence.
	   AITrigger is probability to notice (and, therefore, switch AITriggered)
	   player if is in monster's FOV. */
	if IsInFOV(b, p.X, p.Y, c.X, c.Y) == true && RandInt(100) <= AITrigger {
		c.AITriggered = true
	}
}

func HandleAI(b Board, cs Creatures, o Objects, c *Creature) {
	/* HandleAI is robust function that takes Board, Creatures, Objects,
	   and specific Creature as arguments. The most notable argument is
	   the last one - behavior of this entity will be decided in function body.
	   Its behavior will be decided regarding to AIType.
	   This function is very big and *wet*, but it is here to stay, for a while,
	   at least. I thought about code duplication removal by introducing one
	   generic function that would take Creature as argument, and - after
	   AIType check - would use proper HandleMeleeDumbAI (etc.) functions; or
	   would start with available weapons check. (One may want to peek at
	   issue #98 in repo - https://github.com/VedVid/RAWIG/issues/98 ).
	   But, on the other hand, ai has so many variations and edge cases that
	   unifying monster's behavior would result in smaller flexibility. */
	ai := c.AIType
	switch ai {
	case MeleeDumbAI:
		if c.AITriggered == true {
			if c.DistanceTo(cs[0].X, cs[0].Y) > 1 {
				c.MoveTowards(b, cs, cs[0].X, cs[0].Y, ai)
			} else {
				c.AttackTarget(cs[0])
			}
		} else {
			dx := RandRange(-1, 1)
			dy := RandRange(-1, 1)
			c.Move(dx, dy, b)
		}
	case MeleePatherAI:
		// The same set of functions as for DumbAI.
		// Just for clarity.
		if c.AITriggered == true {
			if c.DistanceTo(cs[0].X, cs[0].Y) > 1 {
				c.MoveTowards(b, cs, cs[0].X, cs[0].Y, ai)
			} else {
				c.AttackTarget(cs[0])
			}
		} else {
			dx := RandRange(-1, 1)
			dy := RandRange(-1, 1)
			c.Move(dx, dy, b)
		}
	case RangedDumbAI:
		if c.AITriggered == true {
			if c.Equipment[SlotWeaponPrimary] != nil {
				// Use primary ranged weapon.
				if c.DistanceTo(cs[0].X, cs[0].Y) >= FOVLength-1 {
					// TODO:
					// For now, every ranged skill has range equal to FOVLength-1
					// but it should change in future.
					c.MoveTowards(b, cs, cs[0].X, cs[0].Y, ai)
				} else {
					// DumbAI will not check if target is valid
					vec, err := NewVector(c.X, c.Y, cs[0].X, cs[0].Y)
					if err != nil {
						fmt.Println(err)
					}
					_ = ComputeVector(vec)
					_, _, target, _ := ValidateVector(vec, b, cs, o)
					if target != nil {
						c.AttackTarget(target)
					}
				}
			} else if c.Equipment[SlotWeaponSecondary] != nil {
				// Use secondary ranged weapon.
				if c.DistanceTo(cs[0].X, cs[0].Y) >= FOVLength-1 {
					// TODO:
					// For now, every ranged skill has range equal to FOVLength-1
					// but it should change in future.
					c.MoveTowards(b, cs, cs[0].X, cs[0].Y, ai)
				} else {
					// DumbAI will not check if target is valid
					vec, err := NewVector(c.X, c.Y, cs[0].X, cs[0].Y)
					if err != nil {
						fmt.Println(err)
					}
					_ = ComputeVector(vec)
					_, _, target, _ := ValidateVector(vec, b, cs, o)
					if target != nil {
						c.AttackTarget(target)
					}
				}
			} else {
				if c.DistanceTo(cs[0].X, cs[0].Y) > 1 {
					c.MoveTowards(b, cs, cs[0].X, cs[0].Y, ai)
				} else {
					c.AttackTarget(cs[0])
				}
			}
		} else {
			dx := RandRange(-1, 1)
			dy := RandRange(-1, 1)
			c.Move(dx, dy, b)
		}
	case RangedPatherAI: // It will depend on ranged weapons and equipment implementation
		if c.AITriggered == true {
			if c.Equipment[SlotWeaponPrimary] != nil {
				if c.DistanceTo(cs[0].X, cs[0].Y) >= FOVLength-1 {
					// TODO:
					// For now, every ranged skill has range equal to FOVLength-1
					// but it should change in future.
					c.MoveTowards(b, cs, cs[0].X, cs[0].Y, ai)
				} else {
					vec, err := NewVector(c.X, c.Y, cs[0].X, cs[0].Y)
					if err != nil {
						fmt.Println(err)
					}
					_ = ComputeVector(vec)
					_, _, target, _ := ValidateVector(vec, b, cs, o)
					if target != cs[0] {
						c.MoveTowards(b, cs, cs[0].X, cs[0].Y, ai)
					} else {
						c.AttackTarget(target)
					}
				}
			} else if c.Equipment[SlotWeaponSecondary] != nil {
				if c.DistanceTo(cs[0].X, cs[0].Y) >= FOVLength-1 {
					// TODO:
					// For now, every ranged skill has range equal to FOVLength-1
					// but it should change in future.
					c.MoveTowards(b, cs, cs[0].X, cs[0].Y, ai)
				} else {
					vec, err := NewVector(c.X, c.Y, cs[0].X, cs[0].Y)
					if err != nil {
						fmt.Println(err)
					}
					_ = ComputeVector(vec)
					_, _, target, _ := ValidateVector(vec, b, cs, o)
					if target != cs[0] {
						c.MoveTowards(b, cs, cs[0].X, cs[0].Y, ai)
					} else {
						c.AttackTarget(target)
					}
				}
			} else {
				if c.DistanceTo(cs[0].X, cs[0].Y) > 1 {
					c.MoveTowards(b, cs, cs[0].X, cs[0].Y, ai)
				} else {
					c.AttackTarget(cs[0])
				}
			}
		} else {
			dx := RandRange(-1, 1)
			dy := RandRange(-1, 1)
			c.Move(dx, dy, b)
		}

	}
}
