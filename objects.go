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
	"fmt"
	"unicode/utf8"
)

const (
	/* Slots for inventory handling.
	   Their order here is important, because it
	   will be order of slots in Equipemnt menu. */
	SlotNA = iota-1

	SlotWeaponPrimary
	SlotWeaponSecondary

	SlotMax
)

var SlotStrings = map[int]string{
	SlotWeaponPrimary: "weapon1",
	SlotWeaponSecondary: "weapon2",
}

const (
	// Use cases, mostly for consumables.
	UseNA = iota

	UseHeal
)

const (
	// Values for handling inventory actions.
	ItemPass   = "pass"
	ItemDrop   = "drop"
	ItemEquip  = "equip"
	ItemDequip = "dequip"
	ItemUse    = "use"
)

type Object struct {
	/* Objects are every other things on map;
	   statues, tables, chairs; but also weapons,
	   armor parts, etc. */
	BasicProperties
	VisibilityProperties
	CollisionProperties
	ObjectProperties
}

// Objects holds every object on map.
type Objects []*Object

func NewObject(layer, x, y int, character, name, color, colorDark string,
	alwaysVisible, blocked, blocksSight bool, pickable, equippable, consumable bool, slot, use int) (*Object, error) {
	/* Function NewObject takes all values necessary by its struct,
	   and creates then returns Object. */
	var err error
	if layer < 0 {
		txt := LayerError(layer)
		err = errors.New("Object layer is smaller than 0. " + txt)
	}
	if x < 0 || x >= MapSizeX || y < 0 || y >= MapSizeY {
		txt := CoordsError(x, y)
		err = errors.New("Object coords is out of window range." + txt)
	}
	if utf8.RuneCountInString(character) != 1 {
		txt := CharacterLengthError(character)
		err = errors.New("Object character string length is not equal to 1." + txt)
	}
	if consumable == true && use == UseNA {
		txt := ConsumableWithoutUseError()
		err = errors.New("Object is consumable, but has undefined use case." + txt)
	}
	if (equippable == false && slot != SlotNA) || (equippable == true && slot == SlotNA) {
		txt := EquippableSlotError(equippable, slot)
		err = errors.New("'equippable' and 'slot' values does not match." + txt)
	}
	if equippable == true && consumable == true {
		//TODO: temporary
		err = errors.New("For now, <equippable> and <consumable> should not exists at the same time.")
	}
	objectBasicProperties := BasicProperties{layer, x, y, character, name,color,
		colorDark}
	objectVisibilityProperties := VisibilityProperties{alwaysVisible}
	objectCollisionProperties := CollisionProperties{blocked, blocksSight}
	objectProperties := ObjectProperties{pickable, equippable, consumable,
	slot, use}
	objectNew := &Object{objectBasicProperties, objectVisibilityProperties,
		objectCollisionProperties, objectProperties}
	return objectNew, err
}

func GatherItemOptions(o *Object) ([]string, error) {
	/* Function GatherItemOptions takes pointer to specific Object
	   as argument and returns slice of strings that is list of
	   possible actions. ItemBack that is necessary, yet last value
	   to include, to provide way to close menu. */
	var options = []string{}
	var err error
	if o.Equippable == true {
		options = append(options, ItemDequip)
	}
	if o.Use != UseNA {
		options = append(options, ItemUse)
	}
	if o.Pickable == true {
		options = append(options, ItemDrop)
	}
	if len(options) == 0 {
		txt := ItemOptionsEmptyError()
		err = errors.New("Object " + o.Name + " has no valid properties." + txt)
	}
	return options, err
}

func GatherEquipmentOptions(o *Object) []string {
	/* GatherEquipmentOptions is function that takes object as argument
	   and returns slice of string.
	   If object is not nil, it calls GatherItemOptions to create
	   list of options based on object properties.
	   If object is nil, it creates slice with two options:
	   one for equipping item in empty slot, and one for going
	   back to previous menu. */
	var options = []string{}
	if o != nil {
		var err error
		options, err = GatherItemOptions(o)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		// there is no object in slot, so
		// print list of equippables
		options = append(options, ItemEquip)
	}
	return options
}

func HandleEquipping(c *Creature, o *Object, slot int) bool {
	/* HandleEquipping is helper function that is called in "Equipment menu".
	   If specific weapon slot (c.Equipment[slot], or o) is nil, it means
	   that slot is empty and game should provide list of all equippable items
	   from inventory. Else, it unequips this item. */
	turnSpent := false
	if o == nil {
		turnSpent = c.HandleEquippables()
		return turnSpent
	} else {
		//dequip item from slot
		turnSpent = c.DequipItem(o, slot)
	}
	return turnSpent
}

func (o *Object) UseItem(c *Creature) (bool, error) {
	/* Method UseItem has Object as receiver and takes Creature as argument.
	   It uses Use value of receiver to determine what action will be performed.
	   If there is no valid o.Use, it breaks switch statement (need proper
	   error handling).
	   It tries to remove item from inventory by calling DestroyItem function,
	   but item will be removed only if its Consumable is set to true.
	   Returns turnSpent that is true, unless o.Use is invalid. */
	turnSpent := false
	var err error
	switch o.Use {
	case UseHeal:
		c.HPCurrent = c.HPMax
		turnSpent = true
	default:
		txt := UseItemError()
		err = errors.New("Item has wrong use case specified." + txt)
		break
	}
	if err == nil {
		err2 := DestroyItem(o, c)
		if err2 != nil {
			fmt.Println(err2)
		}
	}
	return turnSpent, err
}

func DestroyItem(o *Object, c *Creature) error {
	/* Function DestroyItem takes Object and Creature as arguments, and returns error.
	   At first, it iterates through Creature's Inventory, and creates an error if
	   proper index is not found. Otherwise, it removes item from inventory. */
	var err error
	if o.Consumable == true {
		index := -1
		for i := 0; i < len(c.Inventory); i++ {
			if c.Inventory[i] == o {
				index = i
				break
			}
		}
		if index < 0 {
			txt := ItemToDestroyNotFoundError()
			err = errors.New("Consumable to destroy is not found in inventory." + txt)
		} else {
			copy(c.Inventory[index:], c.Inventory[index+1:])
			c.Inventory[len(c.Inventory)-1] = nil
			c.Inventory = c.Inventory[:len(c.Inventory)-1]
		}
	}
	return err
}
