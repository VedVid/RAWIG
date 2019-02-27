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

import (
	blt "bearlibterminal"
)

func Controls(k int, p *Creature, b *Board, c *Creatures, o *Objects) bool {
	/* Function Controls is input handler.
	   It takes integer k (key codes are basically numbers,
	   but creating new "type key int" is not convenient)
	   and Creature p (which is player).
	   Controls handle input, then returns integer value that depends
	   if player spent turn by action or not. */
	turnSpent := false
	switch k {
	case blt.TK_UP:
		turnSpent = p.MoveOrAttack(0, -1, *b, *c)
	case blt.TK_RIGHT:
		turnSpent = p.MoveOrAttack(1, 0, *b, *c)
	case blt.TK_DOWN:
		turnSpent = p.MoveOrAttack(0, 1, *b, *c)
	case blt.TK_LEFT:
		turnSpent = p.MoveOrAttack(-1, 0, *b, *c)

	case blt.TK_F:
		turnSpent = p.Target(*b, *o, *c)
	case blt.TK_L:
		p.Look(*b, *o, *c) // Looking is free action.
	case blt.TK_G:
		turnSpent = p.PickUp(o)
	case blt.TK_I:
		turnSpent = p.InventoryMenu(o)
	case blt.TK_E:
		turnSpent = p.EquipmentMenu(o)
	}
	return turnSpent
}
