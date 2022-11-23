package main

import (
	"fmt"
	"gifhelper"
)

func main() {

	// initialRoad := make(MultiRoad, laneNum)
	// for i := range initialRoad {
	// 	initialRoad[i] = make(Road, roadLength)
	// }

	initialRoad := make(Road, roadLength)

	// Set Traffic light
	// trafficLightPos := roadLength / 2
	// for i := 0; i < laneNum; i++ {
	// 	initialRoad[i][trafficLightPos].kind = 3
	// }

	trafficLightPos := roadLength / 2
	initialRoad[trafficLightPos].kind = 3

	trafficLightTime := make([]int, 3)
	trafficLightTime[0] = 30 // red light
	trafficLightTime[1] = 5  // yellow light
	trafficLightTime[3] = 30 // green light
	numGens := 1000
	cellWidth := 50

	// timePoints := PlayMultiLaneModel(initialRoad, numGens, trafficLightPos, trafficLightTime)
	// fmt.Println("finish")
	// fmt.Println(timePoints)

	timePointsSingle := PlaySingleLaneModel(initialRoad, numGens, trafficLightPos, trafficLightTime)
	fmt.Println("finish")
	fmt.Println(timePointsSingle)
	var timePoints []MultiRoad
	timePoints = make([]MultiRoad, numGens+1)
	for i := range timePoints {
		timePoints[i] = make(MultiRoad, 1)
		timePoints[i][0] = timePointsSingle[i]
	}

	imageList := BoardsToImages(timePoints, cellWidth)
	gifhelper.ImagesToGIF(imageList, "cars")
}

func PlaySingleLaneModel(initialRoad Road, numGens, lightPos int, trafficLightTime []int) []Road {
	roads := make([]Road, numGens+1)
	roads[0] = initialRoad
	oneRound := trafficLightTime[0] + trafficLightTime[1] + trafficLightTime[2]
	for i := 1; i <= numGens; i++ {
		t := i % oneRound
		if 1 <= t && t <= 30 {
			for k := 0; k < laneNum; k++ {
				initialRoad[lightPos].kind = 3 // red light
			}
		} else if 30 < t && t <= 35 {
			for k := 0; k < laneNum; k++ {
				initialRoad[lightPos].kind = 4 // yellow light
			}
		} else {
			initialRoad[lightPos].kind = 5 // green light
		}
		roads[i] = SingleLaneSimulation(roads[i-1])
		//fmt.Println(roads[i])
	}

	return roads
}

func PlayMultiLaneModel(initialRoad MultiRoad, numGens, lightPos int, trafficLightTime []int) []MultiRoad {
	roads := make([]MultiRoad, numGens+1)
	roads[0] = initialRoad
	oneRound := trafficLightTime[0] + trafficLightTime[1] + trafficLightTime[2]
	for i := 1; i <= numGens; i++ {
		t := i % oneRound
		if 1 <= t && t <= 30 {
			for k := 0; k < laneNum; k++ {
				initialRoad[k][lightPos].kind = 3 // red light
			}
		} else if 30 < t && t <= 35 {
			for k := 0; k < laneNum; k++ {
				initialRoad[k][lightPos].kind = 4 // yellow light
			}
		} else {
			initialRoad[k][lightPos].kind = 5 // green light
		}
		roads[i] = MultiLaneSimulation(roads[i-1])

	}

	return roads
}
