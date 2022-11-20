package main

import "math/rand"

func MultiLaneSimulation(currentRoad MultiRoad) {
	// newRoad := make(MultiRoad, laneNum)
	// for i := 0; i < laneNum; i++ {
	// 	newRoad[i] = make([]Car, roadLength)
	// }

	var probOfTurn float64
	for curLane := 0; curLane < laneNum-dedicatedLane; curLane++ {
		for j := 0; j < roadLength; j++ {
			currentCar := currentRoad[curLane][j]
			kind := currentCar.kind
			speed := currentCar.speed
			curAheadCarIndex := GetPrev(currentRoad[curLane], j)
			curAheadCar := currentRoad[curLane][curAheadCarIndex]
			delta_d := curAheadCarIndex - j

			if kind == 1 {
				// change the turning light of NSDV and change the lane
				turningLight := ChangeNSDVTurningLight(currentRoad, curLane, j)
				if turningLight != 0 {
					probOfTurn = 0
					aimLane := curLane + turningLight
					if curAheadCar.backlight == -1 && safeSpaceMin[speed] <= delta_d &&
						safeSpaceMax[speed] > delta_d {
						probOfTurn = cp1
					} else if curAheadCar.backlight >= 0 && safeSpaceMin[speed] <= delta_d &&
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
			}
			if kind == 2 {
				turninglight := ChangeLaneSDVCondition(currentRoad, curLane, j)

				if turninglight == -1 {
					currentcar := currentRoad[curLane][j]
					leftroadIndex := curLane - 1

					if currentRoad[leftroadIndex][j].kind != 0 {
						panic("SDV crahes during changing lane.")
					} else {

						currentRoad[leftroadIndex][j].speed = currentcar.speed
						currentRoad[leftroadIndex][j].kind = currentcar.kind
						currentRoad[leftroadIndex][j].backlight = currentcar.backlight
						currentRoad[leftroadIndex][j].accel = currentcar.accel
						currentRoad[leftroadIndex][j].turninglight = 0

						currentcar.speed = 0
						currentcar.kind = 0
						currentcar.backlight = 0
						currentcar.accel = 0
						currentcar.turninglight = 0

					}
				} else if turninglight == 1 {
					currentcar := currentRoad[curLane][j]
					rightroadIndex := curLane + 1

					if currentRoad[rightroadIndex][j].kind != 0 {
						panic("SDV crahes during changing lane.")
					} else {

						currentRoad[rightroadIndex][j].speed = currentcar.speed
						currentRoad[rightroadIndex][j].kind = currentcar.kind
						currentRoad[rightroadIndex][j].backlight = currentcar.backlight
						currentRoad[rightroadIndex][j].accel = currentcar.accel
						currentRoad[rightroadIndex][j].turninglight = 0

						currentcar.speed = 0
						currentcar.kind = 0
						currentcar.backlight = 0
						currentcar.accel = 0
						currentcar.turninglight = 0

					}
				}
			}
		}
	}
}

func ChangeNSDVTurningLight(currentRoad MultiRoad, curLane, curCarIndex int) int {
	currentCar := currentRoad[curLane][curCarIndex]
	kind := currentCar.kind
	// speed := currentCar.speed
	turningLight := currentCar.turninglight
	// curAheadCarIndex := GetPrev(currentRoad[curLane], curCarIndex)
	// curAheadCar := currentRoad[curLane][curAheadCarIndex]
	// delta_d := curAheadCarIndex - curCarIndex
	// NSDV situation
	if kind == 1 {
		if turningLight == 0 {
			for aimLane := curLane - 1; aimLane < curLane+2; aimLane++ {
				if !ValidLane(aimLane) {
				}

				lcm := LCMforNSDV(currentRoad, curLane, aimLane, curCarIndex)
				lcs := LCSforNSDV(currentRoad, curLane, aimLane, curCarIndex)

				if lcm && lcs {
					turningLight = aimLane - curLane
				}
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
	curAheadCarIndex := GetPrev(road[curLane], curCarIndex)
	// curNextCarIndex := GetNext(road[curLane], curCarIndex)
	aimAheadCarIndex := GetPrev(road[aimLane], curCarIndex-1)
	aimNextCarIndex := GetNext(road[aimLane], curCarIndex+1)
	curAheadDelta_d := curAheadCarIndex - curCarIndex
	aimAheadDelta_d := aimAheadCarIndex - curCarIndex
	// curBackDelta_d := curCarIndex - curNextCarIndex
	aimBackDelta_d := curCarIndex - aimNextCarIndex

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
	if aimLane >= laneNum-dedicatedLane {
		return false
	}
	var res bool
	// kind := road[curLane][curCarIndex].kind
	speed := road[curLane][curCarIndex].speed
	// turningLight := road[curLane][curCarIndex].turninglight
	curAheadCarIndex := GetPrev(road[curLane], curCarIndex)
	// curNextCarIndex := GetNext(road[curLane], curCarIndex)
	aimAheadCarIndex := GetPrev(road[aimLane], curCarIndex)
	// aimNextCarIndex := GetNext(road[aimLane], curCarIndex)
	curAheadDelta_d := curAheadCarIndex - curCarIndex
	aimAheadDelta_d := aimAheadCarIndex - curCarIndex
	// curBackDelta_d := curCarIndex - curNextCarIndex
	// aimBackDelta_d := curCarIndex - aimNextCarIndex
	curAheadSpeed := road[curLane][curAheadCarIndex].speed
	aimAheadSpeed := road[aimLane][aimAheadCarIndex].speed

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

// Whether the lane is exist
func ValidLane(lane int) bool {
	return (lane >= 0 && lane < laneNum)
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
	prevCarIndex := GetPrev(road[currentRoadIndex], currentcarIndex)
	prevCar = road[currentRoadIndex][prevCarIndex]
	delta_dAB := prevCarIndex - currentcarIndex

	// get the previous car from the lane in the left direction
	var prevNeighborCar Car
	prevNeighborCarIndex := GetPrev(road[currentRoadIndex-1], currentcarIndex)
	prevNeighborCar = road[currentRoadIndex-1][prevNeighborCarIndex]
	delta_dAE := prevNeighborCarIndex - currentcarIndex

	currentCar := road[currentRoadIndex][currentcarIndex]
	Dmax_vA := safeSpaceMax[currentCar.speed]

	if delta_dAB < Dmax_vA && currentCar.speed > prevCar.speed && (delta_dAE > Dmax_vA || prevNeighborCar.speed > currentCar.speed) || prevCar.turninglight == -1 && prevCar.kind == 2 || (Checktrain(currentcarIndex) == true && prevCar.turninglight == -1) {
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
	prevCarIndex := GetPrev(road[currentRoadIndex], currentcarIndex)
	prevCar = road[currentRoadIndex][prevCarIndex]
	delta_dAB := prevCarIndex - currentcarIndex

	// get the previous car from the lane in the left direction
	var prevNeighborCar Car
	prevNeighborCarIndex := GetPrev(road[currentRoadIndex+1], currentcarIndex)
	prevNeighborCar = road[currentRoadIndex+1][prevNeighborCarIndex]
	delta_dAE := prevNeighborCarIndex - currentcarIndex

	currentCar := road[currentRoadIndex][currentcarIndex]
	Dmax_vA := safeSpaceMax[currentCar.speed]

	if delta_dAB < Dmax_vA && currentCar.speed > prevCar.speed && (delta_dAE > Dmax_vA || prevNeighborCar.speed > currentCar.speed) || prevCar.turninglight == 1 && prevCar.kind == 2 || (Checktrain(currentcarIndex) == true && prevCar.turninglight == 1) {
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
	leftprevIndex := GetPrev(roads[currentRoadIndex], currentcarIndex)
	delta_dsamel := leftprevIndex - currentcarIndex
	delta_differlprev := GetPrev(roads[leftroadIndex], currentcarIndex) - currentcarIndex
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
	rightprevIndex := GetPrev(roads[currentRoadIndex], currentcarIndex)
	delta_dsamel := rightprevIndex - currentcarIndex
	delta_differlprev := GetPrev(roads[rightroadIndex], currentcarIndex) - currentcarIndex
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
