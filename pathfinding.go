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
	return nodes
}

func (c *Creature) MoveTowardsPath(b Board, tx, ty int) bool {
	finished := false
	nodes := TilesToNodes(b)
	goal := &Node{tx, ty, nodeGoalWeight}
	var adjacent = []*Node{goal} //find neightbours of these nodes
	var adjacent2 = []*Node{}    //neightbours; change values
	w := 0
	for {
		w++
		for i := 0; i < len(adjacent); i++ {
			//find all adjacent nodes
			n := adjacent[i]
			for x := (n.X - 1); x <= (n.X + 1); x++ {
				for y := (n.Y - 1); y <= (n.Y + 1); y++ {
					if x == n.X && y == n.Y {
						//pass if newNode is currentNode
						continue
					} else {
						newNode := nodes[x][y]
						newNode.Weight = w
						adjacent2 = append(adjacent2, newNode)
						if x == c.X && y == c.Y {
							//if current node is start point
							finished = true
							goto IterationEnd
						}
					}
				}
			}
		}
		//copy just traversed nodes to the first slice, so
		//they will be used to find their neightbours
		adjacent = adjacent[:0]
		for j := 0; j < len(adjacent2); j++ {
			adjacent = append(adjacent, adjacent2[j])
		}
		adjacent2 = adjacent2[:0]
		if len(adjacent) == 0 {
			break
		}
	}
IterationEnd:
	return finished
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
