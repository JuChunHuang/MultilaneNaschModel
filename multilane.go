package main

import (
	"fmt"
	"math/rand"
)

// highest level function to simulate multilane running situations
func MultiLaneSimulation(currentRoad MultiRoad, i, laneNum int, sdvPercentage float64) (MultiRoad, int) {
	// whether to produce a new car at the beginning of each road
	// if i%2 == 0 {
	// 	ProduceMulti(&currentRoad, sdvPercentage)
	// }
	ProduceMulti(&currentRoad, sdvPercentage)
	// return a MultiRoad that all cars have been determined changing lane or not
	currentRoad = ChangeLane(&currentRoad, laneNum)
	// return a MultiRoad that all cars move to their new index based on newspeed
	newRoad, carCnt := ChangeSpeed(currentRoad, laneNum)

	return newRoad, carCnt
}

// fucntion changes spped of the car on the raod
// based on functions ChangeSpeedNSDV and ChangeSpeedSDV
func ChangeSpeed(currentRoad MultiRoad, laneNum int) (MultiRoad, int) {
	var carCnt int
	var kind int
	var prevCarIndex int

	// make a new empty multi-lane road
	newRoad := make(MultiRoad, laneNum)
	for i := 0; i < laneNum; i++ {
		newRoad[i] = make(Road, roadLength)
	}

	// carCnt tracks the number of cars that has travelled along the road
	carCnt = 0

	// go over each grid on the roads, find SDVs and NSDVs and change their speeds
	for curLane := 0; curLane < laneNum; curLane++ {
		for j := roadLength - 1; j >= 0; j-- {
			// get the nearest car and traffic light ahead
			prevCarIndex = GetPrevCar(currentRoad[curLane], j)
			kind = currentRoad[curLane][j].kind

			// if current car is a NSDV
			if kind == 1 {
				// if current car is an NSDV
				ChangeSpeedNSDV(&currentRoad, &newRoad, &carCnt, curLane, j, prevCarIndex, laneNum)
			} else if kind == 2 {
				// if current car is a SDV
				ChangeSpeedSDV(&currentRoad, &newRoad, &carCnt, curLane, j, prevCarIndex)
			} else {
				//empty or light
				newRoad[curLane][j] = currentRoad[curLane][j]
			}
		}
	}
	return newRoad, carCnt
}

// function changes speed for SDV
func ChangeSpeedSDV(currentRoad, newRoad *MultiRoad, carCnt *int, curLane, currentIndex, prevCarIndex int) {
	var newSpeed int
	var newLight int
	var prevCar Car
	//(*currentRoad)[curLane][currentIndex].passingTime += 1
	currentCar := (*currentRoad)[curLane][currentIndex]
	speed := currentCar.speed

	//distance between currentcar and previous car
	deltaD := prevCarIndex - currentIndex

	if deltaD >= roadLength {
		newSpeed = speed + 1
		newLight = 0
	} else {
		prevCar = (*currentRoad)[curLane][prevCarIndex]
		if prevCar.kind == 1 || prevCar.kind == 3 {
			// previous car is NSDV or light
			if deltaD >= safeSpaceMax[speed] {
				// safe to speed up
				newSpeed = speed + 1
				newLight = 0
			} else if deltaD < safeSpaceMax[speed] && deltaD >= safeSpaceMin[speed] {
				if prevCar.backlight == 0 {
					newSpeed = speed
					newLight = 0
				} else {
					newSpeed = speed - 1
					newLight = 1
				}

			} else if deltaD < safeSpaceMin[speed] {
				// need to decelerate for safety
				newSpeed = min(speed-1, deltaD-2)
				newLight = 1
			}
		} else if prevCar.kind == 2 {
			// previous car is SDV
			SDVminDis := GetSDVmindis(currentIndex, prevCarIndex, (*currentRoad)[curLane])
			if deltaD >= safeSpaceMax[speed] {
				// safe to speed up
				newSpeed = speed + 1
				newLight = 0
			} else if deltaD > SDVminDis && deltaD < safeSpaceMax[speed] {
				//situation: SDV in a SDV train
				if CheckTrain((*currentRoad)[curLane], prevCarIndex) >= trainLength {
					if prevCar.speed > speed {
						//keep speed
						newSpeed = speed
						newLight = 0
					} else {
						//keep same speed as the previous car
						// prevent a crash
						newSpeed = prevCar.speed
						newLight = 1
					}
				} else {
					// situationL SDV not in a SDV train
					if prevCar.speed > speed {
						newSpeed = speed + 1
						newLight = 0
					} else {
						newSpeed = prevCar.speed + 1
						newLight = 1
					}
				}

			} else {
				newSpeed = prevCar.speed
				newLight = prevCar.backlight

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
		//fmt.Println("passing Time", currentCar.passingTime, currentCar.kind)
	} else if newIndex < roadLength && (*newRoad)[curLane][newIndex].kind != 0 {
		fmt.Println("SDV crashes something.", (*newRoad)[curLane][newIndex].kind)
	} else {
		(*newRoad)[curLane][newIndex].speed = newSpeed
		(*newRoad)[curLane][newIndex].backlight = newLight
		(*newRoad)[curLane][newIndex].kind = currentCar.kind
		(*newRoad)[curLane][newIndex].turninglight = currentCar.turninglight
		(*newRoad)[curLane][newIndex].passingTime = currentCar.passingTime + 1
		(*currentRoad)[curLane][currentIndex].speed = newSpeed
		(*currentRoad)[curLane][currentIndex].backlight = newLight
	}

}

// fucntion will check if car changes left or right situation
func CheckNearbyLaneChange(currentRoad *MultiRoad, curLane, currentIndex, laneNum int) bool {
	var leftPrevIndex int
	var deltaLeft int
	var rightPrevIndex int
	var deltaRight int
	var aimLane int
	var speed int
	var leftNear, rightNear bool

	speed = (*currentRoad)[curLane][currentIndex].speed

	if ValidLane(curLane-1, laneNum) {
		// left lane exists
		aimLane = curLane - 1
		leftPrevIndex = GetPrevCar((*currentRoad)[aimLane], currentIndex)
		deltaLeft = leftPrevIndex - currentIndex
		if leftPrevIndex <= roadLength {
			if deltaLeft < safeSpaceMax[speed] && (*currentRoad)[aimLane][leftPrevIndex].backlight == 1 {
				// can change to left
				leftNear = true
			} else {
				leftNear = false
			}
		}

	}

	// right lane exists
	if ValidLane(curLane+1, laneNum) {
		aimLane = curLane + 1
		rightPrevIndex = GetPrevCar((*currentRoad)[aimLane], currentIndex)
		deltaRight = leftPrevIndex - currentIndex
		if rightPrevIndex <= roadLength {
			if deltaRight < safeSpaceMax[speed] && (*currentRoad)[aimLane][rightPrevIndex].backlight == -1 {
				//can change to right
				rightNear = true
			} else {
				rightNear = false
			}
		}

	}

	if leftNear || rightNear {
		// there is a possible lane to change
		return true
	} else {
		return false
	}
}

// fucntion change speed for NSDV
func ChangeSpeedNSDV(currentRoad, newRoad *MultiRoad, carCnt *int, curLane, currentIndex, prevCarIndex, laneNum int) {
	var newSpeed int
	var newLight int
	var prevCar Car
	currentCar := (*currentRoad)[curLane][currentIndex]
	speed := currentCar.speed

	var deltaD int

	deltaD = prevCarIndex - currentIndex

	if deltaD >= safeSpaceMax[speed] {
		if CheckNearbyLaneChange(currentRoad, curLane, currentIndex, laneNum) == false {
			// there is no possilble lane to change
			newSpeed = speed + 1
			newLight = 0
		} else {
			newSpeed = speed
			newLight = 0
		}

	} else if deltaD < safeSpaceMax[speed] && deltaD >= safeSpaceMin[speed] {
		// based on driver changing lane possibility factor
		prevCar = (*currentRoad)[curLane][currentIndex]
		if prevCar.backlight == 0 {
			newSpeed, newLight = Randomdeceleraion(p2, speed)
		} else {
			newSpeed, newLight = Randomdeceleraion(p1, speed)
		}
		if newSpeed == speed {
			if CheckNearbyLaneChange(currentRoad, curLane, currentIndex, laneNum) {
				// there is possilble lane to change
				newSpeed, newLight = Randomdeceleraion(p1, speed)
			}
		}
	} else if deltaD < safeSpaceMin[speed] {
		// unsafe - slow down
		prevCar = (*currentRoad)[curLane][currentIndex]
		newSpeed = min(speed-1, deltaD-speed-1)
		newLight = 1
	}

	if newSpeed < 0 {
		newSpeed = 0
	} else if newSpeed > maxSpeed {
		newSpeed = maxSpeed
	}

	// change the position and speed of current car on the new road
	newIndex := currentIndex + newSpeed
	if newIndex >= roadLength {
		(*carCnt)++
	} else if newIndex < roadLength && (*newRoad)[curLane][newIndex].kind != 0 {
		fmt.Println("NSDV crashes something.", newIndex)
	} else {
		(*newRoad)[curLane][newIndex].speed = newSpeed
		(*newRoad)[curLane][newIndex].kind = currentCar.kind
		(*newRoad)[curLane][newIndex].backlight = newLight
		(*newRoad)[curLane][newIndex].turninglight = currentCar.turninglight
	}
}

// randome number simulats changing lane possibility
func RandomTurn(p float64, turningLight int) int {
	var newturningLight int

	thresToTurn := rand.Float64()
	if thresToTurn <= p {
		newturningLight = turningLight
	} else {
		newturningLight = 0
	}

	return newturningLight

}

// apply NSDV change lane
func ChangeLaneNSDV(currentRoad, newRoad *MultiRoad, curLane, currentIndex, prevCarIndex, laneNum int) {
	var prevCar Car
	var newturningLight int

	currentCar := (*currentRoad)[curLane][currentIndex]

	if prevCarIndex >= roadLength {
		prevCar.kind = 0
		newturningLight = 0
	} else {
		//get original change lane direction
		prevCar = (*currentRoad)[curLane][prevCarIndex]
		newturningLight = ChangeNSDVTurningLight(currentRoad, curLane, currentIndex, laneNum)
	}

	speed := currentCar.speed
	deltaD := prevCarIndex - currentIndex
	aimLane := curLane + newturningLight
	// decide change lane direction
	if newturningLight != 0 {
		if prevCar.backlight == 1 && safeSpaceMin[speed] <= deltaD &&
			safeSpaceMax[speed] > deltaD {
			newturningLight = RandomTurn(cp1, newturningLight)
		} else if prevCar.backlight == 0 && safeSpaceMin[speed] <= deltaD &&
			safeSpaceMax[speed] > deltaD {
			newturningLight = RandomTurn(cp2, newturningLight)
		} else if deltaD < safeSpaceMin[speed] {
			newturningLight = RandomTurn(cp3, newturningLight)
		} else {
		}
	}

	if currentCar.turninglight == newturningLight {
		currentCar.turninglightTime += 1
	} else {
		currentCar.turninglightTime = 0
	}

	if currentCar.turninglightTime == 1 {
		if (*newRoad)[aimLane][currentIndex].kind != 0 {
			fmt.Println("NSDV crashes during changing lane.")
			(*newRoad)[curLane][currentIndex].kind = currentCar.kind
			(*newRoad)[curLane][currentIndex].speed = speed
			(*newRoad)[curLane][currentIndex].turninglight = 0
			(*newRoad)[curLane][currentIndex].backlight = currentCar.backlight
			(*newRoad)[curLane][currentIndex].turninglightTime = 0
		} else {
			if newturningLight != 0 {
				//fmt.Println("NSDVchange!")
			}

			(*newRoad)[aimLane][currentIndex].kind = currentCar.kind
			(*newRoad)[aimLane][currentIndex].speed = speed
			(*newRoad)[aimLane][currentIndex].turninglight = newturningLight
			(*newRoad)[aimLane][currentIndex].backlight = currentCar.backlight
			(*newRoad)[curLane][currentIndex].turninglightTime = 0
		}
	} else {
		// there is not a car on aimlane of same position
		(*newRoad)[curLane][currentIndex].kind = currentCar.kind
		(*newRoad)[curLane][currentIndex].speed = speed
		(*newRoad)[curLane][currentIndex].turninglight = newturningLight
		(*newRoad)[curLane][currentIndex].backlight = currentCar.backlight
		(*newRoad)[curLane][currentIndex].turninglightTime = currentCar.turninglightTime
	}

}

// highlevel changelane on SDVs and NSDVs
// call function ChangeLaneNSDV pr ChangeLaneSDV
func ChangeLane(currentRoad *MultiRoad, laneNum int) MultiRoad {
	var newRoad MultiRoad
	var kind int
	var prevCarIndex int

	// make a new empty multi-lane road
	newRoad = make(MultiRoad, len(*currentRoad))
	for i := range newRoad {
		newRoad[i] = make(Road, roadLength)
	}

	for curLane := 0; curLane < laneNum; curLane++ {
		for j := roadLength - 1; j >= 0; j-- {
			prevCarIndex = GetPrevCar((*currentRoad)[curLane], j)
			kind = (*currentRoad)[curLane][j].kind
			if kind == 1 {
				ChangeLaneNSDV(currentRoad, &newRoad, curLane, j, prevCarIndex, laneNum)
			} else if kind == 2 {
				// 	if CheckTrain((*currentRoad)[curLane], j) <= trainLength {
				// 		fmt.Println("EE")
				// 		ChangeLaneSDVTrain(currentRoad, &newRoad, curLane, j, laneNum)
				// 	} else if GetTrainTail((*currentRoad)[curLane], j) != j {
				// 		fmt.Println("EE")
				// 		ChangeLaneSDVTrain(currentRoad, &newRoad, curLane, j, laneNum)
				// 	} else {
				// 		ChangeLaneSDV(currentRoad, &newRoad, curLane, j, prevCarIndex, laneNum)
				// 	}
				// } else if kind == 2 && CheckTrain((*currentRoad)[curLane], j) == 1 {
				ChangeLaneSDV(currentRoad, &newRoad, curLane, j, prevCarIndex, laneNum)
			} else {
				newRoad[curLane][j] = (*currentRoad)[curLane][j]
			}
		}
	}

	return newRoad
}

// function judge if change lane and return aimlane number
func ChangeNSDVTurningLight(currentRoad *MultiRoad, curLane, curCarIndex, laneNum int) int {
	var aimLane int
	var lcm, lcs, leftLane, rightLane bool
	leftLane = false
	rightLane = false
	lanePreference := 0.5

	if ValidLane(curLane-1, laneNum) {
		aimLane = curLane - 1
		lcm = LCMforNSDV(currentRoad, curLane, aimLane, curCarIndex, laneNum)
		lcs = LCSforNSDV(currentRoad, curLane, aimLane, curCarIndex)
		if lcm && lcs {
			leftLane = true
		}
	}

	if ValidLane(curLane+1, laneNum) {
		aimLane = curLane + 1
		lcm = LCMforNSDV(currentRoad, curLane, aimLane, curCarIndex, laneNum)
		lcs = LCSforNSDV(currentRoad, curLane, aimLane, curCarIndex)
		if lcm && lcs {
			rightLane = true
		}
	}

	return DecisionLaneChange(leftLane, rightLane, lanePreference)
}

// fucntion let current car change lane on a direction based on trainhead
func ChangeLaneSDVTrain(currentRoad, newRoad *MultiRoad, curLane, currentIndex, laneNum int) {
	var turningLight int
	var trainhead int

	trainhead = GetTrainHead((*currentRoad)[curLane], currentIndex)
	// get direction of trainhead changing lane direction
	turningLight = ChangeSDVTurningLight(currentRoad, curLane, trainhead, laneNum)

	currentCar := (*currentRoad)[curLane][currentIndex]
	speed := currentCar.speed

	aimLane := curLane + turningLight
	// chnage lane
	if (*newRoad)[aimLane][currentIndex].kind != 0 {
		fmt.Printf("SDV crashes during changing lane.")
		(*newRoad)[curLane][currentIndex].kind = currentCar.kind
		(*newRoad)[curLane][currentIndex].speed = speed
		(*newRoad)[curLane][currentIndex].turninglight = 0
		(*newRoad)[curLane][currentIndex].backlight = currentCar.backlight
	} else {
		(*newRoad)[aimLane][currentIndex].kind = currentCar.kind
		(*newRoad)[aimLane][currentIndex].speed = speed
		(*newRoad)[aimLane][currentIndex].turninglight = turningLight
		(*newRoad)[aimLane][currentIndex].backlight = currentCar.backlight

		if turningLight != 0 {
			//fmt.Printf("SDV changed lane from %v to %v at %v\n", curLane, aimLane, currentIndex)
		}
	}
}

// fucntion let current car change lane
func ChangeLaneSDV(currentRoad, newRoad *MultiRoad, curLane, currentIndex, prevCarIndex, laneNum int) {
	var prevCar Car
	var turningLight int

	if prevCarIndex >= roadLength {
		prevCar.kind = 0
		turningLight = 0
	} else {
		// check direction to change lane or not change
		prevCar = (*currentRoad)[curLane][prevCarIndex]
		turningLight = ChangeSDVTurningLight(currentRoad, curLane, currentIndex, laneNum)
	}
	currentCar := (*currentRoad)[curLane][currentIndex]
	speed := currentCar.speed

	aimLane := curLane + turningLight

	if (*newRoad)[aimLane][currentIndex].kind != 0 {
		fmt.Printf("SDV crashes during changing lane.")
		(*newRoad)[curLane][currentIndex].kind = currentCar.kind
		(*newRoad)[curLane][currentIndex].speed = speed
		(*newRoad)[curLane][currentIndex].turninglight = 0
		(*newRoad)[curLane][currentIndex].backlight = currentCar.backlight
	} else {
		(*newRoad)[aimLane][currentIndex].kind = currentCar.kind
		(*newRoad)[aimLane][currentIndex].speed = speed
		(*newRoad)[aimLane][currentIndex].turninglight = turningLight
		(*newRoad)[aimLane][currentIndex].backlight = currentCar.backlight

		if turningLight != 0 {
			//fmt.Printf("SDV changed lane from %v to %v at %v\n", curLane, aimLane, currentIndex)
		}
	}
}

// function return if currentcar will changelane
// return changing lane direction
func DecisionLaneChange(leftLane, rightLane bool, lanePreference float64) int {
	var alter int
	if leftLane && rightLane {
		p := rand.Float64()
		if p < lanePreference {
			alter = -1
		} else {
			alter = 1
		}
	} else if leftLane {
		alter = -1
	} else if rightLane {
		alter = 1
	} else {
		alter = 0
	}
	return alter
}

// define Lane-Changing Safety(LCS) for NSDV
// Check if current NSDV satisfy LCS rule
func LCSforNSDV(road *MultiRoad, curLane, aimLane, curCarIndex int) bool {
	var res bool
	res = true
	if curCarIndex == 0 {
		res = false
	}
	speed := (*road)[curLane][curCarIndex].speed

	curAheadCarIndex := GetPrevCar((*road)[curLane], curCarIndex)
	curNextCarIndex := GetNext((*road)[curLane], curCarIndex)
	aimAheadCarIndex := GetPrevCar((*road)[aimLane], curCarIndex-1)
	aimNextCarIndex := GetNext((*road)[aimLane], curCarIndex+1)
	curAheadDelta_d := curAheadCarIndex - curCarIndex
	curNextDelta_d := curCarIndex - curNextCarIndex
	aimAheadDelta_d := aimAheadCarIndex - curCarIndex
	aimBackDelta_d := curCarIndex - aimNextCarIndex

	if aimAheadDelta_d >= safeSpaceMin[speed] && aimNextCarIndex < 0 {
		res = true
	} else if aimAheadDelta_d >= safeSpaceMin[speed] &&
		aimBackDelta_d >= safeSpaceMin[(*road)[aimLane][aimNextCarIndex].speed] {
		res = true
	} else {
		res = false
	}

	if curAheadCarIndex <= roadLength {
		if curAheadDelta_d <= safeSpaceMin[speed] && (*road)[curLane][curAheadCarIndex].turninglight == aimLane-curLane {
			res = false
		}
	}

	if curNextCarIndex >= 0 {
		if curNextDelta_d <= safeSpaceMin[(*road)[curLane][curNextCarIndex].speed] && (*road)[curLane][curNextCarIndex].turninglight == aimLane-curLane {
			res = false
		}
	}

	return res
}

// define Lane-Changing Safety(LCS) for SDV
// Check if current SDV satisfy LCS rule
func LCSforSDV(road *MultiRoad, curLane, aimLane, curCarIndex int) bool {
	var curLanePrev Car
	var curLaneNext Car
	// define decision parameters t1, t2, t3
	t1 := false
	t2 := false
	t3 := false
	speed := (*road)[curLane][curCarIndex].speed
	curAheadCarIndex := GetPrevCar((*road)[curLane], curCarIndex)
	curNextCarIndex := GetNext((*road)[curLane], curCarIndex)
	aimAheadCarIndex := GetPrevCar((*road)[aimLane], curCarIndex-1)
	aimNextCarIndex := GetNext((*road)[aimLane], curCarIndex+1)
	curAheadDelta_d := curAheadCarIndex - curCarIndex
	curNextDelta_d := curCarIndex - curNextCarIndex
	aimAheadDelta_d := aimAheadCarIndex - curCarIndex
	aimBackDelta_d := curCarIndex - aimNextCarIndex

	// t1 curlane ahead car no turning light
	if curAheadCarIndex == 2*roadLength && curNextCarIndex == -100 {
		t1 = true
	} else if curAheadCarIndex != 2*roadLength {
		curLanePrev = (*road)[curLane][curAheadCarIndex]
		if curLanePrev.kind == 1 || curLanePrev.kind == 3 {
			if curAheadDelta_d <= safeSpaceMin[speed] && curLanePrev.turninglight == aimLane-curLane {
				t1 = false
			} else {
				t1 = true
			}
		} else {
			if curAheadDelta_d <= GetSDVmindis(curCarIndex, curAheadCarIndex, (*road)[curLane]) && curLanePrev.turninglight == aimLane-curLane {
				t1 = false
			} else {
				t1 = true
			}
		}
	} else {
		curLaneNext = (*road)[curLane][curNextCarIndex]
		if curLanePrev.kind == 1 || curLanePrev.kind == 3 {
			if curNextDelta_d <= safeSpaceMin[speed] && curLaneNext.turninglight == aimLane-curLane {
				//safe condition
				t1 = false
			} else {
				t1 = true
			}
		} else {
			if curNextDelta_d <= GetSDVmindis(curNextCarIndex, curCarIndex, (*road)[curLane]) && curLaneNext.turninglight == aimLane-curLane {
				t1 = false
			} else {
				t1 = true
			}
		}
	}

	if aimAheadCarIndex == 2*roadLength {
		t2 = true
	} else if aimAheadCarIndex != 2*roadLength {
		aimAheadCar := (*road)[aimLane][aimAheadCarIndex]
		if aimAheadCar.kind == 1 || aimAheadCar.kind == 3 {
			if aimAheadDelta_d >= safeSpaceMin[speed] {
				t2 = true
			} else {
				t2 = false
			}

		} else {
			//aimAheadDelta_d >= (safeSpaceMin[speed] - safeSpaceMin[(*road)[aimLane][aimAheadCarIndex].speed] + 1 + 2*(*road)[aimLane][aimAheadCarIndex].speed + 1)
			if aimAheadDelta_d >= safeSpaceMin[speed] {
				t2 = true
			} else {
				t2 = false
			}
		}

	}

	if aimNextCarIndex < 0 {
		t3 = true
	} else if aimNextCarIndex >= 0 {
		aimNextCar := (*road)[aimLane][aimNextCarIndex]
		if aimNextCar.kind == 1 || aimNextCar.kind == 3 {
			if aimBackDelta_d >= safeSpaceMin[aimNextCar.speed] {
				t3 = true
			} else {
				t3 = false
			}
		} else {
			if aimBackDelta_d > safeSpaceMin[aimNextCar.speed] || (speed >= aimNextCar.speed && CheckTrain((*road)[aimLane], aimNextCarIndex) < trainLength) {
				t3 = true
			} else {
				t3 = false
			}
		}

	}
	// satisfy all parameters true
	if (t1 && t2 && t3) == true {
		return true
	} else {
		return false
	}

}

// function judge if change lane and return aimlane number
func ChangeSDVTurningLight(currentRoad *MultiRoad, curLane, curCarIndex, laneNum int) int {
	var aimLane int
	var lcm, lcs, leftLane, rightLane bool
	leftLane = false
	rightLane = false
	lanePreference := 0.5

	if ValidLane(curLane-1, laneNum) {
		aimLane = curLane - 1
		lcm = LCMforSDV(currentRoad, curLane, aimLane, curCarIndex)
		lcs = LCSforSDV(currentRoad, curLane, aimLane, curCarIndex)
		if lcm && lcs {
			leftLane = true
		}
	}

	if ValidLane(curLane+1, laneNum) {
		aimLane = curLane + 1
		lcm = LCMforSDV(currentRoad, curLane, aimLane, curCarIndex)
		lcs = LCSforSDV(currentRoad, curLane, aimLane, curCarIndex)
		if lcm && lcs {
			rightLane = true
		}
	}

	return DecisionLaneChange(leftLane, rightLane, lanePreference)
}

// define Lane-Changing Motivation(LCM) for NSDV
// Check if current NSDV satisfy LCM rule
func LCMforNSDV(road *MultiRoad, curLane, aimLane, curCarIndex, laneNum int) bool {
	var curAheadDelta_d int
	var aimAheadDelta_d int
	if aimLane >= laneNum {
		return false
	}
	var res bool
	res = false
	speed := (*road)[curLane][curCarIndex].speed
	curAheadCarIndex := GetPrevCar((*road)[curLane], curCarIndex)
	aimAheadCarIndex := GetPrevCar((*road)[aimLane], curCarIndex)

	// if there is no car ahead in the curLane
	if curAheadCarIndex > roadLength {
		curAheadDelta_d = 0
	} else {
		curAheadDelta_d = curAheadCarIndex - curCarIndex
	}

	// if there is no car ahead in the aimlane
	if aimAheadCarIndex > roadLength {
		aimAheadDelta_d = roadLength
	} else {
		aimAheadDelta_d = aimAheadCarIndex - curCarIndex
	}

	if curAheadDelta_d < safeSpaceMax[speed] || CheckNearbyLaneChange(road, curLane, curCarIndex, laneNum) {
		if aimAheadDelta_d > safeSpaceMax[speed] {
			res = true
		} else {
			p := rand.Float64()
			if p > 0.8 {
				res = true
			}
		}
	} else {
		if aimAheadDelta_d > safeSpaceMax[speed] {
			p := rand.Float64()
			if p > 0.3 {
				res = true
			}
		} else {
			res = false
		}

	}

	return res
}

// define Lane-Changing Motivation(LCM) for SDV
// Check if current SDV satisfy LCM rule
func LCMforSDV(road *MultiRoad, currentRoadIndex, aimRoadIndex int, currentCarIndex int) bool {
	var prevCar Car
	var aimprevCar Car
	var res bool

	res = false
	speed := (*road)[currentRoadIndex][currentCarIndex].speed

	prevCarIndex := GetPrevCar((*road)[currentRoadIndex], currentCarIndex)
	delta := prevCarIndex - currentCarIndex
	if prevCarIndex > roadLength {
		prevCar.kind = 0
	} else {
		prevCar = (*road)[currentRoadIndex][prevCarIndex]
	}

	prevAimCarIndex := GetPrevCar((*road)[aimRoadIndex], currentCarIndex)
	deltaAimCar := prevAimCarIndex - currentCarIndex
	if prevAimCarIndex > roadLength {
		aimprevCar.kind = 0
	} else {
		aimprevCar = (*road)[aimRoadIndex][prevAimCarIndex]
	}

	if prevCar.kind == 0 {
		res = false
	} else if prevCar.kind == 1 || prevCar.kind == 3 {
		if delta <= safeSpaceMax[speed] {
			if deltaAimCar > safeSpaceMax[speed] {
				res = true
			} else {
				res = false
			}
		} else {
			res = false
		}

	} else {
		// if delta <= safeSpaceMax[speed] || (prevCar.speed < speed && CheckTrain((*road)[currentRoadIndex], prevCarIndex) >= trainLength) {
		// 	if aimprevCar.kind == 0 {
		// 		res = true
		// 	} else if deltaAimCar > safeSpaceMax[speed] || (deltaAimCar > GetSDVmindis(currentCarIndex, prevAimCarIndex, (*road)[aimRoadIndex]) && aimprevCar.speed > speed && CheckTrain((*road)[aimRoadIndex], prevAimCarIndex) < trainLength) {
		// 		res = true
		// 	} else {
		// 		res = false
		// 	}
		// } else {
		// 	res = false
		// }

	}
	return res

}
