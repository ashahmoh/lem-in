package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type antFarm struct {
	Rooms []*room
}

type room struct { //vertex
	Name     string
	Links    []*room //links
	Id       int
	visited  bool
	occupied bool
	Tunnel   []*room //adjacent vertex
	ants     int
}

var (
	from string
	to   string
	//pathSlice [][]*room
	thePaths [][]*room
	//endRoom  = getEnd()
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
	//DFS(ourFarm.startRoom(), *ourFarm)
	//ants := theAnts{}
	//ants.antOutput()
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
		for _, neighbour := range current.Links {
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

var requiredSteps int

func (f *antFarm) FindPaths(start, end *room) [][]*room {
	// Create a list to keep track of the rooms that have been visited
	visited := make([]bool, len(f.Rooms))

	// Initialize an empty path
	path := []*room{}

	// Initialize an empty list of paths
	paths := [][]*room{}

	// Call the depth-first search function to find all paths from the start room to the end room
	f.dfs(start, end, visited, path, &paths)

	// Return the list of all paths
	return paths
}

// findCompatiblePaths is a function that takes a 2D array of paths and returns a 2D array of compatible paths.
func FindCompatiblePaths(paths [][]*room) [][]int {
	// Initialize a 2D slice to store the compatible paths.
	var compatiblePaths [][]int
	// Loop through each path in the array, and compare it to every subsequent path in the array.
	for i, path1 := range paths {
		// Add the index of the current path to the compatiblePaths slice as a new array containing only that index.
		compatiblePaths = append(compatiblePaths, []int{i})
		// Create a map to keep track of the rooms in the current path.
		roomMap := make(map[int]struct{})
		// Loop through each room in the current path and add it to the roomMap.
		for _, room := range path1[1 : len(path1)-1] {
			roomMap[room.Id] = struct{}{}
		}
		// Loop through each subsequent path and compare it to the current path.
		for j, path2 := range paths[i+1:] {
			// Assume that the two paths are compatible.
			isCompatible := true
			// Loop through each room in the current path and check if it appears in the roomMap of the other path.
			for _, room := range path2[1 : len(path2)-1] {
				if _, ok := roomMap[room.Id]; ok {
					// If a room appears in both paths, the paths are not compatible.
					isCompatible = false
					break
				}
			}
			// If the paths are compatible, add the index of the other path to the compatiblePaths slice for the current path.
			if isCompatible {
				compatiblePaths[i] = append(compatiblePaths[i], i+1+j)
				// Loop through each room in the other path and add it to the roomMap.
				for _, room := range path2[1 : len(path2)-1] {
					roomMap[room.Id] = struct{}{}
				}
			}
		}
	}
	// Return the compatiblePaths slice.
	return compatiblePaths
}

// pathAssign is a function that takes the 2D array of paths and the compatible paths, along with the number of ants, and assigns a path to each ant.
func PathAssign(paths [][]*room, validPaths [][]int, antNbr int) []string {
	// Initialize variables to keep track of the best assigned path and its maximum step length.
	var bestAssignedPath []string
	bestMaxStepLength := math.MaxInt32
	// Loop through each valid path.
	for _, validPath := range validPaths {
		// Initialize a slice to store the step lengths of each path in the current valid path.
		var stepLength []int
		// Initialize a slice to store the assigned path for each ant.
		var assignedPath []string
		// Loop through each index in the current valid path and add the step length of the corresponding path to the stepLength slice.
		for _, pathIndex := range validPath {
			path := paths[pathIndex]
			stepLength = append(stepLength, len(path)-1)
		}
		// Loop through each ant.
		for i := 1; i <= antNbr; i++ {
			// Find the path in the valid path with the shortest step length and assign the ant to that path
			minStepsIndex := 0
			for j, steps := range stepLength {
				if steps <= stepLength[minStepsIndex] {
					minStepsIndex = j
				}
			}
			assignedPath = append(assignedPath, fmt.Sprintf("%d-%d", i, validPath[minStepsIndex]))
			stepLength[minStepsIndex]++
		}
		// Calculate the maximum step length in the assigned path.
		maxStepLength := 0
		for _, steps := range stepLength {
			if steps > maxStepLength {
				maxStepLength = steps
			}
		}
		// If the maximum step length in the assigned path is less than the best maximum step length so far, update the bestAssignedPath and bestMaxStepLength.
		if maxStepLength < bestMaxStepLength {
			bestAssignedPath = assignedPath
			bestMaxStepLength = maxStepLength
		}
	}
	// Store the required number of steps as the best maximum step length.
	requiredSteps = bestMaxStepLength
	// Return the best assigned path.
	return bestAssignedPath

}

func PrintAntSteps(filteredPaths [][]*room, pathStrings []string) {
	// Initialize a 2D slice to store the steps taken by each ant in order.
	var antSteps [][]string
	// Calculate the number of turns required to complete the path.
	arrayLen := requiredSteps - 1
	// Initialize a slice to store the steps taken by each ant in order.
	orderedSteps := make([][]string, arrayLen)
	// Loop through each assigned path.
	for _, antPath := range pathStrings {
		// Initialize a slice to store the steps taken by the current ant.
		var steps []string
		// Split the antPath string into its ant number and path index components.
		parts := strings.SplitN(antPath, "-", 2)
		antStr := parts[0]
		antPath, _ := strconv.Atoi(string(parts[1]))
		// Loop through each room in the path and add a step string to the steps slice for each room.
		for i := 1; i < len(filteredPaths[antPath]); i++ {
			path := filteredPaths[antPath][i]
			temp := "L" + antStr + "-" + path.Name
			steps = append(steps, temp)
		}
		// Add the steps slice to the antSteps slice.
		antSteps = append(antSteps, steps)
	}
	// Loop through each step in each ant's path and add it to the orderedSteps slice in order.
	for i := 0; i < len(antSteps); i++ {
		slice := antSteps[i]
		var row int
		for j := 0; j < len(slice); j++ {
			temp := slice[j]
			if j == 0 {
				// Split the step string to get the room name and use getRow to find the first row in orderedSteps that does not contain the room name.
				parts := strings.SplitN(temp, "-", 2)
				row = getRow(orderedSteps, "-"+parts[1])
			}
			// Add the step string to the row in orderedSteps.
			orderedSteps[j+row] = append(orderedSteps[j+row], temp)
		}
		row = 0
	}
	// Loop through each step in the orderedSteps slice and print it.
	for i, printline := range orderedSteps {
		fmt.Println(strings.Trim(fmt.Sprint(printline), "[]"))
		if i == len(orderedSteps)-1 {
			fmt.Println()
			fmt.Printf("Number of turns: %v\n", i+1)
		}
	}
}

func getRow(tocheck [][]string, value string) int {
	// Loop through each row in the slice.
	for i, row := range tocheck {
		found := false
		// Loop through each item in the current row and check if it contains the value.
		for _, item := range row {
			if strings.Contains(item, value) {
				found = true
				break
			}
		}
		// If the value is not found in the current row, return its index.
		if !found {
			return i
		}
	}
	// If the value is found in every row, return 0.
	return 0
}
