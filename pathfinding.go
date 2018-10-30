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

import "math"
import "strconv"
import blt "bearlibterminal"
import "fmt"

const nodeBaseWeight = -1
const nodeGoalWeight = 0

type Node struct {
	X, Y   int
	Weight int
}

func TilesToNodes(b Board) [][]*Node {
	_ = b //TODO: use Board to check for blocked tiles etc.
	nodes := make([][]*Node, WindowSizeX)
	for i := range nodes {
		nodes[i] = make([]*Node, WindowSizeY)
	}
	for x := 0; x < WindowSizeX; x++ {
		for y := 0; y < WindowSizeY; y++ {
			nodes[x][y] = &Node{x, y, nodeBaseWeight}
		}
	}
	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			fmt.Println(nodes[x][y].Weight)
		}
	}
	return nodes
}

func (c *Creature) MoveTowardsPath(b Board, tx, ty int) {
	/*nodes := TilesToNodes(b)
	goal := nodes[tx][ty]
	goal.Weight = nodeGoalWeight
	//nodes that were weighted in previous iteration
	var traversed = []*Node{goal}
	w := nodeGoalWeight //weight
	for {
		fmt.Println()
		var adjacent = []*Node{}
		//at the end of iteration, traversed = list-of-adjacent-but-unweighted-tiles
		//if empty, break the loop - whole map is traversed already
		w++ //increase weight
		if len(traversed) == 0 {
			break
		} else { //not necessary, but it'll make code more readable
			//create list of tiles that are adjacent to lastly traversed ones
			for i := 0; i < len(traversed); i++ {
				t := traversed[i]
				for x := (t.X - 1); x <= (t.X + 1); x++ {
					for y := (t.Y - 1); y <= (t.Y + 1); y++ {
						if x == t.X && y == t.Y {
							//skip "t"
							continue
						} else {
							//add neightbours to the adjacent slice
							if nodes[x][y].Weight == nodeBaseWeight {
								// it means that node was not traversed yet
								adjacent = append(adjacent, nodes[x][y])
							}
						}
					}
				}
			}
			for i := 0; i < len(adjacent); i++ {
				adjacent[i].Weight = w
			}
			fmt.Println(len(traversed))
			fmt.Println(len(adjacent))
			//perfecly cloned slice :3
			traversed = nil
			traversed = append(adjacent[:0:0], adjacent...)
			adjacent = nil
			fmt.Println(len(traversed))
			fmt.Println(len(adjacent))
		}
	}*/
}

func RenderPath(nodes [][]*Node) {
	blt.Layer(5)
	for x := 0; x < WindowSizeX; x++ {
		for y := 0; y < WindowSizeY; y++ {
			glyph := strconv.Itoa(nodes[x][y].Weight)
			if nodes[x][y].Weight < 0 {
				glyph = "X"
			}
			blt.Print(x, y, glyph)
			fmt.Println(x, y, glyph)
		}
	}
}

func (c *Creature) MoveTowardsDumb(b Board, tx, ty int) {
	/*MoveTowardsDumb is Creature method;
	  it is main part of creature pathfinding. It is very simple algorithm that
	  is not supposed to replace good, old A-Star.*/
	dx := tx - c.X
	dy := ty - c.Y
	ddx, ddy := 0, 0
	if dx > 0 {
		ddx = 1
	} else if dx < 0 {
		ddx = (-1)
	}
	if dy > 0 {
		ddy = 1
	} else if dy < 0 {
		ddy = (-1)
	}
	if b[c.X+ddx][c.Y+ddy].Blocked == false {
		c.Move(ddx, ddy, b)
	} else {
		if ddx != 0 {
			if b[c.X+ddx][c.Y].Blocked == false {
				c.Move(ddx, 0, b)
			}
		} else if ddy != 0 {
			if b[c.X][c.Y+ddy].Blocked == false {
				c.Move(0, ddy, b)
			}
		}
	}
}

func (c *Creature) DistanceTo(tx, ty int) int {
	/*DistanceTo is Creature method. It takes target x and target y as args;
	  computes then returns distance from receiver to target.*/
	dx := float64(tx - c.X)
	dy := float64(ty - c.Y)
	return RoundFloatToInt(math.Sqrt(math.Pow(dx, 2) + math.Pow(dy, 2)))
}
