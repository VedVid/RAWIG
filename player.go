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
	   so, it's basically NewCreature function. */
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
	/* Inventory menu is method of Creature (that is supposed to be a player).
	   It calls PrintInventoryMenu (that have much better docstring).
	   It returns boolean value that depends if real action (like using /
	   dropping item) was performed. */
	PrintInventoryMenu(UIPosX, UIPosY, "Inventory:", p.Inventory)
	turnSpent := p.HandleInventory(o, KeyToOrder(blt.Read())) //it's ugly one-liner
	return turnSpent
}

func (p *Creature) HandleInventory(o *Objects, option int) bool {
	/* HandleInventory is method that has pointer to Creature as receiver,
	   but it is supposed to be player every time. It takes
	   slice of game objects and chosen option (that is index of item in Inventory)
	   as arguments.
	   If option is valid index, ie is not out of Inventory bounds, it calls
	   InventoryActions method for handling actions that are possible for
	   this specific item. */
	turnSpent := false
	if option <= len(p.Inventory)-1 { //valid input
		turnSpent = p.InventoryActions(o, option)
	}
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
	   - gets player input and checks if it's valid;
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
		options := GatherItemOptions(object)
		PrintMenu(UIPosX, UIPosY, object.Name, options)
		var chosenStr string
		chosenInt := KeyToOrder(blt.Read())
		if chosenInt > len(options)-1 {
			chosenStr = ItemPass
		} else {
			chosenStr = options[chosenInt]
		}
		switch chosenStr {
		case ItemEquip:
			fmt.Println("Equipping items is not implemented yet. ")
			break Loop
		case ItemDrop:
			turnSpent = p.Drop(o, option)
			break Loop
		case ItemUse:
			turnSpent = object.UseItem(p)
			break Loop
		case ItemBack:
			break Loop
		default:
			continue Loop
		}
	}
	return turnSpent
}

func (p *Creature) EquipmentMenu() bool {
	/* EquipmentMenu is method of Creature (that is supposed to be player)
	   that prints menu with all equipped objects.
	   Currently it returns false all the time, because there is no
	   other things to do with it than printing menu. */
	eq := GetAllSlots(p)
	PrintEquipmentMenu(UIPosX, UIPosY, "Equipment:", eq)
	turnSpent := p.HandleEquipment(KeyToOrder(blt.Read()))
	return turnSpent
}

func (p *Creature) HandleEquipment(option int) bool {
	turnSpent := false
	option++ // Minimal default option is 0; minimal proper slot iota is 1.
	Loop:
		for {
			switch option {
			case SlotWeapon:
				//turnSpent = p.EquipmentActions(p.SlotWeapon, SlotWeapon)
				break Loop
			default:
				continue Loop
			}
		}
	return turnSpent
}
