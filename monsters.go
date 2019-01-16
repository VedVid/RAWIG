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

import (
	"errors"
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

func NewCreature(layer, x, y int, character, name, color, colorDark string,
	alwaysVisible, blocked, blocksSight bool, ai, hp, attack,
	defense int, equipment EquipmentComponent) (*Creature, error) {
	/* Function NewCreature takes all values necessary by its struct,
	   and creates then returns pointer to Creature. */
	var err error
	if layer < 0 {
		txt := LayerError(layer)
		err = errors.New("Creature layer is smaller than 0." + txt)
	}
	if x < 0 || x >= MapSizeX || y < 0 || y >= MapSizeY {
		txt := CoordsError(x, y)
		err = errors.New("Creature coords is out of window range." + txt)
	}
	if utf8.RuneCountInString(character) != 1 {
		txt := CharacterLengthError(character)
		err = errors.New("Creature character string length is not equal to 1." + txt)
	}
	if hp < 0 {
		txt := InitialHPError(hp)
		err = errors.New("Creature HPMax is smaller than 0." + txt)
	}
	if attack < 0 {
		txt := InitialAttackError(attack)
		err = errors.New("Creature attack value is smaller than 0." + txt)
	}
	if defense < 0 {
		txt := InitialDefenseError(defense)
		err = errors.New("Creature defense value is smaller than 0." + txt)
	}
	creatureBasicProperties := BasicProperties{layer, x, y, character, name, color,
		colorDark}
	creatureVisibilityProperties := VisibilityProperties{alwaysVisible}
	creatureCollisionProperties := CollisionProperties{blocked, blocksSight}
	creatureFighterProperties := FighterProperties{ai, hp, hp, attack, defense}
	creatureNew := &Creature{creatureBasicProperties,
		creatureVisibilityProperties, creatureCollisionProperties,
		creatureFighterProperties, equipment}
	return creatureNew, err
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
	var turnSpent bool
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
	for i := 0; i <= len(obj); i++ {
		if obj[i].X == c.X && obj[i].Y == c.Y && obj[i].Pickable == true {
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
	var err error
	if o == nil {
		txt := EquipNilError(c)
		err = errors.New("Creature tried to equip *Object that was nil." + txt)
	}
	if c.Equipment[slot] != nil {
		txt := EquipSlotNotNilError(c, slot)
		err = errors.New("Creature tried to equip item into already occupied slot." + txt)
	}
	turnSpent := false
	// Equip item...
	c.Equipment[slot] = o
	// ...then remove it from inventory.
	//copy(c.Inventory[index:], c.Inventory[index+1:])
	//c.Inventory[len(c.Inventory)-1] = nil
	//c.Inventory = c.Inventory[:len(c.Inventory)-1]
	turnSpent = true
	return turnSpent, err
}

func (c *Creature) DequipItem(o *Object, slot int) (bool, error) {
	/* DequipItem is method of Creature. It is called when receiver is about
	   to dequip weapon from "ready" equipment slot.
	   At first, weapon is added to Inventory, then Equipment slot is set to nil. */
	var err error
	if o == nil {
		txt := DequipNilError(c, slot)
		err = errors.New("Creature tried to DequipItem that was nil." + txt)
	}
	turnSpent := false
	c.Inventory = append(c.Inventory, o) //adding items to inventory should have own function, that will check "bounds" of inventory
	c.Equipment[slot] = nil
	turnSpent = true
	return turnSpent, err
}

func (c *Creature) Die() {
	/* Method Die is called when Creature's HP drops below zero.
	   Die() has *Creature as receiver.
	   Receiver properties changes to fit better to corpse. */
	c.Layer = DeadLayer
	c.Color = "dark red"
	c.ColorDark = "dark red"
	c.Char = CorpseChar
	c.Blocked = false
	c.BlocksSight = false
	c.AIType = NoAI
}
