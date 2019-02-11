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

type BasicProperties struct {
	/* BasicProperties is struct that aggregates
	   all widely used data, necessary for every
	   map tile and object representation. */
	Layer     int
	X, Y      int
	Char      string
	Name      string
	Color     string
	ColorDark string
}

type VisibilityProperties struct {
	/* VisibilityProperties is simple struct
	   for checking if object is always visible,
	   regardless of player's fov. */
	AlwaysVisible bool
}

type CollisionProperties struct {
	/* CollisionProperties is struct filled with
	   boolean values, for checking several
	   collision conditions: if cell is blocked,
	   if it blocks creature sight, etc. */
	Blocked     bool
	BlocksSight bool
}

type FighterProperties struct {
	/* FighterProperties stores information about
	   things that can live and fight (ie fighters);
	   it may be used for destructible environment
	   elements as well.
	   AI types are iota (integers) defined
	   in creatures.go. */
	AIType      int
	AITriggered bool
	HPMax       int
	HPCurrent   int
	Attack      int
	Defense     int
}

type ObjectProperties struct {
	/* Not every Object can be picked up - like tables;
	   also, not every Object can be equipped - like cheese.
	   It's place for other properties - like slot it will
	   occupy, use cases, etc.
	   Note that currently Equippable can not be Consumable,
	   due to removing from Inventory / Equipment problems. */
	Pickable   bool
	Equippable bool
	Consumable bool
	Slot       int
	Use        int
}

type EquipmentComponent struct {
	/* EquipmentComponent helps with inventory management.
	   It's part of Creature. Slot is generic place for
	   equipped items (it could be "head" for helmets,
	   "feet" for boots, etc).
	   Inventory is list of items in backpack. */
	Equipment Objects
	Inventory Objects
}
