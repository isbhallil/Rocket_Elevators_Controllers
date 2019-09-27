import math;


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
        response = {}

        self.columns

elevators = [
    [{
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
        self.lowestFloor = lowestFloor
        self.highestFloor = highestFloor

    def initButtons(self, lowestFloor, highestFloor, columnId):
        buttons = []

        floor = lowestFloor
        maxFloor = highestFloor

        buttons.insert(0, FloorButton('up', 1, self.id))

        for i in range(floor, highestFloor + 1):
            buttons.append(FloorButton('up', i, self.id));
            buttons.append(FloorButton('down', i, self.id))

        return buttons
    
    # def getColumnButton(floor, direction):
    #     response = {}

        


# // **************************************************************************

# //      FLOOR_BUTTON

# // **************************************************************************
class FloorButton:
    def __init__(self, direction, number, columnId):
        self.columnId = columnId
        self.FloorNumber = number
        self.direction = direction
        self.light = False

    def requestElevator(self):
        # requestElevator(self.FloorNumber, self.direction)
        self.toggleLight()
        print('light')

    def toggleLight(self):
        self.light = not self.light
        print('light: ' + str(self.light))



column = Column(1, 2, 10)
# button = FloorButton('up', 10, 1);
# # button.requestElevator()       
        

# battery = Battery(10,1)

# print(battery.numberOfFloors);


# def requestElevator(floor, direction):
#     query = battery.getColumnButton(floor, direction)

#     if query:
#         query.column.requestElevator(floor, direction)
#     else:
#         print('the request you trying to make is not possible because only God can')



# def resquestFloor(elevator, requestedFloor):
#     for col in battery.columns:
#         for elevatorItem in col.elevators:
#             if (elevator.id == elevator):
#                 elevatorItem.addTask(requestedFloor)
#                 print('Elevator found,' +  str(elevatorItem))
#                 break