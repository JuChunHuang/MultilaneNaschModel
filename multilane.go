package main

import (
	"fmt"
	"math/rand"
)

func MultiLaneSimulation(currentRoad MultiRoad, i, laneNum int, sdvPercentage float64) (MultiRoad, int) {
	// whether to produce a new car at the beginning of each road
	if i%2 == 0 {
		ProduceMulti(&currentRoad, sdvPercentage)
	}
	// return a MultiRoad that all cars have been determined changing lane or not
	currentRoad = ChangeLane(&currentRoad, laneNum)
	// return a MultiRoad that all cars move to their new index based on newspeed
	newRoad, carCnt := ChangeSpeed(currentRoad, laneNum)

	return newRoad, carCnt
}

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
				ChangeSpeedNSDV(&currentRoad, &newRoad, &carCnt, curLane, j, prevCarIndex, laneNum)
			} else if kind == 2 {
				ChangeSpeedSDV(&currentRoad, &newRoad, &carCnt, curLane, j, prevCarIndex)
			} else {
				newRoad[curLane][j] = currentRoad[curLane][j]
			}
		}
	}
	return newRoad, carCnt
}

func ChangeSpeedSDV(currentRoad, newRoad *MultiRoad, carCnt *int, curLane, currentIndex, prevCarIndex int) {
	var newSpeed int
	var newLight int
	var prevCar Car
	currentCar := (*currentRoad)[curLane][currentIndex]
	speed := currentCar.speed

	deltaD := prevCarIndex - currentIndex

	if deltaD >= roadLength {
		newSpeed = speed + 1
		newLight = 0
	} else {
		prevCar = (*currentRoad)[curLane][prevCarIndex]
		if prevCar.kind == 1 || prevCar.kind == 3 {
			if deltaD >= safeSpaceMax[speed] {
				newSpeed = speed + 1
				newLight = 0
			} else if deltaD < safeSpaceMax[speed] && deltaD >= safeSpaceMin[speed] {
				if prevCar.backlight == 0 {
					newSpeed = speed + 1
					newLight = 0
				} else {
					newSpeed = speed - 1
					newLight = 1
				}

			} else if deltaD < safeSpaceMin[speed] {
				newSpeed = min(speed-1, deltaD-1)
				newLight = 1
			}
		} else if prevCar.kind == 2 {
			SDVminDis := GetSDVmindis(currentIndex, prevCarIndex, (*currentRoad)[curLane])
			if deltaD >= safeSpaceMax[speed] {
				newSpeed = speed + 1
				newLight = 0
			} else if deltaD > SDVminDis && deltaD < safeSpaceMax[speed] {
				if CheckTrain((*currentRoad)[curLane], prevCarIndex) >= 5 {
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
				if CheckTrain((*currentRoad)[curLane], prevCarIndex) >= 5 {
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
	} else if newIndex < roadLength && (*newRoad)[curLane][newIndex].kind != 0 {
		fmt.Println("SDV crashes something.", newIndex)
	} else {
		(*newRoad)[curLane][newIndex].speed = newSpeed
		(*newRoad)[curLane][newIndex].backlight = newLight
		(*newRoad)[curLane][newIndex].kind = currentCar.kind
		(*newRoad)[curLane][newIndex].turninglight = currentCar.turninglight
	}

}

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
		aimLane = curLane - 1
		leftPrevIndex = GetPrevCar((*currentRoad)[aimLane], currentIndex)
		deltaLeft = leftPrevIndex - currentIndex
		if leftPrevIndex <= roadLength {
			if deltaLeft < safeSpaceMax[speed] && (*currentRoad)[aimLane][leftPrevIndex].backlight == 1 {
				leftNear = true
			} else {
				leftNear = false
			}
		}

	}

	if ValidLane(curLane+1, laneNum) {
		aimLane = curLane + 1
		rightPrevIndex = GetPrevCar((*currentRoad)[aimLane], currentIndex)
		deltaRight = leftPrevIndex - currentIndex
		if rightPrevIndex <= roadLength {
			if deltaRight < safeSpaceMax[speed] && (*currentRoad)[aimLane][rightPrevIndex].backlight == -1 {
				rightNear = true
			} else {
				rightNear = false
			}
		}

	}

	if leftNear || rightNear {
		return true
	} else {
		return false
	}
}

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
			newSpeed = speed + 1
			newLight = 0
		} else {
			newSpeed = speed
			newLight = 0
		}

	} else if deltaD < safeSpaceMax[speed] && deltaD >= safeSpaceMin[speed] {
		prevCar = (*currentRoad)[curLane][currentIndex]
		if prevCar.backlight == 0 {
			newSpeed, newLight = Randomdeceleraion(p2, speed)
		} else {
			newSpeed, newLight = Randomdeceleraion(p1, speed)
		}
		if newSpeed == speed {
			if CheckNearbyLaneChange(currentRoad, curLane, currentIndex, laneNum) {
				newSpeed, newLight = Randomdeceleraion(p1, speed)
			}
		}
	} else if deltaD < safeSpaceMin[speed] {
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

func ChangeLaneNSDV(currentRoad, newRoad *MultiRoad, curLane, currentIndex, prevCarIndex, laneNum int) {
	var prevCar Car
	var newturningLight int

	currentCar := (*currentRoad)[curLane][currentIndex]

	if prevCarIndex >= roadLength {
		prevCar.kind = 0
		newturningLight = 0
	} else {
		prevCar = (*currentRoad)[curLane][prevCarIndex]
		newturningLight = ChangeNSDVTurningLight(currentRoad, curLane, currentIndex, laneNum)
	}

	speed := currentCar.speed
	deltaD := prevCarIndex - currentIndex
	aimLane := curLane + newturningLight

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
			fmt.Println("change!")
			(*newRoad)[aimLane][currentIndex].kind = currentCar.kind
			(*newRoad)[aimLane][currentIndex].speed = speed
			(*newRoad)[aimLane][currentIndex].turninglight = newturningLight
			(*newRoad)[aimLane][currentIndex].backlight = currentCar.backlight
			(*newRoad)[curLane][currentIndex].turninglightTime = 0
		}
	} else {
		(*newRoad)[curLane][currentIndex].kind = currentCar.kind
		(*newRoad)[curLane][currentIndex].speed = speed
		(*newRoad)[curLane][currentIndex].turninglight = newturningLight
		(*newRoad)[curLane][currentIndex].backlight = currentCar.backlight
		(*newRoad)[curLane][currentIndex].turninglightTime = currentCar.turninglightTime
	}

}

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
			} else if kind == 2 && CheckTrain((*currentRoad)[curLane], j) > 1 {
				if CheckTrain((*currentRoad)[curLane], j) <= trainLength {
					ChangeLaneSDVTrain(currentRoad, &newRoad, curLane, j, laneNum)
				} else if GetTrainTail((*currentRoad)[curLane], j) != j {
					ChangeLaneSDVTrain(currentRoad, &newRoad, curLane, j, laneNum)
				} else {
					ChangeLaneSDV(currentRoad, &newRoad, curLane, j, prevCarIndex, laneNum)
				}
			} else if kind == 2 && CheckTrain((*currentRoad)[curLane], j) == 1 {
				ChangeLaneSDV(currentRoad, &newRoad, curLane, j, prevCarIndex, laneNum)
			} else {
				newRoad[curLane][j] = (*currentRoad)[curLane][j]
			}
		}
	}

	return newRoad
}

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

func ChangeLaneSDVTrain(currentRoad, newRoad *MultiRoad, curLane, currentIndex, laneNum int) {
	var turningLight int
	var trainhead int

	trainhead = GetTrainHead((*currentRoad)[curLane], currentIndex)

	turningLight = ChangeSDVTurningLight(currentRoad, curLane, trainhead, laneNum)

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
			fmt.Printf("SDV changed lane from %v to %v at %v\n", curLane, aimLane, currentIndex)
		}
	}
}

func ChangeLaneSDV(currentRoad, newRoad *MultiRoad, curLane, currentIndex, prevCarIndex, laneNum int) {
	var prevCar Car
	var turningLight int

	if prevCarIndex >= roadLength {
		prevCar.kind = 0
		turningLight = 0
	} else {
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
			fmt.Printf("SDV changed lane from %v to %v at %v\n", curLane, aimLane, currentIndex)
		}
	}
}

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

func LCSforSDV(road *MultiRoad, curLane, aimLane, curCarIndex int) bool {
	var curLanePrev Car
	var curLaneNext Car

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
			if aimAheadDelta_d >= safeSpaceSDVMin[speed] {
				t2 = true
			} else {
				t2 = false
			}

		} else {
			//aimAheadDelta_d >= (safeSpaceMin[speed] - safeSpaceMin[(*road)[aimLane][aimAheadCarIndex].speed] + 1 + 2*(*road)[aimLane][aimAheadCarIndex].speed + 1)
			if aimAheadDelta_d >= GetSDVmindis(curCarIndex, aimAheadCarIndex, (*road)[aimLane]) {
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
			if aimBackDelta_d > GetSDVmindis(aimNextCarIndex, curCarIndex, (*road)[aimLane]) || (speed >= aimNextCar.speed && CheckTrain((*road)[aimLane], aimNextCarIndex) < trainLength) {
				t3 = true
			} else {
				t3 = false
			}
		}

	}

	if (t1 && t2 && t3) == true {
		return true
	} else {
		return false
	}

}

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
			res = false
		}
	} else {
		if aimAheadDelta_d > safeSpaceMax[speed] {
			p := rand.Float64()
			if p > 0.5 {
				res = true
			}
		} else {
			res = false
		}

	}

	return res
}

func LCMforSDV(road *MultiRoad, currentRoadIndex, aimRoadIndex int, currentCarIndex int) bool {
	var prevCar Car
	var aimprevCar Car
	var res bool

	res = true
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
		if delta <= safeSpaceMax[speed] || (prevCar.speed < speed && CheckTrain((*road)[currentRoadIndex], prevCarIndex) >= trainLength) {
			if aimprevCar.kind == 0 {
				res = true
			} else if deltaAimCar > safeSpaceMax[speed] || (deltaAimCar > GetSDVmindis(currentCarIndex, prevAimCarIndex, (*road)[aimRoadIndex]) && aimprevCar.speed > speed && CheckTrain((*road)[aimRoadIndex], prevAimCarIndex) < trainLength) {
				res = true
			} else {
				res = false
			}
		} else {
			res = false
		}

	}
	return res

}
