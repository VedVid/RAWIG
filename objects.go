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
	/* Slots for inventory handling.
	   Their order here is important, because it
	   will be order of slots in Equipment menu. */
	SlotNA = iota - 1

	SlotWeaponPrimary
	SlotWeaponSecondary
	SlotWeaponMelee

	SlotMax
)

var SlotStrings = map[int]string{
	SlotWeaponPrimary:   "weapon1",
	SlotWeaponSecondary: "weapon2",
	SlotWeaponMelee:     "weapon3",
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

func NewObject(x, y int, objectPath string) (*Object, error) {
	/* NewObject is function that returns new Creature from
	   json file passed as argument. It replaced old code that
	   was encouraging hardcoding data in go files.
	   Errors returned by json package are not very helpful, and
	   hard to work with, so there is lazy panic for them. */
	var object = &Object{}
	err := ObjectFromJson(ObjectsPathJson+objectPath, object)
	if err != nil {
		fmt.Println(err)
		panic(-1)
	}
	object.X, object.Y = x, y
	var err2 error
	if object.Layer < 0 {
		txt := LayerError(object.Layer)
		err = errors.New("Object layer is smaller than 0. " + txt)
	}
	if object.Layer != ObjectsLayer {
		txt := LayerWarning(object.Layer, ObjectsLayer)
		err2 = errors.New("Creature layer is not equal to CreaturesLayer constant." + txt)
	}
	if object.X < 0 || object.X >= MapSizeX || object.Y < 0 || object.Y >= MapSizeY {
		txt := CoordsError(object.X, object.Y)
		err = errors.New("Object coords is out of window range." + txt)
	}
	if utf8.RuneCountInString(object.Char) != 1 {
		txt := CharacterLengthError(object.Char)
		err = errors.New("Object character string length is not equal to 1." + txt)
	}
	if object.Consumable == true && object.Use == UseNA {
		txt := ConsumableWithoutUseError()
		err = errors.New("Object is consumable, but has undefined use case." + txt)
	}
	if (object.Equippable == false && object.Slot != SlotNA) ||
		(object.Equippable == true && object.Slot == SlotNA) {
		txt := EquippableSlotError(object.Equippable, object.Slot)
		err = errors.New("'equippable' and 'slot' values does not match." + txt)
	}
	if object.Equippable == true && object.Consumable == true {
		//TODO: temporary
		err = errors.New("For now, <equippable> and <consumable> should not exists at the same time.")
	}
	return object, err2
}

func GatherItemOptions(o *Object) ([]string, error) {
	/* Function GatherItemOptions takes pointer to specific Object
	   as argument and returns slice of strings that is list of
	   possible actions. ItemBack that is necessary, yet last value
	   to include, to provide way to close menu. */
	var options = []string{}
	var err error
	if o.Equippable == true {
		options = append(options, ItemEquip)
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

func GatherEquipmentOptions(o *Object) ([]string, error) {
	/* Function GatherEquipmentOptions takes pointer to specific Object
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

func GetEquippablesFromInventory(c *Creature, slot int) Objects {
	/* GetEquippablesFromInventory is function that takes *Creature as arguments
	   and returns []*Object.
	   It creates empty slice of Object pointers, then adds every *Object
	   from *Creature's Inventory that is not nil, and has Equippable bool set to true.
	   This function is used to create list of all equippable items from
	   someone's Inventory. */
	var eq = Objects{}
	for i := 0; i < len(c.Inventory); i++ {
		item := c.Inventory[i]
		if item != nil && item.Equippable == true && item.Slot == slot {
			eq = append(eq, item)
		}
	}
	return eq
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
		AddMessage("You used " + o.Name + ".")
		err2 := DestroyItem(o, c)
		if err2 != nil {
			fmt.Println(err2)
			// It could be case to set turnSpent to false again.
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
		index, err_ := FindObjectIndex(o, c.Inventory)
		if err_ != nil {
			err = err_ // It looks like ugly hack.
			txt := ItemToDestroyNotFoundError()
			fmt.Println(txt)
		} else {
			copy(c.Inventory[index:], c.Inventory[index+1:])
			c.Inventory[len(c.Inventory)-1] = nil
			c.Inventory = c.Inventory[:len(c.Inventory)-1]
		}
	}
	return err
}
