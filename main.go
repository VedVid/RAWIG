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
	blt "bearlibterminal"
	"fmt"
	"math/rand"
	"os"
	"time"
)

var MsgBuf = []string{}
var LastTarget *Creature

func main() {
	var cellsPtr = new(Board)
	var objsPtr = new(Objects)
	var actorsPtr = new(Creatures)
	StartGame(cellsPtr, actorsPtr, objsPtr)
	cells := *cellsPtr
	objs := *objsPtr
	actors := *actorsPtr
	for {
		RenderAll(cells, objs, actors)
		key := blt.Read()
		if key == blt.TK_ESCAPE {
			err20 := SaveGame(cells, actors, objs)
			if err20 != nil {
				fmt.Println(err20)
			}
			break
		} else if actors[0].HPCurrent <= 0 {
			err30 := DeleteSaves()
			if err30 != nil {
				fmt.Println(err30)
				panic(-1)
			}
			break
		} else {
			turnSpent := Controls(key, actors[0], cells, actors, &objs)
			if turnSpent == true {
				CreaturesTakeTurn(cells, actors, objs)
			}
		}
	}
	blt.Close()
}

func NewGame(b *Board, c *Creatures, o *Objects) {
	slot, _ := NewObject(ObjectsLayer, 0, 0, "}", "weapon", "red", "dark red", true,
		false, false, true, true, false, SlotWeaponPrimary, UseHeal)
	slot2, _ := NewObject(ObjectsLayer, 0, 0, "{", "weapon2", "green", "dark green", true,
		false, false, true, true, false, SlotWeaponSecondary, UseNA)
	slot3, _ := NewObject(ObjectsLayer, 0, 0, "|", "melee", "yellow", "dark yellow", true,
		false, false, true, true, false, SlotWeaponMelee, UseNA)
	item, _ := NewObject(ObjectsLayer, 0, 0, "O", "heal", "blue", "dark blue", true,
		false, false, true, false, true, SlotNA, UseHeal)
	var playerEq = EquipmentComponent{Objects{slot, slot2, slot3}, Objects{item}}
	player, err := NewPlayer(PlayerLayer, 1, 1, "@", "player", "white", "white", true,
		true, false, false, PlayerAI, 999, 5, 2, playerEq)
	if err != nil {
		fmt.Println(err)
	}
	var enemyEq = EquipmentComponent{Objects{nil, nil, nil}, Objects{}}
	enemy, err := NewCreature(CreaturesLayer, 10, 10, "T", "enemy", "green", "green",
		false, true, false, false, MeleePatherAI, 10, 4, 1, enemyEq)
	if err != nil {
		fmt.Println(err)
	}
	slot7, _ := NewObject(ObjectsLayer, 0, 0, "}", "weapon", "red", "dark red", true,
		false, false, true, true, false, SlotWeaponPrimary, UseHeal)
	slot8, _ := NewObject(ObjectsLayer, 0, 0, "{", "weapon2", "green", "dark green", true,
		false, false, true, true, false, SlotWeaponSecondary, UseNA)
	slot9, _ := NewObject(ObjectsLayer, 0, 0, "|", "melee", "yellow", "dark yellow", true,
		false, false, true, true, false, SlotWeaponMelee, UseNA)
	var enemyEq2 = EquipmentComponent{Objects{slot7, slot8, slot9}, Objects{}}
	enemy2, err2 := NewCreature(CreaturesLayer, 11, 11, "T", "enemy", "red", "red",
		false, true, false, false, RangedPatherAI, 10, 4, 1, enemyEq2)
	if err2 != nil {
		fmt.Println(err)
	}
	*c = Creatures{player, enemy, enemy2}
	obj, err := NewObject(ObjectsLayer, 3, 3, "(", "heal2", "blue", "dark blue", true,
		false, false, true, false, false, SlotNA, UseHeal)
	*o = Objects{obj}
	if err != nil {
		fmt.Println(err)
	}
	*b = InitializeEmptyMap()
}

func StartGame(b *Board, c *Creatures, o*Objects) {
	_, errBoard := os.Stat("./map.gob")
	_, errCreatures := os.Stat("./monsters.gob")
	_, errObjects := os.Stat("./objects.gob")
	if errBoard == nil && errCreatures == nil && errObjects == nil {
		LoadGame(b, c, o)
	} else if errBoard != nil && errCreatures != nil && errObjects != nil {
		NewGame(b, c, o)
	} else {
		txt := CorruptedSaveError(errBoard, errCreatures, errObjects)
		fmt.Println("Error: save files are corrupted: " + txt)
		panic(-1)
	}
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
	InitializeFOVTables()
	InitializeBLT()
}
