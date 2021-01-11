package main

import (
	"fmt"
	"math"
)

// screen with 2d projection of a 3d object
type screen struct {
	width            int
	height           int
	screen           [][]string
	zBuffer          [][]float64
	distanceToEye    float64
	distanceToObject float64
}

// newScreen initialises a screen with given width, heigh and tunable distance parameters
func newScreen(width, height int, distanceToEye, distanceToObject float64) *screen {
	screen := screen{
		width:            width,
		height:           height,
		screen:           make([][]string, width),
		zBuffer:          make([][]float64, width),
		distanceToEye:    distanceToEye,
		distanceToObject: distanceToObject,
	}
	for i := range screen.screen {
		screen.screen[i] = make([]string, height)
		for j := range screen.screen[i] {
			screen.screen[i][j] = " "
		}
	}
	for i := range screen.zBuffer {
		screen.zBuffer[i] = make([]float64, height)
	}
	return &screen
}

// render prints the screen on the terminal
func (s screen) render() {
	fmt.Print("\x1b[H\n")
	for j := 0; j < int(s.height); j++ {
		for i := 0; i < int(s.width); i++ {
			fmt.Print(s.screen[i][j])
		}
		fmt.Println()
	}
}

// clear clears the screens points and the zBuffer
func (s screen) clear() {
	for i := range s.screen {
		s.screen[i] = make([]string, s.height)
		for j := range s.screen[i] {
			s.screen[i][j] = " "
		}
	}
	for i := range s.zBuffer {
		s.zBuffer[i] = make([]float64, s.height)
	}
}

// vector3d represents a 3-dimensional vector
type vector3d struct {
	x float64
	y float64
	z float64
}

// dotProduct calculates the dot product of two 3-dimensional vectors
func (v vector3d) dotProduct(w vector3d) float64 {
	return v.x*w.x + v.y*w.y + v.z*w.z
}

// point3d represents a 3-dimensional point belonging to a solid, along with the normal in that point
type point3d struct {
	x      float64
	y      float64
	z      float64
	normal vector3d
}

// rotation3d represents a 3-dimensional rotation
type rotation3d [3][3]float64

// compose composes 2 rotations into one
func (a rotation3d) compose(b rotation3d) rotation3d {
	return rotation3d{
		{
			a[0][0]*b[0][0] + a[0][1]*b[1][0] + a[0][2]*b[2][0],
			a[0][0]*b[0][1] + a[0][1]*b[1][1] + a[0][2]*b[2][1],
			a[0][0]*b[0][2] + a[0][1]*b[1][2] + a[0][2]*b[2][2],
		},
		{
			a[1][0]*b[0][0] + a[1][1]*b[1][0] + a[1][2]*b[2][0],
			a[1][0]*b[0][1] + a[1][1]*b[1][1] + a[1][2]*b[2][1],
			a[1][0]*b[0][2] + a[1][1]*b[1][2] + a[1][2]*b[2][2],
		},
		{
			a[2][0]*b[0][0] + a[2][1]*b[1][0] + a[2][2]*b[2][0],
			a[2][0]*b[0][1] + a[2][1]*b[1][1] + a[2][2]*b[2][1],
			a[2][0]*b[0][2] + a[2][1]*b[1][2] + a[2][2]*b[2][2],
		},
	}
}

// rotate rotates a 3-dimensional point and its normal vector using a given rotation
func (p point3d) rotate(rotation rotation3d) point3d {
	return point3d{
		x: rotation[0][0]*p.x + rotation[1][0]*p.y + rotation[2][0]*p.z,
		y: rotation[0][1]*p.x + rotation[1][1]*p.y + rotation[2][1]*p.z,
		z: rotation[0][2]*p.x + rotation[1][2]*p.y + rotation[2][2]*p.z,
		normal: vector3d{
			x: rotation[0][0]*p.normal.x + rotation[1][0]*p.normal.y + rotation[2][0]*p.normal.z,
			y: rotation[0][1]*p.normal.x + rotation[1][1]*p.normal.y + rotation[2][1]*p.normal.z,
			z: rotation[0][2]*p.normal.x + rotation[1][2]*p.normal.y + rotation[2][2]*p.normal.z,
		},
	}
}

// brightness calculates the brightness of a 3-dimensional point for a given light source
func (p point3d) brightness(lightSource vector3d) (string, error) {
	var brightnessSymbols [13]string = [13]string{".", ",", "-", "~", ":", ";", "=", "!", "*", "#", "$", "@", "@"}

	brightness := p.normal.dotProduct(lightSource)
	if brightness < 0 {
		return "", fmt.Errorf("negative brightness")
	}
	brightnessIndex := int(12 * brightness)
	return brightnessSymbols[brightnessIndex], nil
}

// addToScreen adds a point to a given screen for a given light source
func (p point3d) addToScreen(screen *screen, lightSource vector3d) {
	eyeScreenDist := screen.distanceToEye
	screenObjectDist := screen.distanceToObject
	xProjection := screen.width/2 + int(eyeScreenDist*p.x/(p.z+screenObjectDist))
	yProjection := screen.height/2 - int(eyeScreenDist*p.y/(p.z+screenObjectDist))
	if 1/(p.z+screenObjectDist) <= screen.zBuffer[xProjection][yProjection] {
		// skip points behind other points
		return
	}
	screen.zBuffer[xProjection][yProjection] = 1 / (p.z + screenObjectDist)
	brightness, err := p.brightness(lightSource)
	if err != nil {
		// skip invisible points
		return
	}
	screen.screen[xProjection][yProjection] = brightness
}

// createDonut creates a discrete, static donut for given radii r1 and r2, theta spacing and phi spacing
func createDonut(r1, r2, thetaSpacing, phiSpacing float64) (points []point3d) {
	for theta := 0.0; theta < 2*math.Pi; theta += thetaSpacing {
		var (
			cosTheta float64 = math.Cos(theta)
			sinTheta float64 = math.Sin(theta)
			circleX  float64 = r2 + r1*cosTheta
			circleY  float64 = r1 * sinTheta
		)

		for phi := 0.0; phi < 2*math.Pi; phi += phiSpacing {
			var (
				cosPhi float64 = math.Cos(phi)
				sinPhi float64 = math.Sin(phi)
				x      float64 = circleX * cosPhi
				y      float64 = circleY
				z      float64 = circleX * sinPhi
				xn     float64 = cosTheta * cosPhi
				yn     float64 = sinTheta
				zn     float64 = cosTheta * sinPhi
			)
			points = append(
				points, point3d{
					x: x,
					y: y,
					z: z,
					normal: vector3d{
						x: xn,
						y: yn,
						z: zn,
					},
				},
			)
		}
	}

	return points
}

// main renders a rotating donut
func main() {
	const (
		aSpacing float64 = 0.08
		bSpacing float64 = 0.03
	)

	screen := newScreen(48, 48, 60, 5)
	donut := createDonut(0.5, 1.0, 0.07, 0.02)
	lightSource := vector3d{
		x: 0,
		y: 1 / math.Sqrt(2),
		z: -1 / math.Sqrt(2),
	}
	for _, point := range donut {
		point.addToScreen(screen, lightSource)
	}
	screen.render()
	screen.clear()
	for a := 0.0; a < 2*math.Pi; a += aSpacing {
		rotationA := rotation3d{
			{1, 0, 0},
			{0, math.Cos(a), math.Sin(a)},
			{0, -math.Sin(a), math.Cos(a)},
		}
		for b := 0.0; b < 2*math.Pi; b += bSpacing {
			rotationB := rotation3d{
				{math.Cos(b), math.Sin(b), 0},
				{-math.Sin(b), math.Cos(b), 0},
				{0, 0, 1},
			}
			rotation := rotationA.compose(rotationB)
			for _, point := range donut {
				newPoint := point.rotate(rotation)
				newPoint.addToScreen(screen, lightSource)
			}
			screen.render()
			screen.clear()
		}
	}
}
