package main

import (
	"encoding/json"
	"os"
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
	   encoded into json file. Due to problems with nil using json, all nil
	   weapons are changed to placeholder "object". */
	for i := 0; i < len(c.Equipment); i++ {
		if c.Equipment[i] == nil {
			c.Equipment[i] = NilToObject()
		}
	}
	err := writeJson(path, c)
	return err
}

func CreatureFromJson(path string, c *Creature) error {
	/* Function CreatureFromJson decodes specific json file into Creature,
	   passed as argument. After unmarshalling, function replaces all
	   placeholder weapons with proper nil values. */
	err := readJson(path, c)
	for i := 0; i < len(c.Equipment); i++ {
			if c.Equipment[i].Name == ObjectNilPlaceholder {
				c.Equipment[i] = nil
		}
	}
	return err
}
