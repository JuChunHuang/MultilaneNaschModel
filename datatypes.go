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

const k = 1

var safeSpaceMin = [maxSpeed + 1]int{1 * k, 5 * k, 10 * k, 15 * k, 20 * k, 25 * k, 30 * k, 35 * k, 40 * k, 45 * k, 50 * k}
var safeSpaceMax = [maxSpeed + 1]int{10 * k, 15 * k, 20 * k, 25 * k, 30 * k, 35 * k, 40 * k, 45 * k, 50 * k, 55 * k, 60 * k}
