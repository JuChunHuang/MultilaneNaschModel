package main

import (
	"fmt"
	"os"
	"strconv"
)

func runGofunc(sdvPercentage float64, numGens, laneNum, pos int) int {
	var cellWidth int
	var NSDVdensity float64

	// GIF generation settings

	//drawing settings
	cellWidth = 21

	//SingleLane ==============================================================

	// set beginning NSDV density
	// NSDVdensity = 0.0

	if laneNum == 1 {
		//set incident positions
		incidentPos := make([]int, 0)
		if pos == -1 {
		} else {
			incidentPos[0] = pos
		}
		// initialize single-lane
		initialSingleRoad := initialSingleLane(incidentPos, NSDVdensity)

		// play NaschModel
		timePointsSingle, totalCnt := PlaySingleLaneModel(initialSingleRoad, numGens, sdvPercentage)
		fmt.Println("Finish running single lane model!")

		// converting results from playing single lane model to MultiRoad type
		var timePoints []MultiRoad
		timePoints = make([]MultiRoad, numGens+1)
		for i := range timePoints {
			timePoints[i] = make(MultiRoad, 1)
			timePoints[i][0] = timePointsSingle[i]
		}

		// //generate SingleRoad pattern
		// DrawBoardSingle(timePoints, numGens, "1lane100.png")
		// fmt.Println("Finish running single lane pattern!")

		//generae GIF for singlelane results
		imageList := BoardsToImages(timePoints, cellWidth)
		ImagesToGIF(imageList, "output")
		fmt.Println("Finish drawing single lane model results!")
		return totalCnt
	} else {
		//MultiLane ==============================================================
		// set traffic lights position
		// trafficLightLane = []int{1, 3}
		// trafficLightPos = roadLength / 2
		// trafficLightTime = make([]int, 3)
		// trafficLightTime[0] = 30 // red light
		// trafficLightTime[1] = 5  // yellow light
		// trafficLightTime[2] = 30 // green light

		//set incident positions
		incidentPos := make([][]int, laneNum)
		for i := 0; i < laneNum; i++ {
			incidentPos[i] = make([]int, 0)
		}
		if pos == -1 {
		} else {
			incidentPos[laneNum/2] = append(incidentPos[laneNum/2], pos)
		}

		// initialize multiple roads
		initialMultiRoad := initialMultiRoad(incidentPos, NSDVdensity, laneNum)

		// play NaschModel
		timePointsMulti, totalCnt := PlayMultiLaneModel(initialMultiRoad, numGens, laneNum, sdvPercentage)
		fmt.Println("Finish running multilane model!")
		//fmt.Println(totalCnt)
		//DrawBoardMulti(timePointsMulti, numGens, laneNum, "MultiRoadNSDVPattern.png")
		//fmt.Println("Finish running multilane pattern!")

		// generae GIF for multilane results
		imageListMul := BoardsToImages(timePointsMulti, cellWidth)
		ImagesToGIF(imageListMul, "output")
		fmt.Println("Finish drawing multi lane model results!")
		return totalCnt
	}

}

func main() {
	a, _ := strconv.ParseFloat(os.Args[1], 64) // enter sdv percentage
	b, _ := strconv.Atoi(os.Args[2])           // enter simulation generation
	c, _ := strconv.Atoi(os.Args[3])           // enter lane number
	d, _ := strconv.Atoi(os.Args[4])           // enter incident position
	runGofunc(a, b, c, d)
}

// PlaySingleLaneModel takes the initial road, run Nasch Model for numGens times and return the results
// Input: a Road object initialRoad, an int object numGens, an int object of traffic light position, a slice of int of time for each status in a light cycle
// Output: a slice of Road objects of length numGens+1

func PlaySingleLaneModel(initialRoad Road, numGens int, sdvPercentage float64) ([]Road, int) {
	roads := make([]Road, numGens+1)
	roads[0] = initialRoad
	totalCnt := 0
	cnt := 0

	for i := 1; i <= numGens; i++ {
		roads[i], cnt = SingleLaneSimulation(roads[i-1], sdvPercentage, i)
		totalCnt += cnt
	}
	return roads, totalCnt
}

// PlayMultiLaneModel takes the initial Multiroad, run Nasch Model for numGens times and return the results
// Input: a Road object initialRoad, an int object numGens, an int object of traffic light position, a slice of int of time for each status in a light cycle
// Output: a slice of MultiRoad objects of length numGens+1
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
