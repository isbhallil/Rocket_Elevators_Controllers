package Commercial_Controller;
import java.util.*;

public class Main {

    public static void main(String[] args)
    {
        System.out.println("Rocket Elevator Controller - Corporate");

        // init number of columns by providing the number of elevators in each columns
        ArrayList<Integer> numberOfELevatorsPerColumn = new ArrayList<Integer>(Arrays.asList(5, 5, 5, 5));
        (new Battery(85, numberOfELevatorsPerColumn))
                .initTestElevator(0, 1, new ArrayList<Integer>(Arrays.asList(24)))
                .initTestElevator(1, 23, new ArrayList<Integer>(Arrays.asList(28)))
                .initTestElevator(2, 33, new ArrayList<Integer>(Arrays.asList(1)))
                .initTestElevator(3, 40, new ArrayList<Integer>(Arrays.asList(24)))
                .initTestElevator(4, 42, new ArrayList<Integer>(Arrays.asList(1)))
                .requestElevator(1, "up", 36).assignElevator(36).operate();


        System.out.println("Thank you for trusting Rocket Elevators");
        new Scanner(System.in).nextLine();
    }
}


class Battery
{
    public ArrayList<Column> columns;
    public int numberOfFloors;

    public Battery(int aNumberOfFloors, ArrayList<Integer> aNumberOfElevatorsPerColumn) // contructor method
    {
        numberOfFloors = aNumberOfFloors;
        columns = initColumns(aNumberOfFloors, aNumberOfElevatorsPerColumn);
    }
    private ArrayList<Column> initColumns(int aNumberOfFloors, ArrayList<Integer> aNumberOfElevatorsPerColumn)
    {
        ArrayList<Column> columns = new ArrayList<Column>();

        int floor = 1;
        int columnId = 0;
        int range = (int)Math.floor((double)(aNumberOfFloors / aNumberOfElevatorsPerColumn.size()));
        do
        {
            int lowestFloorLevel;
            if (floor == 1)
            {
                lowestFloorLevel = 2;
            }
            else
            {
                lowestFloorLevel = floor + 1;
            }

            int highestFloorLevel;
            if ((floor + range) > aNumberOfFloors)
            {
                highestFloorLevel = aNumberOfFloors;
            }
            else
            {
                highestFloorLevel = floor + range;
            }

            columns.add(new Column(columnId, lowestFloorLevel, highestFloorLevel, aNumberOfElevatorsPerColumn.get(columnId)));

            floor += range;
            columnId++;
        } while (floor < aNumberOfFloors);

        return columns;
    }
    public Battery initTestElevator(int elevatorId, int currentFloor, ArrayList<Integer> floorsToVisit)
    {

        Column bestColumn = this.columns.stream()
                .filter(column -> ((column.lowest <= currentFloor && column.highest >= currentFloor) || (column.lowest <= floorsToVisit.get(0) && column.highest >= floorsToVisit.get(0))))
                .findFirst()
                .get();

        Elevator elevatorTest = bestColumn.elevators.get(elevatorId);

        elevatorTest.currentFloor = currentFloor;
        elevatorTest.initTasksList(floorsToVisit);

        return this;
    }

    public Elevator requestElevator(int floor, String direction, int target)
    {
         Column requestedColumn = this.columns.stream()
                .filter(column -> ((column.lowest <= floor && column.highest >= floor) || (column.lowest <= target && column.highest >= target)))
                .findFirst()
                .get();

        ElevatorRequest task = new ElevatorRequest(direction, floor, requestedColumn.id);

        return requestedColumn.getBestElevator(task).assignElevator(floor);
    }
}

class Column
{
    public int id;
    public int lowest;
    public int highest;
    public ArrayList<FloorButton> floorButtons;
    public ArrayList<Elevator> elevators;
    public Column(int aId, int aLowestFloor, int aHighestFloor, Object aNumberOfElevators)
    {
        id = aId;
        lowest = aLowestFloor;
        highest = aHighestFloor;
        floorButtons = initFloorButtons(aLowestFloor, aHighestFloor);
        elevators = initElevators((Integer)aNumberOfElevators, aLowestFloor, aHighestFloor);

        System.out.println("Column " + id + " has just been created !");
    }
    public Elevator getBestElevator(ElevatorRequest task)
    {
        ArrayList<Elevator> comingElevators = new ArrayList<>();
        ArrayList<Elevator> idleElevators = new ArrayList<>();
        ArrayList<Elevator> otherElevators = new ArrayList<>();

        elevators.forEach(elevator ->
        {
            if (elevator.tasksList.size() == 0)
            {
                idleElevators.add(elevator);
            }
            else if (elevator.isComming(task))
            {
                comingElevators.add(elevator);
            }
            else
            {
                otherElevators.add(elevator);
            }
        });

        if (comingElevators.size() > 0){
            Elevator elevator = getBestElevatorFrom(comingElevators, task);
            System.out.println();
            System.out.println("requestElevator with " + elevator.id);
            return elevator;
        } else if (idleElevators.size() > 0) {
            Elevator elevator = getBestElevatorFrom(idleElevators, task);
            System.out.println();
            System.out.println("requestElevator with " + elevator.id);
            return elevator;
        } else {
            Elevator elevator = getBestElevatorFrom(otherElevators, task);
            System.out.println();
            System.out.println("requestElevator with " + elevator.id);
            return elevator;
        }
    }

    private Elevator getBestElevatorFrom(ArrayList<Elevator> list, ElevatorRequest task)
    {
        Elevator bestElevator = list.get(0);
        var bestTravelSteps = bestElevator.getStepsToReach(task);

        for ( Elevator elevator: list) {
            int elevatorSteps = elevator.getStepsToReach(task);
            if (elevatorSteps < bestTravelSteps)
            {
                bestElevator = elevator;
                bestTravelSteps = elevatorSteps;
            }
        };

        return bestElevator;
    }
    private ArrayList<FloorButton> initFloorButtons(int aLowestFloor, int aHighestFloor)
    {
        ArrayList<FloorButton> list = new ArrayList<FloorButton>();
        list.add(new FloorButton(id, 1, "up"));
        list.add(new FloorButton(id, 1, "down"));

        int floor = aLowestFloor;
        while (floor <= aHighestFloor)
        {
            list.add(new FloorButton(id, floor, "up"));
            list.add(new FloorButton(id, floor, "down"));
            floor++;
        }

        return list;
    }
    private ArrayList<Elevator> initElevators(int aNumberOfElevators, int aStartingFloor, int aHighestFloor)
    {
        ArrayList<Elevator> list = new ArrayList<Elevator>();
        for (int i = 0; i < aNumberOfElevators; i++)
        {
            list.add(new Elevator(i, aStartingFloor, aHighestFloor));
        }

        return list;
    }
}

class Elevator
{
    public int id;
    public int currentFloor;
    public boolean isDoorOpen;
    public boolean isSafe;
    public int maxWeight;
    public ArrayList<ElevatorButton> buttonsList;
    public ArrayList<ElevatorTask> tasksList;

    public Elevator(int aId, int aStartingFloor, int aHighestFloor)
    {
        id = aId;
        currentFloor = aStartingFloor;
        isDoorOpen = false;
        isSafe = true;
        maxWeight = 3500;
        tasksList = new ArrayList<ElevatorTask>(Arrays.asList( ));
        buttonsList = initElevatorButtons(aStartingFloor, aHighestFloor);
        System.out.println("Elevator created " + id + " has " + tasksList.size());
    }

    public Elevator assignElevator(int floor)
    {
        System.out.println("request floor from elevator " + id + " to go to floor " + floor);
        addTask(new ElevatorTask(floor)).operate();

        return this;
    }

    public Elevator addTask(ElevatorTask task)
    {
        tasksList.add(task);

        String direction = this.getDirection();
        if (direction.equals("up"))
        {
                tasksList.sort(Comparator.comparingInt(ElevatorTask::getFloor));
        }
        else
        {
            tasksList.sort(Comparator.comparingInt(ElevatorTask::getFloor).reversed());
        }

        return this;
    }

    public boolean isComming(ElevatorRequest task)
    {
        boolean isComming = false;
        String direction = this.getDirection();

        if (direction.equals(task.direction))
        {
            if (direction.equals("up"))
            {
                for (ElevatorTask taskItem : this.tasksList)
                {
                    if (this.currentFloor <= taskItem.floor && currentFloor < taskItem.floor)
                    {
                        return true;
                    }
                }
            }

            else if (direction.equals("down"))
            {
                boolean again = true;
                for ( ElevatorTask taskItem : this.tasksList)
                {
                    if (this.currentFloor >= taskItem.floor && this.currentFloor > taskItem.floor)
                    {
                        return true ;
                    }
                };
            }
        }

        return false;
    }

    public int getStepsToReach(ElevatorRequest task)
    {
        String direction = this.getDirection();
        int taskFloor = task.floor;
        int stepsToCome = 0;
        int previousTask = this.currentFloor;

        if (this.tasksList.size() == 0)
        {
            stepsToCome += Math.abs(previousTask - taskFloor);
        }
        else
        {
            int index = 0;

            for ( ElevatorTask elevatorTask : this.tasksList)
            {
                int nextTaskFloor;
                if (tasksList.get(index) != null)
                {
                    nextTaskFloor = tasksList.get(index).floor;
                }
                else
                {
                    nextTaskFloor = taskFloor;
                }


                if (direction.equals("up") && taskFloor >= previousTask && taskFloor <= nextTaskFloor)
                {
                    stepsToCome += Math.abs(previousTask - taskFloor);
                }
                else if (direction.equals("down") && taskFloor <= previousTask && taskFloor >= nextTaskFloor)
                {
                    stepsToCome += Math.abs(previousTask - taskFloor);
                }
                else
                {
                    stepsToCome += Math.abs(previousTask - nextTaskFloor);
                }

                previousTask = nextTaskFloor;
                index++;
            };

            // stepsToCome += Math.Abs(previousTask - taskFloor);
        }

        System.out.println("Elevator " + id + " have next step as " + previousTask);
        System.out.println("Elevator " + id + " gonna do " + stepsToCome + " to reach " + task.floor);
        return stepsToCome;
    }

    public String getDirection()
    {
        boolean movingDown = false;
        boolean movingUp = false;
        String direction = "null";

        movingDown = tasksList.size() > 0 && currentFloor > tasksList.get(0).floor;
        movingUp = tasksList.size() > 0 && currentFloor < tasksList.get(0).floor;

        if (movingDown)
        {
            direction = "down";
        }
        else if (movingUp)
        {
            direction = "up";
        }
        else
        {
            direction = "idle";
        }

        return direction;
    }

    public Elevator operate()
    {
        while (tasksList.size() > 0)
        {
            if (this.tasksList.size() > 0 )
            {
                writeScreen("Elevator is at floor " + currentFloor);
                this.move();
                if (this.currentFloor == this.tasksList.get(0).getFloor()){
                    writeScreen("you arrived at floor" + this.currentFloor);
                    this.openDoor();
                    this.removeNextTask();
                    this.closeDoor();
                }
            }
        }
        return this;
    }

    public Elevator move()
    {
           if (this.currentFloor < this.tasksList.get(0).getFloor())
           {
               this.moveup();
           }
           else if ( this.currentFloor > this.tasksList.get(0).getFloor())
           {
               this.moveDown();
           }

           return this;
    }

    public void writeScreen(String message)
    {
        System.out.println();
        System.out.println("================================");
        System.out.println("=> " + message);
        System.out.println();
        System.out.println();
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
        tasksList.remove(0);
        writeScreen("Task removed !");
    }

    public void moveup()
    {
        writeScreen("Elevator " + id + " is moving up to floor " + (currentFloor + 1));

        currentFloor += 1;
        writeScreen("Elevator " + id + " is at floor " + currentFloor);
    }

    public void moveDown()
    {
        writeScreen("Elevator " + id + " is moving down to floor " + (currentFloor - 1));

        currentFloor -= 1;
        writeScreen("Elevator " + id + " is at floor " + currentFloor);
    }

    public ArrayList<ElevatorButton> initElevatorButtons(int aStartingFloor, int aHighestFloor)
    {
        ArrayList<ElevatorButton> list = new ArrayList<ElevatorButton>();

        int nextFloor = aStartingFloor;
        while (nextFloor <= aHighestFloor)
        {
            list.add(new ElevatorButton(nextFloor, id));
            nextFloor++;
        }

        return list;
    }

    public Elevator initTasksList(ArrayList<Integer> predefinedTasks)
    {
        ArrayList<ElevatorTask> list = new ArrayList<ElevatorTask>();

        int maxLoop = predefinedTasks.size();
        int index = 0;
        while (index < maxLoop)
        {
            int floor = predefinedTasks.get(index);
            ElevatorTask task = new ElevatorTask(floor);
            list.add(task);
            index++;
        }

        tasksList = list;

        return this;
    }

}

class ElevatorTask
{
    public int floor;
    public Date createdAt;

    public ElevatorTask(int aFloor)
    {
        floor = aFloor;
        createdAt = new Date();
    }

    public int getFloor()
    {
        return this.floor;
    }
}

class ElevatorButton
{

    public boolean isActive;
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

        System.out.println("assignElevator");
    }
}

class FloorButton
{
    public String direction;

    public boolean isActive;
    public int floor;
    public int columnId;
    public FloorButton(int aColumnId, int aFloor, String aDirection)
    {
        isActive = false;
        floor = aFloor;
        columnId = aColumnId;
        direction = aDirection;
    }

    public void requestElevator()
    {
        System.out.println("requestElevator");
    }

}

class ElevatorRequest {
    public String direction;
    public int floor;
    public int columnId;

    public ElevatorRequest(String aDirection, int aFloor, int aColumnId) {
        direction = aDirection;
        floor = aFloor;
        columnId = aColumnId;
    }
}
