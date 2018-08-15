package main

type Creature struct {
	/*Creatures are living objects that
	  moves, attacks, dies, etc.*/
	Block Basic
}

/*Monsters holds every monster on map.*/
type Monsters []Creature
