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
	height := len(r)
	width := len(r[0])

	c := canvas.CreateNewCanvas(3*width*scalingFactor/3, 2*height*scalingFactor)
	// c.SetFillColor(canvas.MakeColor(0, 0, 0))
	// c.ClearRect(0, 0, rows, cols*scalingFactor)
	// c.Fill()

	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			if r[i][j].kind == 0 {
				c.SetFillColor(canvas.MakeColor(255, 255, 255))
			} else if r[i][j].kind == 1 {
				if r[i][j].turninglight == -1 || r[i][j].turninglight == 1 {
					c.SetFillColor(canvas.MakeColor(0, 255, 0))
				} else {
					c.SetFillColor(canvas.MakeColor(244, 114, 208))
				}
			} else if r[i][j].kind == 2 {
				if r[i][j].turninglight == -1 || r[i][j].turninglight == 1 {
					c.SetFillColor(canvas.MakeColor(0, 255, 0))
				} else {
					c.SetFillColor(canvas.MakeColor(0, 0, 255))
				}
			} else if r[i][j].kind == 3 {
				c.SetFillColor(canvas.MakeColor(255, 0, 0)) // red light
			} else if r[i][j].kind == 4 {
				c.SetFillColor(canvas.MakeColor(255, 255, 0)) // yellow light
			} else if r[i][j].kind == 5 {
				c.SetFillColor(canvas.MakeColor(0, 255, 0)) // green light
			}

			// x1, y1 := i, scalingFactor*cols
			// x2, y2 := i+1, scalingFactor*(cols+1)
			x1, y1 := (3*j)*scalingFactor/3, 2*i*scalingFactor
			x2, y2 := (3*j+3)*scalingFactor/3, (2*i+2)*scalingFactor
			c.ClearRect(x1, y1, x2, y2)
			c.Fill()

			// c.SetFillColor(canvas.MakeColor(0, 255, 0))
			// if r[i][j].turninglight == -1 {
			// 	x1, y1 = (3*j)*scalingFactor/3, 2*i*scalingFactor
			// 	x2, y2 = (3*j+1)*scalingFactor/3, (2*i+1)*scalingFactor
			// 	c.ClearRect(x1, y1, x2, y2)
			// 	c.Fill()
			// } else if r[i][j].turninglight == 1 {
			// 	x1, y1 = (3*j)*scalingFactor/3, (2*i+1)*scalingFactor
			// 	x2, y2 = (3*j+1)*scalingFactor/3, (2*i+2)*scalingFactor
			// 	c.ClearRect(x1, y1, x2, y2)
			// 	c.Fill()
			// }

		}
	}

	return c.GetImage()

}

func DrawBoardSingle(t []MultiRoad, numGens int, filename string) {
	c := canvas.CreateNewCanvas(roadLength, numGens+1)
	c.SetLineWidth(1)

	laneNum := 1
	for i := 0; i < numGens+1; i++ {
		for k := 0; k < laneNum; k++ {
			for j := 0; j < roadLength; j++ {
				if t[i][0][j].kind == 0 {
					c.SetFillColor(canvas.MakeColor(255, 255, 255))
					drawSquare(c, j, i)
				} else if t[i][0][j].kind == 1 {
					if t[i][0][j].turninglight == -1 || t[i][0][j].turninglight == 1 {
						c.SetFillColor(canvas.MakeColor(0, 255, 0))
						drawSquare(c, j, i)
					} else {
						c.SetFillColor(canvas.MakeColor(244, 114, 208))
						drawSquare(c, j, i)
					}
				} else if t[i][0][j].kind == 2 {
					if t[i][0][j].turninglight == -1 || t[i][0][j].turninglight == 1 {
						c.SetFillColor(canvas.MakeColor(0, 255, 0))
						drawSquare(c, j, i)
					} else {
						c.SetFillColor(canvas.MakeColor(0, 0, 255))
						drawSquare(c, j, i)
					}
				} else if t[i][0][j].kind == 3 {
					c.SetFillColor(canvas.MakeColor(255, 0, 0)) // red light
					drawSquare(c, j, i)
				} else if t[i][0][j].kind == 4 {
					c.SetFillColor(canvas.MakeColor(255, 255, 0)) // yellow light
					drawSquare(c, j, i)
				} else if t[i][0][j].kind == 5 {
					c.SetFillColor(canvas.MakeColor(0, 255, 0)) // green light
					drawSquare(c, j, i)
				}
			}

		}
	}
	c.SaveToPNG(filename)
}

func DrawBoardMulti(t []MultiRoad, numGens, laneNum int, filename string) {
	width := (numGens + 1) / 20 * laneNum
	c := canvas.CreateNewCanvas(roadLength, width)
	c.SetLineWidth(1)

	for i := 0; i < (numGens+1)/20; i++ {
		for k := 0; k < laneNum; k++ {
			for j := 0; j < roadLength; j++ {
				if t[i][k][j].kind == 0 {
					c.SetFillColor(canvas.MakeColor(255, 255, 255))
					drawSquare(c, j, i*k)
				} else if t[i][k][j].kind == 1 {
					if t[i][k][j].turninglight == -1 || t[i][k][j].turninglight == 1 {
						c.SetFillColor(canvas.MakeColor(0, 255, 0))
						drawSquare(c, j, i*k)
					} else {
						c.SetFillColor(canvas.MakeColor(244, 114, 208))
						drawSquare(c, j, i*k)
					}
				} else if t[i][k][j].kind == 2 {
					if t[i][k][j].turninglight == -1 || t[i][k][j].turninglight == 1 {
						c.SetFillColor(canvas.MakeColor(0, 255, 0))
						drawSquare(c, j, i*k)
					} else {
						c.SetFillColor(canvas.MakeColor(0, 0, 255))
						drawSquare(c, j, i*k)
					}
				} else if t[i][k][j].kind == 3 {
					c.SetFillColor(canvas.MakeColor(255, 0, 0)) // red light
					drawSquare(c, j, i*k)
				} else if t[i][k][j].kind == 4 {
					c.SetFillColor(canvas.MakeColor(255, 255, 0)) // yellow light
					drawSquare(c, j, i*k)
				} else if t[i][k][j].kind == 5 {
					c.SetFillColor(canvas.MakeColor(0, 255, 0)) // green light
					drawSquare(c, j, i*k)
				}
			}
		}
	}
	c.SaveToPNG(filename)
}

func DrawPoint(a canvas.Canvas, r, c int) {
	a.ClearRect(r, c, r+1, c+1)
}
func drawSquare(a canvas.Canvas, r, c int) {
	x1, y1 := float64(r), float64(c)
	x2, y2 := float64(r+1), float64(c+1)
	a.MoveTo(x1, y1)
	a.LineTo(x1, y2)
	a.LineTo(x2, y2)
	a.LineTo(x2, y1)
	a.LineTo(x1, y1)
	a.Fill()
}
