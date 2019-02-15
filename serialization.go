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
/*
func CreatureToJson(c *Creature) {
	var jsonData []byte
	jsonData, err := json.Marshal(c)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(jsonData))
}
*/

func CreatureToJson2(c *Creature) error {
	err := writeJson("./player.json", c)
	return err
}

