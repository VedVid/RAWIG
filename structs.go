package main

type Basic struct {
	/*Basic is struct that aggregates all
	  widely used data, necessary for every
	  map tile and object representation*/
	Layer int
	X, Y  int
	Char  string
	Color string
}
