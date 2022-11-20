package main

import (
	"fmt"
	"gifhelper"
)

func main() {
	initialRoad := make(Road, roadLength)
	for i := range initialRoad {
		initialRoad[i].accel = 0
		initialRoad[i].kind = 0
		initialRoad[i].backlight = 0
		initialRoad[i].speed = 0
	}

	numGens := 1000
	cellWidth := 50
	timePoints := PlayNaschModel(initialRoad, numGens)
	fmt.Println("finish")
	imageList := BoardsToImages(timePoints, cellWidth)
	gifhelper.ImagesToGIF(imageList, "cars")
}

func PlayNaschModel(initialRoad Road, numGens int) []Road {
	roads := make([]Road, numGens+1)
	roads[0] = initialRoad
	for i := 1; i <= numGens; i++ {
		roads[i] = SingleLaneSimulation(roads[i-1])
		//fmt.Println(roads[i])
	}

	return roads
}
