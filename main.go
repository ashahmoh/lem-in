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
	Name     string
	path     []*room
	visited  bool
	occupied bool
	Tunnel   []*room
}

var (
	from      string
	to        string
	pathSlice [][]*room
)

func main() {

	ourFarm := &antFarm{}

	file := readFile(os.Args[1])

	var therooms string
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
			//fmt.Println("from:", from, "	to:", to)

		}

	}

	ourFarm.addTunnel(from, to)
	ourFarm.showRooms()
	//DFS(ourFarm.startRoom(), *ourFarm)

}

func (f *antFarm) showRooms() {
	allAnts()

	for _, val := range f.Rooms {

		for _, val := range val.Tunnel {
			//add condition to check it doesn't print same room twice
			fmt.Printf(" %v", val.Name)
		}
		//print out the name for each room
		fmt.Printf("\nRoom: %v", val.Name)
	}
	fmt.Println()
	f.startRoom()
	f.endRoom()
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
		if val.Name == from || val.Name == to {
			return f.Rooms[index]
		}
	}
	return nil
}

func (f *antFarm) addTunnel(from, to string) {

	fromRoom := f.getRoom(from)
	toRoom := f.getRoom(to)

	if fromRoom == nil || toRoom == nil {
		err := fmt.Errorf("invalid tunnel (%v --> %v)", from, to)
		fmt.Println(err.Error())
	} else if contains(fromRoom.Tunnel, to) {
		err := fmt.Errorf("tunnel already exists (%v --> %v)", from, to)
		fmt.Println(err.Error())
	} else if fromRoom == toRoom {
		err := fmt.Errorf("cannot connect room to itself (%v --> %v)", from, to)
		fmt.Println(err.Error())
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

func allAnts() {
	file := readFile(os.Args[1])
	ant := file[0]
	if file[0] <= "0" {
		err := fmt.Errorf("invalid number of ants")
		fmt.Println(err.Error())
	}
	theAnts, _ := strconv.Atoi(string(ant))
	fmt.Println("\nNumber of Ants:", theAnts)
	fmt.Println()

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

func DFS(r *room, f antFarm) {
	visitedList := []string{}
	startRoom := f.startRoom()
	endRoom := f.endRoom()
	if r.Name != endRoom.Name {
		r.visited = true
		visitedList = append(visitedList, r.Name)
		for _, neighbour := range r.Tunnel {
			if !neighbour.visited {
				neighbour.path = append(r.path, neighbour)
				if contains(neighbour.path, endRoom.Name) {

					pathSlice = append(pathSlice, neighbour.path)
				}
				visitedList = append(visitedList, neighbour.Name)
				DFS(neighbour, antFarm{f.Rooms})
			}
		}
	} else {
		if len(startRoom.Tunnel) > 1 && !contains(startRoom.Tunnel, endRoom.Name) {
			visitedList = append(visitedList, r.Name)
			fmt.Printf("1stvisitedList %v", visitedList)
			startRoom.Tunnel = f.startRoom().Tunnel[1:][:]
			DFS(startRoom, antFarm{f.Rooms})
		} else {
			visitedList = append(visitedList, r.Name)
			fmt.Printf("2ndvisitedList %v\n", visitedList)
			fmt.Println("start", startRoom.Name, "end", endRoom.Name)
		}
	}
}
