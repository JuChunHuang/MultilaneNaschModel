package main

import (
	"fmt"
	"math/rand"
)

func MultiLaneSimulation(currentRoad MultiRoad, i int) MultiRoad {
	// whether to produce a new car at the beginning of each road
	if i%2 == 0 {
		ProduceMulti(&currentRoad, 0.5)
	}

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
					// panic("NSDV crashes something.")
				} else {
					newRoad[curLane][newIndex].speed = newSpeed
					newRoad[curLane][newIndex].kind = kind
					newRoad[curLane][newIndex].backlight = newLight
				}
			} else if kind == 2 && (roadLength/2 <= j || prevLight.kind == 5 || prevLight.kind == -1) {
				if delta_d >= safeSpaceSDVMin[speed] {
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
					if delta_d < safeSpaceSDVMin[0] {
						newSpeed = 0
						newLight = 0
						newAccel = 0
					}
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

				if newSpeed < 0 {
					newSpeed = 0
				} else if newSpeed > 10 {
					newSpeed = 10
				}
				newIndex := j + newSpeed

				if newIndex >= roadLength {
					carCnt++
				} else if newIndex < roadLength && newRoad[curLane][newIndex].kind != 0 {
					// panic("SDV crashes something.")
				} else {
					newRoad[curLane][newIndex].speed = newSpeed
					newRoad[curLane][newIndex].backlight = newLight
					newRoad[curLane][newIndex].accel = newAccel
					newRoad[curLane][newIndex].kind = kind
				}

			} else if kind == 2 && roadLength/2 > j && (prevLight.kind == 3 || prevLight.kind == 4) {
				if prevCarIndex > prevLightIndex {
					delta_d = prevLightIndex - j
					prevCarIndex = prevLightIndex
				}
				if delta_d >= safeSpaceMin[speed] {
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
					if delta_d < safeSpaceSDVMin[0] {
						newSpeed = 0
						newLight = 0
						newAccel = 0
					}
					newSpeed = speed - 1
					newLight = -1
					newAccel = 0
				} else if prevLight.kind >= 3 && deltaDLight <= safeSpaceMin[0] {
					newSpeed = 0
					newLight = 0
					newAccel = 0
				} else {
					newLight = 0
					newAccel = 0
				}

				if delta_d < safetraffic[speed] {
					newSpeed = 0
					newLight = -1
				}

				newIndex := j + newSpeed
				if newIndex > roadLength/2 {
					newSpeed = 0
				}
				if newSpeed < 0 {
					newSpeed = 0
				} else if newSpeed > 10 {
					newSpeed = 10
				}

				newRoad[curLane][newIndex].speed = newSpeed
				newRoad[curLane][newIndex].backlight = newLight
				newRoad[curLane][newIndex].accel = newAccel
				newRoad[curLane][newIndex].kind = kind

				if delta_d >= safeSpaceMin[speed] {
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
				} else if prevCar.kind == 2 && delta_d <= GetSDVmindis(j, prevCarIndex, currentRoad[curLane]) && prevCar.speed != 0 {
					if delta_d < safeSpaceSDVMin[0] {
						newSpeed = 0
						newLight = 0
						newAccel = 0
					}
					newSpeed = speed - 1
					newLight = -1
					newAccel = 0
				} else if prevCar.kind == 2 && delta_d <= GetSDVmindis(j, prevCarIndex, currentRoad[curLane]) && prevCar.speed == 0 {
					if delta_d < safeSpaceSDVMin[0] {
						newSpeed = 0
						newLight = 0
						newAccel = 0
					}
					newSpeed = 0
					newLight = 0
					newAccel = 0

				} else if prevLight.kind > 3 && deltaDLight <= safeSpaceMin[0] {
					newSpeed = 0
					newLight = 0
					newAccel = 0
				} else {
					newLight = 0
					newAccel = 0
				}

				if delta_d < safetraffic[speed] {
					newSpeed = 0
					newLight = -1
				}

				newIndex = j + newSpeed
				if newIndex > roadLength/2 {
					newSpeed = 0
				}
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

	return newRoad, carCnt
}

func ChangeLane(currentRoad MultiRoad) MultiRoad {
	var probOfTurn float64
	var turninglight int
	var prevCar Car

	for curLane := 0; curLane < laneNum; curLane++ {
		turninglight = -1
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
			} else {
				if kind == 1 {
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
								currentRoad[aimLane][j].turninglight = turninglight
								currentRoad[aimLane][j].backlight = currentCar.backlight
								fmt.Println(turninglight)

								currentRoad[curLane][j].kind = 0
								currentRoad[curLane][j].speed = 0
								currentRoad[curLane][j].turninglight = 0
								currentRoad[curLane][j].backlight = 0

								fmt.Printf("NSDV changed lane from %v to %v at %v\n", curLane, aimLane, j)
							}
						}
					}
				} else if kind == 2 {
					turninglight := ChangeSDVTurningLight(currentRoad, curLane, j)
					fmt.Println(curLane, turninglight)

					if turninglight == -1 {
						leftroadIndex := curLane - 1

						if currentRoad[leftroadIndex][j].kind != 0 {
							panic("SDV crahes during changing lane.")
						} else {
							currentRoad[leftroadIndex][j].speed = currentCar.speed
							currentRoad[leftroadIndex][j].kind = currentCar.kind
							currentRoad[leftroadIndex][j].backlight = currentCar.backlight
							currentRoad[leftroadIndex][j].accel = currentCar.accel
							currentRoad[leftroadIndex][j].turninglight = turninglight

							currentCar.speed = 0
							currentCar.kind = 0
							currentCar.backlight = 0
							currentCar.accel = 0
							currentCar.turninglight = 0
							fmt.Printf("SDV changed lane from %v to %v at %v\n", curLane, leftroadIndex, j)
						}
					} else if turninglight == 1 {
						rightroadIndex := curLane + 1
						// fmt.Println(rightroadIndex)

						if currentRoad[rightroadIndex][j].kind != 0 {
							// panic("SDV crahes during changing lane.")
						} else {
							currentRoad[rightroadIndex][j].speed = currentCar.speed
							currentRoad[rightroadIndex][j].kind = currentCar.kind
							currentRoad[rightroadIndex][j].backlight = currentCar.backlight
							currentRoad[rightroadIndex][j].accel = currentCar.accel
							currentRoad[rightroadIndex][j].turninglight = turninglight

							currentCar.speed = 0
							currentCar.kind = 0
							currentCar.backlight = 0
							currentCar.accel = 0
							currentCar.turninglight = 0

							fmt.Printf("SDV changed lane from %v to %v at %v\n", curLane, rightroadIndex, j)
						}
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
	if curCarIndex == 0 {
		res = false
	}
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

func LCSforSDV(road MultiRoad, curLane, aimLane, curCarIndex int) bool {
	t1 := false
	t2 := false
	t3 := false
	speed := road[curLane][curCarIndex].speed
	curAheadCarIndex := GetPrevCar(road[curLane], curCarIndex)
	// curNextCarIndex := GetNext(road[curLane], curCarIndex)
	aimAheadCarIndex := GetPrevCar(road[aimLane], curCarIndex-1)
	aimNextCarIndex := GetNext(road[aimLane], curCarIndex+1)
	curAheadDelta_d := curAheadCarIndex - curCarIndex
	aimAheadDelta_d := aimAheadCarIndex - curCarIndex
	// curBackDelta_d := curCarIndex - curNextCarIndex
	aimBackDelta_d := curCarIndex - aimNextCarIndex

	if curAheadCarIndex == 2*roadLength {
		t1 = true
	} else if curAheadCarIndex != 2*roadLength {
		if curAheadDelta_d >= safeSpaceMin[speed] {
			t1 = true
		} else {
			t1 = false
		}

		if curAheadDelta_d >= (safeSpaceMin[speed] - safeSpaceMin[road[curLane][curAheadCarIndex].speed] + 1 + 2*road[curLane][curAheadCarIndex].speed) {
			t1 = true
		} else {
			t1 = false
		}

	}

	if aimAheadCarIndex == 2*roadLength {
		t2 = true
	} else if aimAheadCarIndex != 2*roadLength {
		if aimAheadDelta_d >= safeSpaceMin[speed] {
			t2 = true
		} else {
			t2 = false
		}

		if aimAheadDelta_d >= (safeSpaceMin[speed] - safeSpaceMin[road[aimLane][aimAheadCarIndex].speed] + 1 + 2*road[aimLane][aimAheadCarIndex].speed + 1) {
			t2 = true
		} else {
			t2 = false
		}

	}

	if aimNextCarIndex == -1 {
		t3 = true
	} else if aimNextCarIndex != -1 {
		if aimBackDelta_d >= safeSpaceMin[road[aimLane][aimNextCarIndex].speed] {
			t3 = true
		} else {
			t3 = false
		}

		if aimBackDelta_d >= (safeSpaceMin[road[aimLane][aimNextCarIndex].speed] - safeSpaceMin[speed] + 2*speed + 1) {
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
func ChangeSDVTurningLight(currentRoad MultiRoad, curLane, curCarIndex int) int {
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
			lcm := LCMforSDV(currentRoad, curLane, aimLane, curCarIndex)
			lcs := LCSforSDV(currentRoad, curLane, aimLane, curCarIndex)

			if lcm && lcs {
				turningLight = aimLane - curLane
			}
		}

	}

	return turningLight
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
func LCMforSDV(road []Road, currentRoadIndex, aimRoadIndex int, currentCarIndex int) bool {

	speed := road[currentRoadIndex][currentCarIndex].speed

	prevCarIndex := GetPrevCar(road[currentRoadIndex], currentCarIndex)
	delta := prevCarIndex - currentCarIndex

	prevAimCarIndex := GetPrevCar(road[aimRoadIndex], currentCarIndex)
	deltaAimCar := prevAimCarIndex - currentCarIndex

	if delta > safeSpaceMax[speed] {
		return false
	}

	if road[currentRoadIndex][prevCarIndex].kind == 1 {
		if (deltaAimCar > safeSpaceMax[speed] || road[aimRoadIndex][prevAimCarIndex].speed >= speed) &&
			(road[currentRoadIndex][prevCarIndex].speed < speed) {
			return true
		} else {
			return false
		}
	} else {
		if road[currentRoadIndex][prevCarIndex].turninglight == aimRoadIndex-currentRoadIndex {
			return true
		} else {
			return false
		}
	}

}
