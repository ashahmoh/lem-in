package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type antFarm struct {
	Rooms []*room
}

type room struct {
	Id     int
	Name   string
	Tunnel []*room //adjacent vertex (neighbours)
}

var (
	from     string
	to       string
	therooms string
	//pathSlice [][]*room
// 	thePaths [][]*room
// 	endRoom  = getEnd()
)

func (f *antFarm) showRooms() {

	for _, val := range f.Rooms {
		fmt.Printf("\nRoom: %v", val.Name)
		for _, val := range val.Tunnel {
			fmt.Printf(" %v", val.Name)
		}

	}
	fmt.Println()
	f.startRoom()
	f.endRoom()
	fmt.Println()
}

func contains(rList []*room, rname string) bool {
	for _, val := range rList {
		if rname == val.Name {
			return true
		}
	}
	return false
}

func (f *antFarm) getRoom(name string) *room {
	for index, val := range f.Rooms {
		if val.Name == name {
			return f.Rooms[index]
		}
	}
	return nil
}

func (f *antFarm) addTunnel(from, to string) {

	//get Room
	fromRoom := f.getRoom(from)
	toRoom := f.getRoom(to)

	//check error
	if fromRoom == nil || toRoom == nil {
		err := fmt.Errorf("ERROR: invalid Tunnel (%v-->%v)", from, to)
		fmt.Println(err.Error())
		os.Exit(0)
	} else if contains(fromRoom.Tunnel, to) {
		err := fmt.Errorf("ERROR: existing Tunnel (%v-->%v)", from, to)
		fmt.Println(err.Error())
		os.Exit(0)
	} else if fromRoom == toRoom {
		err := fmt.Errorf("ERROR: cannot connect room to itself (%v --> %v)", from, to)
		fmt.Println(err.Error())
		os.Exit(0)
	} else if fromRoom.Name == f.endRoom().Name {
		toRoom.Tunnel = append(toRoom.Tunnel, fromRoom)

	} else if toRoom.Name == f.startRoom().Name {
		toRoom.Tunnel = append(toRoom.Tunnel, fromRoom)
	} else {
		fromRoom.Tunnel = append(fromRoom.Tunnel, toRoom)
	}

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

func printAnts(fname []string) int {
	file := readFile(os.Args[1])
	ant := file[0]
	if file[0] <= "0" {
		err := fmt.Errorf("invalid number of ants")
		fmt.Println(err.Error())
	}
	theAnts, _ := strconv.Atoi(string(ant))
	fmt.Println("\nNumber of Ants:", theAnts)
	fmt.Println()

	return theAnts
}

func (f *antFarm) startRoom() *room {
	var start string

	file := readFile(os.Args[1])

	for i := range file[1:] {
		if file[i] == "##start" {
			start = file[i+1]
			if strings.Contains(start, " ") {
				words := strings.Fields(start)
				start = words[0]
				fmt.Printf("\nStart: %v\n", start)
			}

		}

	}
	return f.getRoom(start)

}

func (f *antFarm) endRoom() *room {

	var end string
	file := readFile(os.Args[1])
	for i := range file[1:] {
		if file[i] == "##end" {
			end = file[i+1]
			if strings.Contains(end, " ") {
				words := strings.Fields(end)
				end = words[0]
				fmt.Printf("End: %v\n", end)
			}

		}

	}
	return f.getRoom(end)

}

func getEnd() string {
	var end string
	file := readFile(os.Args[1])

	for i := range file {
		if file[i] == "##end" {
			end = strings.Split(string(file[i+1]), " ")[0]
		}
	}
	return end
}
func (f *antFarm) dfs(current, end *room, visited []bool, path []*room, paths *[][]*room) {
	// Mark the current room as visited
	visited[current.Id] = true

	// Add the current room to the path
	path = append(path, current)

	// If the current room is the end room, add the path to the list of all paths
	if current == end {
		pathCopy := make([]*room, len(path))
		copy(pathCopy, path)
		*paths = append(*paths, pathCopy)
	} else {
		// Recursively search the neighbours of the current room
		for _, neighbour := range current.Tunnel {
			if !visited[neighbour.Id] {
				f.dfs(neighbour, end, visited, path, paths)
			}
		}
	}

	// Backtrack by removing the current room from the path
	path = path[:len(path)-1]

	// Mark the current room as not visited
	visited[current.Id] = false
}

func main() {

	ourFarm := &antFarm{}

	file := readFile(os.Args[1])
	printAnts(file)

	for _, v := range file[1:] {
		if !strings.Contains(v, "-") && strings.Contains(string(v), " ") {
			//step 2:
			words := strings.Fields(string(v))
			therooms = words[0]
			ourFarm.Rooms = append(ourFarm.Rooms, &room{Name: therooms})

		}

	}

	var remdash string

	for _, v := range file[1:] {
		if strings.Contains(v, "-") {
			remdash = strings.Replace(v, "-", " ", 1)
			words := regexp.MustCompile(" ").Split(remdash, -1)
			from = words[0]
			to = words[1]
			fmt.Println("from:", from, "	to:", to)

		}

	}
	fmt.Println("tunnel:", ourFarm.startRoom().Tunnel)
	ourFarm.addTunnel(from, to)

	ourFarm.showRooms()

}
