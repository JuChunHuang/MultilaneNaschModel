package main

import (
	"math/rand"
	"time"
)

// PlayNaschModel takes the intitial Road, run the model for numGens times, and get road configuration for each time
// Input: a Road object, an int object
// Output: a slice of Roads of numGens length
func PlayNaschModel(initialRoad Road, numGens int) []Road {
	var roads []Road
	roads = make([]Road, numGens+1)
	// give the initial road
	roads[0] = initialRoad

	// update the road for numGens times
	for i := 1; i <= numGens; i++ {
		roads[i] = UpdateRoad(roads[i-1])
	}

	return roads
}

// UpdateRoad takes the currentRoad and update each grid
// Input: a Road object
// Output: a Road object
func UpdateRoad(currentRoad Road) Road {
	var prevCar Car
	var prevCarIndex int
	var newSpeed int
	var newLight int
	var newAccel int
	var probOfDecel float64
	var SDVprobability float64
	var delta_d int

	newRoad := make(Road, roadLength)

	// carCnt tacks the number of cats running across the road
	carCnt := 0

	// generate new cars at the beginning of the road
	SDVprobability = 0.5
	Produce(&currentRoad, SDVprobability)

	// go over each grid and update their status
	for i := range currentRoad {

		// Get the index of previous car
		prevCarIndex = GetPrev(currentRoad, i)
		if prevCarIndex >= roadLength {
			// no cars in front
			prevCar.light = -3
		} else {
			prevCar = currentRoad[prevCarIndex]
		}

		// get the distance between the current grid and the grid containing the previous car
		delta_d = prevCarIndex - i

		// If there's an NSDV car in the current grid
		if currentRoad[i].kind == 1 {
			// the car is a NSDV, change the speed of the car
			currentSpeed := currentRoad[i].speed
			// get the safeSpace of current cat based on its speed
			if prevCar.light == 1 && delta_d > safeSpaceMin[currentSpeed] && delta_d < safeSpaceMax[currentSpeed] {
				probOfDecel = p1
			} else if prevCar.light == 0 && delta_d > safeSpaceMin[currentSpeed] && delta_d < safeSpaceMax[currentSpeed] {
				probOfDecel = p2
			} else if currentSpeed == 0 {
				probOfDecel = p3
			} else {
				probOfDecel = 0
			}

			// Decel the current car stochastically
			thresToDecel := rand.Float64()

			if probOfDecel < thresToDecel {
				if currentSpeed < maxSpeed && delta_d > safeSpaceMax[currentSpeed] {
					// acceleration because no car in front of it
					newSpeed = currentSpeed + 1
					newLight = 1
				} else if prevCar.light == 1 && delta_d > safeSpaceMin[currentSpeed] {
					// acceleration because the front car is accelerated
					newSpeed = currentSpeed + 1
					newLight = 1
				}
			} else {
				if delta_d < safeSpaceMin[currentSpeed] {
					// deceleration
					newSpeed = currentSpeed - 1
					newLight = -1
				} else if delta_d == safeSpaceMin[currentSpeed] {
					// on hold case, speed not changed
					newSpeed = currentSpeed
					newLight = 0
				}
			}

			// Get the grid our current car is going to appear at next time point
			newIndex := i + newSpeed
			if newIndex >= roadLength {
				// car run out of the road
				// add carCnt
				carCnt++
			} else if newIndex < roadLength && newRoad[newIndex].kind != 0 {
				// two cars appear in the same grid
				panic("NSDV crashes something.")
			} else {
				// limit the car speed
				if newSpeed < 0 {
					newSpeed = 0
				} else if newSpeed > 10 {
					newSpeed = 10
				}
				// update the newroad
				newRoad[newIndex].speed = newSpeed
				newRoad[newIndex].kind = currentRoad[i].kind
				newRoad[newIndex].light = newLight
			}
		}

		// If there's a SDV car in the current grid
		if currentRoad[i].kind == 2 {

			currentSpeed := currentRoad[i].speed

			// get the minimum safe distance between current cat and previous car if the previous car is an SDV
			SDVMinDistance := GetSDVmindis(i, prevCarIndex, currentRoad)

			if delta_d >= safeSpaceMax[currentSpeed] {
				// the previous cat is in safe space
				newSpeed = currentSpeed + 1
				newLight = 1
				newAccel = 1

			} else if prevCar.kind == 1 && prevCar.light != -1 && delta_d >= safeSpaceMin[currentSpeed] {
				// the previous car is an SDV with negative accelaration
				newSpeed = currentSpeed + 1
				newLight = 1
				newAccel = 1

			} else if prevCar.kind == 2 && delta_d > SDVMinDistance {
				// the previous car is an NSDV and current distance is larger than safe distance
				newSpeed = currentSpeed + 1
				newLight = 1
				newAccel = 1

			} else if prevCar.kind == 2 && delta_d <= SDVMinDistance {
				// the previous car is an NSDV and current distance is smaller than safe distance
				newSpeed = currentSpeed - 1
				newLight = -1
				newAccel = 0

			} else {
				newLight = 0
				newAccel = 0
			}

			// Get the index of previous car
			newIndex := i + newSpeed
			if newIndex >= roadLength {
				// car ran out of the road
				carCnt++
			} else if newIndex < roadLength && newRoad[newIndex].kind != 0 {
				// two cars appear in the same grid
				panic("SDV crashes something.")
			} else {
				// limit the speed range of cars
				if newSpeed < 0 {
					newSpeed = 0
				} else if newSpeed > 10 {
					newSpeed = 10
				}
				newRoad[newIndex].speed = newSpeed
				newRoad[newIndex].light = newLight
				newRoad[newIndex].accel = newAccel
				newRoad[newIndex].kind = currentRoad[i].kind
			}

		}
	}

	return newRoad
}

// GetPrev get the index nearest car before the current grid of given index
// Input: a Road object, an index
// Output: an int object
func GetPrev(currentRoad Road, index int) int {
	for c := index + 1; c < roadLength; c++ {
		if currentRoad[c].kind != 0 {
			return c
		}
	}

	return 2 * roadLength
}

// GetPrev get the index nearest car after the current grid of given index
// Input: a Road object, an index
// Output: an int object
func GetNext(currentRoad Road, index int) int {
	for c := index - 1; c <= 0; c-- {
		if currentRoad[c].kind != 0 {
			return c
		}
	}

	return 0
}

// Produce generate cars at the begining of the road according to the possibilitiy of SDVs and NSDVs
// Input: a pointer to Road, a float object of possibility
// Output: a bool object
func Produce(currentRoad *Road, kindPossiblity float64) bool {
	// Determine the kind of next car based on kind possibility
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	p := r.Float64()
	var kind int
	var initSpeedBound int
	if p < kindPossiblity {
		kind = 1
	} else {
		kind = 2
	}

	// determin the speed range that the car can obtain
	if (*currentRoad)[0].kind == 0 {
		initSpeedBound = 0
		prevCar := GetPrev((*currentRoad), 0)

		// if no car before
		if prevCar > roadLength {
			initSpeedBound = maxSpeed
		} else if (kind == 1) || (kind == 2 && (*currentRoad)[prevCar].kind == 1) {
			// if the new car is a NSDV or the new car is SDV and prevCar is NSDV
			for i := 0; i < maxSpeed; i++ {
				if prevCar < safeSpaceMin[i] {
					initSpeedBound = i - 1
					break
				}
			}
		} else if kind == 2 && (*currentRoad)[prevCar].kind == 2 {
			// if the new car is a SDV and prevCar is a SDV
			for i := 0; i < maxSpeed; i++ {
				minSpace := safeSpaceMin[i] - safeSpaceMin[(*currentRoad)[prevCar].speed] + 2*(*currentRoad)[prevCar].speed + 1
				if prevCar < minSpace {
					initSpeedBound = i - 1
					break
				}
			}
		}
	}

	if initSpeedBound <= 0 {
		// no car produced
		return false
	} else {
		(*currentRoad)[0].speed = 1 + rand.Intn(initSpeedBound)
		(*currentRoad)[0].kind = kind
		(*currentRoad)[0].light = 0
		(*currentRoad)[0].accel = 0
	}
	return true
}

// GetSDVmindis get the minimum safe distance between SDVs
// Input: a Road Object, two int objects
// Output: a int object
func GetSDVmindis(k, m int, currentroad Road) int {
	var BrakingDistance int
	vm := currentroad[m].speed
	va := currentroad[k].speed
	maxv := max(vm-2, 0)
	BrakingDistance = safeSpaceMin[maxv]
	maxva := max(safeSpaceMin[va]-BrakingDistance+1, 1)

	return maxva

}

// max return the maximum of two integars
// Input: two int objects
// Output: a int object

func max(k, m int) int {
	if k > m {
		return k
	} else {
		return m
	}
}
