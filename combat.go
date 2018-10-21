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

func (c *Creature) AttackTarget(t *Creature) {
	/*Method Attack is, for now, only placeholder. It's called by
	CreaturesTakeTurn, for Creature that can attack player.*/
	att := RandInt(c.Attack) //basic attack roll
	att2 := 0                //critical bonus
	def := t.Defense         //opponet's defense
	dmg := 0                 //dmg delivered
	crit := false            //was it critical hit?
	if att == c.Attack {     //critical hit!
		crit = true
		att2 = RandInt(c.Attack)
	}
	if att < def {
		if crit == false {
			fmt.Println("Attack deflected!")
		} else {
			att = att2 //critical hit, but against heavily armored enemy
			fmt.Println("Critical hit! <heavily armored enemy>")
		}
	} else if att == def {
		if crit == false {
			att = 1 //just a scratch...
			fmt.Println("Attack successful, but it is just a scratch...")
		} else {
			//att = att
			fmt.Println("Critical hit, but it barely bypassed opponent's armor.")
		}
	} else {
		if crit == false {
			//att = att, normal attack
			fmt.Println("Successful attack!")
		} else {
			att += att2 //critical attack!
			fmt.Println("Critical attack!")
		}
	}
	dmg = att
	_ = dmg
}
