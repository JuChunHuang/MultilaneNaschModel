package main

const roadLength = 1000
const laneNum = 1
const maxSpeed = 10

// const car_num = 10
const p1 = 0.94
const p2 = 0.5
const p3 = 0.2

type Car struct {
	speed int
	kind  int
	light int
	accel int
}

type Road []Car

var safeSpaceMin = [maxSpeed + 1]int{1, 5, 10, 15, 20, 25, 30, 35, 40, 45, 50}
var safeSpaceMax = [maxSpeed + 1]int{10, 15, 20, 25, 30, 35, 40, 45, 50, 55, 60}
