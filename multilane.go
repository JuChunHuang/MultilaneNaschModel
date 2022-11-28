package main

import (
	"fmt"
	"math/rand"
)

func MultiLaneSimulation(currentRoad MultiRoad) MultiRoad {
	// whether to produce a new car at the beginning of each road
	ProduceMulti(&currentRoad, 0.5)

	// return a MultiRoad that all cars have been determined changing lane or not
	currentRoad = ChangeLane(currentRoad)

	// return a MultiRoad that all cars move to their new index based on newspeed
	newRoad, carCnt := ChangeSpeed(currentRoad)

	if carCnt != 0 {
		fmt.Println("CarCnt:", carCnt)
	}

	return newRoad
}

func ChangeSpeed(currentRoad MultiRoad) (MultiRoad, int) {
	var carCnt int
	var newSpeed int
	var newLight int
	var newAccel int
	var probOfDecel float64
	var thresToDecel float64
	var prevCar Car
	var prevLight Car

	newRoad := make(MultiRoad, laneNum)
	for i := 0; i < laneNum; i++ {
		newRoad[i] = make(Road, roadLength)
	}

	carCnt = 0
	for curLane := 0; curLane < laneNum; curLane++ {
		for j := roadLength - 1; j >= 0; j-- {
			currentCar := currentRoad[curLane][j]
			kind := currentCar.kind
			speed := currentCar.speed
			prevCarIndex := GetPrevCar(currentRoad[curLane], j)
			prevLightIndex := GetPrevLight(currentRoad[curLane], j)

			if prevCarIndex >= roadLength {
				prevCar.backlight = -2
			} else {
				prevCar = currentRoad[curLane][prevCarIndex]
			}

			if prevLightIndex >= roadLength {
				prevLight.kind = -1
			} else {
				prevLight = currentRoad[curLane][prevLightIndex]
			}

			delta_d := prevCarIndex - j
			deltaDLight := prevLightIndex - j

			if kind == 1 {
				// the car is a NSDV, change the speed of the car
				if prevCar.backlight == -1 && delta_d > safeSpaceMin[speed] && delta_d < safeSpaceMax[speed] &&
					(deltaDLight > safetraffic[speed] || deltaDLight < 0) {
					probOfDecel = p1
				} else if prevCar.backlight >= 0 && delta_d > safeSpaceMin[speed] && delta_d < safeSpaceMax[speed] &&
					(deltaDLight > safetraffic[speed] || deltaDLight < 0) {
					probOfDecel = p2
				} else if speed == 0 {
					probOfDecel = p3
				} else {
					probOfDecel = 0
				}

				thresToDecel = rand.Float64()

				if probOfDecel < thresToDecel {
					if speed < maxSpeed && delta_d > safeSpaceMax[speed] {
						// acceleration because no car in front of it
						newSpeed = speed + 1
						newLight = 1
					} else if prevCar.backlight == 1 && delta_d > safeSpaceMin[speed] {
						// acceleration because the front car is accelerated
						newSpeed = speed + 1
						newLight = 1
					}
				} else {
					if delta_d < safeSpaceMin[speed] {
						// deceleration
						newSpeed = speed - 1
						newLight = -1
					} else if delta_d == safeSpaceMin[speed] {
						// on hold case, speed not changed
						newSpeed = speed
						newLight = 0
					}
				}

				if deltaDLight >= 0 && deltaDLight < safetraffic[speed] && (prevLight.kind == 3 || prevLight.kind == 4) {
					newSpeed = 0
					newLight = -1
				}

				if newSpeed < 0 {
					newSpeed = 0
				} else if newSpeed > 10 {
					newSpeed = 10
				}

				newIndex := j + newSpeed
				if newIndex >= roadLength {
					carCnt++
				} else if newIndex < roadLength && newRoad[curLane][newIndex].kind != 0 {
					//panic("NSDV crashes something.")
				} else {
					newRoad[curLane][newIndex].speed = newSpeed
					newRoad[curLane][newIndex].kind = kind
					newRoad[curLane][newIndex].backlight = newLight
				}
			} else if kind == 2 {

				if delta_d >= safeSpaceMax[speed] {
					newSpeed = speed + 1
					newLight = 1
					newAccel = 1
				} else if prevCar.kind == 1 && prevCar.backlight != -1 && delta_d >= safeSpaceMin[speed] {
					newSpeed = speed + 1
					newLight = 1
					newAccel = 1
				} else if prevCar.kind == 2 && delta_d > GetSDVmindis(j, prevCarIndex, currentRoad[curLane]) {
					newSpeed = speed + 1
					newLight = 1
					newAccel = 1
				} else if prevCar.kind == 2 && delta_d <= GetSDVmindis(j, prevCarIndex, currentRoad[curLane]) {
					newSpeed = speed - 1
					newLight = -1
					newAccel = 0
				} else if prevLight.kind > 3 && deltaDLight <= safeSpaceMin[0] {
					newSpeed = speed - 1
					newLight = -1
					newAccel = 0
				} else {
					newLight = 0
					newAccel = 0
				}

				trainHead := GetTrainHead(currentRoad[curLane], j)
				if CheckTrain(currentRoad[curLane], j) == true && trainHead != j {
					if delta_d != GetSDVmindis(j, prevCarIndex, currentRoad[curLane]) {
						panic("not SDV Train")
					}
					newSpeed = currentRoad[curLane][trainHead].speed
					newLight = currentRoad[curLane][trainHead].backlight
					newAccel = currentRoad[curLane][trainHead].accel
				}

				newIndex := j + newSpeed
				if newIndex >= roadLength {
					carCnt++
					// } else if newIndex < roadLength && newRoad[curLane][newIndex].kind != 0 {
					// 	panic("SDV crashes something.")
				} else {
					if newSpeed < 0 {
						newSpeed = 0
					} else if newSpeed > 10 {
						newSpeed = 10
					}
					newRoad[curLane][newIndex].speed = newSpeed
					newRoad[curLane][newIndex].backlight = newLight
					newRoad[curLane][newIndex].accel = newAccel
					newRoad[curLane][newIndex].kind = kind
				}
			}
		}
	}

	return newRoad, carCnt
}

func ChangeLane(currentRoad MultiRoad) MultiRoad {
	var probOfTurn float64
	var turninglight int
	var prevCar Car

	turninglight = -1
	for curLane := 0; curLane < laneNum; curLane++ {
		for j := 0; j < roadLength; j++ {
			currentCar := currentRoad[curLane][j]
			kind := currentCar.kind
			speed := currentCar.speed
			/////// should consider if there is no car in front of it
			prevCarIndex := GetPrevCar(currentRoad[curLane], j)
			if prevCarIndex >= roadLength {
				// no car in front of currentCar, no need to change lane
				turninglight = 0
			} else {
				prevCar = currentRoad[curLane][prevCarIndex]
			}
			delta_d := prevCarIndex - j

			// change lane
			if turninglight == 0 {
				// no car in front of currentCar, no need to change lane
			} else if kind == 1 {
				// change the turning light of NSDV and change the lane
				turninglight := ChangeNSDVTurningLight(currentRoad, curLane, j)
				if turninglight != 0 {
					probOfTurn = 0
					aimLane := curLane + turninglight
					if prevCar.backlight == -1 && safeSpaceMin[speed] <= delta_d &&
						safeSpaceMax[speed] > delta_d {
						probOfTurn = cp1
					} else if prevCar.backlight >= 0 && safeSpaceMin[speed] <= delta_d &&
						safeSpaceMax[speed] > delta_d {
						probOfTurn = cp2
					} else if speed == 0 {
						probOfTurn = cp3
					} else {
						probOfTurn = 0
					}

					thresToTurn := rand.Float64()

					if currentRoad[aimLane][j].kind != 0 {
						panic("NSDV crashes during changing lane.")
					} else {
						if probOfTurn >= thresToTurn {
							currentRoad[aimLane][j].kind = kind
							currentRoad[aimLane][j].speed = speed
							currentRoad[aimLane][j].turninglight = 0
							currentRoad[aimLane][j].backlight = currentCar.backlight

							currentRoad[curLane][j].kind = 0
							currentRoad[curLane][j].speed = 0
							currentRoad[curLane][j].turninglight = 0
							currentRoad[curLane][j].backlight = 0
						}
					}
				}
			} else if kind == 2 {
				turninglight := ChangeLaneSDVCondition(currentRoad, curLane, j)

				if turninglight == -1 {
					leftroadIndex := curLane - 1

					if currentRoad[leftroadIndex][j].kind != 0 {
						panic("SDV crahes during changing lane.")
					} else {
						currentRoad[leftroadIndex][j].speed = currentCar.speed
						currentRoad[leftroadIndex][j].kind = currentCar.kind
						currentRoad[leftroadIndex][j].backlight = currentCar.backlight
						currentRoad[leftroadIndex][j].accel = currentCar.accel
						currentRoad[leftroadIndex][j].turninglight = 0

						currentCar.speed = 0
						currentCar.kind = 0
						currentCar.backlight = 0
						currentCar.accel = 0
						currentCar.turninglight = 0
					}
				} else if turninglight == 1 {
					rightroadIndex := curLane + 1

					if currentRoad[rightroadIndex][j].kind != 0 {
						panic("SDV crahes during changing lane.")
					} else {
						currentRoad[rightroadIndex][j].speed = currentCar.speed
						currentRoad[rightroadIndex][j].kind = currentCar.kind
						currentRoad[rightroadIndex][j].backlight = currentCar.backlight
						currentRoad[rightroadIndex][j].accel = currentCar.accel
						currentRoad[rightroadIndex][j].turninglight = 0

						currentCar.speed = 0
						currentCar.kind = 0
						currentCar.backlight = 0
						currentCar.accel = 0
						currentCar.turninglight = 0
					}
				}
			}
		}
	}

	return currentRoad
}

func ChangeNSDVTurningLight(currentRoad MultiRoad, curLane, curCarIndex int) int {
	currentCar := currentRoad[curLane][curCarIndex]
	// kind := currentCar.kind
	// speed := currentCar.speed
	turningLight := currentCar.turninglight
	// curAheadCarIndex := GetPrevCar(currentRoad[curLane], curCarIndex)
	// curAheadCar := currentRoad[curLane][curAheadCarIndex]
	// delta_d := curAheadCarIndex - curCarIndex
	// NSDV situation
	for aimLane := curLane - 1; aimLane < curLane+2; aimLane++ {
		if !ValidLane(aimLane) {
		} else {
			lcm := LCMforNSDV(currentRoad, curLane, aimLane, curCarIndex)
			lcs := LCSforNSDV(currentRoad, curLane, aimLane, curCarIndex)

			if lcm && lcs {
				turningLight = aimLane - curLane
			}
		}

	}

	return turningLight
}

func LCSforNSDV(road MultiRoad, curLane, aimLane, curCarIndex int) bool {
	var res bool
	// kind := road[curLane][curCarIndex].kind
	speed := road[curLane][curCarIndex].speed
	turningLight := road[curLane][curCarIndex].turninglight
	curAheadCarIndex := GetPrevCar(road[curLane], curCarIndex)
	// curNextCarIndex := GetNext(road[curLane], curCarIndex)
	aimAheadCarIndex := GetPrevCar(road[aimLane], curCarIndex-1)
	aimNextCarIndex := GetNext(road[aimLane], curCarIndex+1)
	curAheadDelta_d := curAheadCarIndex - curCarIndex
	aimAheadDelta_d := aimAheadCarIndex - curCarIndex
	// curBackDelta_d := curCarIndex - curNextCarIndex
	aimBackDelta_d := curCarIndex - aimNextCarIndex

	if aimNextCarIndex < 0 {
		aimNextCarIndex = 0
	}

	if curAheadDelta_d >= safeSpaceMin[speed] && aimAheadDelta_d >= safeSpaceMin[speed] &&
		aimBackDelta_d >= safeSpaceMin[road[aimLane][aimNextCarIndex].speed] {
		if turningLight != aimLane-curLane {
			res = true
		} else {
			res = false
		}
	}

	return res
}

func LCMforNSDV(road MultiRoad, curLane, aimLane, curCarIndex int) bool {
	if aimLane >= laneNum {
		return false
	}
	var res bool
	// kind := road[curLane][curCarIndex].kind
	speed := road[curLane][curCarIndex].speed
	// turningLight := road[curLane][curCarIndex].turninglight
	curAheadCarIndex := GetPrevCar(road[curLane], curCarIndex)
	// curNextCarIndex := GetNext(road[curLane], curCarIndex)
	aimAheadCarIndex := GetPrevCar(road[aimLane], curCarIndex)
	// aimNextCarIndex := GetNext(road[aimLane], curCarIndex)

	// curBackDelta_d := curCarIndex - curNextCarIndex
	// aimBackDelta_d := curCarIndex - aimNextCarIndex
	var curAheadDelta_d int
	var aimAheadDelta_d int
	var curAheadSpeed int
	var aimAheadSpeed int

	// if there is no car ahead in the currentlane
	if curAheadCarIndex > roadLength {
		curAheadDelta_d = 0
		curAheadSpeed = 0
	} else {
		curAheadDelta_d = curAheadCarIndex - curCarIndex
		curAheadSpeed = road[curLane][curAheadCarIndex].speed
	}

	// if there is no car ahead in the aimlane
	if aimAheadCarIndex > roadLength {
		aimAheadDelta_d = roadLength
		aimAheadSpeed = maxSpeed + 1
	} else {
		aimAheadDelta_d = aimAheadCarIndex - curCarIndex
		aimAheadSpeed = road[aimLane][aimAheadCarIndex].speed
	}

	if curAheadDelta_d > safeSpaceMax[speed] {
		res = false
	}

	if aimAheadDelta_d > safeSpaceMax[speed] || (aimAheadSpeed >= speed && curAheadSpeed < speed) {
		res = true
	} else {
		res = false
	}

	return res
}

func LCMforSDVLeft(road []Road, currentRoadIndex int, currentcarIndex int) bool {
	// if currentlane is the edge of the whole lane
	edge1 := 0
	edge2 := len(road) - 1
	if currentRoadIndex == edge1 || currentRoadIndex == edge2 {
		return false
	}

	// get the previous car from currentlane
	var prevCar Car
	prevCarIndex := GetPrevCar(road[currentRoadIndex], currentcarIndex)
	prevCar = road[currentRoadIndex][prevCarIndex]
	delta_dAB := prevCarIndex - currentcarIndex

	// get the previous car from the lane in the left direction
	var prevNeighborCar Car
	prevNeighborCarIndex := GetPrevCar(road[currentRoadIndex-1], currentcarIndex)
	prevNeighborCar = road[currentRoadIndex-1][prevNeighborCarIndex]
	delta_dAE := prevNeighborCarIndex - currentcarIndex

	currentCar := road[currentRoadIndex][currentcarIndex]
	Dmax_vA := safeSpaceMax[currentCar.speed]

	if delta_dAB < Dmax_vA && currentCar.speed > prevCar.speed && (delta_dAE > Dmax_vA || prevNeighborCar.speed > currentCar.speed) || prevCar.turninglight == -1 && prevCar.kind == 2 {
		return true
	}

	return false

}

func LCMforSDVRight(road []Road, currentRoadIndex int, currentcarIndex int) bool {
	// if currentlane is the edge of the whole lane
	edge1 := 0
	edge2 := len(road) - 1
	if currentRoadIndex == edge1 || currentRoadIndex == edge2 {
		return false
	}

	// get the previous car from currentlane
	var prevCar Car
	prevCarIndex := GetPrevCar(road[currentRoadIndex], currentcarIndex)
	prevCar = road[currentRoadIndex][prevCarIndex]
	delta_dAB := prevCarIndex - currentcarIndex

	// get the previous car from the lane in the left direction
	var prevNeighborCar Car
	prevNeighborCarIndex := GetPrevCar(road[currentRoadIndex+1], currentcarIndex)
	prevNeighborCar = road[currentRoadIndex+1][prevNeighborCarIndex]
	delta_dAE := prevNeighborCarIndex - currentcarIndex

	currentCar := road[currentRoadIndex][currentcarIndex]
	Dmax_vA := safeSpaceMax[currentCar.speed]

	if delta_dAB < Dmax_vA && currentCar.speed > prevCar.speed && (delta_dAE > Dmax_vA || prevNeighborCar.speed > currentCar.speed) || prevCar.turninglight == 1 && prevCar.kind == 2 {
		return true
	}

	return false

}

func ChangeLaneSDVCondition(road MultiRoad, currentlane, num int) int {
	// check if the SDV is already to turn left or right,
	//  if it is left, return -1 ; if it is right, return 1 ; otherwise return 0
	if LCMforSDVLeft(road, currentlane, num) == true && LCSforSDVLeft(road, currentlane, num) == true {
		return -1
	} else if LCMforSDVRight(road, currentlane, num) == true && LCSforSDVRight(road, currentlane, num) == true {
		return 1
	} else {
		return 0
	}
}

func LCSforSDVLeft(roads MultiRoad, currentRoadIndex int, currentcarIndex int) bool {

	leftroadIndex := currentRoadIndex - 1
	leftprevIndex := GetPrevCar(roads[currentRoadIndex], currentcarIndex)

	delta_dsamel := leftprevIndex - currentcarIndex
	delta_differlprev := GetPrevCar(roads[leftroadIndex], currentcarIndex) - currentcarIndex
	leftnextcarIndex := GetNext(roads[leftroadIndex], currentcarIndex)
	delta_differlnext := currentcarIndex - leftnextcarIndex

	safeDisPrev := safeSpaceMin[roads[currentRoadIndex][currentcarIndex].speed]
	safeDisNext := safeSpaceMin[roads[leftroadIndex][leftnextcarIndex].speed]

	currentcar := roads[currentRoadIndex][currentcarIndex]

	if delta_dsamel > safeDisPrev && delta_differlprev > safeDisPrev && delta_differlnext > safeDisNext && (currentcar.turninglight != -1 || currentcar.kind == 2) {
		return true
	}
	return false

}

func LCSforSDVRight(roads MultiRoad, currentRoadIndex int, currentcarIndex int) bool {

	rightroadIndex := currentRoadIndex + 1
	rightprevIndex := GetPrevCar(roads[currentRoadIndex], currentcarIndex)

	delta_dsamel := rightprevIndex - currentcarIndex
	delta_differlprev := GetPrevCar(roads[rightroadIndex], currentcarIndex) - currentcarIndex
	rightnextcarIndex := GetNext(roads[rightroadIndex], currentcarIndex)
	delta_differlnext := currentcarIndex - rightnextcarIndex

	safeDisPrev := safeSpaceMin[roads[currentRoadIndex][currentcarIndex].speed]
	safeDisNext := safeSpaceMin[roads[rightroadIndex][rightnextcarIndex].speed]

	currentcar := roads[currentRoadIndex][currentcarIndex]

	if delta_dsamel > safeDisPrev && delta_differlprev > safeDisPrev && delta_differlnext > safeDisNext && (currentcar.turninglight != 1 || currentcar.kind == 2) {
		return true
	}
	return false

}

func CarTurnLeft(roads []Road, currentRoadIndex int, currentcarIndex int) {
	currentcar := roads[currentRoadIndex][currentcarIndex]
	leftroadIndex := currentRoadIndex - 1
	if LCMforSDVLeft(roads, currentRoadIndex, currentcarIndex) == true && LCSforSDVLeft(roads, currentRoadIndex, currentcarIndex) == true {
		currentcar.turninglight = -1

		roads[leftroadIndex][currentcarIndex].speed = currentcar.speed
		roads[leftroadIndex][currentcarIndex].kind = currentcar.kind
		roads[leftroadIndex][currentcarIndex].backlight = currentcar.backlight
		roads[leftroadIndex][currentcarIndex].accel = currentcar.accel

		currentcar.speed = 0
		currentcar.kind = 0
		currentcar.turninglight = 0
		currentcar.accel = 0

	}
}

func CarTurnRight(roads []Road, currentRoadIndex int, currentcarIndex int) {
	currentcar := roads[currentRoadIndex][currentcarIndex]
	RightroadIndex := currentRoadIndex + 1
	if LCMforSDVRight(roads, currentRoadIndex, currentcarIndex) == true && LCSforSDVRight(roads, currentRoadIndex, currentcarIndex) == true {
		currentcar.turninglight = -1

		roads[RightroadIndex][currentcarIndex].speed = currentcar.speed
		roads[RightroadIndex][currentcarIndex].kind = currentcar.kind
		roads[RightroadIndex][currentcarIndex].turninglight = currentcar.turninglight
		roads[RightroadIndex][currentcarIndex].accel = currentcar.accel

		currentcar.speed = 0
		currentcarIndex = 0
		currentcarIndex = 0
		currentcarIndex = 0

	}
}
