package main

import (
	"C"
	"fmt"
)

func main() {
	var numGens int
	var cellWidth int
	var trafficLightPos int
	var trafficLightTime []int
	var trafficLightLane []int

	// GIF generation settings
	numGens = 1000

	// drawing settings
	cellWidth = 21

	// SingleLane ==============================================================

	// // set traffic lights position and time
	// trafficLightPos = roadLength / 2
	// trafficLightTime = make([]int, 3)
	// trafficLightTime[0] = 30 // red light
	// trafficLightTime[1] = 5  // yellow light
	// trafficLightTime[2] = 30 // green light

	// // initialize single-lane
	// initialSingleRoad := initialSingleLane(trafficLightPos)

	// // play NaschModel
	// timePointsSingle := PlaySingleLaneModel(initialSingleRoad, numGens, trafficLightPos, trafficLightTime)
	// fmt.Println("Finish running single lane model!")

	// // converting results from playing single lane model to MultiRoad type
	// var timePoints []MultiRoad
	// timePoints = make([]MultiRoad, numGens+1)
	// for i := range timePoints {
	// 	timePoints[i] = make(MultiRoad, 1)
	// 	timePoints[i][0] = timePointsSingle[i]
	// }

	// //generate SingleRoad pattern
	// DrawBoardSingle(timePoints, numGens, "SingleRoadPattern.png")
	// fmt.Println("Finish running single lane pattern!")

	// // generae GIF for singlelane results
	// imageList := BoardsToImages(timePoints, cellWidth)
	// gifhelper.ImagesToGIF(imageList, "SingleLane")
	// fmt.Println("Finish drawing single lane model results!")

	// MultiLane ==============================================================

	// set traffic lights position
	trafficLightLane = []int{1, 3}
	trafficLightPos = roadLength / 2
	trafficLightTime = make([]int, 3)
	trafficLightTime[0] = 30 // red light
	trafficLightTime[1] = 5  // yellow light
	trafficLightTime[2] = 30 // green light

	// set SDVs percentage
	sdvPercentage := 0.5
	nsdvPercentage := 1.0 - sdvPercentage

	// set lane number
	laneNum := 5

	// initialize multiple roads
	initialMultiRoad := initialMultiRoad(trafficLightLane, trafficLightPos, laneNum)

	// play NaschModel
	timePointsMulti, totalCnt := PlayMultiLaneModel(initialMultiRoad, numGens, trafficLightPos, laneNum, trafficLightLane, trafficLightTime, nsdvPercentage)
	fmt.Println("Finish running multilane model!")
	//DrawBoardMulti(timePointsMulti, numGens, laneNum, "MultiRoadNSDVPattern.png")
	//fmt.Println("Finish running multilane pattern!")

	// generae GIF for multilane results
	imageListMul := BoardsToImages(timePointsMulti, cellWidth)
	ImagesToGIF(imageListMul, "Multilane")
	fmt.Println("Finish drawing multi lane model results!")
	fmt.Println(totalCnt)

}

// PlaySingleLaneModel takes the initial road, run Nasch Model for numGens times and return the results
// Input: a Road object initialRoad, an int object numGens, an int object of traffic light position, a slice of int of time for each status in a light cycle
// Output: a slice of Road objects of length numGens+1
//
//export PlaySingleLaneModel
func PlaySingleLaneModel(initialRoad Road, numGens, lightPos int, trafficLightTime []int) []Road {
	roads := make([]Road, numGens+1)
	roads[0] = initialRoad
	oneRound := trafficLightTime[0] + trafficLightTime[1] + trafficLightTime[2]

	for i := 1; i <= numGens; i++ {
		//Set traffic light status for each generation
		t := i % oneRound
		if 1 <= t && t <= 30 {
			roads[i-1][lightPos].kind = 3 // red light
		} else if 30 < t && t <= 60 {
			roads[i-1][lightPos].kind = 5 // green light
		} else {
			roads[i-1][lightPos].kind = 4 // yellow light
		}
		//Run Nasch Model
		roads[i] = SingleLaneSimulation(roads[i-1])
	}

	return roads
}

// PlayMultiLaneModel takes the initial Multiroad, run Nasch Model for numGens times and return the results
// Input: a Road object initialRoad, an int object numGens, an int object of traffic light position, a slice of int of time for each status in a light cycle
// Output: a slice of MultiRoad objects of length numGens+1
//
//export PlayMultiLaneModel
func PlayMultiLaneModel(initialRoad MultiRoad, numGens, lightPos, laneNum int, lightLane, trafficLightTime []int, nsdvPercentage float64) ([]MultiRoad, int) {
	roads := make([]MultiRoad, numGens+1)
	roads[0] = initialRoad
	//oneRound := trafficLightTime[0] + trafficLightTime[1] + trafficLightTime[2]
	totalCnt := 0
	cnt := 0

	if lightPos > 0 {
		//Set traffic light status for each generation
		for i := 1; i <= numGens; i++ {
			// t := i % oneRound
			// for _, val := range lightLane {
			// 	if 1 <= t && t <= 30 {
			// 		roads[i-1][val][lightPos].kind = 3 // red light
			// 	} else if 30 < t && t <= 60 {
			// 		roads[i-1][val][lightPos].kind = 5 // green light
			// 	} else {
			// 		roads[i-1][val][lightPos].kind = 4 // yellow light
			// 	}
			// }

			roads[i], cnt = MultiLaneSimulation(roads[i-1], i, laneNum, nsdvPercentage)
			totalCnt += cnt
		}
	} else {
		for i := 1; i <= numGens; i++ {
			roads[i], cnt = MultiLaneSimulation(roads[i-1], i, laneNum, nsdvPercentage)
			totalCnt += cnt
		}
	}

	return roads, totalCnt
}
