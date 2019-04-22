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

func (c *Creature) AttackTarget(t *Creature, o *Objects) {
	/* Method Attack handles damage rolls for combat. Receiver "c" is attacker,
	   argument "t" is target. Critical hit is if attack roll is the same as receiver
	   attack attribute.
	   Result of attack is displayed in combat log, but messages need more polish. */
	att := RandInt(c.Attack) //basic attack roll
	att2 := 0                //critical bonus
	def := t.Defense         //opponent's defense
	dmg := 0                 //dmg delivered
	crit := false            //was it critical hit?
	if att == c.Attack {     //critical hit!
		crit = true
		att2 = RandInt(c.Attack)
	}
	switch {
	case att < def: // Attack score if lower than target defense.
		if crit == false {
			AddMessage("Attack deflected!")
		} else {
			dmg = att2 // Critical hit, but against heavily armored enemy.
			AddMessage("Critical hit! <heavily armored enemy>")
		}
	case att == def: // Attack score is equal to target defense.
		if crit == false {
			dmg = 1 // It's just a scratch...
			AddMessage("Attack successful, but it is just a scratch...")
		} else {
			dmg = att
			AddMessage("Critical hit, but it barely bypassed opponent's armor.")
		}
	case att > def: // Attack score is bigger than target defense.
		if crit == false {
			dmg = att
			AddMessage("Successful attack!")
		} else {
			dmg = att + att2 // Critical attack!
			AddMessage("Critical attack!")
		}
	}
	t.TakeDamage(dmg, o)
}

func (c *Creature) TakeDamage(dmg int, o *Objects) {
	/* Method TakeDamage has *Creature as receiver and takes damage integer
	   as argument. dmg value is deducted from Creature current HP.
	   If HPCurrent is below zero after taking damage, Creature dies. */
	c.HPCurrent -= dmg
	if c.HPCurrent <= 0 {
		c.Die(o)
	}
}
