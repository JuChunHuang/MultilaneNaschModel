package main

import (
	"fmt"
	"testing"
)

func TestCheckPreviousTrain(t *testing.T) {
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

	if CheckPreviousTrain(road, 1) == 2 {
		fmt.Println("Pass the CheckPreviousTrain test!")
	} else {
		t.Error("Can't pass the CheckPreviousTrain test!")
	}
}

func TestCheckNextTrain(t *testing.T) {
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

	if CheckNextTrain(road, 4) == 2 {
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
