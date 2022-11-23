package main

const roadLength = 1000
const laneNum = 5
const maxSpeed = 10
const dedicatedLane = 0

const p1 = 0.94
const p2 = 0.5
const p3 = 0.2
const p4 = 0.7
const cp1 = 0.8
const cp2 = 0.5
const cp3 = 0.1

type Car struct {
	speed        int
	kind         int
	backlight    int
	turninglight int
	accel        int
}

type Road []Car

type MultiRoad []Road

const k = 1

var safeSpaceMin = [maxSpeed + 1]int{1 * k, 5 * k, 10 * k, 15 * k, 20 * k, 25 * k, 30 * k, 35 * k, 40 * k, 45 * k, 50 * k}
var safeSpaceMax = [maxSpeed + 1]int{10 * k, 15 * k, 20 * k, 25 * k, 30 * k, 35 * k, 40 * k, 45 * k, 50 * k, 55 * k, 60 * k}
