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

import "fmt"
import "math"
import "strconv"
import blt "bearlibterminal"

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

func FindAdjacent(nodes [][]*Node, frontiers []*Node, start *Node, w int) ([]*Node, bool) {
	var adjacent = []*Node{}
	startFound := false
	for i := 0; i < len(frontiers); i++ {
		for x := (frontiers[i].X - 1); x <= (frontiers[i].X + 1); x++ {
			for y := (frontiers[i].Y - 1); y <= (frontiers[i].Y + 1); y++ {
				if x < 0 || x >= WindowSizeX || y < 0 || y >= WindowSizeY {
					continue
				}
				if nodes[x][y].Weight != (-1) {
					continue
				}
				if x == frontiers[i].X && y == frontiers[i].Y {
					continue
				}
				nodes[x][y].Weight = w
				adjacent = append(adjacent, nodes[x][y])
				if x == start.X && y == start.Y {
					startFound = true
					goto End
				}
			}
		}
	}
End:
	return adjacent, startFound
}

func (c *Creature) MoveTowardsPath(b Board, tx, ty int) {
	nodes := TilesToNodes(b) //convert tiles to nodes
	start := nodes[c.X][c.Y]
	startFound := false
	goal := nodes[tx][ty]
	goal.Weight = 0
	var frontiers = []*Node{goal}
	w := 0
	for {
		w++
		if len(frontiers) == 0 || startFound == true {
			break
		}
		frontiers, startFound = FindAdjacent(nodes, frontiers, start, w)
	}
	RenderWeights(nodes)
	dx, dy := BacktrackPath(nodes, start)
	c.Move(dx, dy, b)
}

func BacktrackPath(nodes [][]*Node, start *Node) (int, int) {
	direction := *start
	for x := (start.X - 1); x <= (start.X + 1); x++ {
		for y := (start.Y - 1); y <= (start.Y + 1); y++ {
			if x < 0 || x >= WindowSizeX || y < 0 || y >= WindowSizeY {
				continue
			}
			if x == start.X && y == start.Y {
				continue
			}
			if nodes[x][y].Weight < 0 {
				continue
			}
			if nodes[x][y].Weight < direction.Weight {
				direction = *nodes[x][y]
			}
		}
	}
	//needs error checking if couldn't find tile with smaller weight than start
	dx := direction.X - start.X
	dy := direction.Y - start.Y
	return dx, dy
}

func RenderWeights(nodes [][]*Node) { //for debugging purposes
	for x := 0; x < WindowSizeX; x++ {
		for y := 0; y < WindowSizeY; y++ {
			fmt.Println(nodes[x][y].X, nodes[x][y].Y, nodes[x][y].Weight)
			glyph := strconv.Itoa(nodes[x][y].Weight)
			if nodes[x][y].Weight < 0 {
				glyph = "-"
			} else if nodes[x][y].Weight > 9 {
				glyph = "+"
			}
			blt.Print(x, y, glyph)
		}
	}
	blt.Refresh()
	blt.Read()
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
