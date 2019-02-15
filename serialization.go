package main

import (
	"encoding/json"
	"os"
)

func writeJson(path string, thing interface{}) error {
	/* Function writeGob takes path-to-file, and any object (as interface{})
	   as arguments, then encodes it to gob file. Returns error - unfortunately,
	   errors that are built in gob package are not very helpful, and whole process
	   is hard to debug. */
	f, err := os.Create(path)
	defer f.Close()
	if err == nil {
		encoder := json.NewEncoder(f)
		encoder.Encode(thing)
	}
	return err
}

func readJson(path string, thing interface{}) error {
	f, err := os.Open(path)
	defer f.Close()
	if err == nil {
		decoder := json.NewDecoder(f)
		err = decoder.Decode(thing)
	}
	return err
}

func CreatureToJson(c *Creature) error {
	err := writeJson("./player.json", c)
	return err
}

func CreatureFromJson(c *Creature) error {
	err := readJson("./player.json", c)
	return err
}