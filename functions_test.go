package main

import (
	"fmt"
	"testing"
)

func TestCheckPreviousTrain(t *testing.T) {
	//change roadLength = 8
	var road Road

	var a Car
	a.kind = 0

	var b Car
	b.kind = 2
	b.speed = 1

	var c Car
	c.kind = 0

	var d Car
	d.kind = 2
	d.speed = 3

	var e Car
	e.kind = 2
	e.speed = 4

	var emp Car
	emp.kind = 0

	road = append(road, a, d, emp, emp, emp, emp, emp, e)
	road[1] = d
	road[7] = e

	if CheckPreviousTrain(road, 1) == 1 {
		fmt.Print(GetSDVmindis(1, 7, road))
		fmt.Println("Pass the CheckPreviousTrain test!")
	} else {
		t.Error("Can't pass the CheckPreviousTrain test!")
	}
}

func TestCheckNextTrain(t *testing.T) {
	//change roadLength = 8
	var road Road

	var a Car
	a.kind = 0

	var b Car
	b.kind = 2
	b.speed = 1

	var c Car
	c.kind = 0

	var d Car
	d.kind = 2
	d.speed = 3

	var e Car
	e.kind = 2
	e.speed = 4

	var emp Car
	emp.kind = 0

	road = append(road, a, d, emp, emp, emp, emp, emp, e)
	road[1] = d
	road[7] = e

	if CheckNextTrain(road, 7) == 1 {
		// fmt.Print(GetSDVmindis(1, 7, road))
		fmt.Println("Pass the CheckPreviousTrain test!")
	} else {
		t.Error("Can't pass the CheckPreviousTrain test!")
	}
}

func TestCheckChain(t *testing.T) {
	//change roadLength = 5
	var road Road

	var a Car
	a.kind = 0

	var b Car
	b.kind = 2

	var c Car
	c.kind = 0

	var d Car
	d.kind = 2

	var e Car
	e.kind = 2

	road = append(road, a, b, c, d, e)

	if CheckTrain(road, 1) == true {
		fmt.Println("pass the test for second car!")
	} else if CheckTrain(road, 3) == true {
		fmt.Println("pass the test for third car!")
	} else if CheckTrain(road, 4) == true {
		fmt.Println("pass the test for fourth car!")
	} else {
		t.Error("Can't pass the test!")
	}

}

func TestGetTrainHead(t *testing.T) {
	//change roadLength = 5
	var road Road

	var a Car
	a.kind = 0

	var b Car
	b.kind = 2

	var c Car
	c.kind = 0

	var d Car
	d.kind = 2

	var e Car
	e.kind = 2

	road = append(road, a, b, c, d, e)

	if GetTrainHead(road, 1) == 4 {
		fmt.Println("pass the test for second car!")
	} else if GetTrainHead(road, 3) == 4 {
		fmt.Println("pass the test for fourth car!")
	} else if GetTrainHead(road, 4) == 4 {
		fmt.Println("pass the test for fifth car!")
	} else {
		t.Error("Can't pass the test!")
	}

}

func SDVchangeLane(t *testing.T) {
	var a1, a2, a3, a4, a5 Car
	a1.kind = 0
	a2.kind = 0
	a3.kind = 0
	a4.kind = 0
	a5.kind = 0

	var b1, b2, b3, b4, b5 Car
	b1.kind = 0
	b2.kind = 2
	b3.kind = 0
	b4.kind = 0
	b5.kind = 0

}
