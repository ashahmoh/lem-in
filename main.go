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
	Name   string
	Tunnel []*room
}

var (
	from string
	to   string
)

func main() {

	ourFarm := &antFarm{}

	file := readFile(os.Args[1])

	allAnts()

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
			fmt.Println("from:", from, "	to:", to)

		}

	}

	ourFarm.addTunnel(from, to)
	ourFarm.showRooms()

}

func (f *antFarm) showRooms() {
	for _, val := range f.Rooms {
		//print out the name for each room
		fmt.Printf("\nRoom: %v", val.Name)
		for _, val := range val.Tunnel {
			fmt.Printf(" %v", val.Name)
		}
	}
	fmt.Println()
	//allAnts()
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
				fmt.Printf("\nstart Room: %v\n", start)
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
				fmt.Printf("End Room: %v\n", end)
			}

		}

	}
	return f.getRoom(end)

}
