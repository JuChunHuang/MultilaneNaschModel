package main

func initialSingleLane(trafficLightPos int) Road {
	initialRoad := make(Road, roadLength)

	// set the trafficLight location
	initialRoad[trafficLightPos].kind = 3

	return initialRoad
}

func initialMultiRoad(trafficLightLane []int, trafficLightPos, laneNum int) MultiRoad {
	initialRoad := make(MultiRoad, laneNum)

	//generate laneNum number of initial SingleRoad
	for i := range initialRoad {
		initialRoad[i] = make(Road, roadLength)
	}

	if trafficLightPos > 0 {
		// set the trafficLight locations
		for _, val := range trafficLightLane {
			initialRoad[val][trafficLightPos].kind = 3
		}
	}

	return initialRoad
}
