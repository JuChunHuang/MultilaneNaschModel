package main

import (
	"canvas"
	"image"
)

func BoardsToImages(roads []MultiRoad, cellWidth int) []image.Image {
	imageList := make([]image.Image, len(roads))
	for i := range roads {
		imageList[i] = roads[i].BoardToImage(cellWidth)
	}
	return imageList
}

// BoardToImage converts a GameBoard to an image, in which
// each cell has a cell width given by a parameter
func (r MultiRoad) BoardToImage(scalingFactor int) image.Image {
	rows := len(r[0])
	cols := len(r)

	c := canvas.CreateNewCanvas(rows, cols)
	c.SetFillColor(canvas.MakeColor(255, 255, 255))
	c.ClearRect(0, 0, rows, cols)
	c.Fill()

	for j := 0; j < cols; j++ {
		for i := 0; i < rows; i++ {
			if r[j][i].kind == 0 {
				c.SetFillColor(canvas.MakeColor(0, 0, 0))
			} else if r[j][i].kind == 1 {
				c.SetFillColor(canvas.MakeColor(100, 10, 90))
			} else if r[j][i].kind == 2 {
				c.SetFillColor(canvas.MakeColor(10, 180, 60))
			} else if r[j][i].kind == 3 {
				c.SetFillColor(canvas.MakeColor(0, 255, 0))
			} else if r[j][i].kind == 4 {
				c.SetFillColor(canvas.MakeColor(0, 255, 255)) //应该是黄色,rgb是我猜的
			} else if r[j][i].kind == 4 {
				c.SetFillColor(canvas.MakeColor(255, 0, 0)) //应该是黄色
			}

			x1, y1 := i, scalingFactor*cols
			x2, y2 := i+1, scalingFactor*(cols+1)

			c.ClearRect(x1, y1, x2, y2)

			c.Fill()
		}

	}

	return c.GetImage()

}
