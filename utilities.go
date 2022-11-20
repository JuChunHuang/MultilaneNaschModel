package main

import (
	"math/rand"
	"time"
)

func GetPrev(currentRoad Road, index int) int {
	for c := index + 1; c < roadLength; c++ {
		if currentRoad[c].kind != 0 {
			return c
		}
	}

	return 2 * roadLength
}

func GetNext(currentRoad Road, index int) int {
	for c := index - 1; c >= 0; c-- {
		if currentRoad[c].kind != 0 {
			return c
		}
	}

	return -1
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

func GenRandom(n int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	rn := r.Intn(n)
	return rn
}

// Whether the lane is exist
func ValidLane(lane int) bool {
	return (lane >= 0 && lane < laneNum)
}

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
		(*currentRoad)[0].backlight = 0
		(*currentRoad)[0].accel = 0
	}
	return true
}

func ProduceMulti(currentRoads *MultiRoad, kindPossiblity float64) bool {
	for i := range *currentRoads {
		a := Produce(&((*currentRoads)[i]), kindPossiblity)
		if a == false {
			return false
		}
	}

	return true
}

func GetTrainHead(road Road, carIndex int) int {
	trainHeadIndex := carIndex
	index := CheckPreviousTrain(road, carIndex)

	if index == 0 {
		return trainHeadIndex
	} else {
		for i := carIndex + 1; i < roadLength-1; i++ {
			if road[i].kind == 0 {
				trainHeadIndex++

			} else if road[i].kind == 2 {
				index -= 1
				if index == 0 {
					return trainHeadIndex
				}
			}
		}
		return trainHeadIndex
	}

}

// Only checktrain if the car is an SDV.
func CheckTrain(road Road, carIndex int) bool {
	var sum int
	sum = 1 + CheckPreviousTrain(road, carIndex) + CheckNextTrain(road, carIndex)

	if sum >= 3 {
		return true
	} else {
		return false
	}
}

func CheckPreviousTrain(road Road, carIndex int) int {
	var sum int

	prevIndex := GetPrev(road, carIndex)
	if prevIndex > roadLength {
		return sum
	} else {
		if road[prevIndex].kind == 1 {
			return sum
		} else if road[prevIndex].kind == 2 {
			sum += 1 + CheckPreviousTrain(road, prevIndex)
		}
	}

	return sum
}

func CheckNextTrain(road Road, carIndex int) int {
	var sum int

	nextIndex := GetNext(road, carIndex)
	if nextIndex < 0 {
		return sum
	} else {
		if road[nextIndex].kind == 1 {
			return sum
		} else if road[nextIndex].kind == 2 {
			sum += 1 + CheckNextTrain(road, nextIndex)
		}
	}

	return sum
}