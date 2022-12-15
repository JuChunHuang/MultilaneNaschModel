package main

import (
	"fmt"
)

func main() {
	var numGens int
	var cellWidth int
	var NSDVdensity float64
	var sdvPercentage float64

	// GIF generation settings
	numGens = 500

	// drawing settings
	cellWidth = 21

	// //SingleLane ==============================================================

	// // set incident positions

	// // set NSDV density
	// NSDVdensity = 0.07
	sdvPercentage = 0.6

	// // initialize single-lane
	// initialSingleRoad := initialSingleLane(incidentPos, NSDVdensity)

	// // play NaschModel
	// timePointsSingle := PlaySingleLaneModel(initialSingleRoad, numGens, sdvPercentage)
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

	//MultiLane ==============================================================

	// set traffic lights position
	// trafficLightLane = []int{1, 3}
	// trafficLightPos = roadLength / 2
	// trafficLightTime = make([]int, 3)
	// trafficLightTime[0] = 30 // red light
	// trafficLightTime[1] = 5  // yellow light
	// trafficLightTime[2] = 30 // green light

	// set lane number
	laneNum := 5

	incidentPos := make([][]int, laneNum)
	for i := 0; i < laneNum; i++ {
		incidentPos[i] = make([]int, 0)
	}
	incidentPos[2] = append(incidentPos[2], roadLength/2)
	// initialize multiple roads
	initialMultiRoad := initialMultiRoad(incidentPos, NSDVdensity, laneNum)

	// play NaschModel
	timePointsMulti, totalCnt := PlayMultiLaneModel(initialMultiRoad, numGens, laneNum, sdvPercentage)
	fmt.Println("Finish running multilane model!")
	DrawBoardMulti(timePointsMulti, numGens, laneNum, "MultiRoadNSDVPattern.png")
	fmt.Println("Finish running multilane pattern!")

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
func PlaySingleLaneModel(initialRoad Road, numGens int, sdvPercentage float64) []Road {
	roads := make([]Road, numGens+1)
	roads[0] = initialRoad

	for i := 1; i <= numGens; i++ {
		roads[i] = SingleLaneSimulation(roads[i-1], sdvPercentage, i)
	}

	return roads
}

// PlayMultiLaneModel takes the initial Multiroad, run Nasch Model for numGens times and return the results
// Input: a Road object initialRoad, an int object numGens, an int object of traffic light position, a slice of int of time for each status in a light cycle
// Output: a slice of MultiRoad objects of length numGens+1
//
//export PlayMultiLaneModel
func PlayMultiLaneModel(initialRoad MultiRoad, numGens, laneNum int, sdvPercentage float64) ([]MultiRoad, int) {
	roads := make([]MultiRoad, numGens+1)
	roads[0] = initialRoad
	//oneRound := trafficLightTime[0] + trafficLightTime[1] + trafficLightTime[2]
	totalCnt := 0
	cnt := 0

	for i := 1; i <= numGens; i++ {
		roads[i], cnt = MultiLaneSimulation(roads[i-1], i, laneNum, sdvPercentage)
		totalCnt += cnt
	}

	return roads, totalCnt
}
