/*
Copyright (c) 2018, Tomasz "VedVid" Nowakowski
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this
   list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

package main

import (
	"encoding/gob"
	"fmt"
	"os"
)

const (
	// Constant values for save files manipulation.
	MapNameGob       = "map.gob"
	MapPathGob       = "./" + MapNameGob
	CreaturesNameGob = "monsters.gob"
	CreaturesPathGob = "./" + CreaturesNameGob
	ObjectsNameGob   = "objects.gob"
	ObjectsPathGob   = "./" + ObjectsNameGob
)

const (
	// Unique name that serves as identifier to values
	// that should be converted from nil to object or from object to nil.
	ObjectNilPlaceholder = "ObjectNilPlaceholder"
)

func NilToObject() *Object {
	/* Function nilToObject returns *Object with >>placeholder<< identifier.
	   It serves to find data that is nil in game - but format gob does not
	   work well with nil values (and interfaces).
	   It is ugly hack, but works. */
	placeholder := &Object{BasicProperties{0, 0, "o", ObjectNilPlaceholder,
		"black", "black"},
		VisibilityProperties{0, false},
		CollisionProperties{false, false},
		ObjectProperties{false, false, false, 0, 0}}
	return placeholder
}

func writeGob(path string, thing interface{}) error {
	/* Function writeGob takes path-to-file, and any object (as interface{})
	   as arguments, then encodes it to gob file. Returns error - unfortunately,
	   errors that are built in gob package are not very helpful, and whole process
	   is hard to debug. */
	f, err := os.Create(path)
	defer f.Close()
	if err == nil {
		encoder := gob.NewEncoder(f)
		encoder.Encode(thing)
	}
	return err
}

func readGob(path string, thing interface{}) error {
	/* Function readGob takes path-to-file, and any object (as interface{})
	   as arguments, then decodes file to interface. Returns error - Decoding has
	   as unhelpful errors as Encoding. */
	f, err := os.Open(path)
	defer f.Close()
	if err == nil {
		decoder := gob.NewDecoder(f)
		err = decoder.Decode(thing)
	}
	return err
}

func saveBoard(b Board) error {
	/* Function saveBoard is helper function that takes game map
	   as argument and encodes it to save file. */
	err := writeGob(MapPathGob, b)
	return err
}

func loadBoard(b *Board) error {
	/* Function loadBoard is helper function that decodes saved data
	   to game map. */
	err := readGob(MapPathGob, b)
	return err
}

func saveCreatures(c Creatures) error {
	/* Function saveBoard is helper function that takes monsters
	   as argument and encodes it to save file.
	   Unfortunately, gob format/package does not work well with
	   nil values. To encode it properly, there is placeholder
	   object created for every nil object; these false objects
	   should be decode to nil by loadCreatures. */
	for i := 0; i < len(c); i++ {
		for j := 0; j < len(c[i].Equipment); j++ {
			if c[i].Equipment[j] == nil {
				c[i].Equipment[j] = NilToObject()
			}
		}
		for k := 0; k < len(c[i].Inventory); k++ {
			if c[i].Inventory[k] == nil {
				c[i].Inventory[k] = NilToObject()
			}
		}
	}
	err := writeGob(CreaturesPathGob, c)
	return err
}

func loadCreatures(c *Creatures) error {
	/* Function loadCreatures is helper function that decodes saved data
	   to slice of creatures. Gob package has troubles with handling nil
	   values, so every nil is represented as placeholder object.
	   During decoding, every placeholder becomes nil again. */
	err := readGob(CreaturesPathGob, c)
	for i := 0; i < len(*c); i++ {
		objs := (*c)[i].Equipment
		for j := 0; j < len(objs); j++ {
			if objs[j].Name == ObjectNilPlaceholder {
				objs[j] = nil
			}
		}
		inv := (*c)[i].Inventory
		for k := 0; k < len(inv); k++ {
			if inv[k].Name == ObjectNilPlaceholder {
				inv[k] = nil
			}
		}
	}
	return err
}

func saveObjects(o Objects) error {
	/* Function saveObjects is helper function that takes game objects
	   as argument and encodes it to save file. */
	err := writeGob(ObjectsPathGob, o)
	return err
}

func loadObjects(o *Objects) error {
	/* Function loadObjects is helper function that decodes saved data
	   to slice of objects. */
	err := readGob(ObjectsPathGob, o)
	return err
}

func SaveGame(b Board, c Creatures, o Objects) error {
	/* Function SaveGame encodes game map, monsters, objects into
	   save files, using Go's gob format. This function may need better
	   error handling - it relies on gob's built-in errors that are
	   not very helpful. */
	var err error
	err = saveBoard(b)
	if err != nil {
		fmt.Println(err)
	}
	err = saveCreatures(c)
	if err != nil {
		fmt.Println(err)
	}
	err = saveObjects(o)
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func LoadGame(b *Board, c *Creatures, o *Objects) error {
	/* Function LoadGame decoded save files (their names and paths are
	   specified as constants on the top of this file) into
	   game map, monsters and objects. As SaveGame, it may need
	   better error handling due to unhelpful gob's error messages. */
	var err error
	err = loadBoard(b)
	if err != nil {
		fmt.Println(err)
	}
	err = loadCreatures(c)
	if err != nil {
		fmt.Println(err)
	}
	err = loadObjects(o)
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func DeleteSaves() {
	/* Function DeleteSaves sereves, well, deleting saves (mostly upon death).
	   It checks if certain save file exists. If so, removes it.
	   Returns the first encountered error. */
	var err error
	_, err = os.Stat(MapPathGob)
	if err == nil {
		os.Remove(MapPathGob)
	}
	_, err = os.Stat(CreaturesPathGob)
	if err == nil {
		os.Remove(CreaturesPathGob)
	}
	_, err = os.Stat(ObjectsPathGob)
	if err == nil {
		os.Remove(ObjectsPathGob)
	}
}
