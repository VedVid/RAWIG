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
	blt "bearlibterminal"
	"errors"
	"fmt"
	"sort"
)

func (c *Creature) Look(b Board, o Objects, cs Creatures) {
	/* Look is method of Creature (that is supposed to be player).
	   It has to take Board, "global" Objects and Creatures as arguments,
	   because function PrintVector need to call RenderAll function.
	   At first, Look creates new para-vector, with player coords as
	   starting point, and dynamic end position.
	   Then ComputeVector checks what tiles are present
	   between Start and End, and adds their coords to vector values.
	   Line from Vector is drawn, then game waits for player input,
	   that will change position of "looking" cursors.
	   Loop breaks with Escape key as input. */
	startX, startY := c.X, c.Y
	targetX, targetY := startX, startY
	for {
		vec, err := NewVector(startX, startY, targetX, targetY)
		if err != nil {
			fmt.Println(err)
		}
		ComputeVector(vec)
		_ = ValidateVector(vec, b, cs)
		PrintVector(vec, VectorColorNeutral, VectorColorNeutral, b, o, cs)
		key := blt.Read()
		if key == blt.TK_ESCAPE {
			break
		}
		switch key {
		case blt.TK_UP:
			targetY--
		case blt.TK_RIGHT:
			targetX++
		case blt.TK_DOWN:
			targetY++
		case blt.TK_LEFT:
			targetX--
		}
	}
}

func (c *Creature) Target(b Board, o Objects, cs Creatures) {
	var target *Creature
	targets := c.FindTargets(FOVLength, b, cs)
	if LastTarget != nil && LastTarget != c &&
		IsInFOV(b, c.X, c.Y, LastTarget.X, LastTarget.Y) == true {
		target = LastTarget
	} else {
		var err error
		target, err = c.FindTarget(targets)
		if err != nil {
			fmt.Println(err)
		}
	}
	targetX, targetY := target.X, target.Y
	for {
		vec, err := NewVector(c.X, c.Y, targetX, targetY)
		if err != nil {
			fmt.Println(err)
		}
		ComputeVector(vec)
		_ = ValidateVector(vec, b, targets)
		PrintVector(vec, VectorColorGood, VectorColorBad, b, o, cs)
		key := blt.Read()
		if key == blt.TK_ESCAPE {
			break
		}
		if key == blt.TK_F {
			monster := FindMonsterByXY(targetX, targetY, cs)
			if monster != nil {
				LastTarget = monster
			}
			break //fire!
		} else if key == blt.TK_TAB {
			monster := FindMonsterByXY(targetX, targetY, cs)
			if monster != nil {
				target = NextTarget(monster, targets)
			} else {
				target = NextTarget(target, targets)
			}
			targetX, targetY = target.X, target.Y
			continue //switch target
		}
		switch key {
		case blt.TK_UP:
			targetY--
		case blt.TK_RIGHT:
			targetX++
		case blt.TK_DOWN:
			targetY++
		case blt.TK_LEFT:
			targetX--
		}
	}
}

func (c *Creature) FindTargets(length int, b Board, cs Creatures) Creatures {
	targets := c.MonstersInFov(b, cs)
	targetable, unreachable := c.MonstersInRange(b, targets, length) //use ValidateVector
	sort.Slice(targetable, func(i, j int) bool {
		return targetable[i].DistanceBetweenCreatures(c) <
			targetable[j].DistanceBetweenCreatures(c)
	})
	sort.Slice(unreachable, func(i, j int) bool {
		return unreachable[i].DistanceBetweenCreatures(c) <
			unreachable[j].DistanceBetweenCreatures(c)
	})
	targets = nil
	targets = append(targets, targetable...)
	targets = append(targets, unreachable...)
	return targets
}

func (c *Creature) FindTarget(targets Creatures) (*Creature, error) {
	var target *Creature
	if len(targets) == 0 {
		target = c
	} else {
		if LastTarget != nil && CreatureIsInSlice(LastTarget, targets) {
			target = LastTarget
		} else {
			target = targets[0]
			LastTarget = target
		}
	}
	var err error
	if target == nil {
		txt := TargetNilError(c, targets)
		err = errors.New("Could not find target, even the 'self' one." + txt)
	}
	return target, err
}

func NextTarget(target *Creature, targets Creatures) *Creature {
	i, err := FindCreatureIndex(target, targets)
	if err != nil {
		fmt.Println(err)
	}
	var t *Creature
	if len(targets) > i+1 {
		t = targets[i+1]
	} else {
		t = targets[0] //player?
	}
	return t
}

func (c *Creature) MonstersInRange(b Board, cs Creatures, length int) (Creatures, Creatures) {
	var inRange = Creatures{}
	var outOfRange = Creatures{}
	for i, v := range cs {
		vec, err := NewVector(c.X, c.Y, v.X, c.Y)
		if err != nil {
			fmt.Println(err)
		}
		if DistanceBetween(c.X, c.Y, v.X, v.Y) <= length {
			if ValidateVector(vec, b, cs) == true {
				inRange = append(inRange, cs[i])
			} else {
				outOfRange = append(outOfRange, cs[i])
			}
		}
	}
	return inRange, outOfRange
}

func ZeroLastTarget(c *Creature) {
	if LastTarget == c {
		LastTarget = nil
	}
}
