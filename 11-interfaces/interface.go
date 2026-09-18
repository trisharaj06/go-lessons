package main

import "fmt"

type geometry interface {
	area() float64
	peri() float64
}

// * interface embedding
type geometry3D interface {
	geometry
	volume() float64
}

type rect struct {
	height, width float64
}

// * struct embedding
type cuboid struct {
	rect
	bredth float64
}

func area(r rect) float64 {
	return r.height * r.width
}

func peri(r rect) float64 {
	return (2 * r.height) + (2 * r.width)
}

func cuboid_area(c cuboid) float64 {
	return c.height * c.width * c.bredth
}

func main() {

	r := rect{
		height: 4,
		width:  5,
	}

	c := cuboid{
		height: 4,
		width:  5,
		bredth: 6,
	}
	fmt.Println(area(r))
	fmt.Println(peri(r))

	fmt.Println(cuboid_area(c))
}
