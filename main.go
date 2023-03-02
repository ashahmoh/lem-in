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
	Id      int
	Name    string
	Tunnel  []*room //adjacent vertex (neighbours)
	visited bool
	Path    []*room
	Ants    int
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
	//file := readFile(os.Args[1])
	f.startRoom()
	f.endRoom()
	fmt.Println()

	for _, v := range f.Rooms {
		fmt.Printf("\nRoom %v: ", v.Name)
		for _, v := range v.Tunnel {
			fmt.Printf("%v ", v.Name)
		}
	}

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
				//fmt.Printf("\nStart: %v\n", start)
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
				//fmt.Printf("End: %v\n", end)
			}

		}

	}
	return f.getRoom(end)

}

func (r *room) dfs(end *room, path []*room, paths map[int][]*room, visited map[*room]bool) {
	// v == current
	visited[r] = true      // marks current vertex as visited
	path = append(path, r) // append current to path

	if r == end {
		// Found a path from start to end
		length := len(paths)
		paths[length] = path
	} else {
		for _, tunnel := range r.Tunnel {
			if !visited[tunnel] {
				tunnel.dfs(end, path, paths, visited)
			}
		}
	}
	// Remove v from the current path and visited set to backtrack
	delete(visited, r)
}

func main() {
	file := readFile(os.Args[1])
	printAnts(file)
	ourFarm := &antFarm{}

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
			ourFarm.addTunnel(from, to)

		}

	}

	// start := ourFarm.startRoom()
	// end := ourFarm.endRoom()
	// start.dfs(end)

	ourFarm.showRooms()
	fmt.Println()

}
