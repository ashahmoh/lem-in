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

// TO DO: put the add room bits from func main into here.
func (f *antFarm) addRoom(rname string) {
	//first place statement to check whether a room in the antFarm has is a certain room with a certain name(
	///check a certain room is in the antFarm under a certian name:
	//if statement is true, means antFarm already has that certain room
	// file := readFile(os.Args[1]

	if contains(f.Rooms, rname) {
		// 	err := fmt.Errorf("room %v not added as room already exists", rname)
		// 	fmt.Println(err.Error())

		// } else {
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
	file := readFile(os.Args[1])

	ourFarm := &antFarm{}
	var therooms string
	//step 1: read from file,
	//step 2: remove coordinates
	//step three save each name at index as room name.
	//step 4: add room to antfarm!
	for _, v := range file[1:] {
		if !strings.Contains(v, "-") && strings.Contains(string(v), " ") {
			//step 2:
			words := strings.Fields(string(v))
			therooms = words[0]
			ourFarm.Rooms = append(ourFarm.Rooms, &room{Name: therooms})

		}

	}

	ourFarm.showRooms()
}

func (f *antFarm) showRooms() {
	for _, val := range f.Rooms {
		//print out the name for each room
		fmt.Printf("\nRoom: %v", val.Name)
		//another loop that can print out ALL the rooms inside the adjacency lsit
		for _, val := range val.Tunnel {
			fmt.Printf(" %v", val.Name)
		}
	}
	fmt.Println()
}

func allAnts() {
	file := readFile(os.Args[1])
	ant := file[0]
	theAnts, _ := strconv.Atoi(string(ant))
	fmt.Println(theAnts)

}
