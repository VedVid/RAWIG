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
	"strconv"
	"unicode/utf8"

	blt "bearlibterminal"
)

func NewPlayer(layer, x, y int, character, name, color, colorDark string,
	alwaysVisible, blocked, blocksSight, triggered bool, ai, hp, attack,
	defense int, equipment EquipmentComponent) (*Creature, error) {
	/* Function NewPlayer takes all values necessary by its struct,
	   and creates then returns pointer to Creature;
	   so, it is basically NewCreature function. */
	var err error
	if layer < 0 {
		txt := LayerError(layer)
		err = errors.New("Player layer is smaller than 0." + txt)
	}
	if x < 0 || x >= MapSizeX || y < 0 || y >= MapSizeY {
		txt := CoordsError(x, y)
		err = errors.New("Player coords is out of window range." + txt)
	}
	if utf8.RuneCountInString(character) != 1 {
		txt := CharacterLengthError(character)
		err = errors.New("Player character string length is not equal to 1." + txt)
	}
	if triggered != false {
		err = errors.New("Warning: Player should not be triggered!")
	}
	if ai != PlayerAI {
		txt := PlayerAIError(ai)
		err = errors.New("Warning: Player AI is supposed to be " +
			strconv.Itoa(PlayerAI) + "." + txt)
	}
	if hp < 0 {
		txt := InitialHPError(hp)
		err = errors.New("Player HPMax is smaller than 0." + txt)
	}
	if attack < 0 {
		txt := InitialAttackError(attack)
		err = errors.New("Player attack value is smaller than 0." + txt)
	}
	if defense < 0 {
		txt := InitialDefenseError(defense)
		err = errors.New("Player defense value is smaller than 0." + txt)
	}
	playerBasicProperties := BasicProperties{layer, x, y, character, name, color,
		colorDark}
	playerVisibilityProperties := VisibilityProperties{alwaysVisible}
	playerCollisionProperties := CollisionProperties{blocked, blocksSight}
	playerFighterProperties := FighterProperties{ai, triggered, hp, hp, attack, defense}
	playerNew := &Creature{playerBasicProperties, playerVisibilityProperties,
		playerCollisionProperties, playerFighterProperties,
		equipment}
	return playerNew, err
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
	eq := GetEquippablesFromInventory(p)
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
	turnSpent, err := p.EquipItem(eq[option], slot)
	if err != nil {
		fmt.Println(err)
	}
	return turnSpent
}
