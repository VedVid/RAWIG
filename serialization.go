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
	"encoding/json"
	"os"
)

const (
	// Constant values for data files manipulation.
	CreaturesPathJson = "./data/monsters/"
	ObjectsPathJson   = "./data/objects/"
	MapsPathJson      = "./data/maps/"
)

func writeJson(path string, thing interface{}) error {
	/* Function writeJson takes path-to-file, and any object (as interface{})
	   as arguments, then encodes it to json file. Returns error - built-in json package. */
	f, err := os.Create(path)
	defer f.Close()
	if err == nil {
		encoder := json.NewEncoder(f)
		encoder.Encode(thing)
	}
	return err
}

func readJson(path string, thing interface{}) error {
	/* Function readJson takes path-to-file, and any object (as interface{})
	   as arguments, then decodes file to interface. Returns error - built-in json package. */
	f, err := os.Open(path)
	defer f.Close()
	if err == nil {
		decoder := json.NewDecoder(f)
		err = decoder.Decode(thing)
	}
	return err
}

func CreatureToJson(path string, c *Creature) error {
	/* Function CreatureToJson takes Creature as argument that will be
	   encoded into json file. */
	err := writeJson(path, c)
	return err
}

func CreatureFromJson(path string, c *Creature) error {
	/* Function CreatureFromJson decodes specific json file into Creature,
	   passed as argument. */
	err := readJson(path, c)
	return err
}

func ObjectToJson(path string, o *Object) error {
	/* Function ObjectToJson takes Object as argument that will be
	   encoded into json file. */
	err := writeJson(path, o)
	return err
}

func ObjectFromJson(path string, o *Object) error {
	/* Function ObjectFromJson decodes specific json file into Object,
	   passed as argument. */
	err := readJson(path, o)
	return err
}

func MapFromJson(path string, m *MapJson) error {
	err := readJson(path, m)
	return err
}
