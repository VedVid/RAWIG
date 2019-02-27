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
	"errors"
	"fmt"
	"unicode/utf8"
)

const (
	// Special characters.
	CorpseChar = "%"
)

type Creature struct {
	/* Creatures are living objects that
	   moves, attacks, dies, etc. */
	BasicProperties
	VisibilityProperties
	CollisionProperties
	FighterProperties
	EquipmentComponent
}

// Creatures holds every creature on map.
type Creatures []*Creature

func NewCreature(x, y int, monsterFile string) (*Creature, error) {
	/* NewCreature is function that returns new Creature from
	   json file passed as argument. It replaced old code that
	   was encouraging hardcoding data in go files.
	   Errors returned by json package are not very helpful, and
	   hard to work with, so there is lazy panic for them. */
	var monster = &Creature{}
	err := CreatureFromJson(CreaturesPathJson+monsterFile, monster)
	if err != nil {
		fmt.Println(err)
		panic(-1)
	}
	monster.X, monster.Y = x, y
	var err2 error
	if monster.Layer < 0 {
		txt := LayerError(monster.Layer)
		err2 = errors.New("Creature layer is smaller than 0." + txt)
	}
	if monster.Layer != CreaturesLayer {
		txt := LayerWarning(monster.Layer, CreaturesLayer)
		err2 = errors.New("Creature layer is not equal to CreaturesLayer constant." + txt)
	}
	if monster.X < 0 || monster.X >= MapSizeX || monster.Y < 0 || monster.Y >= MapSizeY {
		txt := CoordsError(monster.X, monster.Y)
		err2 = errors.New("Creature coords is out of window range." + txt)
	}
	if utf8.RuneCountInString(monster.Char) != 1 {
		txt := CharacterLengthError(monster.Char)
		err2 = errors.New("Creature character string length is not equal to 1." + txt)
	}
	if monster.HPMax < 0 {
		txt := InitialHPError(monster.HPMax)
		err2 = errors.New("Creature HPMax is smaller than 0." + txt)
	}
	if monster.Attack < 0 {
		txt := InitialAttackError(monster.Attack)
		err2 = errors.New("Creature attack value is smaller than 0." + txt)
	}
	if monster.Defense < 0 {
		txt := InitialDefenseError(monster.Defense)
		err = errors.New("Creature defense value is smaller than 0." + txt)
	}
	return monster, err2
}

func (c *Creature) MoveOrAttack(tx, ty int, b Board, all Creatures) bool {
	/* Method MoveOrAttack decides if Creature will move or attack other Creature;
	   It has *Creature receiver, and takes tx, ty (coords) integers as arguments,
	   and map of current level, and list of all Creatures.
	   Starts by target that is nil, then iterates through Creatures. If there is
	   Creature on targeted tile, that Creature becomes new target for attack.
	   Otherwise, Creature moves to specified Tile.
	   It's supposed to take player as receiver (attack / moving enemies is
	   handled differently - check ai.go and combat.go). */
	var target *Creature
	turnSpent := false
	for i, _ := range all {
		if all[i].X == c.X+tx && all[i].Y == c.Y+ty {
			if all[i].HPCurrent > 0 {
				target = all[i]
				break
			}
		}
	}
	if target != nil {
		c.AttackTarget(target)
		turnSpent = true
	} else {
		turnSpent = c.Move(tx, ty, b)
	}
	return turnSpent
}

func (c *Creature) Move(tx, ty int, b Board) bool {
	/* Move is method of Creature; it takes target x, y as arguments;
	   check if next move won't put Creature off the screen, then updates
	   Creature coords. */
	turnSpent := false
	newX, newY := c.X+tx, c.Y+ty
	if newX >= 0 &&
		newX <= MapSizeX-1 &&
		newY >= 0 &&
		newY <= MapSizeY-1 {
		if b[newX][newY].Blocked == false {
			c.X = newX
			c.Y = newY
			turnSpent = true
		}
	}
	return turnSpent
}

func (c *Creature) PickUp(o *Objects) bool {
	/* PickUp is method that has *Creature as receiver
	   and slice of *Object as argument.
	   Creature tries to pick object up.
	   If creature stands on object that is possible to pick,
	   object is added to c's inventory, and removed
	   from "global" slice of objects.
	   Picking objects up takes turn only if it is
	   successful attempt. */
	turnSpent := false
	obj := *o
	for i := 0; i < len(obj); i++ {
		if obj[i].X == c.X && obj[i].Y == c.Y && obj[i].Pickable == true {
			if c.AIType == PlayerAI {
				AddMessage("You found " + obj[i].Name + ".")
			}
			c.Inventory = append(c.Inventory, obj[i])
			copy(obj[i:], obj[i+1:])
			obj[len(obj)-1] = nil
			*o = obj[:len(obj)-1]
			turnSpent = true
			break
		}
	}
	return turnSpent
}

func (c *Creature) DropFromInventory(objects *Objects, index int) bool {
	/* Drop is method that has Creature as receiver and takes
	   "global" list of objects as main argument, and additional
	   integer that is index of item to be dropped from c's Inventory.
	   At first, turnSpent is set to false, to make it true
	   at the end of function. It may be considered as obsolete WET,
	   because 'return true' would be sufficient, but it is
	   a bit more readable now.
	   Objs is dereferenced objects and it is absolutely necessary
	   to do any actions on these objects.
	   Drop do two things:
	   at first, it adds specific item to the game map,
	   then it removes this item from its owner Inventory. */
	turnSpent := false
	objs := *objects
	if c.AIType == PlayerAI {
		AddMessage("You dropped " + c.Inventory[index].Name + ".")
	}
	// Add item to the map.
	object := c.Inventory[index]
	object.X, object.Y = c.X, c.Y
	objs = append(objs, object)
	*objects = objs
	// Then remove item from inventory.
	copy(c.Inventory[index:], c.Inventory[index+1:])
	c.Inventory[len(c.Inventory)-1] = nil
	c.Inventory = c.Inventory[:len(c.Inventory)-1]
	turnSpent = true
	return turnSpent
}

func (c *Creature) DropFromEquipment(objects *Objects, slot int) bool {
	/* DropFromEquipment is method of *Creature that takes "global" objects,
	   and int (as index) as arguments, and returns bool (result depends if
	   action was successful, therefore if took a turn).
	   This function is very similar to DropFromInventory, but is kept
	   due to explicitness.
	   The difference is that Equipment checks Equipment index, not
	   specific object, so additionally checks for nils, and instead of
	   removing item from slice, makes it nil.
	   This behavior is important, because while Inventory is "dynamic"
	   slice, Equipment is supposed to be "fixed size" - slots are present
	   all the time, but the can be empty (ie nil) or occupied (ie object). */
	turnSpent := false
	objs := *objects
	object := c.Equipment[slot]
	if object == nil {
		return turnSpent // turn is not spent because there is no object to drop
	}
	// else {
	if c.AIType == PlayerAI {
		AddMessage("You removed and dropped " + object.Name + ".")
	}
	// add item to map
	object.X, object.Y = c.X, c.Y
	objs = append(objs, object)
	*objects = objs
	// then remove from slot
	c.Equipment[slot] = nil
	turnSpent = true
	return turnSpent
}

func (c *Creature) EquipItem(o *Object, slot int) (bool, error) {
	/* EquipItem is method of *Creature that takes *Object and int (that is
	   indicator to index of Equipment slot) as arguments; it returns
	   bool and error.
	   At first, EquipItem checks for errors:
	    - if object to equip exists
	    - if this equipment slot is not occupied
	   then equips item and removes it from inventory. */
	var err error
	if o == nil {
		txt := EquipNilError(c)
		err = errors.New("Creature tried to equip *Object that was nil." + txt)
	}
	if c.Equipment[slot] != nil {
		txt := EquipSlotNotNilError(c, slot)
		err = errors.New("Creature tried to equip item into already occupied slot." + txt)
	}
	if o.Slot != slot {
		txt := EquipWrongSlotError(o.Slot, slot)
		err = errors.New("Creature tried to equip item into wrong slot." + txt)
	}
	turnSpent := false
	// Equip item...
	c.Equipment[slot] = o
	// ...then remove it from inventory.
	index, err := FindObjectIndex(o, c.Inventory)
	if err != nil {
		fmt.Println(err)
	}
	copy(c.Inventory[index:], c.Inventory[index+1:])
	c.Inventory[len(c.Inventory)-1] = nil
	c.Inventory = c.Inventory[:len(c.Inventory)-1]
	if c.AIType == PlayerAI {
		AddMessage("You equipped " + o.Name + ".")
	}
	turnSpent = true
	return turnSpent, err
}

func (c *Creature) DequipItem(slot int) (bool, error) {
	/* DequipItem is method of Creature. It is called when receiver is about
	   to dequip weapon from "ready" equipment slot.
	   At first, weapon is added to Inventory, then Equipment slot is set to nil. */
	var err error
	if c.Equipment[slot] == nil {
		txt := DequipNilError(c, slot)
		err = errors.New("Creature tried to DequipItem that was nil." + txt)
	}
	if c.AIType == PlayerAI {
		AddMessage("You dequipped " + c.Equipment[slot].Name + ".")
	}
	turnSpent := false
	c.Inventory = append(c.Inventory, c.Equipment[slot]) //adding items to inventory should have own function, that will check "bounds" of inventory
	c.Equipment[slot] = nil
	turnSpent = true
	return turnSpent, err
}

func (c *Creature) Die() {
	/* Method Die is called when Creature's HP drops below zero.
	   Die() has *Creature as receiver.
	   Receiver properties changes to fit better to corpse. */
	c.Layer = DeadLayer
	c.Name = "corpse of " + c.Name
	c.Color = "dark red"
	c.ColorDark = "dark red"
	c.Char = CorpseChar
	c.Blocked = false
	c.BlocksSight = false
	c.AIType = NoAI
	ZeroLastTarget(c)
}

func FindMonsterByXY(x, y int, c Creatures) *Creature {
	/* Function FindMonsterByXY takes desired coords and list
	   of all available creatures. It iterates through this list,
	   and returns nil or creature that occupies specified coords. */
	var monster *Creature
	for i := 0; i < len(c); i++ {
		if x == c[i].X && y == c[i].Y {
			monster = c[i]
			break
		}
	}
	return monster
}
