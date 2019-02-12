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
	   Loop breaks with Escape, Space or Enter input. */
	startX, startY := c.X, c.Y
	targetX, targetY := startX, startY
	msg := ""
	i := false
	for {
		vec, err := NewVector(startX, startY, targetX, targetY)
		if err != nil {
			fmt.Println(err)
		}
		_ = ComputeVector(vec)
		_, _, _, _ = ValidateVector(vec, b, cs, o)
		PrintVector(vec, VectorColorNeutral, VectorColorNeutral, b, o, cs)
		if b[targetX][targetY].Explored == true {
			if IsInFOV(b, c.X, c.Y, targetX, targetY) == true {
				s := GetAllStringsFromTile(targetX, targetY, b, cs, o)
				msg = FormatLookingMessage(s, true)
			} else {
				// Skip monsters if tile is out of c's field of view.
				s := GetAllStringsFromTile(targetX, targetY, b, nil, o)
				msg = FormatLookingMessage(s, false)
			}
		} else {
			msg = "You don't know what is here."
		}
		PrintLookingMessage(msg, i)
		key := blt.Read()
		if key == blt.TK_ESCAPE || key == blt.TK_ENTER || key == blt.TK_SPACE {
			break
		}
		CursorMovement(&targetX, &targetY, key)
		i = true
	}
}

func PrintLookingMessage(s string, b bool) {
	/* Function PrintLookingMessage takes string (message) and bool ("is it
	   a first iteration?") as arguments.
	   It is used to provide dynamic printing looking message:
	   player do not need to confirm target to see what is it, but messages
	   will not flood message log. */
	l := len(MsgBuf)
	if s != "" {
		switch {
		case l == 0:
			AddMessage(s)
		case l >= MaxMessageBuffer:
			RemoveLastMessage()
			AddMessage(s)
		case l > 0 && l < MaxMessageBuffer:
			if b == true {
				RemoveLastMessage()
			}
			AddMessage(s)
		}
	}
}

func FormatLookingMessage(s []string, fov bool) string {
	/* FormatLookingMessage is function that takes slice of strings as argument
	   and returns string.
	   Player "see" things in his fov, and "recalls" out of his fov.
	   It is used to format Look() messages properly.
	   If slice is empty, it return empty tile message.
	   If slice contains only one item, it creates simplest message.
	   If slice is longer, it starts to format message - but it is
	   explicitly visible in function body. */
	const inFov = "see"
	const outFov = "recall"
	txt := ""
	if fov == true {
		txt = inFov
	} else {
		txt = outFov
	}
	if len(s) == 0 {
		return "There is nothing here."
	}
	if len(s) == 1 {
		return "You " + txt + " " + s[0] + " here."
	}
	msg := "You " + txt + " "
	for i, v := range s {
		if i < len(s)-2 { // Regular items.
			msg = msg + v + ", "
		} else if i == len(s)-1-1 { // One-before-last item.
			msg = msg + v + " and "
		} else { // Last item.
			msg = msg + v + " here."
		}
	}
	return msg
}

func (c *Creature) Target(b Board, o Objects, cs Creatures) bool {
	/* Target is method of Creature, that takes game map, objects, and
	   creatures as arguments. Returns bool that serves as indicator if
	   action took some time or not.
	   This method is "the big one", general, for handling targeting.
	   In short, player starts targetting, line is drawn from player
	   to monster, then function waits for input (confirmation - "fire",
	   breaking the loop, or continuing).
	   Explicitly:
	   - creates list of all potential targets in fov
	    * tries to automatically last target, but
	    * if fails, it targets the nearest enemy
	   - draws line between source (receiver) and target (coords)
	    * creates new vector
	    * checks if it is valid - monsterHit should not be nil
	    * prints brensenham's line (ie so-called "vector")
	   - waits for player input
	    * if player cancels, function ends
	    * if player confirms, valley is shoot (in target, or empty space)
	    * if valley is shot in empty space, vector is extrapolated to check
	      if it will hit any target
	    * player can switch between targets as well; it targets
	      next target automatically; at first, only monsters that are
	      valid target (ie clean shot is possible), then monsters that
	      are in range and fov, but line of shot is not clear
	    * in other cases, game will try to move cursor; invalid input
	      is ignored */
	turnSpent := false
	var target *Creature
	targets := c.FindTargets(FOVLength, b, cs, o)
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
	i := false
	for {
		vec, err := NewVector(c.X, c.Y, targetX, targetY)
		if err != nil {
			fmt.Println(err)
		}
		_ = ComputeVector(vec)
		_, _, monsterHit, _ := ValidateVector(vec, b, targets, o)
		PrintVector(vec, VectorColorGood, VectorColorBad, b, o, cs)
		if monsterHit != nil {
			msg := "There is " + monsterHit.Name + " here."
			PrintLookingMessage(msg, i)
		}
		key := blt.Read()
		if key == blt.TK_ESCAPE {
			break
		}
		if key == blt.TK_F {
			monsterAimed := FindMonsterByXY(targetX, targetY, cs)
			if monsterAimed != nil && monsterAimed != c && monsterAimed.HPCurrent > 0 {
				LastTarget = monsterAimed
				c.AttackTarget(monsterAimed)
			} else {
				if monsterAimed == c {
					break // Do not hurt yourself.
				}
				if monsterHit != nil {
					if monsterHit.HPCurrent > 0 {
						LastTarget = monsterHit
						c.AttackTarget(monsterHit)
					}
				} else {
					vx, vy := FindVectorDirection(vec)
					v := ExtrapolateVector(vec, vx, vy)
					_, _, monsterHitIndirectly, _ := ValidateVector(v, b, targets, o)
					if monsterHitIndirectly != nil {
						c.AttackTarget(monsterHitIndirectly)
					}
				}
			}
			turnSpent = true
			break
		} else if key == blt.TK_TAB {
			i = true
			monster := FindMonsterByXY(targetX, targetY, cs)
			if monster != nil {
				target = NextTarget(monster, targets)
			} else {
				target = NextTarget(target, targets)
			}
			targetX, targetY = target.X, target.Y
			continue // Switch target
		}
		CursorMovement(&targetX, &targetY, key)
		i = true
	}
	return turnSpent
}

func CursorMovement(x, y *int, key int) {
	/* CursorMovement is function that takes pointers to coords, and
	   int-based user input. It uses MoveCursor function to
	   modify original values. */
	switch key {
	case blt.TK_UP:
		MoveCursor(x, y, 0, -1)
	case blt.TK_RIGHT:
		MoveCursor(x, y, 1, 0)
	case blt.TK_DOWN:
		MoveCursor(x, y, 0, 1)
	case blt.TK_LEFT:
		MoveCursor(x, y, -1, 0)
	}
}

func MoveCursor(x, y *int, dx, dy int) {
	/* Function MoveCursor takes pointers to coords, and
	   two other ints as direction indicators.
	   It adds direction to coordinate, checks if it is in
	   map bounds, and modifies original values accordingly.
	   This function is called by CursorMovement. */
	newX, newY := *x+dx, *y+dy
	if newX < 0 || newX >= MapSizeX {
		newX = *x
	}
	if newY < 0 || newY >= MapSizeY {
		newY = *y
	}
	*x, *y = newX, newY
}

func (c *Creature) FindTargets(length int, b Board, cs Creatures, o Objects) Creatures {
	/* FindTargets is method of Creature that takes several arguments:
	   length (that is supposed to be max range of attack), and: map, creatures,
	   objects. Returns list of creatures.
	   At first, method creates list of all monsters im c's field of view.
	   Then, this list is divided to two: first, with all "valid" targets
	   (clean (without obstacles) line between c and target) and second,
	   with all other monsters that remains in fov.
	   Both slices are sorted by distance from receiver, then merged.
	   It is necessary for autotarget feature - switching between targets
	   player will start from the nearest valid target, to the farthest valid target;
	   THEN, it will start to target "invalid" targets - again,
	   from nearest to farthest one. */
	targets := c.MonstersInFov(b, cs)
	targetable, unreachable := c.MonstersInRange(b, targets, o, length)
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
	/* FindTarget is method of Creature that takes Creatures as arguments.
	   It returns specific Creature and error.
	   "targets" is supposed to be slice of Creature in player's fov,
	   sorted as explained in FindTargets docstring.
	   If this slice is empty, the target is set to receiver. If not,
	   it tries to target lastly targeted Creature. If it is not possible,
	   it targets first element of slice, and marks it as LastTarget.
	   This method throws an error if it can not find any target,
	   even including receiver. */
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
	/* Function NextTarget takes specific creature (target) and slice of creatures
	   (targets) as arguments. It tries to find the *next* target (used
	   with switching between targets, for example using Tab key).
	   At the end, it returns the next creature. */
	i, _ := FindCreatureIndex(target, targets)
	var t *Creature
	length := len(targets)
	if length > i+1 {
		t = targets[i+1]
	} else if length == 0 {
		t = target
	} else {
		t = targets[0]
	}
	return t
}

func (c *Creature) MonstersInRange(b Board, cs Creatures, o Objects,
	length int) (Creatures, Creatures) {
	/* MonstersInRange is method of Creature. It takes global map, Creatures
	   and Objects, and length (range indicator) as its arguments. It returns
	   two slices - one with monsters that are in range, and one with
	   monsters out of range.
	   At first, two empty slices are created, then function starts iterating
	   through Creatures from argument. It creates new vector from source (c)
	   to target, adds monster to proper slice. It also validates vector
	   (ie, won't add monster hidden behind wall) and skips all dead monsters. */
	var inRange = Creatures{}
	var outOfRange = Creatures{}
	for i, v := range cs {
		vec, err := NewVector(c.X, c.Y, v.X, v.Y)
		if err != nil {
			fmt.Println(err)
		}
		if ComputeVector(vec) <= length+1 { // "+1" is necessary due Vector values.
			valid, _, _, _ := ValidateVector(vec, b, cs, o)
			if cs[i].HPCurrent <= 0 {
				continue
			}
			if valid == true {
				inRange = append(inRange, cs[i])
			} else {
				outOfRange = append(outOfRange, cs[i])
			}
		}
	}
	return inRange, outOfRange
}

func ZeroLastTarget(c *Creature) {
	/* LastTarget is global variable (will be incorporated into
	   player struct in future). Function ZeroLastTarget changes
	   last target to nil, is last target matches creature
	   passed as argument. */
	if LastTarget == c {
		LastTarget = nil
	}
}
