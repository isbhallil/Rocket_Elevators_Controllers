package main

import (
	"fmt"
	"strconv"
	"time"
	"sort"
)

//\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\

//							BATTERY

//\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\

// Battery type
type Battery struct {
	numberOfFloors int
	columns        []Column
}

// get the best elevator for the request and move the elevator to the requested floors
func (b *Battery) requestElevator(floorNumber, target int, direction string) *Elevator{
	c := b.getRequestedColumn(floorNumber, target)
	e := c.getBestElevator(floorNumber, direction)
	go e.assignElevator(floorNumber)
	go e.assignElevator(target)
	
	return e
}

// get the column that fits the request
func (b*Battery)getRequestedColumn(floorNumber, target int) Column {
	for _, c := range b.columns {
		if keep(c, floorNumber, target) {
			return c
		}
	}

	return b.columns[0]
}

// provide de condition to select the column that fits the request
func keep(c Column, floor , target int) bool {
	if (floor != 1){
		return c.lowest <= floor && c.highest >= floor
	} 

	return c.lowest <= target && c.highest >= target
}

// get the highest between two values
func getHighestBetween(floor1, floor2 int) int {
	if floor1 > floor2 {
		return floor1
	}

	return floor2
}

// Battery constructor
func newBattery(numberOfFloors int, numberOfElevatorsPerColumn []int) *Battery {
	b := new(Battery)

	fmt.Println("[newBattery]")
	b.columns = initColumns(numberOfFloors, numberOfElevatorsPerColumn)
	b.numberOfFloors = numberOfFloors

	return b
}

// Provide columns to the Battery
func initColumns(numberOfFloors int, numberOfElevatorsPerColumn []int) []Column {
	list := []Column{}

	var floor int = 1
	var columnHeight = 21
	index := 0
	for _, numberOfElevators := range numberOfElevatorsPerColumn {
		var firstFloor = floor + 1
		var lastFloor = columnHeight + floor
		column := newColumn(index, firstFloor, lastFloor, numberOfElevators)
		list = append(list, *column)

		floor = columnHeight + floor
		index++
	}

	return list
}





//\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\

//							COLUMN

//\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\

// Column type
type Column struct {
	id, lowest, highest	int
	floorButtons 		[]*FloorButton
	elevators    		[]*Elevator
}

// get the quickest elevator for the request
func (c*Column) getBestElevator(floor int, direction string) (*Elevator){
	commingElevators := []*Elevator{}
	idlingElevators := []*Elevator{}
	otherElevators := []*Elevator{}

	for _, e := range c.elevators {
		if (len(e.tasksList) == 0){
			fmt.Println(strconv.Itoa(e.id) + "is idling")
			idlingElevators = append(idlingElevators, e)
		} else if ( e.isComming(floor, direction) ){
			fmt.Println(strconv.Itoa(e.id) + "is comming")
			commingElevators = append(commingElevators, e)
		} else {
			fmt.Println(strconv.Itoa(e.id) + "is other")
			otherElevators = append(otherElevators, e)
		}
	}

	return c.selectElevator(commingElevators, idlingElevators, otherElevators, floor, direction)
}

// return the best elevator from list according to priority
func (c *Column)selectElevator(priority1, priority2, priority3 []*Elevator , floor int, direction string ) *Elevator{
	if (len(priority1) > 0) {
		return c.getBestElevatorFrom(priority1, floor, direction);
	} else if (len(priority2) > 0) {
		return c.getBestElevatorFrom(priority2, floor, direction);
	} 

	return c.getBestElevatorFrom(priority3, floor, direction);
}

// draw between elevators to select the best fit in the list
func (c*Column) getBestElevatorFrom(list []*Elevator, floor int, direction string) *Elevator {
	bestElevator := list[0]
	bestGap := bestElevator.getGapToReach(floor, direction)

	for _, e := range list {
		gap := e.getGapToReach(floor, direction) 

		if ( gap < bestGap ){
			bestElevator = e
			bestGap = gap
		}
	}

	return bestElevator
}

// Column consutructor
func newColumn(ID, lowestFloor, highestFloor, numberOfElevators int) *Column {
	c := new(Column)
	c.id = ID
	c.lowest = lowestFloor
	c.highest = highestFloor
	c.elevators = initElevators(numberOfElevators, lowestFloor, highestFloor)
	c.floorButtons = initFloorButtons(lowestFloor, highestFloor)

	return c
}

// Provide elevators to the column
func initElevators(numberOfElevators, startingFloor, highestFloor int) []*Elevator {
	list := []*Elevator{}

	for index := 0; index < numberOfElevators; index++ {
		list = append(list, newElevator(index, startingFloor, highestFloor))
	}

	return list
}

// Provide the buttons to the column
func initFloorButtons(lowestFloor, highestFloor int) []*FloorButton {
	list := []*FloorButton{}

	for floor := lowestFloor; floor <= highestFloor; floor++ {
		list = append(list, newFloorButton("up", floor), newFloorButton("down", floor))
	}

	return list
}





//\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\

//							ELEVATOR

//\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\

// Elevator type
type Elevator struct {
	id, currentFloor, maxWeight	int
	isDoorOpen,	isSafe  		bool
	buttonsList  				[]*ElevatorButton
	tasksList    				[]int
}

// return the gap distance between the elevator and the request
func (e*Elevator) getGapToReach(floor int, direction string) int {
	requestDirection := direction
	elevatorDirection := e.getDirection()
	requestFloor := floor
	gapToReach := 0
	previousReach := e.currentFloor

	if len(e.tasksList) == 0 {
		gapToReach = gapToReach + abs(previousReach - requestFloor)
	} else {

		for _, taskFloor := range e.tasksList {
			if (elevatorDirection == "up" && requestFloor >= previousReach && requestFloor <= taskFloor){
				gapToReach = gapToReach + abs(previousReach - requestFloor)
			} else if (requestDirection == "down" && requestFloor <= previousReach && requestFloor >= taskFloor) {
				gapToReach = gapToReach + abs(previousReach - requestFloor)
			} else {
				gapToReach = gapToReach + abs(previousReach - taskFloor)
			}

			previousReach = taskFloor
		}

	}

	return gapToReach + abs(previousReach - floor)
}

// return if elevator is reaching the request or will be able to reach
func (e*Elevator) isComming(requestfloor int, direction string) bool {
	if (e.getDirection() == direction){
		for range e.tasksList {
			if((e.currentFloor <= requestfloor && direction == "up") || (e.currentFloor > requestfloor && direction == "down") ){
				return true
			} 
		}
	}

	return false
}

// return direction of the elevator
func (e*Elevator) getDirection() string {
	if (len(e.tasksList) > 0 && e.currentFloor < e.tasksList[0]){
		return "up"
	} else if ( len(e.tasksList) > 0 && e.currentFloor > e.tasksList[0] ){
		return "down"
	}

	return "idle"
}

// add the request to the elevator's tasks and command to operate
func (e *Elevator) assignElevator(floor int){
	e.addTask(floor).operate()
}

// command to operate and behave according to rocket elevator rules
func (e *Elevator) operate() *Elevator{	
	for len(e.tasksList) > 0 {
		time.Sleep(250)
		e.move()
		if (e.isArrived()){
			e.openDoor()
			screen("arrived at floor " + strconv.Itoa(e.currentFloor))
			e.removeTask()
			time.Sleep(1000)
			e.closeDoor()
		}

	return e
}

// decide witch way to move
func (e *Elevator) move() *Elevator{
	if ( e.isNextTaskAbove() ){
		e.moveUp()
	} else if ( e.isNextTaskBeneath() ) {
		e.moveDown()
	}

	return e
}

// add a floor to the task list
func (e *Elevator) addTask(floor int) *Elevator{
	e.tasksList = append(e.tasksList, floor)

	return e
}

// sort list ascending or descending
func arrange(slice []int, order string) {

	if (order == "ASC"){
		sort.Slice(slice, func(i, j int) bool {
			return slice[i] < slice[j]
		})
	} else if ( order == "DESC"){
		sort.Slice(slice, func(i, j int) bool {
			return slice[i] > slice[j]
		})
	}
	
}

// remove first element of the tasksLiost
func (e *Elevator) removeTask() *Elevator{
	e.tasksList = e.tasksList[1:]
	return e
} 

// return true if the elevator is at a requested floor
func (e *Elevator) isArrived() bool {
	for _, floor := range e.tasksList {
		if e.currentFloor == floor {
			return true
		}
	}
	return false
}

// return true if the nest floor to visit is above the elevator
func (e *Elevator) isNextTaskAbove() bool {
	return e.currentFloor < e.tasksList[0]
}

// return true if the nest floor to visit is beneathe the elevator
func (e *Elevator) isNextTaskBeneath() bool {
	return e.currentFloor > e.tasksList[0]
}

// command the elevator to move 1 floor up
func (e *Elevator) moveDown() *Elevator {
	screen("[Elevator " + strconv.Itoa(e.id) + "]  currentFloor :>" + strconv.Itoa(e.currentFloor) + " | nextFloor :> " + strconv.Itoa(e.tasksList[0]))
	e.currentFloor = e.currentFloor - 1
	return e
}

// command the elevator to move 1 floor down
func (e *Elevator) moveUp() *Elevator {
	screen("[Elevator " + strconv.Itoa(e.id) + "]  currentFloor :>" + strconv.Itoa(e.currentFloor) + " | nextFloor :> " + strconv.Itoa(e.tasksList[0]))
	e.currentFloor = e.currentFloor + 1
	return e
}

// commamd the elevator to open doors
func (e *Elevator) openDoor() *Elevator {
	e.isDoorOpen = true
	screen("[Elevator " + strconv.Itoa(e.id) + "]  currentFloor :>" + strconv.Itoa(e.currentFloor) + " | Door is now open ")
	return e
}

// command the elevator to close doors
func (e *Elevator) closeDoor() *Elevator{
	e.isDoorOpen = false
	screen("[Elevator " + strconv.Itoa(e.id) + "]  currentFloor :>" + strconv.Itoa(e.currentFloor) + " | Door is now closed ")
	return e
}

// Elevator constructor
func newElevator(ID, startingFloor, highestFloor int) *Elevator {
	e := new(Elevator)
	e.id = ID
	e.currentFloor = startingFloor
	e.isDoorOpen = false
	e.isSafe = true
	e.maxWeight = 3400
	e.buttonsList = initElevatorButtons(startingFloor, highestFloor)
	e.tasksList = []int{}

	return e
}

// Provide the Elevator buttons to the elevator
func initElevatorButtons(startingFloor, highestFloor int) []*ElevatorButton {
	list := []*ElevatorButton{}

	for floor := startingFloor; floor <= highestFloor; floor++ {
		list = append(list, newElevatorButton(startingFloor))
		startingFloor++
	}

	return list
}


//\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\

//							FLOOR BUTTON

//\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\

// FloorButton type
type FloorButton struct {
	direction string
	isActive  bool
	floor     int
}

// FloorButton contructor
func newFloorButton(direction string, floor int) *FloorButton {
	f := new(FloorButton)
	f.direction = direction
	f.floor = floor
	f.isActive = false

	return f
}

//\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\

//							ELEVATOR BUTTON

//\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\

// ElevatorButton type
type ElevatorButton struct {
	isActive bool
	floor    int
}

// ElevatorButton contructor
func newElevatorButton(floor int) *ElevatorButton {
	f := new(ElevatorButton)
	f.floor = floor
	f.isActive = false

	return f
}

//\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\

//							MAIN

//\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\

// print on screen a message
func screen(message string){
	fmt.Println("======================================")
	fmt.Println(message)
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
}

// return the absolute value of int value
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// TestElevator Type
type TestElevator struct {
	ID, currentFloor  int
	currentTasks []int
}

// TestElevator constructor
func (b*Battery) newTestElevator(ID ,currentFloor int, currentTasks []int) TestElevator{
	n := TestElevator{}
	n.ID = ID
	n.currentFloor = currentFloor
	n.currentTasks = currentTasks

	fmt.Println("newTestElevator")
	return n
}

// Initialize the scene to perform tests on the system
func (b*Battery) initTest(testElevators []TestElevator, numberOfElevatorsPerColumn []int) *Battery {
	
	i := 0
	Loop:
		for i < len(numberOfElevatorsPerColumn) {
			c := b.columns[i]
			for index, elevator := range testElevators {
				fmt.Printf(strconv.Itoa(index))
				if ((c.lowest <= elevator.currentFloor && c.highest >= elevator.currentFloor) || (c.lowest <= elevator.currentTasks[0] && c.highest >= elevator.currentTasks[0])){
					c.elevators[index].tasksList = elevator.currentTasks
					c.elevators[index].currentFloor = elevator.currentFloor
					if ( index + 1 == len(testElevators)){
						break Loop
					} 
				}

			}
			i = i + 1
		}

	return b 
}


func main() {

	numberOfElevatorsPerColumn := []int{5, 5, 5, 5}
	b := newBattery(85, numberOfElevatorsPerColumn)
	
	w := []TestElevator {
		b.newTestElevator(0, 45, []int{1}),
		b.newTestElevator(1, 2, []int{45}),
		b.newTestElevator(2, 55, []int{61}),
		b.newTestElevator(3, 50, []int{1}),
		b.newTestElevator(4, 63, []int{46}),
	}
	
	b.initTest(w, numberOfElevatorsPerColumn)
	b.requestElevator(1, 64, "up")

	

	screen("THANK FOR TRUSTING ROCKET ELEVATOR | ELEVATE, SPEED AND STYLE !")
}