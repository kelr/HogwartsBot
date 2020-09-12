package main

import (
	"encoding/gob"
	"fmt"
	"os"
)

// Write houseCount info to the file
func writeFile() {
	f, err := os.OpenFile("houses.gob", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	e := gob.NewEncoder(f)

	err = e.Encode(&houseCount)
	if err != nil {
		panic(err)
	}
}

// Load the houseCount info from the file
func loadFile() {
	f, err := os.OpenFile("houses.gob", os.O_RDWR, 0755)
	if err != nil {
		// Handle PathError specifically as it indicates the file does not exist
		if _, ok := err.(*os.PathError); ok {
			fmt.Println("File did not exist, creating now")
			writeFile()
			return
		} else {
			fmt.Println(err)
		}
	}
	defer f.Close()
	d := gob.NewDecoder(f)

	err = d.Decode(&houseCount)
	if err != nil {
		panic(err)
	}
}
