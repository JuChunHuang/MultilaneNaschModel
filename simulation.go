package main

import "math/rand"

func PlayNaschModel(initial_road Road, num_gens, time int) []Road {
	roads := make([]Road, num_gens+1)
	roads[0] = initial_road
	for i := 1; i <= num_gens; i++ {
		roads[i] = UpdateRoad(roads[i-1], time)
	}

	return roads
}

// calculate the next position of each car
func CarNextStep(current_road Road) Road{
	var new_road
	var prob_of_decel float64

	for c := range current_road{
		cur := current_road[c]
		if cur.kind == 1{
			// the car is a NSDV, change the speed of the car
			pre_car_index := GetPrev(current_road, c)
			prev := current_road[pre_car_index]
			delta_d := pre_car_index - c
			if prev.light == 1 && delta_d > safe_space_min[cur.speed] && delta_d < safe_space_max[cur.speed]{
				prob_of_decel = p1
			}else if prev.light == 0 && delta_d > safe_space_min[cur.speed] && delta_d < safe_space_max[cur.speed]{
				prob_of_decel = p2
			}else if cur.speed == 0{
				prob_of_decel = p3
			}else{
				prob_of_decel = 0
			}
			thres_to_decel := rand.Float64()

			if cur.speed < max_speed && delta_d > safe_space_max[cur.speed]{
				// acceleration because no car in front of it
				cur.speed += 1
				cur.light = 1
			}else if prev.light == 1 && delta_d > safe_space_min[cur.speed]{
				// acceleration because the front car is accelerated
				cur.speed += 1
				cur.light = 1
			}else if prob_of_decel >= thres_to_decel || delta_d < safe_space_min[cur.speed]{
				// deceleration
				cur.speed -= 1
				cur.light = -1
			}else if delta_d == safe_space_min[cur.speed]{
				// on hold case, speed not changed
				cur.light = 0
			}
		}
	}

}

func GetPrev(current_road Road, index_of_car int) int {
	for c := index_of_car + 1; c < road_length; c++ {
		if current_road[c].kind != 0 {
			return c
		}
	}

	return 2 * road_length
}

func GetNext(current_road Road, index_of_car int) int {
	for c := index_of_car - 1; c <= 0; c-- {
		if current_road[c].kind != 0 {
			return c
		}
	}

	return 0
}

func Produce(current_road Road) bool {
	var kind int
	kind = 1
	if current_road[0].kind == 0 {
		init_speed_bound := 0
		prev_car := GetPrev(current_road, 0)

		// determin the speed range that the car can obtain
		if prev_car > road_length {
			// if no car before
			init_speed_bound = max_speed
		} else if kind == 1 {
			// if the new car is a NSDV
			for i := 0; i < max_speed; i++ {
				if prev_car < safe_space_min[i] {
					break
				}
				init_speed_bound = i
			}
		} else if current_road[prev_car].kind == 1 {
			// if the new car is a SDV and prev_car is a NSDV
			for i := 0; i < max_speed; i++ {
				if prev_car < safe_space_min[i] {
					break
				}
				init_speed_bound = i
			}
		} else if current_road[prev_car].kind == 2 {
			// if the new car is a SDV and prev_car is a SDV
			for i := 0; i < max_speed; i++ {
				min_space := safe_space_min[i] - safe_space_min[current_road[prev_car].speed] + 2*current_road[prev_car].speed + 1
				if prev_car < min_space {
					break
				}
				init_speed_bound = i
			}
		}

		if init_speed_bound <= 0 {
			// no car produced
			return false
		} else {
			current_road[0].speed = 1 + rand.Intn(init_speed_bound)
			current_road[0].kind = kind
			current_road[0].light = false
			current_road[0].accel = false
		}
	}

	return true
}

//HelloWorld!
