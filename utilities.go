package main

import (
	"C"
	"math/rand"
)
import (
	"time"
)

func GetPrevCar(currentRoad Road, index int) int {
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

	return -100
}

func GetSDVmindis(k, m int, currentroad Road) int {
	var BrakingDistance int
	vm := currentroad[m].speed
	vk := currentroad[k].speed
	maxv := max(vk-vm, 0)
	BrakingDistance = maxv
	maxva := max(safeSpaceSDVMin[maxv]-BrakingDistance+1, 2)

	return maxva

}

func max(k, m int) int {
	if k > m {
		return k
	} else {
		return m
	}
}

func min(k, m int) int {
	if k < m {
		return k
	} else {
		return m
	}
}

// Whether the lane is exist
func ValidLane(lane, laneNum int) bool {
	return (lane >= 0 && lane < laneNum)
}

// Produce generate cars at the first grid according to car type possibility
// Input: a pointer to a Road object of currentRoad, a float object representing the possibility to produce NSDV
// Output: a boolean object
func Produce(currentRoad *Road, kindPossiblity float64) bool {
	// Determine the kind of next car based on kind possibility
	rand.Seed(time.Now().UnixNano())
	p := rand.Float64()
	var kind int
	var initSpeedBound int
	if p < kindPossiblity {
		kind = 2
	} else {
		kind = 1
	}

	// determin the speed range that the car can obtain
	if (*currentRoad)[0].kind == 0 {
		initSpeedBound = 0
		prevCar := GetPrevCar((*currentRoad), 0)

		// if no car before
		if prevCar > roadLength {
			initSpeedBound = maxSpeed + 1
		} else if (kind == 1) || (kind == 2 && (*currentRoad)[prevCar].kind == 1) {
			// if the new car is a NSDV or the new car is SDV and prevCar is NSDV
			for i := maxSpeed - 1; i >= 0; i-- {
				if prevCar > safeSpaceMin[i] {
					initSpeedBound = i
					break
				}
			}
		} else if kind == 2 && (*currentRoad)[prevCar].kind == 2 {
			// if the new car is a SDV and prevCar is a SDV

			for i := maxSpeed - 1; i >= 0; i-- {
				minSpace := max(safeSpaceMin[i]-safeSpaceMin[(*currentRoad)[prevCar].speed]+2*(*currentRoad)[prevCar].speed+1, 1)
				if prevCar > minSpace {
					initSpeedBound = i
					break
				}

			}
		}
	}

	if initSpeedBound <= 0 {
		// no car produced
		return false
	} else {
		rand.Seed(time.Now().UnixNano())
		(*currentRoad)[0].speed = rand.Intn(initSpeedBound)
		(*currentRoad)[0].kind = kind
		(*currentRoad)[0].backlight = 0
	}
	return true
}

func ProduceMulti(currentRoads *MultiRoad, kindPossiblity float64) bool {
	count := 0
	for i := range *currentRoads {
		a := Produce(&((*currentRoads)[i]), kindPossiblity)
		if a == false {
			count += 1
		}
	}

	if count == 5 {
		return false
	} else {
		return true
	}

}

func GetTrainHead(road Road, carIndex int) int {
	trainHeadIndex := carIndex
	index := CheckPreviousTrain(road, carIndex)

	if index == 0 {
		return trainHeadIndex
	} else {
		for i := carIndex + 1; i < roadLength; i++ {
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

func GetTrainTail(road Road, carIndex int) int {
	trainTailIndex := carIndex
	index := CheckNextTrain(road, carIndex)

	if index == 0 {
		return trainTailIndex
	} else {
		for i := carIndex - 1; i >= 0; i-- {
			if road[i].kind == 0 {
				trainTailIndex--

			} else if road[i].kind == 2 {
				index -= 1
				if index == 0 {
					return trainTailIndex
				}
			}
		}
		return trainTailIndex
	}

}

// Only checktrain if the car is an SDV.
func CheckTrain(road Road, carIndex int) int {
	var sum int
	sum = 1 + CheckPreviousTrain(road, carIndex) + CheckNextTrain(road, carIndex)
	return sum
}

func CheckPreviousTrain(road Road, carIndex int) int {
	var sum int

	prevIndex := GetPrevCar(road, carIndex)
	if prevIndex > roadLength {
		return sum
	} else {

		distance := prevIndex - carIndex
		Dmin := GetSDVmindis(carIndex, prevIndex, road)
		// only count the previous SDV car into consideration when the distance bewteen
		// previous SDV and current SDV is equal to SDV min distance
		if road[prevIndex].kind == 2 && distance <= Dmin {
			sum += 1 + CheckPreviousTrain(road, prevIndex)
		} else {
			return sum
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
		distance := carIndex - nextIndex
		Dmin := GetSDVmindis(nextIndex, carIndex, road)

		// only count the next SDV car into consideration when the distance bewteen
		// next SDV and current SDV is equal to SDV min distance
		if road[nextIndex].kind == 2 && distance <= Dmin {
			sum += 1 + CheckNextTrain(road, nextIndex)
		} else {
			return sum
		}
	}

	return sum
}
