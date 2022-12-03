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

	trafficLightPos := roadLength / 2
	initialSingleRoad := initialSingleLane(trafficLightPos)

	// Set Traffic light
	// trafficLightPos := roadLength / 2
	// for i := 0; i < laneNum; i++ {
	// 	initialRoad[i][trafficLightPos].kind = 3
	// }

	trafficLightTime := make([]int, 3)
	trafficLightTime[0] = 30 // red light
	trafficLightTime[1] = 5  // yellow light
	trafficLightTime[2] = 30 // green light
	numGens := 50
	cellWidth := 21

	// timePoints := PlayMultiLaneModel(initialRoad, numGens, trafficLightPos, trafficLightTime)
	// fmt.Println("finish")
	// fmt.Println(timePoints)

	timePointsSingle := PlaySingleLaneModel(initialSingleRoad, numGens, trafficLightPos, trafficLightTime)
	fmt.Println("Finish running single lane model!")

	var timePoints []MultiRoad
	timePoints = make([]MultiRoad, numGens+1)
	for i := range timePoints {
		timePoints[i] = make(MultiRoad, 1)
		timePoints[i][0] = timePointsSingle[i]
	}
	// fmt.Println(timePoints[11][0])
	imageList := BoardsToImages(timePoints, cellWidth)
	gifhelper.ImagesToGIF(imageList, "SingleLane")

	trafficLightLane := []int{1, 2, 3}

	initialMultiRoad := initialMultiRoad(trafficLightLane, trafficLightPos)

	timePointsMulti := PlayMultiLaneModel(initialMultiRoad, numGens, trafficLightPos, trafficLightLane, trafficLightTime)
	fmt.Println("Finish running multilane model!")

	var timePointsMul []MultiRoad
	timePointsMul = make([]MultiRoad, numGens+1)
	for i := range timePointsMul {
		timePointsMul[i] = make(MultiRoad, laneNum)
		for j := range timePointsMulti[i] {
			timePointsMul[i][j] = timePointsMulti[i][j]
		}
	}

	fmt.Println("Finish multiple lane timePoints simulation!")

	imageListMul := BoardsToImages(timePointsMul, cellWidth)
	gifhelper.ImagesToGIF(imageListMul, "Multilane")
}

func PlaySingleLaneModel(initialRoad Road, numGens, lightPos int, trafficLightTime []int) []Road {
	roads := make([]Road, numGens+1)
	roads[0] = initialRoad
	oneRound := trafficLightTime[0] + trafficLightTime[1] + trafficLightTime[2]
	for i := 1; i <= numGens; i++ {
		t := i % oneRound
		if 1 <= t && t <= 30 {
			roads[i-1][lightPos].kind = 3 // red light
		} else if 30 < t && t <= 60 {
			roads[i-1][lightPos].kind = 5 // green light
		} else {
			roads[i-1][lightPos].kind = 4 // yellow light
		}
		roads[i] = SingleLaneSimulation(roads[i-1])
		// fmt.Println(roads[i])
	}

	return roads
}

func PlayMultiLaneModel(initialRoad MultiRoad, numGens, lightPos int, lightLane, trafficLightTime []int) []MultiRoad {

	roads := make([]MultiRoad, numGens+1)
	roads[0] = initialRoad
	oneRound := trafficLightTime[0] + trafficLightTime[1] + trafficLightTime[2]
	for i := 1; i <= numGens; i++ {
		t := i % oneRound
		for _, val := range lightLane {
			if 1 <= t && t <= 30 {
				roads[i-1][val][lightPos].kind = 3 // red light
			} else if 30 < t && t <= 60 {
				roads[i-1][val][lightPos].kind = 3 // green light
			} else {
				roads[i-1][val][lightPos].kind = 3 // yellow light
			}
		}
		roads[i] = MultiLaneSimulation(roads[i-1], i)
		// fmt.Println(roads[i])
	}

	return roads
}
