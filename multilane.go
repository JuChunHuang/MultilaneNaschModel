package main

import (
	"fmt"
	"math/rand"
)

func MultiLaneSimulation(currentRoad MultiRoad, i, laneNum int, nsdvPercentage float64) (MultiRoad, int) {
	// whether to produce a new car at the beginning of each road
	if i%4 == 0 {
		ProduceMulti(&currentRoad, nsdvPercentage)
	}
	//ProduceMulti(&currentRoad, nsdvPercentage)
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
	var prevLightIndex int

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
			prevLightIndex = GetPrevLight(currentRoad[curLane], j)
			kind = currentRoad[curLane][j].kind

			// if current car is a NSDV
			if kind == 1 {
				ChangeSpeedNSDV(&currentRoad, &newRoad, &carCnt, curLane, j, prevCarIndex, prevLightIndex)
			} else if kind == 2 {
				ChangeSpeedSDV(&currentRoad, &newRoad, &carCnt, curLane, j, prevCarIndex, prevLightIndex)
			} else {
				newRoad[curLane][j] = currentRoad[curLane][j]
			}
		}
	}
	return newRoad, carCnt
}

func ChangeSpeedSDV(currentRoad, newRoad *MultiRoad, carCnt *int, curLane, currentIndex, prevCarIndex, prevLightIndex int) {
	var newSpeed int
	var newLight int
	var newAccel int
	var prevCar Car
	var prevLight Car
	currentCar := (*currentRoad)[curLane][currentIndex]
	speed := currentCar.speed

	// make sure the previous car is within the road length
	if prevCarIndex >= roadLength {
		prevCar.kind = 0
	} else {
		prevCar = (*currentRoad)[curLane][prevCarIndex]
	}

	// make sure the previous traffic light is within the road length
	if prevLightIndex >= roadLength {
		prevLight.kind = 0
	} else {
		prevLight = (*currentRoad)[curLane][prevLightIndex]
	}

	deltaD := prevCarIndex - currentIndex
	//deltaDLight := prevLightIndex - currentIndex

	if prevLight.kind == 5 || prevLight.kind == 0 {
		// if there's no traffic light ahead or the traffic light ahead is green

		if deltaD >= safeSpaceSDVMin[speed] {
			newSpeed = speed + 1
			newLight = 1
			newAccel = 1
		} else if prevCar.kind == 1 && prevCar.backlight != -1 {
			newSpeed = speed
			newLight = 0
			newAccel = 0
		} else if prevCar.kind == 1 && prevCar.backlight == -1 {
			newSpeed = speed - 1
			newLight = -1
			newAccel = 0
		} else if prevCar.kind == 2 && deltaD <= GetSDVmindis(currentIndex, prevCarIndex, (*currentRoad)[curLane]) {
			newSpeed = speed - 1
			newLight = -1
			newAccel = 0
		} else if prevCar.kind == 2 && deltaD >= GetSDVmindis(currentIndex, prevCarIndex, (*currentRoad)[curLane]) && CheckTrain((*currentRoad)[curLane], currentIndex) {

			trainHead := GetTrainHead((*currentRoad)[curLane], currentIndex)

			newSpeed = (*currentRoad)[curLane][trainHead].speed
			newLight = (*currentRoad)[curLane][trainHead].backlight
			newAccel = (*currentRoad)[curLane][trainHead].accel
		} else {
			newSpeed = speed
			newLight = 0
			newAccel = 0
		}

		if newSpeed < 0 {
			newSpeed = 0
		} else if newSpeed > maxSpeed {
			newSpeed = maxSpeed
		}

		// if prevLight.kind == 5 && newSpeed == 0 {

		// }

		// if CheckTrain((*currentRoad)[curLane], currentIndex) {
		// 	trainHead := GetTrainHead((*currentRoad)[curLane], currentIndex)
		// 	if trainHead != currentIndex {
		// 		if deltaD > GetSDVmindis(currentIndex, prevCarIndex, (*currentRoad)[curLane]) {
		// 			panic("not SDV Train")
		// 		}
		// 		fmt.Println("IS SDV train!")
		// 		newSpeed = (*currentRoad)[curLane][trainHead].speed
		// 		newLight = (*currentRoad)[curLane][trainHead].backlight
		// 		newAccel = (*currentRoad)[curLane][trainHead].accel
		// 	}
		// }

		newIndex := currentIndex + newSpeed

		if newIndex >= roadLength {
			(*carCnt)++
		} else if newIndex < roadLength && (*newRoad)[curLane][newIndex].kind != 0 {
			//fmt.Println("SDV crashes something.", newIndex)
		} else {
			(*newRoad)[curLane][newIndex].speed = newSpeed
			(*newRoad)[curLane][newIndex].backlight = newLight
			(*newRoad)[curLane][newIndex].accel = newAccel
			(*newRoad)[curLane][newIndex].kind = currentCar.kind
			(*newRoad)[curLane][newIndex].turninglight = currentCar.turninglight
		}

	} else if prevLight.kind == 3 || prevLight.kind == 4 {
		// if the traffic light ahead is yellow or red

		if prevCarIndex > prevLightIndex {
			deltaD = prevLightIndex - currentIndex
			prevCarIndex = prevLightIndex
			speed = 0
			prevCar.kind = 3
		}

		if deltaD >= safeSpaceMin[speed] {
			newSpeed = speed + 1
			newLight = 1
			newAccel = 1
		} else if prevCar.kind == 1 && prevCar.backlight != -1 {
			newSpeed = speed
			newLight = 0
			newAccel = 0
		} else if prevCar.kind == 1 && prevCar.backlight == -1 {
			newSpeed = speed - 1
			newLight = -1
			newAccel = 0
		} else if prevCar.kind == 2 && deltaD <= GetSDVmindis(currentIndex, prevCarIndex, (*currentRoad)[curLane]) {
			newSpeed = speed - 1
			newLight = -1
			newAccel = 1
		} else if prevCar.kind == 2 && deltaD > GetSDVmindis(currentIndex, prevCarIndex, (*currentRoad)[curLane]) && CheckTrain((*currentRoad)[curLane], currentIndex) {

			trainHead := GetTrainHead((*currentRoad)[curLane], currentIndex)

			newSpeed = (*currentRoad)[curLane][trainHead].speed
			newLight = (*currentRoad)[curLane][trainHead].backlight
			newAccel = (*currentRoad)[curLane][trainHead].accel
		} else if prevLight.kind == 3 {
			newSpeed = speed - 1
			newLight = -1
			newAccel = 0
		} else {
			newSpeed = speed
			newLight = 0
			newAccel = 0
		}

		// if deltaD < safetraffic[speed] {
		// 	newSpeed = deltaD / 3
		// 	newLight = -1
		// }

		if newSpeed < 0 {
			newSpeed = 0
		} else if newSpeed > maxSpeed {
			newSpeed = maxSpeed
		}
		// if CheckTrain((*currentRoad)[curLane], currentIndex) {
		// 	trainHead := GetTrainHead((*currentRoad)[curLane], currentIndex)
		// 	if trainHead != currentIndex {
		// 		if deltaD > GetSDVmindis(currentIndex, prevCarIndex, (*currentRoad)[curLane]) {
		// 			panic("not SDV Train")
		// 		}
		// 		fmt.Println("IS SDV train!")
		// 		newSpeed = (*currentRoad)[curLane][trainHead].speed
		// 		newLight = (*currentRoad)[curLane][trainHead].backlight
		// 		newAccel = (*currentRoad)[curLane][trainHead].accel
		// 	}
		// }

		newIndex := currentIndex + newSpeed
		// if newIndex > roadLength/2 {
		// 	newSpeed = 0
		// }

		if newIndex >= roadLength {
			(*carCnt)++
		} else if newIndex < roadLength && (*newRoad)[curLane][newIndex].kind != 0 {
			// fmt.Println("SDV crashes something.", newIndex)
		} else {
			(*newRoad)[curLane][newIndex].speed = newSpeed
			(*newRoad)[curLane][newIndex].backlight = newLight
			(*newRoad)[curLane][newIndex].accel = newAccel
			(*newRoad)[curLane][newIndex].kind = currentCar.kind
			(*newRoad)[curLane][newIndex].turninglight = currentCar.turninglight
		}
	}

	// else if prevLight.kind == 0 {

	// 	if prevCarIndex > roadLength {
	// 		newSpeed = speed + 1
	// 		newLight = 1
	// 		newAccel = 1
	// 	} else {
	// 		if deltaD >= safeSpaceSDVMin[speed] {
	// 			newSpeed = speed + 1
	// 			newLight = 1
	// 			newAccel = 1
	// 		}

	// 		if prevCar.kind == 1 && prevCar.backlight != -1 && deltaD >= safeSpaceMin[speed] {
	// 			newSpeed = speed + 1
	// 			newLight = 1
	// 			newAccel = 1
	// 		} else if prevCar.kind == 2 && deltaD > GetSDVmindis(currentIndex, prevCarIndex, (*currentRoad)[curLane]) {
	// 			newSpeed = deltaD - safeSpaceSDVMin[speed] - 1
	// 			newLight = 1
	// 			newAccel = 1
	// 		} else if prevCar.kind == 2 && deltaD <= GetSDVmindis(currentIndex, prevCarIndex, (*currentRoad)[curLane]) && CheckTrain((*currentRoad)[curLane], currentIndex) {

	// 			trainHead := GetTrainHead((*currentRoad)[curLane], currentIndex)

	// 			newSpeed = (*currentRoad)[curLane][trainHead].speed
	// 			newLight = (*currentRoad)[curLane][trainHead].backlight
	// 			newAccel = (*currentRoad)[curLane][trainHead].accel
	// 		} else if prevLight.kind < 5 && deltaDLight <= safeSpaceMin[0] {
	// 			newSpeed = speed - 1
	// 			newLight = -1
	// 			newAccel = 0
	// 		} else {
	// 			newLight = 0
	// 			newAccel = 0
	// 		}

	// 	}

	// 	if newSpeed < 0 {
	// 		newSpeed = 0
	// 	} else if newSpeed > maxSpeed {
	// 		newSpeed = maxSpeed
	// 	}

	// 	// if CheckTrain((*currentRoad)[curLane], currentIndex) {
	// 	// 	trainHead := GetTrainHead((*currentRoad)[curLane], currentIndex)
	// 	// 	if trainHead != currentIndex {
	// 	// 		if deltaD > GetSDVmindis(currentIndex, prevCarIndex, (*currentRoad)[curLane]) {
	// 	// 			panic("not SDV Train")
	// 	// 		}
	// 	// 		fmt.Println("IS SDV train!")
	// 	// 		newSpeed = (*currentRoad)[curLane][trainHead].speed
	// 	// 		newLight = (*currentRoad)[curLane][trainHead].backlight
	// 	// 		newAccel = (*currentRoad)[curLane][trainHead].accel
	// 	// 	}
	// 	// }

	// 	newIndex := currentIndex + newSpeed

	// 	if newIndex >= roadLength {
	// 		(*carCnt)++
	// 	} else if newIndex < roadLength && (*newRoad)[curLane][newIndex].kind != 0 {
	// 		// fmt.Println("SDV crashes something.", newIndex)
	// 	} else {
	// 		(*newRoad)[curLane][newIndex].speed = newSpeed
	// 		(*newRoad)[curLane][newIndex].backlight = newLight
	// 		(*newRoad)[curLane][newIndex].accel = newAccel
	// 		(*newRoad)[curLane][newIndex].kind = currentCar.kind
	// 		(*newRoad)[curLane][newIndex].turninglight = currentCar.turninglight
	// 	}

	// }
}

func ChangeSpeedNSDV(currentRoad, newRoad *MultiRoad, carCnt *int, curLane, currentIndex, prevCarIndex, prevLightIndex int) {
	var probOfDecel float64
	var thresToDecel float64
	var newSpeed int
	var newLight int
	var prevCar Car
	var prevLight Car
	currentCar := (*currentRoad)[curLane][currentIndex]
	speed := currentCar.speed

	// make sure the previous car is within the road length
	if prevCarIndex >= roadLength {
		prevCar.kind = 0
	} else {
		prevCar = (*currentRoad)[curLane][prevCarIndex]
	}

	// make sure the previous traffic light is within the road length
	if prevLightIndex >= roadLength {
		prevLight.kind = 0
	} else {
		prevLight = (*currentRoad)[curLane][prevLightIndex]
	}

	deltaD := prevCarIndex - currentIndex
	//deltaDLight := prevLightIndex - currentIndex

	if prevLight.kind == 3 || prevLight.kind == 4 {
		if prevCarIndex > prevLightIndex {
			deltaD = prevLightIndex - currentIndex
			prevCarIndex = prevLightIndex
			prevCar.backlight = -1
			speed = 0
			prevCar.kind = 3
		}
	}

	// NSDV may decelerate stochastically, and the probability depends its distance with the previous car

	if deltaD < safeSpaceMin[speed] && prevCar.backlight == -1 {
		probOfDecel = 0.9999
	} else if deltaD < safeSpaceMin[speed] && prevCar.backlight != -1 {
		probOfDecel = 0.9
	} else if prevCar.backlight == -1 && deltaD > safeSpaceMin[speed] && deltaD < safeSpaceMax[speed] {
		probOfDecel = p1
	} else if prevCar.backlight >= 0 && deltaD > safeSpaceMin[speed] && deltaD < safeSpaceMax[speed] {
		probOfDecel = p2
	} else if prevCar.backlight == -1 {
		probOfDecel = p3
	} else {
		probOfDecel = 0
	}

	thresToDecel = rand.Float64()

	if probOfDecel < thresToDecel {
		if speed < maxSpeed && deltaD > safeSpaceMax[speed] {
			// acceleration because no car in front of it
			newSpeed = speed + 1
			newLight = 1
		} else if prevCar.backlight == 1 && deltaD > safeSpaceMin[speed] {
			// acceleration because the front car is accelerated
			newSpeed = speed + 1
			newLight = 1
		}
	} else {
		if deltaD <= safeSpaceMin[speed] {
			// deceleration
			newSpeed = speed - 1
			newLight = -1
		}
		// else if deltaD == safeSpaceMin[speed] {
		// 	// on hold case, speed not changed
		// 	newSpeed = speed
		// 	newLight = 0
		// }
	}

	// // consider the traffic light before
	// if deltaDLight >= 0 && deltaDLight < safetraffic[speed] && (prevLight.kind == 3 || prevLight.kind == 4) {
	// 	newSpeed = speed - 1
	// 	newLight = -1
	// }

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
		// fmt.Println("NSDV crashes something.", newIndex)
	} else {
		(*newRoad)[curLane][newIndex].speed = newSpeed
		(*newRoad)[curLane][newIndex].kind = currentCar.kind
		(*newRoad)[curLane][newIndex].backlight = newLight
		(*newRoad)[curLane][newIndex].turninglight = currentCar.turninglight
	}
}

func ChangeLaneNSDV(currentRoad, newRoad *MultiRoad, curLane, currentIndex, prevCarIndex, prevLightIndex, laneNum int) {
	var prevCar Car
	var turningLight int
	var probOfTurn float64

	probOfTurn = -1

	if prevCarIndex >= roadLength {
		prevCar.kind = 0
		turningLight = 0
	} else {
		prevCar = (*currentRoad)[curLane][prevCarIndex]
		turningLight = ChangeNSDVTurningLight(currentRoad, curLane, currentIndex, laneNum)
	}
	currentCar := (*currentRoad)[curLane][currentIndex]
	speed := currentCar.speed
	deltaD := prevCarIndex - currentIndex
	aimLane := curLane + turningLight

	if turningLight != 0 {
		if prevCar.backlight == -1 && safeSpaceMin[speed] <= deltaD &&
			safeSpaceMax[speed] > deltaD {
			probOfTurn = cp1
		} else if prevCar.backlight >= 0 && safeSpaceMin[speed] <= deltaD &&
			safeSpaceMax[speed] > deltaD {
			probOfTurn = cp2
		} else if speed == 0 {
			probOfTurn = cp3
		} else {
			probOfTurn = 0
		}
	}

	thresToTurn := rand.Float64()

	if (*newRoad)[aimLane][currentIndex].kind != 0 {
		//panic("NSDV crashes during changing lane.")
		(*newRoad)[curLane][currentIndex].kind = currentCar.kind
		(*newRoad)[curLane][currentIndex].speed = speed
		(*newRoad)[curLane][currentIndex].turninglight = 0
		(*newRoad)[curLane][currentIndex].backlight = currentCar.backlight
	} else if thresToTurn <= probOfTurn {
		(*newRoad)[aimLane][currentIndex].kind = currentCar.kind
		(*newRoad)[aimLane][currentIndex].speed = speed
		(*newRoad)[aimLane][currentIndex].turninglight = turningLight
		(*newRoad)[aimLane][currentIndex].backlight = currentCar.backlight

		if turningLight != 0 {
			//fmt.Printf("NSDV changed lane from %v to %v at %v\n", curLane, aimLane, currentIndex)
		}
	} else {
		(*newRoad)[curLane][currentIndex].kind = currentCar.kind
		(*newRoad)[curLane][currentIndex].speed = speed
		(*newRoad)[curLane][currentIndex].turninglight = 0
		(*newRoad)[curLane][currentIndex].backlight = currentCar.backlight

	}
}

func ChangeLane(currentRoad *MultiRoad, laneNum int) MultiRoad {
	var newRoad MultiRoad
	var kind int
	var prevCarIndex int
	var prevLightIndex int

	// make a new empty multi-lane road
	newRoad = make(MultiRoad, len(*currentRoad))
	for i := range newRoad {
		newRoad[i] = make(Road, roadLength)
	}

	for curLane := 0; curLane < laneNum; curLane++ {
		for j := roadLength - 1; j >= 0; j-- {
			prevCarIndex = GetPrevCar((*currentRoad)[curLane], j)
			prevLightIndex = GetPrevLight((*currentRoad)[curLane], j)
			kind = (*currentRoad)[curLane][j].kind
			if kind == 1 {
				ChangeLaneNSDV(currentRoad, &newRoad, curLane, j, prevCarIndex, prevLightIndex, laneNum)
			} else if kind == 2 {
				ChangeLaneSDV(currentRoad, &newRoad, curLane, j, prevCarIndex, prevLightIndex, laneNum)
			} else if kind > 2 {
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

func ChangeLaneSDV(currentRoad, newRoad *MultiRoad, curLane, currentIndex, prevCarIndex, prevLightIndex, laneNum int) {
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
		//panic("SDV crashes during changing lane.")
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
	res = false
	if curCarIndex == 0 {
		res = false
	}
	speed := (*road)[curLane][curCarIndex].speed

	curAheadCarIndex := GetPrevCar((*road)[curLane], curCarIndex)
	curAheadLightIndex := GetPrevLight((*road)[curLane], curCarIndex)
	if curAheadLightIndex <= curAheadCarIndex {
		curAheadCarIndex = curAheadLightIndex
	}
	aimAheadCarIndex := GetPrevCar((*road)[aimLane], curCarIndex-1)
	aimAheadLightIndex := GetPrevLight((*road)[aimLane], curCarIndex)
	if aimAheadLightIndex <= aimAheadCarIndex {
		aimAheadCarIndex = aimAheadLightIndex
	}
	aimNextCarIndex := GetNext((*road)[aimLane], curCarIndex+1)
	if aimNextCarIndex < 0 {
		aimNextCarIndex = 0
	}
	curAheadDelta_d := curAheadCarIndex - curCarIndex
	aimAheadDelta_d := aimAheadCarIndex - curCarIndex
	aimBackDelta_d := curCarIndex - aimNextCarIndex

	if curAheadDelta_d >= safeSpaceMin[speed] && aimAheadDelta_d >= safeSpaceMin[speed] &&
		aimBackDelta_d >= safeSpaceMin[(*road)[aimLane][aimNextCarIndex].speed] {
		res = true
	} else {
		res = false
	}
	return res
}

func LCSforSDV(road *MultiRoad, curLane, aimLane, curCarIndex int) bool {
	t1 := false
	t2 := false
	t3 := false
	speed := (*road)[curLane][curCarIndex].speed
	curAheadCarIndex := GetPrevCar((*road)[curLane], curCarIndex)
	curAheadLightIndex := GetPrevLight((*road)[curLane], curCarIndex)
	if curAheadLightIndex <= curAheadCarIndex {
		curAheadCarIndex = curAheadLightIndex
	}
	aimAheadCarIndex := GetPrevCar((*road)[aimLane], curCarIndex-1)
	aimAheadLightIndex := GetPrevLight((*road)[aimLane], curCarIndex)
	if aimAheadLightIndex <= aimAheadCarIndex {
		aimAheadCarIndex = aimAheadLightIndex
	}
	aimNextCarIndex := GetNext((*road)[aimLane], curCarIndex+1)
	if aimNextCarIndex < 0 {
		aimNextCarIndex = 0
	}
	curAheadDelta_d := curAheadCarIndex - curCarIndex
	aimAheadDelta_d := aimAheadCarIndex - curCarIndex
	aimBackDelta_d := curCarIndex - aimNextCarIndex

	if curAheadCarIndex == 2*roadLength {
		t1 = true
	} else if curAheadCarIndex != 2*roadLength {
		// if curAheadDelta_d >= safeSpaceMin[speed] {
		// 	t1 = true
		// } else {
		// 	t1 = false
		// }

		if curAheadDelta_d >= (safeSpaceMin[speed] - safeSpaceMin[(*road)[curLane][curAheadCarIndex].speed] + 1 + 2*(*road)[curLane][curAheadCarIndex].speed) {
			t1 = true
		} else {
			t1 = false
		}

	}

	if aimAheadCarIndex == 2*roadLength {
		t2 = true
	} else if aimAheadCarIndex != 2*roadLength {
		if aimAheadDelta_d >= safeSpaceSDVMin[speed] {
			t2 = true
		} else {
			t2 = false
		}

		if aimAheadDelta_d >= (safeSpaceMin[speed] - safeSpaceMin[(*road)[aimLane][aimAheadCarIndex].speed] + 1 + 2*(*road)[aimLane][aimAheadCarIndex].speed + 1) {
			t2 = true
		} else {
			t2 = false
		}

	}

	if aimNextCarIndex == -1 {
		t3 = true
	} else if aimNextCarIndex != -1 {
		if aimBackDelta_d >= safeSpaceMin[(*road)[aimLane][aimNextCarIndex].speed] {
			t3 = true
		} else {
			t3 = false
		}

		if aimBackDelta_d >= (safeSpaceMin[(*road)[aimLane][aimNextCarIndex].speed] - safeSpaceMin[speed] + 2*speed + 1) {
			t3 = true
		} else {
			t3 = false
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
		fmt.Println(lcm, lcs)
	}

	if ValidLane(curLane+1, laneNum) {
		aimLane = curLane + 1
		lcm = LCMforSDV(currentRoad, curLane, aimLane, curCarIndex)
		lcs = LCSforSDV(currentRoad, curLane, aimLane, curCarIndex)
		if lcm && lcs {
			rightLane = true
		}
		fmt.Println(lcm, lcs)
	}

	return DecisionLaneChange(leftLane, rightLane, lanePreference)
}

func LCMforNSDV(road *MultiRoad, curLane, aimLane, curCarIndex, laneNum int) bool {
	var curAheadDelta_d int
	var aimAheadDelta_d int
	var curAheadSpeed int
	var aimAheadSpeed int

	if aimLane >= laneNum {
		return false
	}
	var res bool
	res = false
	speed := (*road)[curLane][curCarIndex].speed
	curAheadCarIndex := GetPrevCar((*road)[curLane], curCarIndex)
	curAheadLightIndex := GetPrevLight((*road)[curLane], curCarIndex)
	if curAheadLightIndex <= curAheadCarIndex {
		curAheadCarIndex = curAheadLightIndex
	}

	aimAheadCarIndex := GetPrevCar((*road)[aimLane], curCarIndex)
	aimAheadLightIndex := GetPrevLight((*road)[aimLane], curCarIndex)
	if aimAheadLightIndex <= aimAheadCarIndex {
		aimAheadCarIndex = aimAheadLightIndex
	}

	// if there is no car ahead in the curLane
	if curAheadCarIndex > roadLength {
		curAheadDelta_d = 0
		curAheadSpeed = 0
	} else {
		curAheadDelta_d = curAheadCarIndex - curCarIndex
		curAheadSpeed = (*road)[curLane][curAheadCarIndex].speed
	}

	// if there is no car ahead in the aimlane
	if aimAheadCarIndex > roadLength {
		aimAheadDelta_d = roadLength
		aimAheadSpeed = maxSpeed + 1
	} else {
		aimAheadDelta_d = aimAheadCarIndex - curCarIndex
		aimAheadSpeed = (*road)[aimLane][aimAheadCarIndex].speed
	}

	if curAheadDelta_d < safeSpaceMax[speed] || (aimAheadSpeed >= speed && curAheadSpeed < speed) {
		if aimAheadDelta_d > safeSpaceMin[speed] {
			res = true
		} else {
			res = false
		}
	}

	return res
}

func LCMforSDV(road *MultiRoad, currentRoadIndex, aimRoadIndex int, currentCarIndex int) bool {
	var prevspeed int
	var aimprevspeed int
	var prevkind int
	var aimkind int

	speed := (*road)[currentRoadIndex][currentCarIndex].speed

	prevCarIndex := GetPrevCar((*road)[currentRoadIndex], currentCarIndex)
	prevLightIndex := GetPrevLight((*road)[currentRoadIndex], currentCarIndex)
	if prevLightIndex <= prevCarIndex {
		prevCarIndex = prevLightIndex
	}
	delta := prevCarIndex - currentCarIndex
	if prevCarIndex > roadLength {
		prevkind = 0
	} else {
		prevspeed = (*road)[currentRoadIndex][prevCarIndex].speed
		prevkind = (*road)[currentRoadIndex][prevCarIndex].kind
	}

	prevAimCarIndex := GetPrevCar((*road)[aimRoadIndex], currentCarIndex)
	prevAimLightIndex := GetPrevLight((*road)[aimRoadIndex], currentCarIndex)
	if prevAimLightIndex <= prevAimCarIndex {
		prevAimCarIndex = prevAimLightIndex
	}
	deltaAimCar := prevAimCarIndex - currentCarIndex
	if prevAimCarIndex > roadLength {
		aimkind = 0
	} else {
		aimprevspeed = (*road)[aimRoadIndex][prevAimCarIndex].speed
		aimkind = (*road)[aimRoadIndex][prevAimCarIndex].kind
	}

	if prevkind != 0 && aimkind != 0 {
		if (delta < safeSpaceMax[speed] || aimprevspeed >= speed) && (prevspeed < speed) {
			if deltaAimCar >= safeSpaceMin[speed] {
				return true
			} else {
				return false
			}
		} else {
			return false
		}
	} else if prevkind == 0 {
		return false
	} else {
		return true
	}
}
