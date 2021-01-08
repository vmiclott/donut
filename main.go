package main

import (
	"fmt"
	"math"
)

const (
	thetaSpacing float64 = 0.07
	phiSpacing   float64 = 0.02
	r1           float64 = 0.5
	r2           float64 = 1.0
	k2           float64 = 5.0
	screenWidth  float64 = 48
	screenHeight float64 = 48
	k1           float64 = screenWidth * k2 * 3 / (8 * (r1 + r2))
)

var brightnessSymbols [12]string = [12]string{".", ",", "-", "~", ":", ";", "=", "!", "*", "#", "$", "@"}

func renderFrame(a, b float64) {
	var (
		cosA float64 = math.Cos(a)
		sinA float64 = math.Sin(a)
		cosB float64 = math.Cos(b)
		sinB float64 = math.Sin(b)
	)

	output := make([][]string, int(screenWidth))
	for i := range output {
		output[i] = make([]string, int(screenHeight))
		for j := range output[i] {
			output[i][j] = " "
		}
	}

	zBuffer := make([][]float64, int(screenWidth))
	for i := range zBuffer {
		zBuffer[i] = make([]float64, int(screenHeight))
	}

	for theta := 0.0; theta < 2*math.Pi; theta += thetaSpacing {
		var (
			cosTheta float64 = math.Cos(theta)
			sinTheta float64 = math.Sin(theta)
			circleX  float64 = r2 + r1*cosTheta
			circleY  float64 = r1 * sinTheta
		)

		for phi := 0.0; phi < 2*math.Pi; phi += phiSpacing {
			var (
				cosPhi      float64 = math.Cos(phi)
				sinPhi      float64 = math.Sin(phi)
				x           float64 = circleX*(cosB*cosPhi+sinA*sinB*sinPhi) - circleY*cosA*sinB
				y           float64 = circleX*(sinB*cosPhi-sinA*cosB*sinPhi) + circleY*cosA*cosB
				z           float64 = k2 + cosA*circleX*sinPhi + circleY*sinA
				oneOverZ    float64 = 1 / z
				xProjection int     = int(screenWidth/2 + k1*oneOverZ*x)
				yProjection int     = int(screenHeight/2 - k1*oneOverZ*y)
				brightness  float64 = cosPhi*cosTheta*sinB - cosA*cosTheta*sinPhi - sinA*sinTheta + cosB*(cosA*sinTheta-cosTheta*sinA*sinPhi)
			)

			if brightness > 0 {
				if oneOverZ > zBuffer[xProjection][yProjection] {
					zBuffer[xProjection][yProjection] = oneOverZ
					brightnessIndex := int(brightness * 8)
					output[xProjection][yProjection] = brightnessSymbols[brightnessIndex]
				}
			}
		}
	}

	fmt.Print("\x1b[H\n")
	for j := 0; j < int(screenHeight); j++ {
		for i := 0; i < int(screenWidth); i++ {
			fmt.Print(output[i][j])
		}
		fmt.Println()
	}
}

func main() {
	const (
		aSpacing float64 = 0.08
		bSpacing float64 = 0.03
	)

	for a := 0.0; a < 2*math.Pi; a += aSpacing {
		for b := 0.0; b < 2*math.Pi; b += bSpacing {
			renderFrame(a, b)
		}
	}
}
