package main

import (
	"math/rand"
)

// SingleLaneSimulation update the single Road
// Input: a Road object
// Output: a Road Object
func SingleLaneSimulation(currentRoad Road) Road {
	var prevCar Car
	var prevCarIndex int
	var newSpeed int
	var newLight int
	var newAccel int
	var probOfDecel float64
	var prevLight Car

	// make a new road
	newRoad := make(Road, roadLength)
	carCnt := 0

	//produce cars at the beginning the the road
	Produce(&currentRoad, 0.5)

	for i := roadLength - 1; i >= 0; i-- {
		currentCar := currentRoad[i]
		kind := currentCar.kind
		speed := currentCar.speed
		prevCarIndex = GetPrevCar(currentRoad, i)
		prevLightIndex := GetPrevLight(currentRoad, i)

		// if we do not have a previous car
		if prevCarIndex >= roadLength {
			prevCar.backlight = -2
		} else {
			prevCar = currentRoad[prevCarIndex]
		}

		// if we do not have a previous light
		if prevLightIndex >= roadLength {
			prevLight.kind = 0
		} else {
			prevLight = currentRoad[prevLightIndex]
		}

		// get the distance between the previous and the current car
		delta_d := prevCarIndex - i
		// get the distance between the traffic light and the current car
		deltaDLight := prevLightIndex - i

		// if the current car is an NSDV
		if kind == 1 {
			// Deceleration conditions
			if prevCar.backlight == -1 && delta_d > safeSpaceMin[speed] && delta_d < safeSpaceMin[speed] &&
				(deltaDLight > safetraffic[speed] || deltaDLight < 0) {
				probOfDecel = p1
			} else if prevCar.backlight >= 0 && delta_d > safeSpaceMin[speed] && delta_d < safeSpaceMin[speed] &&
				(deltaDLight > safetraffic[speed] || deltaDLight < 0) {
				probOfDecel = p2
			} else if speed == 0 {
				probOfDecel = p3
			} else {
				probOfDecel = 0
			}

			// get a threshold possibility to decide whether we should decelerate
			thresToDecel := rand.Float64()

			if probOfDecel < thresToDecel {
				if speed < maxSpeed && delta_d > safeSpaceMin[speed] {
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

			// stop in front of a red light or a yellow light
			if deltaDLight >= 0 && deltaDLight <= safetraffic[speed] && (prevLight.kind == 3 || prevLight.kind == 4) {
				newSpeed = 0
				newLight = -1
			}

			// avoid cases when the speed exceeds the limits
			if newSpeed < 0 {
				newSpeed = 0
			} else if newSpeed > 10 {
				newSpeed = 10
			}

			// new position for the current car
			newIndex := i + newSpeed
			if newIndex >= roadLength {
				// count the cars which pass the entire road
				carCnt++
			} else if newIndex < roadLength && newRoad[newIndex].kind != 0 {
				// fmt.Println("NSDV crashes something", newIndex, newRoad[newIndex].kind)
				// panic("NSDV crashes something.")
			} else {
				newRoad[newIndex].speed = newSpeed
				newRoad[newIndex].kind = kind
				newRoad[newIndex].backlight = newLight
			}
			// if the current car is an SDV and has no traffic light ahead
		} else if kind == 2 && prevLight.kind == 0 {
			// if we do not have a previous car, we can accelerate
			if prevCarIndex > roadLength {
				newSpeed = speed + 1
				newLight = 1
				newAccel = 1
			} else {
				// for other condtions that we can accelerate
				if delta_d >= safeSpaceSDVMin[speed] {
					newSpeed = speed + 1
					newLight = 1
					newAccel = 1
				}
				if prevCar.kind == 1 && prevCar.backlight != -1 && delta_d >= safeSpaceMin[speed] {
					newSpeed = speed + 1
					newLight = 1
					newAccel = 1
				} else if prevCar.kind == 2 && delta_d > GetSDVmindis(i, prevCarIndex, currentRoad) {
					newSpeed = delta_d - safeSpaceSDVMin[speed] - 1
					newLight = 1
					newAccel = 1
					// if the previous car is an SDV, we then check if this forms an SDV-train
				} else if prevCar.kind == 2 && delta_d <= GetSDVmindis(i, prevCarIndex, currentRoad) && CheckTrain(currentRoad, i) == true {
					trainHead := GetTrainHead(currentRoad, i)
					// if delta_d > GetSDVmindis(i, prevCarIndex, currentRoad) {
					// 	panic("not SDV Train")
					// }
					newSpeed = currentRoad[trainHead].speed
					newLight = currentRoad[trainHead].backlight
					newAccel = currentRoad[trainHead].accel
					// slow down if there is a red light or yellow light ahead or a car ahead and we are getting close
				} else if prevLight.kind < 5 && deltaDLight <= safeSpaceMin[0] {
					newSpeed = speed - 1
					newLight = -1
					newAccel = 0
				} else {
					newLight = 0
					newAccel = 0
				}

			}

			// avoid cases when the speed exceeds the limits
			if newSpeed < 0 {
				newSpeed = 0
			} else if newSpeed > 10 {
				newSpeed = 10
			}
			// new position for the current car
			newIndex := i + newSpeed

			// count the cars which pass the entire road
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
			// if current car is an SDV and there is a red light or yellow light ahead
		} else if kind == 2 && (prevLight.kind == 3 || prevLight.kind == 4) {
			// if the previous car went over the traffic light
			if prevCarIndex > prevLightIndex {
				delta_d = prevLightIndex - i
				prevCarIndex = prevLightIndex
			}
			if delta_d >= safeSpaceMin[speed] {
				newSpeed = speed + 1
				newLight = 1
				newAccel = 1
			}
			if prevCar.kind == 1 && prevCar.backlight != -1 && delta_d >= safeSpaceMin[speed] {
				newSpeed = speed + 1
				newLight = 1
				newAccel = 1
			} else if prevCar.kind == 2 && delta_d > GetSDVmindis(i, prevCarIndex, currentRoad) {
				newSpeed = delta_d - safeSpaceSDVMin[speed] - 1
				newLight = 1
				newAccel = 1
			} else if prevCar.kind == 2 && delta_d <= GetSDVmindis(i, prevCarIndex, currentRoad) && CheckTrain(currentRoad, i) == true {
				trainHead := GetTrainHead(currentRoad, i)
				// if delta_d > GetSDVmindis(i, prevCarIndex, currentRoad) {
				// 	panic("not SDV Train")
				// }
				newSpeed = currentRoad[trainHead].speed
				newLight = currentRoad[trainHead].backlight
				newAccel = currentRoad[trainHead].accel
				// stop in front of a red light or a yellow light
			} else if prevLight.kind >= 3 && deltaDLight <= safeSpaceMin[0] {
				newSpeed = 0
				newLight = 0
				newAccel = 0
			} else {
				newLight = 0
				newAccel = 0
			}

			// if delta_d <= safetraffic[1] {
			// 	newSpeed = 0
			// 	newLight = -1
			// }

			var newIndex int

			// we do not allow the car goes over the traffic light
			if newIndex > roadLength/2 {
				newSpeed = 0
			}
			// avoid cases when the speed exceeds the limits
			if newSpeed < 0 {
				newSpeed = 0
			} else if newSpeed > 10 {
				newSpeed = 10
			}
			newIndex = i + newSpeed

			// count the cars which pass the entire road
			if newIndex >= roadLength {
				carCnt++
			} else {
				// change the status of the road cells
				newRoad[newIndex].speed = newSpeed
				newRoad[newIndex].backlight = newLight
				newRoad[newIndex].accel = newAccel
				newRoad[newIndex].kind = kind
			}

			// if current car is an SDV and there is a green light ahead
		} else if kind == 2 && prevLight.kind == 5 {
			// if the previous car is far away or we do not have a previous car
			if delta_d >= safeSpaceSDVMin[speed] {
				newSpeed = speed + 1
				newLight = 1
				newAccel = 1
			}
			// other conditions for acceleration
			if prevCar.kind == 1 && prevCar.backlight != -1 && delta_d >= safeSpaceMin[speed] {
				newSpeed = speed + 1
				newLight = 1
				newAccel = 1
			} else if prevCar.kind == 2 && delta_d > GetSDVmindis(i, prevCarIndex, currentRoad) {
				newSpeed = delta_d - safeSpaceSDVMin[speed] - 1
				newLight = 1
				newAccel = 1
			} else if prevCar.kind == 2 && delta_d <= GetSDVmindis(i, prevCarIndex, currentRoad) && CheckTrain(currentRoad, i) == true {
				trainHead := GetTrainHead(currentRoad, i)
				// if delta_d > GetSDVmindis(i, prevCarIndex, currentRoad) {
				// 	panic("not SDV Train")
				// }
				newSpeed = currentRoad[trainHead].speed
				newLight = currentRoad[trainHead].backlight
				newAccel = currentRoad[trainHead].accel
				// slow down in front of a red light or yellow light
			} else if prevLight.kind < 5 && deltaDLight <= safeSpaceMin[0] {
				newSpeed = speed - 1
				newLight = -1
				newAccel = 0
			} else {
				newLight = 0
				newAccel = 0
			}

			// avoid cases when the speed exceeds the limits
			if newSpeed < 0 {
				newSpeed = 0
			} else if newSpeed > 10 {
				newSpeed = 10
			}
			newIndex := i + newSpeed

			// count the cars which pass the entire road
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

		}
	}
	// fmt.Println(carCnt)

	return newRoad
}
