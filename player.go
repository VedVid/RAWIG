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
	alwaysVisible, blocked, blocksSight bool, ai, hp, attack,
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
	playerFighterProperties := FighterProperties{ai,hp, hp, attack, defense}
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
		case ItemDequip:
			fmt.Println("Equipping items is not implemented yet. ")
			break Loop
		case ItemDrop:
			turnSpent = p.DropFromInventory(o, option)
			break Loop
		case ItemUse:
			var err2 error
			turnSpent, err2 = object.UseItem(p)
			if err2 != nil {
				fmt.Println(err2)
			}
			break Loop
		default:
			continue Loop
		}
	}
	return turnSpent
}

func (p *Creature) EquipmentMenu(o *Objects) bool {
	/* EquipmentMenu works as InventoryMenu, but at the start of loop
	   it checks all equipment slots if they are empty or not.
	   It is almost the same function as used in handling inventory,
	   but maybe it is worth to be explicit here. */
	turnSpent := false
	for {
		PrintEquipmentMenu(UIPosX, UIPosY, "Equipment: ", p.Equipment)
		key := blt.Read()
		option := KeyToOrder(key)
		if option == KeyToOrder(blt.TK_ESCAPE) {
			break
		} else if option < SlotMax {
			turnSpent = p.HandleEquipment(o, option)
		} else {
			continue
		}
	}
	return turnSpent
}

func (p *Creature) HandleEquipment(o *Objects, option int) bool {
	/* HandleEquipment is method of Creature (that is supposed to be player)
	   that calls EquipmentActions with proper player Slot, and
	   Slot int indicator, as arguments. */
	eq := p.Equipment[option]
	turnSpent := p.EquipmentActions(o, eq, option)
	return turnSpent
}

func (p *Creature) EquipmentActions(o *Objects, object *Object, slot int) bool {
	/* Method EquipmentActions works as InventoryActions but for Equipment.
	   Refer to InventoryActions for more detailed info, but remember that
	   Inventory and Equipment, even if using the same architecture, may
	   call different functions, for example for dropping stuff. */
	turnSpent := false
Loop:
	for {
		options := GatherEquipmentOptions(object)
		header := "Equipment: "
		if object != nil {
			header = object.Name
		}
		PrintMenu(UIPosX, UIPosY, header, options)
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
			turnSpent = HandleEquipping(p, object, slot)
			break Loop
		case ItemDrop:
			turnSpent = p.DropFromEquipment(o, slot)
			break Loop
		case ItemUse:
			var err error
			turnSpent, err = object.UseItem(p)
			if err != nil {
				fmt.Println(err)
			}
			break Loop
		default:
			continue Loop
		}
	}
	return turnSpent
}
