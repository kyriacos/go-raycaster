package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type LevelData [][]int

// Level - type
type Level struct {
	Data LevelData `json:"map"`
	ID   string    `json:"id"`
}

func (l *Level) At(i, j int) int {
	return l.Data[i][j]
}

// Load level from file
func LoadLevel(filepath string) *Level {
	var l Level

	file, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("Failed to load level file: %s. Error: %s", filepath, err)
	}
	defer file.Close()

	// err := json.Unmarshal(file, &l)
	d := json.NewDecoder(file)
	d.Decode(&l)
	if err != nil {
		log.Fatalf("Couldn't load level: %s. Error: %s", "level-1", err)
	}
	fmt.Println(l.ID)
	return &l
}
