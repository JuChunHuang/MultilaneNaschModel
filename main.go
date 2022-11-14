package main

import (
	"gifhelper"
)

func main() {
	initialRoad := make(Road, roadLength)
	for i := range initialRoad {
		initialRoad[i].accel = GenRandom(2)
		initialRoad[i].kind = GenRandom(3)
		initialRoad[i].light = -1 + GenRandom(3)
		initialRoad[i].speed = GenRandom(10)
	}

	numGens := 300
	cellWidth := 20
	timePoints := PlayNaschModel(initialRoad, numGens)
	imageList := BoardsToImages(timePoints, cellWidth)
	gifhelper.ImagesToGIF(imageList, "prisoners")
}
