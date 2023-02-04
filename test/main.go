package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// graph struct represents an adjacency list graph
type antFarm struct {
	Rooms []*room
}

// room struct represnts the room in the graph
type room struct {
	Name   string  //our key are the room names so it is a string in lem-in, not an int
	Tunnel []*room //adjacent: as this an adjacency list, each room as a list of neighbouring rooms connected by tunnels (list of lists)
}

func readFile(fname string) []string {
	file, _ := os.Open(os.Args[1])
	scanned := bufio.NewScanner(file)
	scanned.Split(bufio.ScanLines)

	var lines []string

	for scanned.Scan() {
		lines = append(lines, scanned.Text())
	}

	return lines
}

func (f *antFarm) addRoom(rname string) {
	//first place statement to check whether a room in the antFarm has is a certain room with a certain name(
	///check a certain room is in the antFarm under a certian name:
	//if statement is true, means antFarm already has that certain room
	// file := readFile(os.Args[1])

	if contains(f.Rooms, rname) {
		err := fmt.Errorf("room %v not added as room already exists", rname)
		fmt.Println(err.Error())

	} else {
		//create a room that has k as the name --> "&room{name: k}"
		//append k to the rooms list in the antFarm (Rooms field) --> f.Rooms = append(f.Rooms, &room{name: k})
		f.Rooms = append(f.Rooms, &room{Name: rname})
	}

}

// take in a room list and this bool will return true if certian room is in list
func contains(rList []*room, rname string) bool {
	//loop through room list and return true if match
	for _, val := range rList {
		if rname == val.Name {
			return true
		}
	}
	return false
}

func main() {
	// file := readFile(os.Args[1])
	allAnts()
	theRooms()

	// t := room{Name: theWords}

}

// read through file and return room names
func theRooms() {
	file := readFile(os.Args[1])
	var theWords string
	for _, v := range file[1:] {
		if !strings.Contains(v, "-") && strings.Contains(string(v), " ") {
			words := strings.Fields(string(v))
			theWords = words[0]

			fmt.Println(string(theWords))

		}

	}

}

func allAnts() {
	file := readFile(os.Args[1])
	ant := file[0]
	theAnts, _ := strconv.Atoi(string(ant))
	fmt.Println(theAnts)

}
