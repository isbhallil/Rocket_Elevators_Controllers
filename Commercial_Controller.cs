using System;
using System.Collections.Generic;
using System.Threading;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace Commercial_Controller
{

    // verfier la possibiliter que le progrome soit bloquer dnas une boucle;

    class Program
    {
        static void Main(string[] args)
        {
            Console.WriteLine("Rocket Elevator Controller - Corporate");

            // init number of columns by providing the number of elevators in each columns



            List<int> numberOfELevatorsPerColumn = new List<int>() { 5, 5, 5, 5 };
            new Battery(85, numberOfELevatorsPerColumn)
                .initTestElevator(0, 1, new List<int> { 24 })
                .initTestElevator(1, 23, new List<int> { 28 })
                .initTestElevator(2, 33, new List<int> { 1 })
                .initTestElevator(3, 40, new List<int> { 24 })
                .initTestElevator(4, 42, new List<int> { 1 })
                .requestElevator(1, "up", 36)
                .assignElevator(36)
                .operate();


            Console.WriteLine("Thank you for trusting Rocket Elevators");
            Console.ReadLine();
        }
    }
    class Battery
    {
        public List<Column> columns;
        public int numberOfFloors;

        public Battery(int aNumberOfFloors, List<int> aNumberOfElevatorsPerColumn) // contructor method
        {
            numberOfFloors = aNumberOfFloors;
            columns = initColumns(aNumberOfFloors, aNumberOfElevatorsPerColumn);
        }
        private List<Column> initColumns(int aNumberOfFloors, List<int> aNumberOfElevatorsPerColumn)
        {
            List<Column> columns = new List<Column>();

            int floor = 1;
            int columnId = 0;
            int range = Convert.ToInt32(Math.Floor(Convert.ToDouble(aNumberOfFloors / aNumberOfElevatorsPerColumn.Count)));
            do
            {
                int lowestFloorLevel;
                if (floor == 1) { lowestFloorLevel = 2; }
                else { lowestFloorLevel = floor + 1; }

                int highestFloorLevel;
                if ((floor + range) > aNumberOfFloors) { highestFloorLevel = aNumberOfFloors; }
                else { highestFloorLevel = floor + range; }

                columns.Add(new Column(columnId, lowestFloorLevel, highestFloorLevel, aNumberOfElevatorsPerColumn[columnId]));

                floor += range;
                columnId++;
            } while (floor < aNumberOfFloors);

            return columns;
        }
        public Battery initTestElevator(int aElevatorId, int aCurrentFloor, List<int> aFloorsToVisit)
        {
            int highestFloor;
            string direction;

            if (aFloorsToVisit.Count == 0 || aCurrentFloor > aFloorsToVisit[0]) highestFloor = aCurrentFloor;
            else highestFloor = aFloorsToVisit[0];

            Elevator elevatorTest = columns
                .Find(column => column.lowest <= highestFloor && column.highest >= highestFloor)   // Exists(button => button.floor == heighestFloor))
                .elevators.Find(elevator => elevator.id == aElevatorId);

            elevatorTest.currentFloor = aCurrentFloor;
            elevatorTest.initTasksList(aFloorsToVisit);

            return this;
        }
        public Elevator requestElevator(int floor, string direction, int target)
        {
            Column requestedColumn = columns.Find(column =>
            (column.lowest <= floor && column.highest >= floor)
            || (column.lowest <= target && column.highest >= target));


            ElevatorRequest task = new ElevatorRequest(direction, floor, requestedColumn.id);

            return requestedColumn
                .getBestElevator(task)
                .assignElevator(floor);
        }
    }
    class Column
    {
        public int id;
        public int lowest;
        public int highest;
        public List<FloorButton> floorButtons;
        public List<Elevator> elevators;
        public Column(int aId, int aLowestFloor, int aHighestFloor, object aNumberOfElevators)
        {
            id = aId;
            lowest = aLowestFloor;
            highest = aHighestFloor;
            floorButtons = initFloorButtons(aLowestFloor, aHighestFloor);
            elevators = initElevators(Convert.ToInt32(aNumberOfElevators), aLowestFloor, aHighestFloor);

            Console.WriteLine("Column " + id + " has just been created !");
        }
        public Elevator getBestElevator(ElevatorRequest task)
        {
            List<Elevator> commingElevators = new List<Elevator>();
            List<Elevator> idleElevators = new List<Elevator>();
            List<Elevator> otherElevators = new List<Elevator>();

            elevators.ForEach(elevator =>
            {
                if (elevator.tasksList.Count == 0)
                {
                    idleElevators.Add(elevator);
                }
                else if (elevator.isComming(task))
                {
                    commingElevators.Add(elevator);
                }
                else
                {
                    otherElevators.Add(elevator);
                }
            });

            Elevator bestElevator = null;
            bool again = true;
            new List<List<Elevator>> { commingElevators, idleElevators, otherElevators }
            .ForEach(list =>
            {
                if (list.Count > 0 && again) bestElevator = getBestElevatorFrom(list, task);
                if (bestElevator != null) again = false;
            });

            Console.WriteLine();
            Console.WriteLine("requestElevator with " + bestElevator.id);

            return bestElevator;
        }
        private Elevator getBestElevatorFrom(List<Elevator> list, ElevatorRequest task)
        {
            Elevator bestElevator = list[0];
            int bestTravelSteps = bestElevator.getStepsToReach(task);

            list.ForEach(elevator =>
            {
                int elevatorSteps = elevator.getStepsToReach(task);

                if (elevatorSteps < bestTravelSteps)
                {
                    bestElevator = elevator;
                    bestTravelSteps = elevatorSteps;
                }

            });

            return bestElevator;
        }
        private List<FloorButton> initFloorButtons(int aLowestFloor, int aHighestFloor)
        {
            List<FloorButton> list = new List<FloorButton>();
            list.Add(new FloorButton(id, 1, "up"));
            list.Add(new FloorButton(id, 1, "down"));

            int floor = aLowestFloor;
            while (floor <= aHighestFloor)
            {
                list.Add(new FloorButton(id, floor, "up"));
                list.Add(new FloorButton(id, floor, "down"));
                floor++;
            }

            return list;
        }
        private List<Elevator> initElevators(int aNumberOfElevators, int aStartingFloor, int aHighestFloor)
        {
            List<Elevator> list = new List<Elevator>();
            for (int i = 0; i < aNumberOfElevators; i++)
            {
                list.Add(new Elevator(i, aStartingFloor, aHighestFloor));
            }

            return list;
        }
    }
    class Elevator
    {
        public int id;
        public int currentFloor;
        public bool isDoorOpen;
        public bool isSafe;
        public int maxWeight;
        public List<ElevatorButton> buttonsList;
        public List<ElevatorTask> tasksList;

        public Elevator(int aId, int aStartingFloor, int aHighestFloor)
        {
            id = aId;
            currentFloor = aStartingFloor;
            isDoorOpen = false;
            isSafe = true;
            maxWeight = 3500;
            tasksList = new List<ElevatorTask> { };
            buttonsList = initElevatorButtons(aStartingFloor, aHighestFloor);
            Console.WriteLine("Elevator created " + id + " has " + tasksList.Count);
        }

        public Elevator assignElevator(int floor)
        {
            Console.WriteLine("request floor from elevator " + id + " to go to floor " + floor);
            addTask(new ElevatorTask(floor)).operate();

            return this;
        }
        public Elevator addTask(ElevatorTask task)
        {
            tasksList.Add(task);

            string direction = getDirection();
            if (direction == "up") tasksList.Sort((a, b) => a.floor.CompareTo(b.floor));
            else tasksList.Sort((a, b) => -1 * a.floor.CompareTo(b.floor));

            return this;
        }
        public bool isComming(ElevatorRequest task)
        {
            bool isComming = false;
            string direction = getDirection();

            if (direction == task.direction)
            {
                if (direction == "up")
                {
                    bool again = true;
                    tasksList.ForEach(taskItem =>
                    {
                        if (currentFloor <= task.floor && currentFloor < task.floor && again == true)
                        {
                            isComming = true;
                            again = false;
                        }
                    });
                }

                else if (direction == "down")
                {
                    bool again = true;
                    tasksList.ForEach(taskItem =>
                    {
                        if (currentFloor >= task.floor && currentFloor > task.floor && again == true)
                        {
                            isComming = true;
                            again = false;
                        }
                    });
                }
            }


            return isComming;
        }
        public int getStepsToReach(ElevatorRequest task)
        {
            string direction = getDirection();
            int taskFloor = task.floor;
            int stepsToCome = 0;
            int previousTask = currentFloor;

            if (tasksList.Count == 0)
            {
                stepsToCome += Math.Abs(previousTask - taskFloor);
            }
            else
            {
                int index = 0;

                tasksList.ForEach(elevatorTask =>
                {
                    int nextTaskFloor;
                    if (tasksList[index] != null) { nextTaskFloor = tasksList[index].floor; }
                    else { nextTaskFloor = taskFloor; }


                    if (direction == "up" && taskFloor >= previousTask && taskFloor <= nextTaskFloor)
                    {
                        stepsToCome += Math.Abs(previousTask - taskFloor);
                    }
                    else if (direction == "down" && taskFloor <= previousTask && taskFloor >= nextTaskFloor)
                    {
                        stepsToCome += Math.Abs(previousTask - taskFloor);
                    }
                    else
                    {
                        stepsToCome += Math.Abs(previousTask - nextTaskFloor);
                    }

                    previousTask = nextTaskFloor;
                    index++;
                });

                // stepsToCome += Math.Abs(previousTask - taskFloor);
            }

            Console.WriteLine("Elevator " + id + " have next step as " + previousTask);
            Console.WriteLine("Elevator " + id + " gonna do " + stepsToCome + " to reach " + task.floor);
            return stepsToCome + abs( previousTask - task.floor);
        }
        public string getDirection()
        {
            bool movingDown = false;
            bool movingUp = false;
            string direction = "null";

            movingDown = tasksList.Count > 0 && currentFloor > tasksList[0].floor; //|| tasksList[tasksList.Count - 2].floor < tasksList[tasksList.Count - 1].floor;

            //if ()


            movingUp = tasksList.Count > 0 && currentFloor < tasksList[0].floor; // || tasksList[tasksList.Count - 2].floor < tasksList[tasksList.Count - 1].floor;


            if (movingDown) { direction = "down"; }
            else if (movingUp) { direction = "up"; }
            else { direction = "idle"; }

            return direction;
        }
        public Elevator operate()
        {
            while (currentFloor != tasksList[0].floor)
            {
                ElevatorTask nextTask = tasksList[0];
                //Thread.Sleep(500);

                Console.WriteLine("Elevator is at floor " + currentFloor);

                if (nextTask.floor > currentFloor)
                {
                    moveup();
                }
                else if (nextTask.floor < currentFloor)
                {
                    moveDown();
                }
                else
                {
                    writeScreen("You arrived at floor " + currentFloor);
                    openDoor();
                    removeNextTask();
                    closeDoor();

                    break;
                }
            };

            if (tasksList.Count != 0) this.operate();

            return this;
        }
        public void writeScreen(string message)
        {
            Console.WriteLine();
            Console.WriteLine("================================");
            Console.WriteLine("=> " + message);
            Console.WriteLine();
            Console.WriteLine();
        }
        public void openDoor()
        {
            writeScreen("Door is openeing...");

            //Thread.Sleep(500);
            isDoorOpen = true;
            writeScreen("Door is now open !");
        }
        public void closeDoor()
        {
            writeScreen("Door is closing...");

            //Thread.Sleep(500);
            isDoorOpen = false;
            writeScreen("Door is now closed !");


        }
        public void removeNextTask()
        {
            tasksList.RemoveAt(0);
            writeScreen("Task removed !");
        }
        public void moveup()
        {
            writeScreen("Elevator " + id + " is moving up to floor " + (currentFloor + 1));

            //Thread.Sleep(500);
            currentFloor += 1;
            writeScreen("Elevator " + id + " is at floor " + currentFloor);
        }
        public void moveDown()
        {
            writeScreen("Elevator " + id + " is moving down to floor " + (currentFloor - 1));

            //Thread.Sleep(500);
            currentFloor -= 1;
            writeScreen("Elevator " + id + " is at floor " + currentFloor);
        }
        public List<ElevatorButton> initElevatorButtons(int aStartingFloor, int aHighestFloor)
        {
            List<ElevatorButton> list = new List<ElevatorButton>();

            int nextFloor = aStartingFloor;
            while (nextFloor <= aHighestFloor)
            {
                list.Add(new ElevatorButton(nextFloor, id));
                nextFloor++;
            }

            return list;
        }
        public Elevator initTasksList(List<int> aPredefinedTasks)
        {
            List<ElevatorTask> list = new List<ElevatorTask>();

            int maxLoop = aPredefinedTasks.Count;
            int index = 0;
            while (index < maxLoop)
            {
                int floor = aPredefinedTasks[index];
                ElevatorTask task = new ElevatorTask(floor);
                list.Add(task);
                index++;
            }

            tasksList = list;

            return this;
        }
    }



    class ElevatorTask
    {
        public int floor;
        public DateTime createdAt;
        public ElevatorTask(int aFloor)
        {
            floor = aFloor;
            createdAt = new DateTime();
        }
    }

    class ElevatorButton
    {

        public bool isActive;
        public int floor;
        public int columnId;
        public ElevatorButton(int aFloor, int aColumnId)
        {
            isActive = false;
            floor = aFloor;
            columnId = aColumnId;
        }
        public void assignElevator(int floor)
        {

            Console.WriteLine("assignElevator");
        }
    }
    class FloorButton
    {
        public string direction;

        public bool isActive;
        public int floor;
        public int columnId;
        public FloorButton(int aColumnId, int aFloor, string aDirection)
        {
            isActive = false;
            floor = aFloor;
            columnId = aColumnId;
            direction = aDirection;
        }

        public void requestElevator()
        {
            Console.WriteLine("requestElevator");
        }

    }

    class ElevatorRequest
    {
        public string direction;
        public int floor;
        public int columnId;

        public ElevatorRequest(string aDirection, int aFloor, int aColumnId)
        {
            direction = aDirection;
            floor = aFloor;
            columnId = aColumnId;
        }
    }

}