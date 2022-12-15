package main

import (
	"C"
	"math/rand"
)
import "fmt"

// SingleLaneSimulation update the single Road
// Input: a Road object
// Output: a Road Object
func SingleLaneSimulation(currentRoad Road, kindPossiblity float64, gen int) Road {
	var prevCarIndex int

	// make a new road
	newRoad := make(Road, roadLength)
	carCnt := 0

	//produce cars at the beginning the the road
	if gen%2 == 0 {
		Produce(&currentRoad, kindPossiblity)
	}

	for i := roadLength - 1; i >= 0; i-- {
		currentCar := currentRoad[i]
		kind := currentCar.kind

		if kind == 0 {
			continue
		}
		if kind == 3 {
			if currentCar.backlight >= 400 {
				currentCar.kind = 0
				currentCar.backlight = 0
				continue
			} else {
				currentCar.backlight += 1
				continue
			}
		}

		prevCarIndex = GetPrevCar(currentRoad, i)

		if kind == 1 {
			SingleLaneNSDVupdate(&currentRoad, &newRoad, &carCnt, i, prevCarIndex)
		} else if kind == 2 {
			SingleLaneSDVupdate(&currentRoad, &newRoad, &carCnt, i, prevCarIndex)
		}
	}

	//fmt.Println(carCnt)
	return newRoad
}

func Randomdeceleraion(p float64, speed int) (int, int) {
	var newSpeed int
	var newLight int

	thresToDecel := rand.Float64()
	if thresToDecel <= p {
		newSpeed = speed - 1
		newLight = 1
	} else {
		newSpeed = speed
		newLight = 0
	}

	return newSpeed, newLight

}

func SingleLaneNSDVupdate(currentRoad, newRoad *Road, carCnt *int, currentIndex, prevCarIndex int) {
	var prevCar Car
	var speed int
	var newSpeed int
	var newLight int

	var deltaD int

	deltaD = prevCarIndex - currentIndex
	speed = (*currentRoad)[currentIndex].speed

	if deltaD >= safeSpaceMax[speed] {
		newSpeed = speed + 1
		newLight = 0
	} else if deltaD < safeSpaceMax[speed] && deltaD >= safeSpaceMin[speed] {
		prevCar = (*currentRoad)[currentIndex]
		if prevCar.backlight == 0 {
			newSpeed, newLight = Randomdeceleraion(p2, speed)
		} else {
			newSpeed, newLight = Randomdeceleraion(p1, speed)
		}
	} else if deltaD < safeSpaceMin[speed] {
		prevCar = (*currentRoad)[currentIndex]
		newSpeed = min(speed-1, deltaD-speed-1)
		newLight = 1
	}

	if newSpeed < 0 {
		newSpeed = 0
	} else if newSpeed > maxSpeed {
		newSpeed = maxSpeed
	}

	newIndex := currentIndex + speed
	if newIndex >= roadLength {
		(*carCnt)++
	} else if newIndex < roadLength && (*newRoad)[newIndex].kind != 0 {
		fmt.Println("NSDV crashes something.", newIndex, (*newRoad)[newIndex].kind)
		// panic("NSDV crashes something.")
	} else {
		(*newRoad)[newIndex].speed = newSpeed
		(*newRoad)[newIndex].kind = 1
		(*newRoad)[newIndex].backlight = newLight
	}
}

func SingleLaneSDVupdate(currentRoad, newRoad *Road, carCnt *int, currentIndex, prevCarIndex int) {
	var prevCar Car
	var speed int
	var newSpeed int
	var newLight int

	var deltaD int

	deltaD = prevCarIndex - currentIndex
	speed = (*currentRoad)[currentIndex].speed

	if deltaD >= roadLength {
		newSpeed = speed + 1
		newLight = 0
	} else {
		prevCar = (*currentRoad)[prevCarIndex]
		if prevCar.kind == 1 {
			if deltaD >= safeSpaceMax[speed] {
				newSpeed = speed + 1
				newLight = 0
			} else if deltaD < safeSpaceMax[speed] && deltaD >= safeSpaceMin[speed] {
				if prevCar.backlight == 0 {
					newSpeed = speed + 1
					newLight = 0
				} else {
					newSpeed = speed
					newLight = 1
				}

			} else if deltaD < safeSpaceMin[speed] {
				newSpeed = min(speed-1, deltaD-1)
				newLight = 1
			}
		} else if prevCar.kind == 2 || prevCar.kind == 3 {
			SDVminDis := GetSDVmindis(currentIndex, prevCarIndex, *currentRoad)
			if deltaD >= safeSpaceMax[speed] {
				newSpeed = speed + 1
				newLight = 0
			} else if deltaD > SDVminDis && deltaD < safeSpaceMax[speed] {
				if CheckTrain(*currentRoad, prevCarIndex) >= 5 {
					if prevCar.speed >= speed {
						newSpeed = speed
						newLight = 0
					} else {
						newSpeed = prevCar.speed
						newLight = 1
					}
				} else {
					if prevCar.speed >= speed {
						newSpeed = speed + 1
						newLight = 0
					} else {
						newSpeed = prevCar.speed + 1
						newLight = 1
					}
				}

			} else {
				if CheckTrain(*currentRoad, prevCarIndex) >= 5 {
					newSpeed = prevCar.speed - 1
					newLight = 1
				} else {
					newSpeed = prevCar.speed
					newLight = prevCar.backlight
				}

			}
		}
	}

	if newSpeed < 0 {
		newSpeed = 0
	} else if newSpeed > maxSpeed {
		newSpeed = maxSpeed
	}
	newIndex := currentIndex + newSpeed

	if newIndex >= roadLength {
		(*carCnt)++
	} else if newIndex < roadLength && (*newRoad)[newIndex].kind != 0 {
		fmt.Println("SDV crashes something.", newIndex, (*newRoad)[newIndex].kind, newSpeed, deltaD, prevCar.speed, (*newRoad)[newIndex].speed)
	} else {
		if prevCarIndex < roadLength {
			(*currentRoad)[currentIndex].speed = newSpeed
		}
		(*newRoad)[newIndex].speed = newSpeed
		(*newRoad)[newIndex].backlight = newLight
		(*newRoad)[newIndex].kind = 2
	}

}
