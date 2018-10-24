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
	/*Method Attack handles damage rolls for combat. Receiver "c" is attacker,
	argument "t" is target. Critical hit is if attack roll is the same as receiver
	attack attribute.
	Println calls will remain here until message log will be implemented.*/
	att := RandInt(c.Attack) //basic attack roll
	att2 := 0                //critical bonus
	def := t.Defense         //opponent's defense
	dmg := 0                 //dmg delivered
	crit := false            //was it critical hit?
	if att == c.Attack {     //critical hit!
		crit = true
		att2 = RandInt(c.Attack)
	}
	if att < def { //if attack score if lower than target defense
		if crit == false {
			fmt.Println("Attack deflected!")
		} else {
			dmg = att2 //critical hit, but against heavily armored enemy
			fmt.Println("Critical hit! <heavily armored enemy>")
		}
	} else if att == def { //if attack score is equal to target defense
		if crit == false {
			dmg = 1 //just a scratch...
			fmt.Println("Attack successful, but it is just a scratch...")
		} else {
			dmg = att
			fmt.Println("Critical hit, but it barely bypassed opponent's armor.")
		}
	} else { //if attack score is bigger than target defense
		if crit == false {
			dmg = att
			fmt.Println("Successful attack!")
		} else {
			dmg = att + att2 //critical attack!
			fmt.Println("Critical attack!")
		}
	}
	t.TakeDamage(dmg)
}

func (c *Creature) TakeDamage(dmg int) {
	/*Method TakeDamage has *Creature as receiver and takes damage integer
	as argument. dmg value is deducted from Creature current HP.
	If HPCurrent is below zero after taking damage, Creature dies.*/
	c.HPCurrent -= dmg
	if c.HPCurrent <= 0 {
		c.Die()
	}
}
