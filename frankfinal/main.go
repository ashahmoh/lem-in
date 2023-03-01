package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type Room struct {
	Id         int
	Name       string
	Neighbours []*Room
}
type Graph struct {
	Rooms map[int]*Room
}

// newGraph creates a new Graph object and initializes its rooms map
func NewGraph() *Graph {
	return &Graph{
		Rooms: make(map[int]*Room),
	}
}

// addRoom adds a new room to the graph with the given id and name
func (g *Graph) AddRoom(id int, name string) *Room {
	// create a new Room object with the given id and name
	r := &Room{
		Id:   id,
		Name: name,
	}
	// add the room to the graph's rooms map
	g.Rooms[id] = r
	// return the room
	return r
}

// addLink creates a two-way link between two rooms in the graph with the given ids
func (g *Graph) AddLink(id1, id2 int) {
	// retrieve the rooms with the given ids from the graph's rooms map
	r1 := g.Rooms[id1]
	r2 := g.Rooms[id2]
	// add each room to the other's neighbours list to create a two-way link
	r1.Neighbours = append(r1.Neighbours, r2)
	r2.Neighbours = append(r2.Neighbours, r1)
}

func FilterData(data []string) (antNbr int, allRooms, links []string) {
	// Convert the first element of data, representing the number of ants, into an integer.
	antNbr, _ = strconv.Atoi(data[0])
	// Find the index of the line containing "##start".
	startIndex := -1
	for i, line := range data {
		if strings.Contains(line, "##start") {
			startIndex = i
			break
		}
	}
	// Extract the room name associated with "##start".
	var startRoom []string
	if startIndex != -1 {
		temp := data[startIndex+1]
		parts := strings.Split(temp, " ")
		startRoom = append(startRoom, parts[0])
	}
	// Find the index of the line containing "##end".
	endIndex := -1
	for i, line := range data {
		if strings.Contains(line, "##end") {
			endIndex = i
			break
		}
	}
	// Extract the room name associated with "##end".
	var endRoom string
	if endIndex != -1 {
		temp := data[endIndex+1]
		parts := strings.Split(temp, " ")
		endRoom = parts[0]
	}
	// Extract the names of all rooms, excluding the start and end rooms.
	var rooms []string
	for i, line := range data {
		if i == startIndex+1 || i == endIndex+1 {
			continue
		}
		if strings.Contains(line, " ") {
			parts := strings.Split(line, " ")
			rooms = append(rooms, parts[0])
		}
	}
	// Combine the start room, all rooms, and end room into a single array.
	allRooms = append(startRoom, rooms...)
	allRooms = append(allRooms, endRoom)
	for _, line := range data {
		linkedrooms := extractLinks(line)
		//links refers to the roomlinks that have now been seperated.
		//so to rooms and from rooms
		links = append(links, linkedrooms...)
	}
	return antNbr, allRooms, links
}
func extractLinks(s string) []string {
	// Split the input string into words using the space character as a delimiter.
	words := strings.Split(s, " ")
	// Create an empty slice to store the links.
	var links []string
	// Iterate over the slice of words.
	for _, word := range words {
		// Check if the word contains a hyphen.
		if strings.Contains(word, "-") {
			// If it does, add the word to the links slice.
			links = append(links, word)
		}
	}
	// Return the slice of links.
	return links
}

func (g *Graph) dfs(currentRoom, end *Room, visited []bool, path []*Room, paths *[][]*Room) {
	// Mark the current room as visited
	visited[currentRoom.Id] = true
	// Add the current room to the path
	path = append(path, currentRoom)
	// If the current room is the end room, add the path to the list of all paths
	if currentRoom == end {
		pathCopy := make([]*Room, len(path))
		copy(pathCopy, path)
		*paths = append(*paths, pathCopy)
	} else {
		// Recursively search the neighbours of the current room
		for _, neighbour := range currentRoom.Neighbours {
			if !visited[neighbour.Id] {
				g.dfs(neighbour, end, visited, path, paths)
			}
		}
	}
	// Backtrack by removing the current room from the path
	path = path[:len(path)-1]
	// Mark the current room as not visited
	visited[currentRoom.Id] = false
}

// findPaths is a function that finds all paths from a start room to an end room in a graph
func (g *Graph) FindPaths(start, end *Room) [][]*Room {
	// Create a list to keep track of the rooms that have been visited
	visited := make([]bool, len(g.Rooms))
	// Initialize an empty path
	path := []*Room{}
	// Initialize an empty list of paths
	paths := [][]*Room{}
	// Call the depth-first search function to find all paths from the start room to the end room
	g.dfs(start, end, visited, path, &paths)
	// Return the list of all paths
	return paths
}

func ReadData() ([]string, error) {
	if len(os.Args) != 2 {
		return nil, fmt.Errorf("ERROR: Incorrect number of arguments.\ninput format: go run . example00.txt")
	}
	path := os.Args[1]
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("ERROR: Failed to open the file: %v", err)
	}
	defer file.Close()
	var data []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data = append(data, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("ERROR: Failed to read the file: %v", err)
	}
	return data, scanner.Err()
}
func PrintOutput(data []string) {
	antNbr, allRooms, allLinks := FilterData(data)
	if antNbr <= 0 {
		fmt.Println("ERROR: Invalid data format. Invalid number of ants. Must be > 0")
		return
	}
	graph := NewGraph()
	//add rooms to the graph
	roomIDs := make(map[string]int)
	for id, room := range allRooms {
		roomIDs[room] = id
		graph.AddRoom(id, room)
	}
	//add links to the graph
	for _, link := range allLinks {
		parts := strings.Split(link, "-")
		id1 := roomIDs[parts[0]]
		id2 := roomIDs[parts[1]]
		graph.AddLink(id1, id2)
	}
	//assign start and end points
	startRoom := graph.Rooms[0]
	endRoom := graph.Rooms[len(graph.Rooms)-1]
	paths := graph.FindPaths(startRoom, endRoom)
	if len(paths) == 0 {
		fmt.Println("ERROR: Invalid data format. No path found or text file is formatted incorrectly")
		return
	}
	validPaths := FindCompatiblePaths(paths)
	bestPath := PathAssign(paths, validPaths, antNbr)
	path := os.Args[1]
	bytes, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	content := string(bytes)
	fmt.Println(content)
	fmt.Println()
	PrintAntSteps(paths, bestPath)
}

var requiredSteps int

// findCompatiblePaths is a function that takes a 2D array of paths and returns a 2D array of compatible paths.
func FindCompatiblePaths(paths [][]*Room) [][]int {
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
func PathAssign(paths [][]*Room, validPaths [][]int, antNbr int) []string {
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

// printAntSteps is a function that takes the filtered paths and the assigned paths and prints the steps taken by each ant.
func PrintAntSteps(filteredPaths [][]*Room, pathStrings []string) {
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

// getRow is a helper function that takes a 2D slice and a value to search for, and returns the index of the first row in the slice that does not contain the value.
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

func main() {
	data, err := ReadData()
	if err != nil {
		fmt.Println(err)
		return
	}
	PrintOutput(data)
}
