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
	err := writeGob("./map.gob", b)
	//fmt.Println(len(b))
	return err
}

func loadBoard(b *Board) error {
	err := readGob("./map.gob", b)
	//fmt.Println(len(*b))
	return err
}

func saveCreatures(c Creatures) error {
	err := writeGob("./monsters.gob", c)
	//fmt.Println(len(c))
	return err
}

func loadCreatures(c *Creatures) error {
	err := readGob("./monsters.gob", c)
	//fmt.Println(len(*c))
	return err
}

func saveObjects(o Objects) error {
	err := writeGob("./objects.gob", o)
	//fmt.Println(len(o))
	return err
}

func loadObjects(o *Objects) error {
	err := readGob("./objects.gob", o)
	//fmt.Println(len(*o))
	return err
}

func SaveGame(b Board, c Creatures, o Objects) error {
	var err error
	fmt.Println(1)
	err = saveBoard(b)
	if err != nil {
		fmt.Println(2)
		fmt.Println(err)
	}
	fmt.Println(3)
	err = saveCreatures(c)
	fmt.Println(4)
	if err != nil {
		fmt.Println(5)
		fmt.Println(err)
	}
	fmt.Println(6)
	err = saveObjects(o)
	fmt.Println(7)
	if err != nil {
		fmt.Println(8)
		fmt.Println(err)
	}
	fmt.Println(9)
	return err
}

func LoadGame(b *Board, c *Creatures, o *Objects) error {
	var err error
	fmt.Println(1)
	err = loadBoard(b)
	if err != nil {
		fmt.Println(2)
		fmt.Println(err)
	}
	fmt.Println(3)
	err = loadCreatures(c)
	if err != nil {
		fmt.Println(6)
		fmt.Println(err)
	}
	err = loadObjects(o)
	if err != nil {
		fmt.Println(4)
		fmt.Println(err)
	}
	fmt.Println(5)
	return err
}