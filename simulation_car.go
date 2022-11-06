package main

import (
	"math/rand"
)

func SDVNNextStep(currentRoad Road) Road {
	var newRoad Road
	var prevCar Car
	var prevCarIndex int
	var newSpeed int
	var newLight int
	var newAccel int
	var probOfDecel float64

	carCnt := 0

	for i := range currentRoad {

		prevCarIndex = GetPrev(currentRoad, i)
		prevCar = currentRoad[prevCarIndex]
		delta_d := prevCarIndex - i

		if currentRoad[i].kind == 1 {
			// the car is a NSDV, change the speed of the car
			if prevCar.light == 1 && delta_d > safeSpaceMin[currentRoad[i].speed] && delta_d < safeSpaceMax[currentRoad[i].speed] {
				probOfDecel = p1
			} else if prevCar.light == 0 && delta_d > safeSpaceMin[currentRoad[i].speed] && delta_d < safeSpaceMax[currentRoad[i].speed] {
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
				} else if prevCar.light == 1 && delta_d > safeSpaceMin[currentRoad[i].speed] {
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
			if newIndex < roadLength && newRoad[newIndex].kind != 0 {
				panic("NSDV crashes something.")
			} else if newIndex >= roadLength {
				carCnt++
			} else {
				newRoad[newIndex].speed = newSpeed
				newRoad[newIndex].kind = currentRoad[i].kind
				newRoad[newIndex].light = newLight
			}
		}

		if currentRoad[i].kind == 2 {
			if delta_d >= safeSpaceMax[currentRoad[i].speed] {
				newSpeed = currentRoad[i].speed + 1
				newLight = 1
				newAccel = 1

			} else if prevCar.kind == 1 && prevCar.light != -1 && delta_d >= safeSpaceMin[currentRoad[i].speed] {
				newSpeed = currentRoad[i].speed + 1
				newLight = 1
				newAccel = 1

			} else if prevCar.kind == 2 && delta_d > GetSDVmindis(i, prevCarIndex, currentRoad) {

				newSpeed = currentRoad[i].speed + 1
				newLight = 1
				newAccel = 1

			} else if Checktrain(i) == true && delta_d == GetSDVmindis(i, prevCarIndex, currentRoad) && currentRoad[GetPrex(currentRoad, prevCarIndex)].accel == 1 {
				newSpeed = currentRoad[i].speed + 1
				newLight = 1
				newAccel = 1

			} else if prevCar.kind == 1 && delta_d <= safeSpaceMin[currentRoad[i].speed] {
				newSpeed = currentRoad[i].speed - 1
				newLight = -1
				newAccel = 0

			} else if prevCar.kind == 2 && Checktrain(i) == false && delta_d <= GetSDVmindis(i, prevCarIndex, currentRoad) {
				newSpeed = currentRoad[i].speed - 1
				newLight = -1
				newAccel = 0
			} else if Checktrain(i) == true && currentRoad[GetPrev(currentRoad, prevCarIndex)].light == -1 {
				newSpeed = currentRoad[i].speed - 1
				newLight = -1
				newAccel = 0
			} else {
				newLight = 0
				newAccel = 0
			}

			newIndex := i + newSpeed
			if newIndex < roadLength && newRoad[newIndex].kind != 0 {
				panic("SDV crashes something.")
			} else if newIndex >= roadLength {
				carCnt++
			} else {
				newRoad[newIndex].speed = newSpeed
				newRoad[newIndex].light = newLight
				newRoad[newIndex].accel = newAccel
				newRoad[newIndex].kind = currentRoad[i].kind
			}

		}
	}

	return newRoad
}

func GetSDVmindis(k, m int, currentroad Road) int {
	var BrakingDistance int
	vm := currentroad[m].speed
	va := currentroad[k].speed
	maxv := max(vm-2, 0)
	BrakingDistance = safeSpaceMin[maxv]
	maxva := max(safeSpaceMin[va]-BrakingDistance+1, 1)

	return maxva

}

func max(k, m int) int {
	if k > m {
		return k
	} else {
		return m
	}
}
