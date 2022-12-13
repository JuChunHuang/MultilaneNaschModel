package main

import "C"

const roadLength = 50
const maxSpeed = 10

const p1 = 0.8
const p2 = 0.4
const p3 = 0.1
const cp1 = 0.8
const cp2 = 0.4
const cp3 = 0.1

type Car struct {
	speed        int
	kind         int
	backlight    int
	turninglight int
	accel        int
}

// kind = 0  <-------> empty road
// kind = 1  <-------> NSDV
// kind = 2  <-------> SDV
// kind = 3  <-------> red traffic lights
// kind = 4  <-------> yellow traffic lights
// kind = 5  <-------> green traffic lights
// backlight = 0  <-------> remain same speed
// backlight = -1 <-------> deceleration
// backlight = 1  <-------> acceleration
// turninglight = 0  <-------> no turning
// turninglight = -1 <-------> turning left
// turninglight = 1  <-------> turning right
// accel = 0  <-------> no acceleration
// accel = 1  <-------> acceleration = 1

type Road []Car

type MultiRoad []Road

const k = 1

var safeSpaceMin = [maxSpeed + 1]int{1 * k, 3 * k, 6 * k, 9 * k, 12 * k, 15 * k, 18 * k, 21 * k, 24 * k, 27 * k, 30 * k}
var safeSpaceMax = [maxSpeed + 1]int{5 * k, 10 * k, 15 * k, 20 * k, 25 * k, 30 * k, 35 * k, 40 * k, 45 * k, 50 * k, 55 * k}
var safetraffic = [maxSpeed + 1]int{1 * k, 3 * k, 6 * k, 9 * k, 12 * k, 15 * k, 18 * k, 21 * k, 24 * k, 27 * k, 30 * k}

var safeSpaceSDVMin = [maxSpeed + 1]int{1 * k, 2 * k, 4 * k, 6 * k, 8 * k, 10 * k, 12 * k, 14 * k, 16 * k, 18 * k, 20 * k}
