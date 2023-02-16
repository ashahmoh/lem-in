package main

import (
	"bufio"
	"fmt"
	"os"
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

//4 functions: 1. addRoom, 2. addTunnel, 3. getRoom, 4. showRooms(contents of graph), 5. contains(if room contains duplicates!)

// this method adds to the graph- creates the graph: has a receiver to the graph(pointer to graph struct) and takes in the name of the room(the key)
func (f *antFarm) addRoom(rname string) {
	//check a certain room is in the antFarm under a certian name:
	//if statement is true, means antFarm already has that certain room
	if contains(f.Rooms, rname) {
		err := fmt.Errorf("room %v not added as room already exists", rname)
		fmt.Println(err.Error())

	} else {
		//create a room that has k as the name --> "&room{name: k}"
		//append k to the rooms list in the antFarm (Rooms field) --> f.Rooms = append(f.Rooms, &room{name: k})
		f.Rooms = append(f.Rooms, &room{Name: rname})
	}

}

// showRooms will print the contents of the graph (adjacency list) for each room of the antFarm
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

//add the tunnels!

func (f *antFarm) addTunnel(from, to string) {
	//this is for DIRECTED GRAPH so it has both FROM and TO
	// So addTunnel(2,1) NOT SAME as addTunnel(1,2)
	//THINK WE NEED TO MAKE THIS UNDIRECTED-- may need refacturing

	//1. get room
	//get room at index??
	fromRoom := f.getRoom(from)
	toRoom := f.getRoom(to)

	// 2. error check to filter (e.g if tunnel already exists or if trying to add to tunnel that doesn't)
	// 2.1 check if room actually exists: so if getRoom method returns nil, we will return error as it means we are trying to add tunnels to rooms that do not exist!:
	if fromRoom == nil || toRoom == nil {
		err := fmt.Errorf("invalid tunnel (%v --> %v)", from, to)
		fmt.Println(err.Error())
	} else if contains(fromRoom.Tunnel, to) { //check duplicatated tunnel
		err := fmt.Errorf("tunnel already exists (%v --> %v)", from, to)
		fmt.Println(err.Error())
	} else {
		// Add tunnel
		fromRoom.Tunnel = append(fromRoom.Tunnel, toRoom)
	}

}

// getRoom return a pointer to the room with a name string
func (f *antFarm) getRoom(name string) *room {
	for index, val := range f.Rooms {
		if val.Name == name {
			return f.Rooms[index]
		}
	}
	return nil
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

//define path and define room

func main() {
	// ourFarm := &antFarm{}

	// readfile from arg
	file := readFile(os.Args[1])

	// 1. print to screen: number of ants (file index[0])
	fmt.Print(string(file[0]), "\n")

}

func remDash(str []string) string {
	var res1 string
	for _, v := range str[1:] {
		if strings.Contains(v, "-") {
			var theRooms []string

			theRooms = append(theRooms, v)
			for _, v := range theRooms {
				res1 = strings.Replace(v, "-", " ", 1)
				fmt.Println(res1)
			}

		}

	}

	return res1
}
