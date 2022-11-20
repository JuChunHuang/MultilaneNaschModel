package main

import (
	"math/rand"
)

func SingleLaneSimulation(currentRoad Road) Road {
	var prevCar Car
	var prevCarIndex int
	var newSpeed int
	var newLight int
	var newAccel int
	var probOfDecel float64

	newRoad := make(Road, roadLength)
	carCnt := 0

	Produce(&currentRoad, 0.5)
	for i := range currentRoad {
		prevCarIndex = GetPrev(currentRoad, i)
		if prevCarIndex >= roadLength {
			prevCar.backlight = -2
		} else {
			prevCar = currentRoad[prevCarIndex]
		}

		delta_d := prevCarIndex - i

		if currentRoad[i].kind == 1 {
			// the car is a NSDV, change the speed of the car
			if prevCar.backlight == -1 && delta_d > safeSpaceMin[currentRoad[i].speed] && delta_d < safeSpaceMax[currentRoad[i].speed] {
				probOfDecel = p1
			} else if prevCar.backlight >= 0 && delta_d > safeSpaceMin[currentRoad[i].speed] && delta_d < safeSpaceMax[currentRoad[i].speed] {
				probOfDecel = p2
			} else if currentRoad[i].speed == 0 {
				probOfDecel = p3
			} else {
				probOfDecel = 0
			}

			thresToDecel := rand.Float64()

			if probOfDecel < thresToDecel {
				if currentRoad[i].speed < maxSpeed && delta_d > safeSpaceMax[currentRoad[i].speed] {
					// acceleration because no car in front of it
					newSpeed = currentRoad[i].speed + 1
					newLight = 1
				} else if prevCar.backlight == 1 && delta_d > safeSpaceMin[currentRoad[i].speed] {
					// acceleration because the front car is accelerated
					newSpeed = currentRoad[i].speed + 1
					newLight = 1
				}
			} else {
				if delta_d < safeSpaceMin[currentRoad[i].speed] {
					// deceleration
					newSpeed = currentRoad[i].speed - 1
					newLight = -1
				} else if delta_d == safeSpaceMin[currentRoad[i].speed] {
					// on hold case, speed not changed
					newSpeed = currentRoad[i].speed
					newLight = 0
				}
			}

			newIndex := i + newSpeed
			if newIndex >= roadLength {
				carCnt++
			} else if newIndex < roadLength && newRoad[newIndex].kind != 0 {
				panic("NSDV crashes something.")
			} else {
				if newSpeed < 0 {
					newSpeed = 0
				} else if newSpeed > 10 {
					newSpeed = 10
				}
				newRoad[newIndex].speed = newSpeed
				newRoad[newIndex].kind = currentRoad[i].kind
				newRoad[newIndex].backlight = newLight
			}
		}

		if currentRoad[i].kind == 2 {

			if delta_d >= safeSpaceMax[currentRoad[i].speed] {
				newSpeed = currentRoad[i].speed + 1
				newLight = 1
				newAccel = 1

			} else if prevCar.kind == 1 && prevCar.backlight != -1 && delta_d >= safeSpaceMin[currentRoad[i].speed] {
				newSpeed = currentRoad[i].speed + 1
				newLight = 1
				newAccel = 1

			} else if prevCar.kind == 2 && delta_d > GetSDVmindis(i, prevCarIndex, currentRoad) {

				newSpeed = currentRoad[i].speed + 1
				newLight = 1
				newAccel = 1

			} else if prevCar.kind == 2 && delta_d <= GetSDVmindis(i, prevCarIndex, currentRoad) {
				newSpeed = currentRoad[i].speed - 1
				newLight = -1
				newAccel = 0

			} else {
				newLight = 0
				newAccel = 0
			}

			newIndex := i + newSpeed
			if newIndex >= roadLength {
				carCnt++
			} else if newIndex < roadLength && newRoad[newIndex].kind != 0 {
				panic("SDV crashes something.")
			} else {
				if newSpeed < 0 {
					newSpeed = 0
				} else if newSpeed > 10 {
					newSpeed = 10
				}
				newRoad[newIndex].speed = newSpeed
				newRoad[newIndex].backlight = newLight
				newRoad[newIndex].accel = newAccel
				newRoad[newIndex].kind = currentRoad[i].kind
			}

		}
	}

	return newRoad
}
