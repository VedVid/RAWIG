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

type Creature struct {
	/*Creatures are living objects that
	  moves, attacks, dies, etc.*/
	Block Basic
}

/*Monsters holds every monster on map.*/
type Monsters []Creature

func NewCreature(layer, x, y int, colour, character string) Creature {
	/*Function NewCreture takes all values necessary by its struct,
	and creates then returns Creature*/
	creatureBlock := Basic{layer, x, y, colour, character}
	creatureNew := Creature{creatureBlock}
	return creatureNew
}
