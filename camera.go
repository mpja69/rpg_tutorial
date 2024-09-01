package main

type Camera struct {
	X float64
	Y float64
}

func NewCamera(x, y float64) *Camera {
	return &Camera{
		X: x,
		Y: y,
	}
}

func (c *Camera) followTarget(targetX, targetY, screenWidth, screenHeight float64) {
	c.X = -targetX + screenWidth/2.0
	c.Y = -targetY + screenHeight/2.0

}

// vim: ts=4
