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

import "errors"
import "math"
import "strconv"
import blt "bearlibterminal"

const (
	//values important for creating and backtracking graph
	nodeInitialWeight = -1 //value of not traversed nodes
	nodeGoalWeight    = 0  //value of goal node
)

type Node struct {
	/*Node is struct that mimics some properties
	of Tile struct (implemented in map.go);
	X, Y are coords of Node, and Weight is value important
	for graph creating, and later - finding shortest path
	from source (creature) to goal (coords);
	Weight is supposed to be set to -1 initially - it
	marks Node as not traversed during
	graph creation process.*/
	X, Y   int
	Weight int
}

func TilesToNodes(b Board) [][]*Node {
	/*TilesToNodes is function that takes Board
	(ie map, or fragment, of level) as argument. It converts
	Tiles to Nodes, and returns 2d array of *Node to mimic
	Board behaviour.
	In future, it may be worth to create new type, ie
	type Nodes [][]*Node;
	During initializatio, every newly created Node has
	its Weight set to -1 to mark it to not traversed.*/
	nodes := make([][]*Node, WindowSizeX)
	for i := range nodes {
		nodes[i] = make([]*Node, WindowSizeY)
	}
	for x := 0; x < WindowSizeX; x++ {
		for y := 0; y < WindowSizeY; y++ {
			nodes[x][y] = &Node{x, y, nodeInitialWeight}
		}
	}
	return nodes
}

func FindAdjacent(b Board, nodes [][]*Node, frontiers []*Node, start *Node, w int) ([]*Node, bool) {
	/*Function FindAdjacent takes Board, Board-like [][]*Node array,
	coords of starting point, and current value to attribute Weight field
	of Node; FindAdjacent returns slice of adjacent tiles and startFound
	bool flag;
	at start, empty slice of *Node is created, and boolean flag startFound
	is set to false; this flag will be set to true, if function will find
	node that is source of path, and it'll break the loops.
	primary for loop uses one of frontiers, and x, y nested loops
	checks for its adjacent tiles (more details in in-line comments);
	if tile met conditions, its Weight is set to current w value, then node
	is added to list of adjacent tiles.*/
	var adjacent = []*Node{}
	startFound := false
	for i := 0; i < len(frontiers); i++ {
		for x := (frontiers[i].X - 1); x <= (frontiers[i].X + 1); x++ {
			for y := (frontiers[i].Y - 1); y <= (frontiers[i].Y + 1); y++ {
				if x < 0 || x >= WindowSizeX || y < 0 || y >= WindowSizeY {
					continue //node is out of map bounds
				}
				if nodes[x][y].Weight != (-1) {
					continue //node is marked as traversed already
				}
				if x == frontiers[i].X && y == frontiers[i].Y {
					continue //it's the current frontier node
				}
				if b[x][y].Blocked == true || b[x][y].BlocksSight == true {
					continue //tile is blocked, or it blocks line of sight
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
	/*MoveTowardsPath is one of main pathfinding methods. It takes
	Board and ints tx, ty (ie target coords) as arguments;
	MoveTowardsPath uses weighted graph to find shortest path
	from goal (tx, ty - it's more universal than passing Node or
	Creature) to source (creature, ie receiver);
	at first, it creates simple graph with all nodes' Weight set to
	-1 as not-yet-traversed; later, it starts potentially infinite loop
	that breaks if starting position is found by FindAdjacent function,
	or when FindAdjacent won't find any proper tiles that are
	adjacent to previously found ones (ie frontiers);
	after every iteration, local variable "w" used to attribute
	node Weight increases by one, to mark that it's another step away
	from goal position; it makes backtracking easy - Creature position
	is end of path / graph, so Creature has only find node with
	Weight set to lesser value that node occupied by Creature;
	effect may be a bit strange as it takes first node that met
	conditions, but works rather well with basic MoveTowards method.*/
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
		frontiers, startFound = FindAdjacent(b, nodes, frontiers, start, w)
	}
	//uncomment line below, if you want to see nodes' weights
	//RenderWeights(nodes)
	dx, dy, err := BacktrackPath(nodes, start)
	if err != nil {
		fmt.Println(err)
	}
	c.Move(dx, dy, b)
}

func BacktrackPath(nodes [][]*Node, start *Node) (int, int, error) {
	/*Function BacktrackPath takes 2d array of *Node, and
	starting *Node as arguments; it returns two ints, that serves
	as coords;
	BacktrackPath is used in pathfinding; it takes weighted graph
	that has some sort of path already created (more in comments for
	MoveTowardsPath and FindAdjacent) as argument; instead of creating
	proper path, or using search algorithm, structure of graph
	allows to use just node with smaller Weight than start node.
	It returns error if can't find proper tile.
	Note: returning three values at once is ugly.*/
	direction := *start
	for x := (start.X - 1); x <= (start.X + 1); x++ {
		for y := (start.Y - 1); y <= (start.Y + 1); y++ {
			if x < 0 || x >= WindowSizeX || y < 0 || y >= WindowSizeY {
				continue //node is out of map bounds
			}
			if x == start.X && y == start.Y {
				continue //this node is the current node
			}
			if nodes[x][y].Weight < 0 {
				continue //node is not part of path
			}
			if nodes[x][y].Weight < direction.Weight {
				direction = *nodes[x][y] //node is closer to goal than current node
			}
		}
	}
	var err error
	if direction == *start {
		//this error doesn't need helper function from err_.go
		err = errors.New("Warning: function BacktrackPath could not find direction that met all requirements." +
			"\n    Returned coords are coords of starting position.")
	}
	dx := direction.X - start.X
	dy := direction.Y - start.Y
	return dx, dy, err
}

func RenderWeights(nodes [][]*Node) {
	/*RenderWeights is created for debugging purposes;
	it clears whole map, and prints Weights of all nodes
	of graph, then waits for user input to reset;
	it's supposed to be called near the end of
	MoveTowardsPath method.*/
	blt.Clear()
	for x := 0; x < WindowSizeX; x++ {
		for y := 0; y < WindowSizeY; y++ {
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

func (c *Creature) MoveTowards(b Board, tx, ty int, ai int) {
	/*MoveTowards is *the* main method for pathfinding;
	it has *Creature as receiver, and takes Board (ie map of level),
	ints tx and ty (ie coords of Node - in that case, it's more
	universal than passing whole Node or Creature), and ai - it's
	style of ai; these style markers are enums declared in ai.go;
	standard behaviour is always the same - check next tile on the single
	ray between source and target; if it's available to pass, make a move;
	if not, behavior is different for every style;
	creatures with DumbAI checks for adjacent tiles - if are available,
	takes a step, otherwise stands still;
	creatures with other styles (currently only PatherAI is implemented)
	calls MoveTowardsPath function, that creates weighted graph and finds
	shortest path from source to goal.*/
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
		if ai == DumbAI {
			if ddx != 0 {
				if b[c.X+ddx][c.Y].Blocked == false {
					c.Move(ddx, 0, b)
				}
			} else if ddy != 0 {
				if b[c.X][c.Y+ddy].Blocked == false {
					c.Move(0, ddy, b)
				}
			}
		} else {
			c.MoveTowardsPath(b, tx, ty)
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
