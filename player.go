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
	"strconv"
	"unicode/utf8"

	blt "bearlibterminal"
)

func NewPlayer(x, y int) (*Creature, error) {
	/* NewPlayer is function that returns new Creature
	   (that is supposed to be player) from json file passed as argument.
	   It replaced old code that was encouraging hardcoding data in go files.
	   Errors returned by json package are not very helpful, and
	   hard to work with, so there is lazy panic for them. */
	const playerPath = "./data/player/player.json"
	var player = &Creature{}
	err := CreatureFromJson(playerPath, player)
	if err != nil {
		fmt.Println(err)
		panic(-1)
	}
	player.X, player.Y = x, y
	var err2 error
	if player.Layer < 0 {
		txt := LayerError(player.Layer)
		err2 = errors.New("Creature layer is smaller than 0." + txt)
	}
	if player.Layer != PlayerLayer {
		txt := LayerWarning(player.Layer, PlayerLayer)
		err2 = errors.New("Creature layer is not equal to CreaturesLayer constant." + txt)
	}
	if player.X < 0 || player.X >= MapSizeX || player.Y < 0 || player.Y >= MapSizeY {
		txt := CoordsError(player.X, player.Y)
		err2 = errors.New("Creature coords is out of window range." + txt)
	}
	if utf8.RuneCountInString(player.Char) != 1 {
		txt := CharacterLengthError(player.Char)
		err2 = errors.New("Creature character string length is not equal to 1." + txt)
	}
	if player.AIType != PlayerAI {
		txt := PlayerAIError(player.AIType)
		err = errors.New("Warning: Player AI is supposed to be " +
			strconv.Itoa(PlayerAI) + "." + txt)
	}
	if player.HPMax < 0 {
		txt := InitialHPError(player.HPMax)
		err2 = errors.New("Creature HPMax is smaller than 0." + txt)
	}
	if player.Attack < 0 {
		txt := InitialAttackError(player.Attack)
		err2 = errors.New("Creature attack value is smaller than 0." + txt)
	}
	if player.Defense < 0 {
		txt := InitialDefenseError(player.Defense)
		err = errors.New("Creature defense value is smaller than 0." + txt)
	}
	return player, err2
}

func (p *Creature) InventoryMenu(o *Objects) bool {
	/* InventoryMenu is method of *Creature that takes *Objects as argument
	   and returns boolean value - indicator if action took turn or not.
	   It starts by loop that prints creature's (player) inventory and
	   waits for input. Then changes input to alphabetic order.
	   Handling input is simple - if input is within inventory range,
	   it invokes HandleInventory; if key is one point larger than
	   length of inventory, it backs menu (by breaking the loop);
	   otherwise, it loops. */
	turnSpent := false
	for {
		PrintInventoryMenu(UIPosX, UIPosY, "Inventory", p.Inventory)
		key := blt.Read()
		option := KeyToOrder(key)
		if option == KeyToOrder(blt.TK_ESCAPE) {
			break
		} else if option < len(p.Inventory) {
			turnSpent = p.HandleInventory(o, option)
		} else {
			continue
		}
	}
	return turnSpent
}

func (p *Creature) HandleInventory(o *Objects, option int) bool {
	/* HandleInventory is method that has pointer to Creature as receiver,
	   but it is supposed to be player every time. It takes
	   slice of game objects and chosen option (that is index of item in Inventory)
	   as arguments.
	   It calls InventoryActions method for handling actions that are possible
	   for specific item. */
	turnSpent := p.InventoryActions(o, option)
	return turnSpent
}

func (p *Creature) InventoryActions(o *Objects, option int) bool {
	/* InventoryActions is method that has *Creature as receiver
	   (that is supposed to be player) and takes *Objects and index
	   of specific item (ie integer) as arguments.
	   It loops rendering menu until proper input is provided.
	   Loop is pretty complicated:
	   - is labelled as Loop to make breaking simpler
	   - at first, it uses GatherItemOptions to group all actions
	     that are possible for this specific item
	   - gets player input and checks if is valid;
	     valid input is binded to chosenStr variable;
	     invalid is, before binding, transformed to ItemPass value
	   - switch expression is called:
	     * all cases are resolved by external functions
	     * equipping items is not implemented yet
	     * with invalid input, loop continues */
	turnSpent := false
	object := p.Inventory[option]
Loop:
	for {
		options, err1 := GatherItemOptions(object)
		if err1 != nil {
			fmt.Println(err1)
		}
		PrintMenu(UIPosX, UIPosY, object.Name, options)
		var chosenStr string
		chosenInt := KeyToOrder(blt.Read())
		if chosenInt == KeyToOrder(blt.TK_ESCAPE) {
			break Loop
		} else if chosenInt > len(options)-1 {
			chosenStr = ItemPass
		} else {
			chosenStr = options[chosenInt]
		}
		switch chosenStr {
		case ItemEquip:
			turnSpent = p.EquipFromInventory(object)
			break Loop
		case ItemDrop:
			turnSpent = p.DropFromInventory(o, option)
			break Loop
		case ItemUse:
			var err2 error
			turnSpent, err2 = object.UseItem(p)
			if err2 != nil {
				fmt.Println(err2)
				turnSpent = false
			}
			break Loop
		default:
			continue Loop
		}
	}
	return turnSpent
}

func (p *Creature) EquipFromInventory(o *Object) bool {
	/* EquipFromInventory is method of Creature (that is supposed to be player)
	   that takes Object (already chosen item from inventory) as argument, and
	   returns true if actions is success.
	   This method is used to equip item directly from inventory. */
	turnSpent := false
	for {
		PrintEquipmentMenu(UIPosX, UIPosY, "Equipment:", p.Equipment)
		key := blt.Read()
		option := KeyToOrder(key)
		if option == KeyToOrder(blt.TK_ESCAPE) {
			break
		} else if option < SlotMax {
			if p.Equipment[option] != nil {
				AddMessage("This slot is already occupied.")
				continue
			} else if option != o.Slot {
				AddMessage("You can't equip this here.")
				continue
			} else {
				var err error
				turnSpent, err = p.EquipItem(o, option)
				if err != nil {
					fmt.Println(err)
				}
			}
			break
		} else {
			continue
		}
	}
	return turnSpent
}

func (p *Creature) EquipmentMenu(o *Objects) bool {
	/* EquipmentMenu start similar to InventoryMenu - it prints Equipment
	   and waits for player input, then checks if input is valid.
	   If test will pass, it tries to dequip item from selected slot;
	   if this slot is already empty, it call EquippablesMenu to
	   provide list of all equippables items from Inventory. */
	turnSpent := false
	for {
		PrintEquipmentMenu(UIPosX, UIPosY, "Equipment: ", p.Equipment)
		key := blt.Read()
		option := KeyToOrder(key)
		if option == KeyToOrder(blt.TK_ESCAPE) {
			break
		} else if option < SlotMax {
			if p.Equipment[option] != nil {
				turnSpent = p.EquipmentActions(o, option)
			} else {
				turnSpent = p.EquippablesMenu(option)
			}
		} else {
			continue
		}
	}
	return turnSpent
}

func (p *Creature) EquipmentActions(o *Objects, slot int) bool {
	/* Method EquipmentActions works as InventoryActions but for Equipment.
	   Refer to InventoryActions for more detailed info, but remember that
	   Inventory and Equipment, even if using the same architecture, may
	   call different functions, for example for dropping stuff. */
	turnSpent := false
	object := p.Equipment[slot]
Loop:
	for {
		options, err1 := GatherEquipmentOptions(object)
		if err1 != nil {
			fmt.Println(err1)
		}
		PrintMenu(UIPosX, UIPosY, object.Name, options)
		var chosenStr string
		chosenInt := KeyToOrder(blt.Read())
		if chosenInt == KeyToOrder(blt.TK_ESCAPE) {
			break Loop
		} else if chosenInt > len(options)-1 {
			chosenStr = ItemPass
		} else {
			chosenStr = options[chosenInt]
		}
		switch chosenStr {
		case ItemDequip:
			var err2 error
			turnSpent, err2 = p.DequipItem(slot)
			if err2 != nil {
				fmt.Println(err2)
				turnSpent = false
			}
			break Loop
		case ItemDrop:
			turnSpent = p.DropFromEquipment(o, slot)
			break Loop
		case ItemUse:
			var err3 error
			turnSpent, err3 = object.UseItem(p)
			if err3 != nil {
				fmt.Println(err3)
				turnSpent = false
			}
			break Loop
		default:
			continue Loop
		}
	}
	return turnSpent
}

func (p *Creature) EquippablesMenu(slot int) bool {
	/* EquippablesMenu is method of Creature (that is supposed to be player).
	   It returns true if action was success, false otherwise.
	   At start, GetEquippablesFromInventory is called to create new slice
	   of equippables separated from inventory. Then function waits for player
	   input and, if possible, calls HandleEquippables to fill empty slot. */
	turnSpent := false
	eq := GetEquippablesFromInventory(p, slot)
	for {
		PrintEquippables(UIPosX, UIPosY, "Equippables: ", eq)
		key := blt.Read()
		option := KeyToOrder(key)
		if option == KeyToOrder(blt.TK_ESCAPE) {
			break
		} else if option < len(eq) {
			turnSpent = p.HandleEquippables(eq, option, slot)
			break
		} else {
			continue
		}
	}
	return turnSpent
}

func (p *Creature) HandleEquippables(eq Objects, option, slot int) bool {
	/* HandleEquippables is method of Creature (player) that takes
	   list of equippables (ie slice of *Object), and two ints as arguments.
	   It returns true if action is success.
	   The body if this function calls EquipItem and handles it error. */
	turnSpent := false
	var err error
	turnSpent, err = p.EquipItem(eq[option], slot)
	if err != nil {
		fmt.Println(err)
	}
	return turnSpent
}
