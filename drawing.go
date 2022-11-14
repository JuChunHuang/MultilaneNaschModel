package main

import (
	"canvas"
	"image"
)

func BoardsToImages(roads []Road, cellWidth int) []image.Image {
	imageList := make([]image.Image, len(roads))
	for i := range roads {
		imageList[i] = roads[i].BoardToImage(cellWidth)
	}
	return imageList
}

//BoardToImage converts a GameBoard to an image, in which
//each cell has a cell width given by a parameter
func (r Road) BoardToImage(scalingFactor int) image.Image {
	rows := len(r)
	cols := 1

	c := canvas.CreateNewCanvas(rows, scalingFactor*5)
	c.SetFillColor(canvas.MakeColor(255, 255, 255))
	c.ClearRect(0, 0, rows, scalingFactor*5)
	c.Fill()

	for i := 0; i < rows; i++ {
		if r[i].kind == 0 {
			c.SetFillColor(canvas.MakeColor(0, 0, 0))
		} else if r[i].kind == 1 {
			c.SetFillColor(canvas.MakeColor(255, 0, 0))
		} else if r[i].kind == 2 {
			c.SetFillColor(canvas.MakeColor(0, 255, 0))
		}

		x1, y1 := i, scalingFactor*cols
		x2, y2 := i+1, scalingFactor*(cols+1)

		c.ClearRect(x1, y1, x2, y2)

		c.Fill()
	}
	return c.GetImage()

}
