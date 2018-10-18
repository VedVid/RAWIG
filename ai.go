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

const (
	//ai types
	PlayerAI = iota
	DumbAI
)

func MonstersTakeTurn(b Board, m Monsters) {
	/*Function MonstersTakeTurn is supposed to handle all enemy monsters
	actions: movement, attacking, etc.
	It takes Board and Monsters as arguments.
	Iterates through all Monsters slice, and handles monster behaviour:
	if distance between monster and player is bigger than 1, monster
	moves towards player.
	It uses switch for matching AIType and behaviour.
	At first, I wanted to use map[int]METHOD, but it's not easy to implement.*/
	for _, v := range m {
		if v.DistanceTo(m[0].X, m[0].Y) > 1 {
			switch v.AIType {
			case DumbAI:
				v.MoveTowardsDumb(b, m[0].X, m[0].Y)
			}
		}
	}
}
