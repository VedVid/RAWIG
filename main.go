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
	"fmt"
	"math/rand"
	"time"

	blt "bearlibterminal"
)

var MsgBuf = []string{}

func main() {
	slot, _ := NewObject(ObjectsLayer, 0, 0, "}", "red", "dark red", true, false, false, true, true)
	item, _ := NewObject(ObjectsLayer, 0, 0, "O", "blue", "dark blue", true, false, false, true, false)
	var playerEq = EquipmentComponent{slot, Objects{item}}
	player, err := NewPlayer(PlayerLayer, 1, 1, "@", "white", "white", true, true, false, PlayerAI, 20, 5, 2, playerEq)
	if err != nil {
		fmt.Println(err)
	}
	var enemyEq = EquipmentComponent{nil, Objects{}}
	enemy, err := NewCreature(CreaturesLayer, 10, 10, "T", "green", "green", false, true, false, PatherAI, 10, 4, 1, enemyEq)
	if err != nil {
		fmt.Println(err)
	}
	var actors = Creatures{player, enemy}
	obj, err := NewObject(ObjectsLayer, 3, 3, "(", "blue", "dark blue", true, false, false, true, true)
	var objs = Objects{obj}
	if err != nil {
		fmt.Println(err)
	}
	cells := InitializeEmptyMap()
	for {
		RenderAll(cells, objs, actors)
		key := blt.Read()
		if key == blt.TK_ESCAPE || actors[0].HPCurrent <= 0 {
			break
		} else {
			turnSpent := Controls(key, player, cells, actors, &objs)
			if turnSpent == true {
				CreaturesTakeTurn(cells, actors)
			}
		}
	}
	blt.Close()
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
	InitializeFOVTables()
	InitializeBLT()
}
