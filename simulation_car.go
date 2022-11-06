package main

import (
	"math/rand"
)

func PlayNaschModel(initialRoad Road, numGens int) []Road {
	roads := make([]Road, numGens+1)
	roads[0] = initialRoad
	for i := 1; i <= numGens; i++ {
		roads[i] = UpdateRoad(roads[i-1])
	}

	return roads
}

func UpdateRoad(currentRoad Road) Road {
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

			} else if Checktrain(i) == true && delta_d == GetSDVmindis(i, prevCarIndex, currentRoad) && currentRoad[GetPrev(currentRoad, prevCarIndex)].accel == 1 {
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

func GetPrev(currentRoad Road, index int) int {
	for c := index + 1; c < roadLength; c++ {
		if currentRoad[c].kind != 0 {
			return c
		}
	}

	return 2 * roadLength
}

func GetNext(currentRoad Road, index int) int {
	for c := index - 1; c <= 0; c-- {
		if currentRoad[c].kind != 0 {
			return c
		}
	}

	return 0
}

func Produce(currentRoad Road, kindPossiblity float64) bool {

	// Determine the kind of next car based on kind possibility
	p := rand.Float64()
	var kind int
	var initSpeedBound int
	if p < kindPossiblity {
		kind = 1
	} else {
		kind = 2
	}

	// determin the speed range that the car can obtain
	if currentRoad[0].kind == 0 {
		initSpeedBound = 0
		prevCar := GetPrev(currentRoad, 0)

		// if no car before
		if prevCar > roadLength {
			initSpeedBound = maxSpeed
		} else if (kind == 1) || (kind == 2 && currentRoad[prevCar].kind == 1) {
			// if the new car is a NSDV or the new car is SDV and prevCar is NSDV
			for i := 0; i < maxSpeed; i++ {
				if prevCar < safeSpaceMin[i] {
					break
				}
				initSpeedBound = i - 1
			}
		} else if kind == 2 && currentRoad[prevCar].kind == 2 {
			// if the new car is a SDV and prevCar is a SDV
			for i := 0; i < maxSpeed; i++ {
				minSpace := safeSpaceMin[i] - safeSpaceMin[currentRoad[prevCar].speed] + 2*currentRoad[prevCar].speed + 1
				if prevCar < minSpace {
					break
				}
				initSpeedBound = i - 1
			}
		}
	}

	if initSpeedBound <= 0 {
		// no car produced
		return false
	} else {
		currentRoad[0].speed = 1 + rand.Intn(initSpeedBound)
		currentRoad[0].kind = kind
		currentRoad[0].light = 0
		currentRoad[0].accel = 0
	}
	return true
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
