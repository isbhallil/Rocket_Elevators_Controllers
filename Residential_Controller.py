import math;
import time


# // **************************************************************************

# //                           BATTERY

# // **************************************************************************

numberOfFloors = 10

class Battery:
    def __init__(self, numberOfFloors, numberOfColumns):
        self.numberOfFloors = numberOfFloors
        self.columns = self.initColumns(numberOfFloors, numberOfColumns)
        print(str(self.columns));

    def initColumns(self, numberOfFloors, numberOfColumns):
        floorRange = math.floor(numberOfFloors / numberOfColumns)
        columns = []
        floor = 1
        maxFloor = numberOfFloors
        columnId = 0

        print(str(floorRange))

        while floor < maxFloor:
            if floor == 1:
                lowestFloor = 2
            else:
                lowestFloor = floor + 1
            

            if floor + floorRange > numberOfFloors:
                highestFloor = numberOfFloors
            else:
                highestFloor = floor + floorRange

            floor = floor + floorRange
            columnId = columnId + 1
            
            columns.append(Column(columnId, lowestFloor, highestFloor))

        return columns

    def getColumnButton(self, floor, direction):
        response = None

        for columnItem in self.columns:
            for floorButton in columnItem.floorButtons:
                print('floorButton ' + str(floorButton.floorNumber) + "  " +str(floorButton.direction))
                if floorButton.floorNumber == floor and floorButton.direction ==  direction:
                    response = {"column": columnItem, "button": floorButton}
                    
        
        print(str(response))
        return response




# // **************************************************************************

# //                           COLUMN

# // **************************************************************************

scene = [
    [{
        "currentFloor": 10,
        "tasksList": []
    }, {
        "currentFloor": 3,
        "tasksList": [6]
    }], [{
        "currentFloor": 10,
        "tasksList": []
    }, {
        "currentFloor": 3,
        "tasksList": [6]
    }]
]   

class Column:
    def __init__(self, columnId, lowestFloor, highestFloor):
        self.id = columnId
        self.floorButtons = self.initButtons(lowestFloor, highestFloor, columnId)
        self.elevators = self.initElevators(scene, columnId)
        self.lowestFloor = lowestFloor
        self.highestFloor = highestFloor
        print('new Column')

    def initButtons(self, lowestFloor, highestFloor, columnId):
        buttons = []

        floor = lowestFloor
        maxFloor = highestFloor

        buttons.insert(0, FloorButton('up', 1, self.id))

        for i in range(floor, highestFloor + 1):
            buttons.append(FloorButton('up', i, self.id));
            buttons.append(FloorButton('down', i, self.id))

        return buttons

    def initElevators(self, scene, columnId):
        elevatorsList = []
        index = 0
        currentScene = scene[columnId]

        for elevator in currentScene:
            elevatorsList.append(Elevator(elevator['currentFloor'], elevator['tasksList'], index))
            index = index + 1
    
        print(elevatorsList)
        return elevatorsList
    
    def requestElevator(self, floor, direction):
        print('Elevator request at : ' + str(floor) + " to go " + str(direction))

        request = ElevatorRequest(floor, direction)
        elevator = self.getBestElevator(request)

        time.sleep(1)
        print("Rocket Elevator hove found you an elevator | id: " + str(elevator.id) + "is comming from floor: " + str(elevator.currentFloor))
        elevator.addTask(floor)
        time.sleep(1)
        elevator.operate()
        
    def getBestElevator(self, request):
        idleElevators = []
        commingElevators = []
        othersElevators = []

        for elevator in self.elevators:
            if len(elevator.tasksList) == 0:
                idleElevators.append(elevator)
                print('elevator appended in idle')
            elif elevator.isComming(request):
                commingElevators.append(elevator)
                print('appende to comming')
            else: 
                othersElevators.append(elevator)
                print('appended to others')

        bestElevator = None
        
        if len(commingElevators) > 0:
            bestElevator = self.getBestFrom(commingElevators, request)
        elif len(idleElevators) > 0:
            bestElevator = self.getBestFrom(idleElevators, request)
        else:
            bestElevator = self.getBestFrom(othersElevators, request)

        return bestElevator

    def getBestFrom(self, list, request):
        bestElevator = None
        bestTravelSteps = None

        for elevator in list:
            elevatorTravelSteps = self.getStepsToCome(elevator, request)

            if bestTravelSteps == None or elevatorTravelSteps < bestTravelSteps:
                bestElevator = elevator
                bestTravelSteps = elevatorTravelSteps

        return bestElevator

    def getStepsToCome(self, elevator, request):
        print('getStepsToCome with : ' + str(elevator))
        direction = elevator.getDirection()
        floor = request.requestedFloor
        stepsToCome = 0
        previousTask = elevator.currentFloor

        if len(elevator.tasksList) == 0:
            stepsToCome = stepsToCome + abs(previousTask - floor)
        else:
            index = 0
            for task in elevator.tasksList:
                nextTask = None

                if elevator.tasksList[index]:
                    nextTask = elevator.tasksList[index]
                else:
                    nextTask = floor

                if direction == "up" and floor >= previousTask and floor <= nextTask:
                    stepsToCome = stepsToCome + abs(previousTask - floor);
                elif direction == "down" and floor <= previousTask and floor >= nextTask:
                    stepsToCome == stepsToCome + abs(previousTask - nextTask)
                else:
                    stepsToCome = stepsToCome + abs(previousTask - nextTask)

                previousTask = nextTask
        
        return stepsToCome
                    







class ElevatorRequest:
    def __init__(self, requestedFloor, direction):
        self.requestedFloor = requestedFloor
        self.direction = direction
        



# // **************************************************************************

# //      ELEVATOR

# // **************************************************************************

class Elevator:
    def __init__(self, currentFloor, tasksList, elevatorId):
        self.id = elevatorId
        self.currentFloor = currentFloor
        self.doorState = 'closed'
        self.isSafe = True
        self.maxWheit = 3500
        self.buttonsList = []
        self.tasksList = tasksList
        print("new Elevator: "+  str(self.id))

    def isComming(self, request):
        print('is comming')
        elevatorDirection = self.getDirection()
        isElevatorComming = False

        if elevatorDirection == request.direction:
            if elevatorDirection == "up":
                for task in self.tasksList:
                    if self.currentFloor <= request.requestedFloor and self.currentFloor < task:
                        isElevatorComming = True
                        break
            elif elevatorDirection == "down":
                for task in self.tasksList:
                    if self.currentFloor >= request.requestedFloor and self.currentFloor > task:
                        isElevatorComming = True
                        break
        
        print("isElevatorComming " + str(isElevatorComming))
        return isElevatorComming

    def getDirection(self):
        movingDown = None
        movingUp = None

        if len(self.tasksList) == 1 and self.currentFloor > self.tasksList[0]:
            movingDown = True
        if len(self.tasksList) > 2:
            if self.tasksList[len(self.tasksList) - 2] > self.tasksList[len(self.tasksList) - 1]:
                movingDown = True

        if len(self.tasksList) == 1 and self.currentFloor < self.tasksList[0]:
            movingUp = True
        if len(self.tasksList) > 2:
            if self.tasksList[len(self.tasksList) - 2] < self.tasksList[len(self.tasksList) - 1]:
                movingUp = True


        if len(self.tasksList) == 0:
            print('none')
            return None
        elif movingDown:
            print('down')
            return 'down'

        else:
            print('up')
            return 'up'

    def operate(self):
        print('operate')

        again = True
        if len(self.tasksList) > 0:
            time.sleep(1)
            while again:
                time.sleep(1)
                if self.tasksList[0] > self.currentFloor:
                    self.moveUp()
                    print('elevator is moving and at floor : ' + str(self.currentFloor))
                elif self.tasksList[0] < self.currentFloor:
                    self.moveDown()
                    print('elevator is moving and at floor requested: ' + str(self.currentFloor))
                else: 
                    self.openDoor()
                    self.removeTask()
                    # time(1)
                    self.closeDoor()
                    again = False

    def moveDown(self):
        self.currentFloor = self.currentFloor - 1
        print('move down')

    def moveUp(self):
        self.currentFloor = self.currentFloor + 1
        print('move up')

    def openDoor(self):
        print('elevator is at floor requested: ' + str(self.currentFloor))
        time.sleep(1)
        self.doorState = "opened"
        time.sleep(1)
        print('open doors')

    def closeDoor(self):
        self.doorState = "closed"
        print("waiting 5 seconds befor closing")
        time.sleep(5)
        print('close doors')

    def addTask(self, task):
        self.tasksList.append(task)
        print('add Task')
        print('YOUR ELEVATOR IS MOVING TO THE REQUESTED FLOOR ' +  str(task))

    def removeTask(self):
        del self.tasksList[0]
        print('remove task')






# // **************************************************************************

# //      FLOOR_BUTTON

# // **************************************************************************
class FloorButton:
    def __init__(self, direction, number, columnId):
        self.columnId = columnId
        self.floorNumber = number
        self.direction = direction
        self.light = False

    def requestElevator(self):
        # requestElevator(self.FloorNumber, self.direction)
        self.toggleLight()
        print('light')

    def toggleLight(self):
        self.light = not self.light
        print('light: ' + str(self.light))






battery = Battery(10,1)

# battery.columns[0].elevators[1].operate()
def requestElevator(floor, direction):
    print('requestElevator at floor: ' + str(floor) + " and direction: " + str(direction))

    query = battery.getColumnButton(floor, direction)

    if query == None:
        print("Your request cannot be handled, maybe the floor exists in your dreams !")
    else:
        query['column'].requestElevator(floor, direction)

    # if query:
    #     query.column.requestElevator(floor, direction)
    # else:
    #     print('the request you trying to make is not possible because only God can')



def requestFloor(elevatorId, requestedFloor):
    for col in battery.columns:
        print('col: ' + str(col))
        for elevatorItem in col.elevators:
            print( elevatorItem.id)
            if elevatorItem.id == elevatorId:
                print('id: ')
                elevatorItem.addTask(requestedFloor)
                elevatorItem.operate()
                print('Elevator found,' +  str(elevatorItem))
                print('YOUR ELEVATOR IS MOVING TO THE REQUESTED FLOOR ' +  str(requestedFloor))
            #     break

requestFloor(1, 7);