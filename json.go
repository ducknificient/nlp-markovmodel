package main

import (
	"fmt"
	"os"
)

func saveJsonToFile(filename string, byteData []byte) {

	path := `json`

	// create directory
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Printf("Creating directory: %v\n", path)
		os.MkdirAll(path, os.ModePerm)
	}

	filename = filename + `.json`
	filepath := path + `/` + filename
	// Create or open a JSON file for writing

	file, err := os.Create(filepath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Write the JSON data to the file
	_, err = file.Write(byteData)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

}
