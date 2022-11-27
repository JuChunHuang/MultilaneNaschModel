package main

import (
	"fmt"
	"math/rand"
)

func SingleLaneSimulation(currentRoad Road) Road {
	var prevCar Car
	var prevCarIndex int
	var newSpeed int
	var newLight int
	var newAccel int
	var probOfDecel float64
	var prevLight Car

	newRoad := make(Road, roadLength)
	carCnt := 0

	Produce(&currentRoad, 0.5)
	for i := roadLength - 1; i >= 0; i-- {
		currentCar := currentRoad[i]
		kind := currentCar.kind
		speed := currentCar.speed
		prevCarIndex = GetPrevCar(currentRoad, i)
		prevLightIndex := GetPrevLight(currentRoad, i)

		if prevCarIndex >= roadLength {
			prevCar.backlight = -2
		} else {
			prevCar = currentRoad[prevCarIndex]
		}

		if prevLightIndex >= roadLength {
			prevLight.kind = -1
		} else {
			prevLight = currentRoad[prevLightIndex]
		}

		delta_d := prevCarIndex - i
		deltaDLight := prevLightIndex - i

		if kind == 1 {
			if prevCar.backlight == -1 && delta_d > safeSpaceMin[speed] && delta_d < safeSpaceMax[speed] &&
				(deltaDLight > safetraffic[speed]) {
				probOfDecel = p1
			} else if prevCar.backlight >= 0 && delta_d > safeSpaceMin[speed] && delta_d < safeSpaceMax[speed] &&
				(deltaDLight > safetraffic[speed]) {
				probOfDecel = p2
			} else if speed == 0 {
				probOfDecel = p3
			} else {
				probOfDecel = 0
			}

			thresToDecel := rand.Float64()

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

			newIndex := i + newSpeed
			if newIndex >= roadLength {
				carCnt++
			} else if newIndex < roadLength && newRoad[newIndex].kind != 0 {
				fmt.Println("NSDV crashes something", newIndex, newRoad[newIndex].kind)
				// panic("NSDV crashes something.")
			} else {
				newRoad[newIndex].speed = newSpeed
				newRoad[newIndex].kind = kind
				newRoad[newIndex].backlight = newLight
			}
		} else if (kind == 2 && roadLength/2 <= i) || prevLight.kind == 5 {

			if delta_d >= safeSpaceMax[speed] {
				newSpeed = speed + 1
				newLight = 1
				newAccel = 1
			} else if prevCar.kind == 1 && prevCar.backlight != -1 && delta_d >= safeSpaceMin[speed] {
				newSpeed = speed + 1
				newLight = 1
				newAccel = 1
			} else if prevCar.kind == 2 && delta_d > GetSDVmindis(i, prevCarIndex, currentRoad) {
				newSpeed = speed + 1
				newLight = 1
				newAccel = 1
			} else if prevCar.kind == 2 && delta_d <= GetSDVmindis(i, prevCarIndex, currentRoad) {
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

			trainHead := GetTrainHead(currentRoad, i)
			if CheckTrain(currentRoad, i) == true && trainHead != i {
				if delta_d != GetSDVmindis(i, prevCarIndex, currentRoad) {
					panic("not SDV Train")
				}
				newSpeed = currentRoad[trainHead].speed
				newLight = currentRoad[trainHead].backlight
				newAccel = currentRoad[trainHead].accel
			}

			if newSpeed < 0 {
				newSpeed = 0
			} else if newSpeed > 10 {
				newSpeed = 10
			}
			newIndex := i + newSpeed

			if newIndex >= roadLength {
				carCnt++
			} else if newIndex < roadLength && newRoad[newIndex].kind != 0 {
				// panic("SDV crashes something.")
			} else {
				newRoad[newIndex].speed = newSpeed
				newRoad[newIndex].backlight = newLight
				newRoad[newIndex].accel = newAccel
				newRoad[newIndex].kind = kind
			}

		} else if kind == 2 && roadLength/2 > i && (prevLight.kind == 3 || prevLight.kind == 4) {
			if prevCarIndex > prevLightIndex {
				delta_d = prevLightIndex - i
				prevCarIndex = prevLightIndex
			}
			if delta_d >= safeSpaceMax[speed] {
				newSpeed = speed + 1
				newLight = 1
				newAccel = 1
			} else if prevCar.kind == 1 && prevCar.backlight != -1 && delta_d >= safeSpaceMin[speed] {
				newSpeed = speed + 1
				newLight = 1
				newAccel = 1
			} else if prevCar.kind == 2 && delta_d > GetSDVmindis(i, prevCarIndex, currentRoad) {
				newSpeed = speed + 1
				newLight = 1
				newAccel = 1
			} else if prevCar.kind == 2 && delta_d <= GetSDVmindis(i, prevCarIndex, currentRoad) {
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

			newIndex := i + newSpeed
			if newIndex > roadLength/2 {
				newSpeed = 0
			}
			if newSpeed < 0 {
				newSpeed = 0
			} else if newSpeed > 10 {
				newSpeed = 10
			}

			newRoad[newIndex].speed = newSpeed
			newRoad[newIndex].backlight = newLight
			newRoad[newIndex].accel = newAccel
			newRoad[newIndex].kind = kind

		}
	}

	return newRoad
}
