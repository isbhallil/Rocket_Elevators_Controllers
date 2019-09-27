// // **************************************************************************

// //                           BATTERY

// // **************************************************************************

class Battery {
    constructor(numberOfFloors, numberOfColumns) {
        this.numberOfFloors = numberOfFloors ? numberOfFloors : 10
        this.columns = this.initColumns(numberOfFloors, numberOfColumns)
    }

    initColumns = (numberOfFloors, numberOfColumns) => {
        let range = Math.floor(numberOfFloors / numberOfColumns);
        let columns = []

        let floor = 1
        let maxFloor = numberOfFloors
        let columnId = 0
        for (; floor < maxFloor; floor += range) {
            let lowestFloorLevel = floor == 1 ? 2 : floor + 1;
            let highestFloorLevel = (floor + range) > numberOfFloors ? numberOfFloors : floor + range

            columns.push(new Column(columnId++, lowestFloorLevel, highestFloorLevel));
        }

        return columns
    }

    getColumnButton = (floor, direction) => {
        let response;

        battery.columns.forEach(columnItem => {
            columnItem.floorButtons.some(buttonItem => {
                if (buttonItem.floorNumber == floor && buttonItem.direction == direction) {

                    response = {
                        column: columnItem,
                        button: buttonItem
                    }

                }
            })
        })

        return response
    }
}









// **************************************************************************

//                              COLUMN

// **************************************************************************

let elevators = [
    [{
        currentFloor: 10,
        tasksList: []
    }, {
        currentFloor: 3,
        tasksList: [6]
    }]
]


class Column {
    constructor(id, lowestFloor, highestFloor) {
        this.id = id
        this.floorButtons = this.initButtons(lowestFloor, highestFloor, id)
        this.elevators = this.initElevators(elevators[id])
        console.log("Column: ", this.floorButtons, this.elevators);
    }

    initButtons = (lowestFloor, highestFloor, columnId) => {
        let buttons = []

        let floor = lowestFloor
        let maxFloor = highestFloor

        buttons.unshift(new FloorButton('up', 1, columnId))

        for (; floor <= maxFloor; floor++) {
            buttons.push(new FloorButton('up', floor, columnId))
            buttons.push(new FloorButton('down', floor, columnId));
        }

        return buttons
    }

    initElevators = (scene) => {
        let elevatorsList = []

        if (scene) scene.forEach((elevator, index) =>
            elevatorsList.push(new Elevator(elevator, index + 1))
        )

        return elevatorsList
    }

    requestElevator = (requestedFloor, direction) => {
        let elevatorRequest = new ElevatorRequest(requestedFloor, direction);
        let elevator = this.getBestElevator(elevatorRequest)
        
        console.log(requestedFloor)
        elevator.addTask(requestedFloor);
    }

    getBestElevator = (request) => {
        let idleElevators = []
        let commingElevators = []
        let othersElevators = []


        this.elevators.forEach(elevator => {
            if (elevator.tasksList == 0) idleElevators.push(elevator)
            else if (elevator.isComming(request)) commingElevators.push(elevator)
            else othersElevators.push(elevator)
        });


        let bestElevator;
        [commingElevators, idleElevators, othersElevators].some(list => {
             
            bestElevator = this.getFastestFrom(list, request)
            return bestElevator
        })
        return bestElevator
    }

    getFastestFrom = (elevatorsList, request) => {
        let bestElevator;
        let bestTravelSteps;

        console.log('getBestFrom: ', elevatorsList, request)
        elevatorsList.forEach(elevator => {
            let elevatorTravelSteps = this.getStepsToCome(elevator, request)

            if (bestTravelSteps == undefined || elevatorTravelSteps < bestTravelSteps) {
                bestElevator = elevator
                bestTravelSteps = elevatorTravelSteps
            }
        })
        return bestElevator
    }


    getStepsToCome = (elevator, request) => {
        let direction = elevator.getDirection()
        let floor = request.requestedFloor
        let idle = elevator.tasksList.length == 0
        let stepsToCome = 0
        let previousTask = elevator.currentFloor

         
        if (idle) {
            stepsToCome += Math.abs(previousTask - request.requestedFloor);
        } else {
            elevator.tasksList.some((task, index) => {
                let nextTask = elevator.tasksList[index] ? elevator.tasksList[index] : floor

                if (direction == 'up' && floor >= previousTask && floor <= nextTask) {
                    return stepsToCome += Math.abs(previousTask - floor)
                } else if (direction == 'down' && floor <= previousTask && floor >= nextTask) {
                    return stepsToCome += Math.abs(previousTask - floor)
                } else stepsToCome += Math.abs(previousTask - nextTask)

                previousTask = nextTask
            })
        }



        return stepsToCome
    }
}

class ElevatorRequest {
    constructor(floor, direction) {
        this.requestedFloor = floor
        this.direction = direction
        this.createdAt = Date.now()
    }
}










// **************************************************************************

//                           FLOOR_BUTTON

// **************************************************************************

class FloorButton {
    constructor(direction, number, columnId) {
        this.columnId = columnId
        this.floorNumber = number
        this.direction = direction
        this.columnId = columnId
        this.buttonsLight = false
    }

    requestElevator = () => {
        requestElevator(this.floorNumber, this.direction)
        this.toggleLight();
    }

    toggleLight = () => {
        this.buttonsLight = !this.buttonsLight
        console.log("FloorButton", this)
    }
}








// // **************************************************************************

// //    ELEVATOR

// // **************************************************************************

class Elevator {
    constructor(elevator, elevatorId) {
        this.id = elevatorId
        this.currentFloor = elevator.currentFloor ? elevator.currentFloor : 1
        this.doorState = elevator.doorState ? elevator.doorState : 'closed'
        this.isSafe = elevator.isSafe ? elevator.isSafe : 'true'
        this.maxWeight = elevator.maxWeight ? elevator.maxWeight : 3500
        this.buttonsList = elevator.buttonsList ? elevator.buttonsList : []
        this.tasksList = elevator.tasksList ? elevator.tasksList : []
    }


    getDirection = () => {
        let movingDown = (this.tasksList.length == 1 && this.currentFloor > this.tasksList[0]) ||
            this.tasksList[this.tasksList.length - 2] > this.tasksList[this.tasksList.length - 1]

        let movingUp = (this.tasksList.length == 1 && this.currentFloor < this.tasksList[0]) ||
            this.tasksList[this.tasksList.length - 2] < this.tasksList[this.tasksList.length - 1]


        if (this.tasksList.length == 0) return null
        else if (movingDown) {
            return 'down'
        } else if (movingUp) {
            return 'up'
        }
    }

    isComming = (request) => {
        let elevatorDirection = this.getDirection()
        let isComming = false

        if (elevatorDirection == request.direction) {
            if (elevatorDirection == "up") {
                isComming = this.tasksList.some(task => {
                    if (task <= request.requestedFloor) return true
                })
            } else if (elevatorDirection == "down") {
                isComming = this.tasksList.some(task => {
                    if (task >= request.requestedFloor) return true
                })
            }
        }

        return isComming
    }

    operate = () => {
        let again = true
        while (again) {
            if (this.tasksList[0] > this.currentFloor){
                 this.moveUp()
            }
            else if (this.tasksList[0] < this.currentFloor) {
                this.moveDown()
            }
            else {
                console.log('you are at your destination: ', this.currentFloor)
                this.openDoor()
                this.removeTask(this.tasksList[0])
                again = false

                setTimeout(()=>{
                   this.closeDoor()
                }, 5000)
            }   
        }
    }

    moveUp = () => {
        this.currentFloor = this.currentFloor + 1
        console.log('moving up to floor: ', this.currentFloor)
    }

    moveDown = () => {
        this.currentFloor = this.currentFloor - 1
        console.log('moving down to floor: ', this.currentFloor)
    }

    openDoor = () => {
        this.doorState == 'opened'
        this.removeTask()
        console.log('door is now open')
    }

    closeDoor = () => {
        this.doorState == "closed"
        console.log('door is closed')
        if (this.tasksList.length > 0) this.operate()
    }

    addTask = (task) => {    
        this.tasksList.push(task)
        console.log(task)

        this.operate()
    }

    removeTask = (task) => {
        this.tasksList.shift()
    }
}














const numberOfFloors = 10
const numberOfColumns = 1
const battery = new Battery(numberOfFloors, numberOfColumns);

console.log(`
**********************************************************************

    TO PUSH A BUTTON SIMPLY CALL requestElevator(floor, direction)
    direction == 'up' || 'down'

**********************************************************************
`)

const requestElevator = (floor, direction) => {
    let query = battery.getColumnButton(floor, direction)
    console.log(query)


    query ? query.column.requestElevator(floor, direction) : console.log('the request you trying to make is not possible because only God can')
}

const requestFloor = (elevator, requestedFloor) => {
    battery.columns.forEach(col => {
        col.elevators.some(elevatorItem => {
            if (elevatorItem.id == elevator) {
                elevatorItem.addTask(requestedFloor)
                console.log('Elevator Found, ', elevatorItem)
                return true
            }
        })
    })

}



























// const battery = {
//     numberOfFloors: 60,
//     numberOfColumns: 4,
// }

// const elevatorsList = [{
//     column: 0,
//     state: 'idle',
//     currentFloor: 2,
//     doorState: 'closed',
//     isSafe: true,
//     buttonsList: [],
//     tasksList: []
// }, {
//     column: 0,
//     state: 'moving',
//     currentFloor: 10,
//     doorState: 'closed',
//     isSafe: true,
//     buttonsList: [],
//     tasksList: []
// }]





// // **************************************************************************

// //    BATTERY

// // **************************************************************************


// class Battery {
//     constructor(numberOfFloors, numberOfColumns) {
//         this.columns = this.initColumns(numberOfFloors, numberOfColumns)

//         console.log(`[Battery] | constructor using numberOfFLoor: ${numberOfFloors} and numberOfColumns: ${numberOfColumns}`)
//     }

//     initColumns = (numberOfFLoors, numberOfColumns) => {
//         let range = Math.floor(numberOfFLoors / numberOfColumns);
//         let columns = []

//         let floor = 1
//         let maxFloor = numberOfFLoors - floor
//         for(; floor <= maxFloor; floor += range){
//             let column = {
//                 lowestFloorLevel: floor ,
//                 highestFloorLevel: floor + range                
//             }

//             columns.push(new Column(Math.random, column));
//             console.log(`[initColumns] | range: ${range}, floor: ${floor}`)
//         }

//         return columns
//     };


//     RequestElevator = (RequestedFloor, Direction, ) => {
//         let request = new Request(RequestedFloor, Direction);
//         console.log('floor requested');
//     }
// }





// const elevators = [

// ]





// // **************************************************************************

// //    COLUMN

// // **************************************************************************

// class Column {
//     constructor(id, column) {
//         this.id = id;
//         this.floorsList = this.initFloors(column.lowestFloorLevel, column.highestFloorLevel)
//         //this.elevatorsList = this.initElevators(column.elevatorsList)
//         console.log(`new Column ${this.id}`)
//     }

//     initFloors = (lowest, highest) => {
//         let floors = [];
//         floors.push(new Floor(0, this.id));

//         let i = lowest
//         while (i <= highest) {
//             floors.push(new Floor(i, this.id))
//             i++
//         }


//         return floors
//     }

//     initElevators = (elevatorsList) => {
//         let elevators = [];

//         let i = 1
//         const iMax = elevatorsList.length
//         for (; i < iMax; i++) {
//             elevators.push(new Elevator(elevatorsList[i]));
//         }

//         return elevators
//     }

// }










// // **************************************************************************

// //    FLOOR

// // **************************************************************************

// class Floor {
//     constructor(floorNumber, columnId) {
//         this.floorNumber = floorNumber;
//         this.floorButtons = this.initFloorButtons(this.floorNumber, columnId)
//         this.columnId = columnId
//         this.requestsList = []

//         console.log(`new Floor ${floorNumber} || on column ${this.columnId}`)
//     }

//     initFloorButtons = (floorNumber, columnId) => {
//         const buttons = []

//         if (floorNumber != battery.numberOfFloors) buttons.push(new FloorButton('up', floorNumber, columnId))
//         if (this.floorNumber != 0) buttons.push(new FloorButton('down', floorNumber, columnId));

//         return buttons
//     }
// }










// // **************************************************************************

// //    FLOOR BUTTON

// // **************************************************************************

// class FloorButton {
//     constructor(direction, number, columnId) {
//         this.floorNumber = number
//         this.direction = direction
//         this.columnId = columnId
//         this.buttonsLight = false

//         console.log(`new FloorButton | direction: ${this.direction}, floorNumber: ${this.floorNumber}, column: ${this.columnId}`)
//     }

//     requestElevator = () => {
//         // let request = new ElevatorRequest(this.floorNumber, this.direction, this.columnId)

//         // this.toggleLight();
//         // this.request = request
//     }

//     toggleLight = () => {
//         this.buttonsLight = !this.buttonsLight
//     }
// }


// class Request {
//     constructor(requestedFloor, direction, columnId) {
//         this.requestedFloor = requestElevator
//         this.direction = direction
//         this.createdAt = Date.now()
//         this.arrivedAt = null
//         this.columnId = columnId
//     }
// }


























// const b1 = new Battery(battery.numberOfFloors, battery.numberOfColumns)