package main

const roadLength = 1000
const laneNum = 5
const maxSpeed = 10
const dedicatedLane = 0

const p1 = 0.94
const p2 = 0.5
const p3 = 0.2
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

const k = 2

var safeSpaceMin = [maxSpeed + 1]int{1 * k, 5 * k, 10 * k, 15 * k, 20 * k, 25 * k, 30 * k, 35 * k, 40 * k, 45 * k, 50 * k}
var safeSpaceMax = [maxSpeed + 1]int{10 * k, 15 * k, 20 * k, 25 * k, 30 * k, 35 * k, 40 * k, 45 * k, 50 * k, 55 * k, 60 * k}
var safetraffic = [maxSpeed + 1]int{10 * k, 15 * k, 20 * k, 25 * k, 30 * k, 35 * k, 40 * k, 45 * k, 50 * k, 55 * k, 60 * k}

const k1 = 1

var safeSpaceSDVMin = [maxSpeed + 1]int{1 * k1, 2 * k1, 3 * k1, 4 * k1, 5 * k1, 6 * k1, 7 * k1, 8 * k1, 9 * k1, 10 * k1, 11 * k1}
