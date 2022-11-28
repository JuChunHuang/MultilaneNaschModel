package main

import (
	"canvas"
	"image"
)

func BoardsToImages(roads []MultiRoad, cellWidth int) []image.Image {
	imageList := make([]image.Image, len(roads))
	for i := range roads {
		imageList[i] = roads[i].BoardToImage(cellWidth)
		// break
	}
	return imageList
}

// BoardToImage converts a GameBoard to an image, in which
// each cell has a cell width given by a parameter
func (r MultiRoad) BoardToImage(scalingFactor int) image.Image {
	height := len(r)
	width := len(r[0])

	c := canvas.CreateNewCanvas(width*scalingFactor/4, height*scalingFactor)
	// c.SetFillColor(canvas.MakeColor(0, 0, 0))
	// c.ClearRect(0, 0, rows, cols*scalingFactor)
	// c.Fill()

	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			if r[i][j].kind == 0 {
				c.SetFillColor(canvas.MakeColor(255, 255, 255))
			} else if r[i][j].kind == 1 {
				c.SetFillColor(canvas.MakeColor(244, 114, 208))
			} else if r[i][j].kind == 2 {
				c.SetFillColor(canvas.MakeColor(0, 0, 255))
			} else if r[i][j].kind == 3 {
				c.SetFillColor(canvas.MakeColor(255, 0, 0)) // red light
			} else if r[i][j].kind == 4 {
				c.SetFillColor(canvas.MakeColor(255, 255, 0)) // yellow light
			} else if r[i][j].kind == 5 {
				c.SetFillColor(canvas.MakeColor(0, 255, 0)) // green light
			}

			// x1, y1 := i, scalingFactor*cols
			// x2, y2 := i+1, scalingFactor*(cols+1)
			x1, y1 := j*scalingFactor/4, i*scalingFactor
			x2, y2 := (j+1)*scalingFactor/4, (i+1)*scalingFactor

			c.ClearRect(x1, y1, x2, y2)

			c.Fill()
		}
	}

	return c.GetImage()

}