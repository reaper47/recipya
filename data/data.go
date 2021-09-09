package data

import (
	"bufio"
	"log"
	"os"
)

// Data is the global variable for the DataReader interface.
var Data DataReader

// DataReader is an interface for the data reader.
//
// It makes testing easier because we can mock.
type DataReader interface {
	ReadBlacklistIngredients() map[string]int8
}

// DataStruct holds functions related to files under the data directory.
type DataStruct struct{}

// ReadBlacklistIngredients reads the list of blacklisted
// ingredients listed in the data/blacklist_ingredients.txt file.
func (d *DataStruct) ReadBlacklistIngredients() map[string]int8 {
	f, err := os.Open("data/blacklist_ingredients.txt")
	if err != nil {
		log.Println("Read blacklist ingredients err:", err)
		return nil
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	items := make(map[string]int8)
	for scanner.Scan() {
		items[scanner.Text()] = 0
	}
	return items
}

func init() {
	Data = &DataStruct{}
}
