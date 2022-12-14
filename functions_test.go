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

	if CheckTrain(road, 1) >= 3 {
		fmt.Println("pass the test for second car!")
	} else if CheckTrain(road, 3) >= 3 {
		fmt.Println("pass the test for third car!")
	} else if CheckTrain(road, 4) >= 3 {
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

func TestGetPrevCar(t *testing.T) {
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
	if GetPrevCar(road, 1) == 7 {
		fmt.Println("pass the test for GetPrevCar!")
	} else if CheckTrain(road, 3) >= 3 {
		fmt.Println("Do not pass the test for GetPrevCar!")
	}
}

func TestGetPrevLight(t *testing.T) {
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
	e.kind = 4

	var emp Car
	emp.kind = 0

	road = append(road, a, d, emp, emp, emp, emp, emp, e)
	if GetPrevCar(road, 1) == 7 {
		fmt.Println("pass the test for GetPrevLight!")
	} else if CheckTrain(road, 3) >= 3 {
		fmt.Println("Do not pass the test for GetPrevLight!")
	}
}
