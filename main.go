package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type antFarm struct {
	Rooms []*room
}

type room struct {
	Name     string
	Tunnel   []*room //adjacent vertex (neighbours)
	Visited  bool
	Path     []*room
	Occupied bool
}

type Ants struct {
	antz []*Ant
}

type Ant struct {
	Name        string
	Path        []*room
	CurrentRoom room
}

var (
	from       string
	to         string
	therooms   string
	validPaths [][]*room
	dfsPaths   [][]*room
	bfsPaths   [][]*room
)

//=====================================BUILD THE ANT FARM ===============================\\

//steps:
//1. read the file
//2. get the rooms
//3. add the tunnels
//4. find the start room
//5. find the end room

func readFile() []string {
	file, _ := os.Open(os.Args[1])
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines
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

func contains(rList []*room, rname string) bool {
	for _, val := range rList {
		if rname == val.Name {
			return true
		}
	}
	return false
}

func (f *antFarm) startRoom() *room {
	var start string

	file := readFile()

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
	file := readFile()
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

//=====================================DFS AND BFS ALGORITHMS ===============================\\

//The two algorithms will be used as a way to compare paths and select the more efficient one of the 2.

func DFS(r *room, f antFarm) {
	eRoom := f.endRoom().Name
	// set the current room being checked's visited status to true
	if r.Name != f.endRoom().Name {
		r.Visited = true
		// range through the neighbours of the room being checked
		for _, nbr := range r.Tunnel {
			if !nbr.Visited {
				/* for each neighbour that hasn't been visited,
				- append their key to the visited slice,
				- then apply dfs to them recursively,
				- then append their key to their path value
				*/
				nbr.Path = append(r.Path, nbr)
				if contains(nbr.Path, eRoom) {
					dfsPaths = append(dfsPaths, nbr.Path)
				}
				DFS(nbr, antFarm{f.Rooms})
			}
		}
	} else {
		if len(f.startRoom().Tunnel) > 1 && !contains(f.startRoom().Tunnel, eRoom) {
			f.startRoom().Tunnel = f.startRoom().Tunnel[1:][:]
			DFS(f.startRoom(), antFarm{f.Rooms})
		} else {
		}
	}
	dfsPaths = PathDupeCheck(dfsPaths)

}

func PathDupeCheck(path [][]*room) [][]*room {

	farmMap := make(map[*room][]*room)

	for _, room := range path {
		if value, ok := farmMap[room[0]]; !ok {
			farmMap[room[0]] = room
		} else {
			if len(room) <= len(value) {
				farmMap[room[0]] = room

			}
		}
	}

	var output [][]*room

	for _, room := range farmMap {
		output = append(output, room)
	}

	return output
}

func BFS(r *room, f antFarm) {

	var vPaths [][]*room

	//queue will add rooms to be visited to the queue in a FiFo order
	var queue []*room

	//set the start room as visited (as it is the current room for all ants at start of traversal)
	r.Visited = true

	//initialise the queue(begin with the start room)
	queue = append(queue, r)

	// this loop checks if there is a direct tunnel between the start room and end room
	for i, v := range f.startRoom().Tunnel {
		if v.Name == f.endRoom().Name {
			f.endRoom().Path = append(f.endRoom().Path, f.startRoom())
			vPaths = append(vPaths, f.endRoom().Path)
			f.startRoom().Tunnel = append(f.startRoom().Tunnel[:i], f.startRoom().Tunnel[i+1:]...)
		}

	}

	//checks the queue for the end room and if the queue is not empty

	for !contains(queue, f.endRoom().Name) && len(queue) >= 1 {
		qfront := queue[0]

		for _, room := range qfront.Tunnel {
			if !room.Visited {
				room.Visited = true
				room.Path = append(qfront.Path, room)
				//
				queue = append(queue, room)
			}

		}

		queue = queue[1:]

		if doesContainRoom(queue, f.endRoom().Name) {

			for _, room := range f.Rooms {
				room.Visited = false
			}
			vPaths = append(vPaths, qfront.Path)

			for _, r := range qfront.Path {
				deleteTunnel(r, f)

			}
			if len(f.startRoom().Tunnel) == 0 {

				break
			}

			if len(f.startRoom().Tunnel) >= 1 {
				for _, froom := range f.startRoom().Tunnel {
					for _, sroom := range froom.Tunnel {
						if sroom.Name != f.endRoom().Name {
							break
						} else {
							BFS(f.startRoom(), antFarm{f.Rooms})
							queue = queue[1:]
						}
					}
				}
			}
			BFS(f.startRoom(), antFarm{f.Rooms})

		}
	}
	for _, v := range vPaths {
		v = append(v, f.endRoom())
		bfsPaths = append(bfsPaths, v)
	}
	bfsPaths = PathDupeCheck(bfsPaths)

}

// delete edge from starting room
func deleteTunnel(r *room, f antFarm) {
	for i := 0; i < len(r.Path); i++ {
		for _, room := range f.Rooms {
			//	for _ , edge := range room.Tunnel
			for j := 0; j < len(room.Tunnel); j++ {
				if room.Tunnel[j] == r.Path[i] {
					room.Tunnel = remove(room.Tunnel, r.Name)
				}
			}
		}
	}
}

// removes a string from a slice (unordered)
func remove(s []*room, k string) []*room {
	for i := 0; i < len(s); i++ {
		if s[i].Name == k {
			s[i] = s[len(s)-1]

		}

	}
	return s[:len(s)-1]
}

func doesContainRoom(sl []*room, s string) bool {

	for _, word := range sl {
		if s == word.Name {
			return true
		}
	}
	return false
}

func lowestInt(a [][]int, b [][]*room) (int, []*room) {

	min := 10000
	var path []*room

	for i := 0; i < len(a); i++ {
		if a[i][0] < min {
			min = a[i][0]
			path = b[i]

		}

	}
	return min, path
}

// increments the zero index for the given array
func Increment(a [][]int, b int) [][]int {

	for _, slice := range a {
		if slice[0] == b {
			slice[0] += 1
			break
		}
	}
	return a

}

func DeleteAnt(a []*Ant, b *Ant) []*Ant {
	ret := make([]*Ant, 0)
	if len(a) == 1 {
		return []*Ant{}
	}
	for i := 0; i < len(a); i++ {
		if a[i].Name == b.Name {
			ret = append(ret, a[:i]...)
			ret = append(ret, a[i+1:]...)
		}
	}
	return ret
}

func pathSlice(a [][]*room) [][]int {
	var slice [][]int
	var s []int

	for i := range a {
		s = append(s, len(a[i]))
		slice = append(slice, s)
		s = []int{}
	}

	return slice
}

func reassign(a [][]*room) [][]*room {

	sort.Slice(a, func(i, j int) bool {
		return len(a[i]) < len(a[j])
	})

	return a

}

// returns the optimal path between bfs & dfs algos
func pathAssign(bfs [][]*room, dfs [][]*room) [][]*room {

	bfsPathNum := len(bfs)
	dfsPathNum := len(dfs)

	if bfsPathNum > dfsPathNum {
		validPaths = append(validPaths, bfsPaths...)
	} else if dfsPathNum > bfsPathNum {
		validPaths = PathDupeCheck(append(validPaths, dfsPaths...))
	} else {

		bfscounter := 0

		dfscounter := 0

		for _, path := range bfs {

			bfscounter += len(path)

		}

		for _, path := range dfs {
			dfscounter += len(path)
		}

		if bfscounter < dfscounter {
			validPaths = append(validPaths, bfs...)
		} else if dfscounter < bfscounter {
			validPaths = append(validPaths, dfs...)
		} else {
			validPaths = append(validPaths, bfs...)
		}

	}
	return validPaths

}

func printAnts(fname []string) int {
	file := readFile()
	ant := file[0]
	if file[0] <= "0" {
		err := fmt.Errorf("invalid number of ants")
		fmt.Println(err.Error())
	}
	theAnts, _ := strconv.Atoi(string(ant))

	return theAnts
}

func main() {
	file := readFile()
	fmt.Println(strings.Repeat("-", 30))
	fmt.Println("The Number of Ants:")
	fmt.Println(strings.Repeat("-", 30))
	fmt.Println(file[0])

	//creat TWO farms:

	//--------------------------------------FARM ONE: BFS FARM ----------------------------\\
	fmt.Println(strings.Repeat("-", 30))
	fmt.Println("The Rooms:")
	fmt.Println(strings.Repeat("-", 30))
	ourBFSFarm := &antFarm{}
	for _, v := range file[1:] {
		if !strings.Contains(v, "-") && strings.Contains(string(v), " ") {
			//step 2:
			words := strings.Fields(string(v))
			therooms = words[0]
			ourBFSFarm.Rooms = append(ourBFSFarm.Rooms, &room{Name: therooms})
			fmt.Printf("Room: %s\n", therooms)
		}

	}

	var BFSremdash string
	fmt.Println(strings.Repeat("-", 30))
	fmt.Println("The Tunnels:")
	fmt.Println(strings.Repeat("-", 30))

	for _, v := range file[1:] {
		if strings.Contains(v, "-") {
			BFSremdash = strings.Replace(v, "-", " ", 1)
			words := regexp.MustCompile(" ").Split(BFSremdash, -1)
			from = words[0]
			to = words[1]
			fmt.Println("from:", from, "	to:", to)
			ourBFSFarm.addTunnel(to, from)

		}

	}
	fmt.Println(strings.Repeat("-", 30))
	fmt.Println(strings.Repeat("=", 30))
	fmt.Printf("##start: %s\n##end  : %s\n", ourBFSFarm.startRoom().Name, ourBFSFarm.endRoom().Name)
	fmt.Println(strings.Repeat("=", 30))

	BFS(ourBFSFarm.startRoom(), *ourBFSFarm)

	//--------------------------------------FARM TWO: DFS FARM --------------------------\\
	ourDFSFarm := &antFarm{}

	for _, v := range file[1:] {
		if !strings.Contains(v, "-") && strings.Contains(string(v), " ") {
			//step 2:
			words := strings.Fields(string(v))
			therooms = words[0]
			ourDFSFarm.Rooms = append(ourDFSFarm.Rooms, &room{Name: therooms})

		}

	}

	var DFSremdash string

	for _, v := range file[1:] {
		if strings.Contains(v, "-") {
			DFSremdash = strings.Replace(v, "-", " ", 1)
			words := regexp.MustCompile(" ").Split(DFSremdash, -1)
			from = words[0]
			to = words[1]
			ourDFSFarm.addTunnel(from, to)

		}

	}

	DFS(ourDFSFarm.startRoom(), *ourDFSFarm)

	// -------------------------------RUN THE ANTS THROUGH BOTH FARMS----------------------------\\
	fmt.Println(strings.Repeat("-", 30))
	fmt.Println("The Path:")
	fmt.Println(strings.Repeat("-", 30))

	arrange := pathSlice(reassign(PathDupeCheck(pathAssign(bfsPaths, dfsPaths))))
	rooms := reassign(PathDupeCheck(pathAssign(bfsPaths, dfsPaths)))

	a := Ants{}
	var unmovedAnts []*Ant
	var movedAnts []*Ant
	counter := 1

	for counter <= printAnts(file) {

		number, _ := lowestInt(arrange, rooms)
		_, route := lowestInt(arrange, rooms)
		a.antz = append(a.antz, &Ant{Name: "L" + strconv.Itoa(counter), Path: route})
		Increment(arrange, number)

		counter++
	}

	unmovedAnts = append(unmovedAnts, a.antz...)

	for len(unmovedAnts) > 0 || len(movedAnts) >= 1 {

		for _, ant := range unmovedAnts {
			if len(ant.Path) == 1 {
				fmt.Print(ant.Name, "-", ant.Path[0].Name, " ")
				ant.Path[0].Occupied = true
				unmovedAnts = DeleteAnt(unmovedAnts, ant)
				break
			}
		}

		for _, ant := range unmovedAnts {

			if !ant.Path[0].Occupied {
				fmt.Print(ant.Name, "-", ant.Path[0].Name, " ")
				ant.Path[0].Occupied = true
				movedAnts = append(movedAnts, ant)
				unmovedAnts = DeleteAnt(unmovedAnts, ant)

			}

		}

		fmt.Println()

		for _, ant := range movedAnts {

			if len(ant.Path) > 1 {
				ant.Path[0].Occupied = false

				ant.Path = ant.Path[1:]
				fmt.Print(ant.Name, "-", ant.Path[0].Name, " ")

			} else {
				movedAnts = DeleteAnt(movedAnts, ant)
				ant.Path = []*room{}
			}

		}

	}

	fmt.Println()

}
