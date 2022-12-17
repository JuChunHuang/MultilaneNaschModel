package main

import (
	"C"
	"math"
	"math/rand"
	"time"
)

//export initialSingleLane
func initialSingleLane(incidentPos []int, NSDVdensity float64) Road {
	var loc int
	var NSDVnumber int
	var posList []int

	initialRoad := make(Road, roadLength)

	// set the incident location
	if len(incidentPos) > 0 {
		for i := range incidentPos {
			loc = incidentPos[i]
			initialRoad[loc].kind = 3
			initialRoad[loc].backlight = 1
		}
	}

	// generate NSDV cars randomly
	NSDVnumber = int(math.Floor(NSDVdensity * roadLength))
	posList = RandomPosGenerate(NSDVnumber, roadLength)
	for i := range posList {
		k := posList[i]
		initialRoad[k].kind = 1
		initialRoad[k].speed = 1
		initialRoad[k].backlight = 0
	}

	return initialRoad
}

// Randomly generate n positions on the road
//
//export RandomPosGenerate
func RandomPosGenerate(n, totalLength int) []int {
	var p int
	var posList []int
	for i := 0; i < n; i++ {
		rand.Seed(time.Now().UnixNano())
		p = rand.Intn(totalLength)
		for CheckRepeat(p, posList) == false {
			p = rand.Intn(totalLength)
		}
		posList = append(posList, p)
	}
	return posList
}

// CheckRepeat checks if there's repeat in n randomly generated numbers
//
//export CheckRepeat
func CheckRepeat(k int, list []int) bool {
	n := len(list)
	for i := 0; i < n; i++ {
		if list[i] == k {
			return false
		}
	}
	return true
}

//export initialMultiRoad
func initialMultiRoad(incidentPos [][]int, NSDVdensity float64, laneNum int) MultiRoad {
	initialRoad := make(MultiRoad, laneNum)

	//generate laneNum number of initial SingleRoad
	for i := range initialRoad {
		initialRoad[i] = initialSingleLane(incidentPos[i], NSDVdensity)
	}

	return initialRoad
}
