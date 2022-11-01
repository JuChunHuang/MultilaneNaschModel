package main

const road_length = 1000
const lane_num = 1
const max_speed = 10

// const car_num = 10
const p1 = 0.94
const p2 = 0.5
const p3 = 0.2

type Car struct {
	speed int
	kind  int
	light int
	accel bool
}

type Road [road_length]Car

var safe_space_min = [max_speed]int{1, 5, 10, 15, 20, 25, 30, 35, 40, 45}
var safe_space_max = [max_speed]int{10, 15, 20, 25, 30, 35, 40, 45, 50, 55}
