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
	/*	fmt.Println("\n\n\nstart of function")
		nodes := TilesToNodes(b)
		for x := 0; x < 10; x++ {
			for y := 0; y < 10; y++ {
				fmt.Println(nodes[x][y].Weight)
			}
		}
		goal := &Node{tx, ty, nodeGoalWeight}
		var adjacent = []*Node{goal} //find neightbours of these nodes
		var adjacent2 = []*Node{}    //neightbours; change values
		w := nodeBaseWeight
		fmt.Println("before for loop")
		for {
			fmt.Println("inside for loop")
			fmt.Println("w++...")
			w++
			fmt.Println("done!")
			fmt.Println("before second for loop")
			for i := 0; i < len(adjacent); i++ {
				//find all adjacent nodes
				fmt.Println("n := adjacent[i]...")
				n := adjacent[i]
				fmt.Println(n.Weight)
				fmt.Println("done!")
				fmt.Println("before x, y for loop")
				for x := (n.X - 1); x <= (n.X + 1); x++ {
					for y := (n.Y - 1); y <= (n.Y + 1); y++ {
						fmt.Println("conditional")
						if x == n.X && y == n.Y {
							fmt.Println("x == n.X && y == n.Y")
							fmt.Println(x, y, " : ", n.X, n.Y)
							//pass if newNode is currentNode
							continue
						} else if n.Weight > nodeBaseWeight {
							fmt.Println("n.Weight > nodeBaseWeight")
							fmt.Println(n.Weight, " : ", nodeBaseWeight)
							//pass if newNode was traversed already
							continue
						} else {
							fmt.Println("else")
							fmt.Println("newNode := nodes[x][y]...")
							newNode := nodes[x][y]
							fmt.Println("done!")
							fmt.Println("newNode.Weight = w...")
							newNode.Weight = w
							fmt.Println("done!")
							fmt.Println("print newNode.Weight and nodes[x][y].Weight...")
							fmt.Println(newNode.Weight, nodes[x][y].Weight)
							fmt.Println("done!")
							fmt.Println("adjacent2 = append(adjacent2, newNode...")
							adjacent2 = append(adjacent2, newNode)
							fmt.Println("done!")
							if x == c.X && y == c.Y {
								fmt.Println("x == c.X && y == c.Y")
								//if current node is start point
								goto IterationEnd
							}
						}
					}
				}
			}
			//copy just traversed nodes to the first slice, so
			//they will be used to find their neightbours
			fmt.Println("end of x, y for loops")
			fmt.Println("adjacent = adjacent[:0]...")
			adjacent = adjacent[:0]
			fmt.Println("done!")
			fmt.Println("j for loop")
			for j := 0; j < len(adjacent2); j++ {
				fmt.Println("adjacent = append(adjacent, adjacent2[j]...")
				adjacent = append(adjacent, adjacent2[j])
				fmt.Println("done!")
			}
			fmt.Println("end of j for loop")
			fmt.Println("adjacent2 = adjacent2[:0]...")
			adjacent2 = adjacent2[:0]
			fmt.Println("done!")
			if len(adjacent) == 0 {
				fmt.Println("len(adjacent)==0")
				break
			}
		}
		fmt.Println("it's goto statement here")
	IterationEnd:
		_ = 5 //RenderPath(nodes)*/
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
