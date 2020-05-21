package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

// Room struct, composed of data found in room.txt
type Room struct {
	Floor        int
	RoomNo       int
	Capacity     int
	TimeSlots    []int
	Availability []int
}

// Input struct, composed of data found in input.txt
type Request struct {
	TeamSize  int
	Floor     int
	StartTime int
	EndTime   int
	Duration  int
}

func main() {
	// Variable to have best-case scenario meeting room stored
	var closestRoom Room

	// Read from rooms.txt and save file contents to a string
	fromOutputFile, err := getRoomsFromOutputFile()
	if err != nil {
		fmt.Println("Error retrieving from file: ", err)
	}

	// Convert the string from the file to a more usable Room struct
	rooms, err := convertToStruct(fromOutputFile)
	if err != nil {
		fmt.Println("Error creating struct: ", err)
	}

	// Convert file content from input.txt to a Request struct
	requestInput, err := getRequestFromInputFile()
	if err != nil {
		fmt.Println("Error creating struct: ", err)
	}

	// Implementing a binary search here allows for both good
	// runtime and to ensure that the room chosen is the
	// closest available distance to the request
	middle := findStartingPoint(rooms, requestInput.Floor)
	right := middle + 1
	left := middle
	isFound := false

	for isFound == false {
		if left >= 0 {
			for i := 0; i < len(rooms[left].TimeSlots)-1; i += 2 {
				if rooms[left].TimeSlots[i] == requestInput.StartTime && // First check for start time compatibility
					rooms[left].Availability[i/2] >= requestInput.Duration && // Next check if there's enough time allotted
					rooms[left].Capacity >= requestInput.TeamSize { // Finally check for space for the team members
					isFound = true
					closestRoom = rooms[left]
				}
			}
			left--
		}
		if right < len(rooms) {
			for i := 0; i < len(rooms[right].TimeSlots)-1; i += 2 {
				if rooms[right].TimeSlots[i] == requestInput.StartTime &&
					rooms[right].Availability[i/2] >= requestInput.Duration &&
					rooms[right].Capacity >= requestInput.TeamSize {
					isFound = true
					closestRoom = rooms[right]
				}
			}
			right++
		}
	}

	// Print closet Meeting room floor and Room # to screen
	fmt.Println("Closet meeting room is on floor:", closestRoom.Floor, "in Room #", closestRoom.RoomNo)

}

// getRoomsFromOutputFile reads rooms.txt and returns the information in a string
func getRoomsFromOutputFile() (string, error) {
	roomBytes, err := ioutil.ReadFile("./rooms.txt")
	if err != nil {
		fmt.Println("Error reading from file: ", err)
		return "", err
	}
	roomInfo := string(roomBytes)
	return roomInfo, nil
}

// Split allows for a string to be split by multiple parameters (, and \n)
func Split(r rune) bool {
	return r == ',' || r == '\n'
}

// convertToStruct assigns multiple arrays of information pulled from the
// rooms.txt file into an array of Room structs
func convertToStruct(fromFile string) ([]Room, error) {
	// count how many lines (rooms) are in input.txt
	roomNumGetter := strings.Split(fromFile, "\n")
	totalRooms := len(roomNumGetter)
	availableRooms := make([]Room, totalRooms)

	floors, rooms, capacities, times, err := convertToArrays(fromFile, totalRooms)
	if err != nil {
		fmt.Println("Error in creating Struct: ", err)
		return nil, err
	}

	for i := 0; i < totalRooms; i++ {
		availableRooms[i].Floor = floors[i]
		availableRooms[i].RoomNo = truncateZero(rooms[i])
		availableRooms[i].Capacity = capacities[i]
		availableRooms[i].TimeSlots, _ = convertTimes(times[i])
		availableRooms[i].Availability = getAvailability(availableRooms[i].TimeSlots)
	}
	return availableRooms, nil
}

// convertToArrays extracts data from a long string and translates to either ints, floats,
// or separate strings to separate information into more useful data structures
func convertToArrays(fromFile string, total int) ([]int, []int, []int, [][]string, error) {
	var floors []int
	var rooms []int
	var capacities []int
	times := make([][]string, total)
	index := -1

	split := strings.FieldsFunc(fromFile, Split) // Split by both , and \n
	for i := 0; i < len(split); i++ {
		if !strings.ContainsAny(split[i], ":") {
			if !strings.ContainsAny(split[i], ".") {
				convertedInt, err := strconv.Atoi(split[i])
				if err != nil {
					fmt.Println("Error converting to int ", err)
					return nil, nil, nil, nil, err
				}
				capacities = append(capacities, convertedInt)
			}
		}
		if strings.ContainsAny(split[i], ".") {
			secondSplit := strings.Split(split[i], ".")
			convertedFloor, err := strconv.Atoi(secondSplit[0])
			if err != nil {
				fmt.Println("Error converting to int ", err)
				return nil, nil, nil, nil, err
			}
			convertedRoomNo, err := strconv.Atoi(secondSplit[1])
			if err != nil {
				fmt.Println("Error converting to int ", err)
				return nil, nil, nil, nil, err
			}

			floors = append(floors, convertedFloor)
			rooms = append(rooms, convertedRoomNo)
			index++
		}
		if strings.ContainsAny(split[i], ":") {
			times[index] = append(times[index], split[i])
		}
	}

	return floors, rooms, capacities, times, nil
}

// truncateZero eliminates the extra zero in 2-digit room numbers after translation
func truncateZero(room int) int {
	if room%10 == 0 {
		room /= 10
	}
	return room
}

// convertTimes converts the string starting and ending times,
// by removing the : and converting the string to an int
func convertTimes(times []string) ([]int, error) {
	var returnInt []int

	for i := 0; i < len(times); i++ {
		split := strings.Split(times[i], ":")
		removedColon := split[0] + split[1]
		convertedInt, err := strconv.Atoi(removedColon)
		if err != nil {
			fmt.Println("Error converting to int ", err)
			return nil, err
		}

		// Translate 15 minute increments into a more workable base-10 model
		lastTwoDigits := convertedInt % 100
		if lastTwoDigits == 15 {
			convertedInt -= 15
			convertedInt += 25
		}
		if lastTwoDigits == 30 {
			convertedInt -= 30
			convertedInt += 50
		}
		if lastTwoDigits == 45 {
			convertedInt -= 45
			convertedInt += 75
		}
		returnInt = append(returnInt, convertedInt)
	}
	return returnInt, nil
}

// getAvailability subtracts the open room's end time from its start time
// allowing for easier comparison between request and room status
func getAvailability(times []int) []int {
	var returnInt []int
	for i := 0; i < len(times)-1; i += 2 {
		returnInt = append(returnInt, times[i+1]-times[i])
	}
	return returnInt
}

// getRequestFromInputFile reads input from input.txt, then
// translates the data into more usable data structures
func getRequestFromInputFile() (Request, error) {
	var returnRequest Request

	requestBytes, err := ioutil.ReadFile("./input.txt")
	if err != nil {
		fmt.Println("Error reading from file: ", err)
		return returnRequest, err
	}
	requestInfo := string(requestBytes)

	split := strings.Split(requestInfo, ",")

	teamSize, err := strconv.Atoi(split[0])
	if err != nil {
		fmt.Println("Error converting to int ", err)
		return returnRequest, err
	}

	Floor, err := strconv.Atoi(split[1])
	if err != nil {
		fmt.Println("Error converting to int ", err)
		return returnRequest, err
	}

	startArr := strings.Split(split[2], ":")

	startString := startArr[0] + startArr[1]

	convertedInt, err := strconv.Atoi(startString)
	if err != nil {
		fmt.Println("Error converting to int ", err)
		return returnRequest, err
	}

	lastTwoDigits := convertedInt % 100
	if lastTwoDigits == 15 {
		convertedInt -= 15
		convertedInt += 25
	}
	if lastTwoDigits == 30 {
		convertedInt -= 30
		convertedInt += 50
	}
	if lastTwoDigits == 45 {
		convertedInt -= 45
		convertedInt += 75
	}

	endArr := strings.Split(split[3], ":")

	endString := endArr[0] + endArr[1]

	convertedEndInt, err := strconv.Atoi(endString)
	if err != nil {
		fmt.Println("Error converting to int ", err)
		return returnRequest, err
	}

	lastTwoDigits = convertedEndInt % 100
	if lastTwoDigits == 15 {
		convertedEndInt -= 15
		convertedEndInt += 25
	}
	if lastTwoDigits == 30 {
		convertedEndInt -= 30
		convertedEndInt += 50
	}
	if lastTwoDigits == 45 {
		convertedEndInt -= 45
		convertedEndInt += 75
	}

	returnRequest.TeamSize = teamSize
	returnRequest.Floor = Floor
	returnRequest.StartTime = convertedInt
	returnRequest.EndTime = convertedEndInt
	returnRequest.Duration = convertedEndInt - convertedInt

	return returnRequest, nil
}

// findMinMax is a precursor to the binary sort, it finds the smallest
// and largest room number of the given Room array
func findMinMax(rooms []Room) (min, max int) {
	min = 9999
	max = -1

	for i := 0; i < len(rooms); i++ {
		if rooms[i].Floor < min {
			min = rooms[i].Floor
		}
		if rooms[i].Floor > max {
			max = rooms[i].Floor
		}
	}
	return min, max
}

// findStartingPoint finds index of a Room array that is the same as the
// floor of the requesting team
func findStartingPoint(rooms []Room, request int) int {
	for i := 0; i < len(rooms); i++ {
		if rooms[i].Floor == request {
			return i
		}
	}
	return request
}
