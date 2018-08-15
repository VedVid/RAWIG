package main

type Object struct {
	/*Objects are every other things on map;
	  statues, tables, chairs; but also weapons,
	  armour parts, etc.*/
	Block Basic
}

/*Objects holds every object on map.*/
type Objects []Object
