package main

type Tile struct {
	/*Tiles are map cells - floors, walls, doors*/
	Block Basic
}

/*Board is map representation, that uses slice
  to hold data of its every cell*/
type Board []Tile
