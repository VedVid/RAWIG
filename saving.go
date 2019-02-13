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
	"errors"
	"encoding/gob"
	"fmt"
	"os"
)

const (
	MapName = "map.gob"
	MapPath = "./" + MapName
	CreaturesName = "monsters.gob"
	CreaturesPath = "./" + CreaturesName
	ObjectsName = "objects.gob"
	ObjectsPath = "./" + ObjectsName
)

const (
	objectNilPlaceholder = "objectNilPlaceholder"
)

func nilToObject() *Object {
	placeholder, _ := NewObject(0, 0, 0, "o", objectNilPlaceholder, "black", "black", false,
		false, false, false, false, false, 0, 0)
	return placeholder
}

func writeGob(path string, thing interface{}) error {
	f, err := os.Create(path)
	defer f.Close()
	if err == nil {
		encoder := gob.NewEncoder(f)
		encoder.Encode(thing)
	}
	return err
}

func readGob(path string, thing interface{}) error {
	f, err := os.Open(path)
	defer f.Close()
	if err == nil {
		decoder := gob.NewDecoder(f)
		err = decoder.Decode(thing)
	}
	return err
}

func saveBoard(b Board) error {
	err := writeGob(MapPath, b)
	return err
}

func loadBoard(b *Board) error {
	err := readGob(MapPath, b)
	return err
}

func saveCreatures(c Creatures) error {
	for i := 0; i < len(c); i++ {
		for j := 0; j < len(c[i].Equipment); j++ {
			if c[i].Equipment[j] == nil {
				c[i].Equipment[j] = nilToObject()
			}
		}
	}
	err := writeGob(CreaturesPath, c)
	return err
}

func loadCreatures(c *Creatures) error {
	err := readGob(CreaturesPath, c)
	for i := 0; i < len(*c); i++ {
		objs := (*c)[i].Equipment
		for j := 0; j < len(objs); j++ {
			if objs[j].Name == objectNilPlaceholder {
				objs[j] = nil
			}
		}
	}
	return err
}

func saveObjects(o Objects) error {
	err := writeGob(ObjectsPath, o)
	return err
}

func loadObjects(o *Objects) error {
	err := readGob(ObjectsPath, o)
	return err
}

func SaveGame(b Board, c Creatures, o Objects) error {
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

func DeleteSaves() error {
	var err error
	_, err = os.Stat(MapPath)
	if err == nil {
		os.Remove(MapPath)
	} else {
		err = errors.New("Error: save file not found: " + MapName + ".")
		return err
	}
	_, err = os.Stat(CreaturesPath)
	if err == nil {
		os.Remove(CreaturesPath)
	} else {
		err = errors.New("Error: save file not found: " + CreaturesName + ".")
		return err
	}
	_, err = os.Stat(ObjectsPath)
	if err == nil {
		os.Remove(ObjectsPath)
	} else {
		err = errors.New("Error: save file not found: " + ObjectsName + ".")
		return err
	}
	return err
}
