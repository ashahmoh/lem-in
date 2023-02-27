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
	path     []*room //the stack
	visited  bool
	occupied bool
	Tunnel   []*room //adjacent vertex
}

type theAnts struct { //ants
	allAnts []*ant //antz
}

type ant struct {
	Name        string  //key
	AntPath     []*room // valid path //path
	CurrentRoom room
	RoomsPassed int
}

var (
	from string
	to   string
	//pathSlice [][]*room
	thePaths [][]*room
	endRoom  = getEnd()
)

func (f *antFarm) showRooms() {
	file := readFile(os.Args[1])
	printAnts(file)

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

// func DFS(r *room, f antFarm) {
// 	visitedList := []string{}
// 	startRoom := f.startRoom()
// 	if r.Name != endRoom {
// 		fmt.Println("pass")
// 		r.visited = true
// 		visitedList = append(visitedList, r.Name)
// 		fmt.Println(r.Name, visitedList)
// 		for _, adjRoom := range r.Tunnel {
// 			fmt.Println("pass2")
// 			if !adjRoom.visited {
// 				adjRoom.path = append(r.path, adjRoom)
// 				if contains(adjRoom.path, endRoom) {

// 					pathSlice = append(pathSlice, adjRoom.path)
// 				}
// 				visitedList = append(visitedList, adjRoom.Name)
// 				fmt.Println(r.visited)
// 				DFS(adjRoom, antFarm{f.Rooms})
// 			}
// 		}
// 	} else {
// 		if len(startRoom.Tunnel) > 1 && !contains(startRoom.Tunnel, endRoom) {
// 			visitedList = append(visitedList, r.Name)
// 			fmt.Printf("Visited List %v", visitedList)
// 			startRoom.Tunnel = startRoom.Tunnel[1:][:]
// 			DFS(startRoom, antFarm{f.Rooms})
// 		} else {
// 			visitedList = append(visitedList, r.Name)
// 			fmt.Printf("Visited %v\n", visitedList)
// 			fmt.Println("start", startRoom.Name, "end", endRoom)
// 		}
// 	}
// }

func DFS(r *room, f antFarm) {

	startRoom := f.startRoom()

	// set the room being checked visited status to true
	if r.Name != endRoom {
		r.visited = true

		// range through the neighbours of the r
		for _, adjRoom := range r.Tunnel {
			if !adjRoom.visited {
				/* for each neighbour that hasn't been visited,
				- append their key to the visited slice,
				- then apply dfs to them recursively,
				- then append their key to their path value
				*/

				adjRoom.path = append(r.path, adjRoom)
				if contains(adjRoom.path, endRoom) {

					thePaths = append(thePaths, adjRoom.path)

				}

				DFS(adjRoom, antFarm{f.Rooms})

			}

		}

	} else {

		if len(startRoom.Tunnel) > 1 && !contains(startRoom.Tunnel, endRoom) {

			startRoom.Tunnel = startRoom.Tunnel[1:][:]

			DFS(startRoom, antFarm{f.Rooms})

		} else {

		}
	}
	thePaths = PathDupeCheck(thePaths)
	//fmt.Println(thePaths)

}

func PathDupeCheck(path [][]*room) [][]*room {

	dataMap := make(map[*room][]*room)

	for _, item := range path {
		if value, ok := dataMap[item[0]]; !ok {
			dataMap[item[0]] = item
		} else {
			if len(item) <= len(value) {
				dataMap[item[0]] = item

			}
		}
	}

	var output [][]*room

	for _, value := range dataMap {
		output = append(output, value)
	}

	return output
}

// func (a *theAnts) antOutput() {

// 	file := readFile(os.Args[1])
// 	ants := printAnts(file)
// 	waitingAnts := []*ant{}
// 	movingAnts := []*ant{}

// 	for i := 1; i <= ants; i++ {
// 		a.allAnts = append(a.allAnts, &ant{Name: "L" + strconv.Itoa(i)})
// 		a.allAnts = append(a.allAnts, &ant{AntPath: pathSlice[0]})
// 	}

// 	waitingAnts = append(waitingAnts, a.allAnts...)

// 	for _, ant := range a.allAnts {
// 		for _, room := range ant.AntPath {
// 			room.occupied = false
// 			if len(waitingAnts)-1 == 0 {
// 				break
// 			}
// 			if !room.occupied {
// 				movingAnts = append(movingAnts, waitingAnts[0])
// 				ant.RoomsPassed = 0
// 				for _, ant := range movingAnts {
// 					ant.CurrentRoom = *room
// 					if !room.occupied && ant.RoomsPassed <= 1 {
// 						ant.RoomsPassed += 1
// 						fmt.Println(ant.Name + "-" + room.Name)
// 					}
// 				}
// 			}
// 		}
// 	}

// }

func getEnd() string {
	var end string
	file := readFile(os.Args[1])

	for i, _ := range file {
		if file[i] == "##end" {
			end = strings.Split(string(file[i+1]), " ")[0]
		}
	}
	return end
}

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
			fmt.Println("from:", from, "	to:", to)

		}

	}

	ourFarm.addTunnel(from, to)
	ourFarm.showRooms()
	DFS(ourFarm.startRoom(), *ourFarm)
	//ants := theAnts{}
	//ants.antOutput()

}
